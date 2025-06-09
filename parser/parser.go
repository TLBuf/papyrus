// Package parser defines a Papyrus parser.
package parser

import (
	"bytes"
	"fmt"
	"iter"
	"strconv"
	"strings"

	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/lexer"
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/token"
	"github.com/TLBuf/papyrus/types"
)

// Option defines an option to configure how parsing is performed.
type Option interface{ apply(*parser) }

type option func(*parser)

// apply implements the [Option] interface.
func (o option) apply(p *parser) {
	o(p)
}

// WithLooseComments controls block and line (i.e. loose) comment processing.
//
// If enabled, loose comments will be attached to the appropriate [ast.Trivia]
// on returned nodes. This is only required when the nodes may need to be
// written back out as source, e.g. when formatting.
func WithLooseComments(enabled bool) Option {
	return option(func(p *parser) {
		p.attachLooseComments = enabled
	})
}

// WithRecovery controls parsing should attempt error recovery and potentially
// include [ast.Error] nodes in the resulting AST.
//
// If enabled, the parser will attempt error recovery if an issue is found and
// instead of immediately failing, it will try to emit an error node instead. It
// is the responsibility of the caller to check for the presensce of [ast.Error]
// nodes.
//
// Enabling this does not guarantee that parsing will never fail with an
// [Error].
func WithRecovery(enabled bool) Option {
	return option(func(p *parser) {
		p.attemptRecovery = enabled
	})
}

// Parser returns the file parsed as an [*ast.Script] or an [Error] if parsing
// encountered one or more issues.
func Parse(file *source.File, opts ...Option) (*ast.Script, error) {
	stream, err := lexer.Lex(file)
	if err != nil {
		return nil, Error{
			Err:      err,
			Location: err.(lexer.Error).Location,
		}
	}
	streamNext, stop := iter.Pull2(stream.All())
	defer stop()
	p := &parser{
		streamNext:          streamNext,
		attachLooseComments: false,
		attemptRecovery:     false,
		prefix:              make(map[token.Kind]prefixParser),
		infix:               make(map[token.Kind]infixParser),
	}
	for _, opt := range opts {
		opt.apply(p)
	}

	registerPrefix(p, p.ParseBoolLiteral, token.True, token.False)
	registerPrefix(p, p.ParseFloatLiteral, token.FloatLiteral)
	registerPrefix(p, p.ParseIdentifier, token.Identifier, token.Self, token.Parent)
	registerPrefix(p, p.ParseIntLiteral, token.IntLiteral)
	registerPrefix(p, p.ParseNoneLiteral, token.None)
	registerPrefix(p, p.ParseParenthetical, token.ParenthesisOpen)
	registerPrefix(p, p.ParseStringLiteral, token.StringLiteral)
	registerPrefix(p, p.ParseUnary, token.Minus, token.LogicalNot)

	registerInfix(p, p.ParseAccess, token.Dot)
	registerInfix(p,
		p.ParseBinary,
		token.LogicalOr,
		token.LogicalAnd,
		token.Equal,
		token.NotEqual,
		token.Greater,
		token.GreaterOrEqual,
		token.Less,
		token.LessOrEqual,
		token.Plus,
		token.Minus,
		token.Divide,
		token.Multiply,
		token.Modulo)
	registerInfix(p, p.ParseCall, token.ParenthesisOpen)
	registerInfix(p, p.ParseCast, token.As)
	registerInfix(p, p.ParseIndex, token.BracketOpen)

	if err := p.next(); err != nil {
		return nil, err
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	script, err := p.ParseScript()
	if err != nil {
		return nil, err
	}

	if p.attachLooseComments {
		if err := attachLooseComments(script, p.comments); err != nil {
			return nil, err
		}
	}

	return script, nil
}

type parser struct {
	streamNext func() (int, token.Token, bool)

	token     *ast.Token
	lookahead *ast.Token

	attachLooseComments bool
	comments            []ast.LooseComment

	attemptRecovery bool
	recovery        bool
	errors          []ast.Error

	prefix map[token.Kind]prefixParser
	infix  map[token.Kind]infixParser
}

// Operator precedence.
const (
	_ int = iota
	lowest
	logicalOr      // ||
	logicalAnd     // &&
	comparison     // ==, !=, >, >=, <, <=
	additive       // +, -
	multiplicitive // *, /, %
	prefix         // -x or !y
	cast           // x As y
	access         // x.y
	call           // x(y)
	index          // x[y]
)

var precedence = map[token.Kind]int{
	token.LogicalOr:       logicalOr,
	token.LogicalAnd:      logicalAnd,
	token.Equal:           comparison,
	token.NotEqual:        comparison,
	token.Greater:         comparison,
	token.GreaterOrEqual:  comparison,
	token.Less:            comparison,
	token.LessOrEqual:     comparison,
	token.Plus:            additive,
	token.Minus:           additive,
	token.Multiply:        multiplicitive,
	token.Divide:          multiplicitive,
	token.Modulo:          multiplicitive,
	token.As:              cast,
	token.Dot:             access,
	token.ParenthesisOpen: call,
	token.BracketOpen:     index,
}

func precedenceOf(t token.Kind) int {
	if p, ok := precedence[t]; ok {
		return p
	}
	return lowest
}

type (
	prefixParser func() (ast.Expression, error)
	infixParser  func(ast.Expression) (ast.Expression, error)
)

// next advances token and lookahead by one
// token while skipping loose comment tokens.
func (p *parser) next() error {
	p.token = p.lookahead
	_, t, ok := p.streamNext()
	if !ok {
		return nil
	}
	p.lookahead = tok(t)
	if p.token == nil {
		return nil
	}
	// Consume loose comments immediately so the rest of the
	// parser never has to deal with them directly.
	if p.token.Kind == token.Semicolon {
		tok, err := p.ParseLineComment()
		if err != nil {
			return err
		}
		if p.attachLooseComments {
			p.comments = append(p.comments, tok)
		}
	}
	if p.token.Kind == token.BlockCommentOpen {
		tok, err := p.ParseBlockComment()
		if err != nil {
			return err
		}
		if p.attachLooseComments {
			p.comments = append(p.comments, tok)
		}
	}
	return nil
}

// tryConsume advances the token position if the current token matches the given
// token type or returns an error.
func (p *parser) tryConsume(t token.Kind, alts ...token.Kind) error {
	if p.token.Kind == t {
		return p.next()
	}
	for _, t := range alts {
		if p.token.Kind == t {
			return p.next()
		}
	}
	return unexpectedTokenError(p.token, t, alts...)
}

func unexpectedTokenError(got *ast.Token, want token.Kind, alts ...token.Kind) error {
	if len(alts) > 0 {
		return newError(got.Location, "expected any of [%s, %s], but found: %s", want, tokensTypesToString(alts...), got.Kind)
	}
	return newError(got.Location, "expected: %s, but found: %s", want, got.Kind)
}

func tokensTypesToString(types ...token.Kind) string {
	if len(types) == 0 {
		return ""
	}
	if len(types) == 1 {
		return types[0].String()
	}
	strs := make([]string, len(types))
	for i, t := range types {
		strs[i] = t.String()
	}
	return strings.Join(strs, ", ")
}

// consumeNewlines advances the token position through the as many newlines as
// possible until a non-newline token is found.
func (p *parser) consumeNewlines() error {
	for p.token.Kind == token.Newline {
		if err := p.next(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) ParseDocComment() (*ast.DocComment, error) {
	node := &ast.DocComment{
		Open: p.token,
	}
	if err := p.tryConsume(token.BraceOpen); err != nil {
		return nil, err
	}
	node.Text = p.token
	if err := p.tryConsume(token.Comment); err != nil {
		return nil, err
	}
	node.Close = p.token
	if err := p.tryConsume(token.BraceClose); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseBlockComment() (*ast.BlockComment, error) {
	node := &ast.BlockComment{
		Open: p.token,
	}
	if err := p.tryConsume(token.BlockCommentOpen); err != nil {
		return nil, err
	}
	node.Text = p.token
	if err := p.tryConsume(token.Comment); err != nil {
		return nil, err
	}
	node.Close = p.token
	if err := p.tryConsume(token.BlockCommentClose); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseLineComment() (*ast.LineComment, error) {
	node := &ast.LineComment{
		Open: p.token,
	}
	if err := p.tryConsume(token.Semicolon); err != nil {
		return nil, err
	}
	node.Text = p.token
	if err := p.tryConsume(token.Comment); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseScript() (*ast.Script, error) {
	var err error
	node := &ast.Script{
		Location: source.Location{
			File:        p.token.Location.File,
			Length:      len(p.token.Location.File.Text),
			StartLine:   1,
			StartColumn: 1,
		},
	}
	node.Keyword = p.token
	if err := p.tryConsume(token.ScriptName); err != nil {
		return nil, err
	}
	if node.Name, err = p.ParseIdentifier(); err != nil {
		return nil, err
	}
	if p.token.Kind == token.Extends {
		node.Extends = p.token
		if err := p.next(); err != nil {
			return nil, err
		}
		if node.Parent, err = p.ParseIdentifier(); err != nil {
			return nil, err
		}
	}
	for p.token.Kind == token.Hidden || p.token.Kind == token.Conditional {
		if p.token.Kind == token.Hidden {
			node.Hidden = append(node.Hidden, p.token)
		} else {
			node.Conditional = append(node.Conditional, p.token)
		}
		if err := p.next(); err != nil {
			return nil, err
		}
	}
	if err := p.tryConsume(token.Newline, token.EOF); err != nil {
		return nil, err
	}
	if p.token.Kind == token.BraceOpen {
		node.Comment, err = p.ParseDocComment()
		if err != nil {
			return nil, err
		}
	}
	if err := p.consumeNewlines(); err != nil {
		return nil, err
	}
	for p.token.Kind != token.EOF {
		stmt, err := p.ParseScriptStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			node.Statements = append(node.Statements, stmt)
		}
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
	}
	return node, nil
}

func (p *parser) ParseScriptStatement() (ast.ScriptStatement, error) {
	start := p.token
	var stmt ast.ScriptStatement
	var err error
	switch p.token.Kind {
	case token.Import:
		stmt, err = p.ParseImport()
	case token.Event:
		stmt, err = p.ParseEvent()
	case token.Auto, token.State:
		stmt, err = p.ParseState()
	case token.Function:
		stmt, err = p.ParseFunction()
	case token.Bool, token.BoolArray, token.Float, token.FloatArray, token.Int, token.IntArray, token.String, token.StringArray, token.Identifier, token.ObjectArray:
		switch p.lookahead.Kind {
		case token.Property:
			stmt, err = p.ParseProperty()
		case token.Function:
			stmt, err = p.ParseFunction()
		case token.Identifier:
			stmt, err = p.ParseScriptVariable()
		default:
			err = unexpectedTokenError(
				p.lookahead,
				token.Property,
				token.Function,
				token.Identifier)
		}
	default:
		err = unexpectedTokenError(
			p.token,
			token.Import,
			token.Event,
			token.Auto,
			token.State,
			token.Function,
			token.Bool,
			token.Float,
			token.Int,
			token.String,
			token.Identifier)
	}
	if err == nil {
		return stmt, nil
	}
	// Error recovery. Attempt to realign to a known statement token and emit an
	// error statement to fill the gap.
	if p.recovery || !p.attemptRecovery {
		// If an error was returned during a recovery operation
		// or we shouldn't even attempt recovery, just propagate it.
		return nil, err
	}
	p.recovery = true
	if err := p.recoverScriptStatement(); err != nil {
		return nil, err
	}
	errStmt := &ast.ErrorScriptStatement{
		Message:  fmt.Sprintf("%v", err),
		Location: source.Span(start.Location, p.token.Location),
	}
	p.errors = append(p.errors, errStmt)
	if err := p.next(); err != nil {
		return nil, err
	}
	p.recovery = false
	return errStmt, nil
}

func (p *parser) recoverScriptStatement() error {
	for {
		switch p.lookahead.Kind {
		case token.EOF:
			// Hit end of file, give up.
			return nil
		case token.Import, token.Event, token.Auto, token.State, token.Function, token.Bool, token.Float, token.Int, token.String, token.Identifier:
			// Next token is the start of a valid statement.
			return nil
		default:
			if err := p.next(); err != nil {
				return err // An error during recovery just fails.
			}
		}
	}
}

func (p *parser) ParseImport() (*ast.Import, error) {
	var err error
	node := &ast.Import{
		Keyword: p.token,
	}
	if err := p.tryConsume(token.Import); err != nil {
		return nil, err
	}
	node.Name, err = p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	node.Location = source.Span(node.Keyword.SourceLocation(), node.Name.Location)
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseState() (ast.ScriptStatement, error) {
	var err error
	node := &ast.State{}
	start := p.token.Location
	if p.token.Kind == token.Auto {
		node.Auto = p.token
		if err := p.next(); err != nil {
			return nil, err
		}
	}
	node.Keyword = p.token
	if err := p.tryConsume(token.State); err != nil {
		return nil, err
	}
	node.Name, err = p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	for p.token.Kind != token.EndState {
		if p.token.Kind == token.EOF {
			// State was never closed, proactively create a
			errStmt := &ast.ErrorScriptStatement{
				Message:  fmt.Sprintf("hit end of file while parsing state %q, did you forget %s?", node.Name.Location.Text(), token.EndState),
				Location: source.Span(start, p.token.Location),
			}
			p.errors = append(p.errors, errStmt)
			return errStmt, nil
		}
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		if p.token.Kind == token.EndState {
			break
		}
		stmt, err := p.ParseInvokable()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			node.Invokables = append(node.Invokables, stmt)
		}
	}
	node.Location = source.Span(start, p.token.Location)
	node.EndKeyword = p.token
	if err := p.tryConsume(token.EndState); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseInvokable() (ast.Invokable, error) {
	start := p.token
	var stmt ast.Invokable
	var err error
	switch p.token.Kind {
	case token.Event:
		stmt, err = p.ParseEvent()
	case token.Function, token.Bool, token.Float, token.Int, token.String, token.Identifier:
		stmt, err = p.ParseFunction()
	default:
		err = unexpectedTokenError(
			p.token,
			token.Event,
			token.Function,
			token.Bool,
			token.Float,
			token.Int,
			token.String,
			token.Identifier)
	}
	if err == nil {
		return stmt, nil
	}
	// Error recovery. Attempt to realign to a known statement token and emit an
	// error statement to fill the gap.
	if p.recovery || !p.attemptRecovery {
		// If an error was returned during a recovery operation
		// or we shouldn't even attempt recovery, just propagate it.
		return nil, err
	}
	p.recovery = true
	if err := p.recoverInvokable(); err != nil {
		return nil, err
	}
	errStmt := &ast.ErrorScriptStatement{
		Message:  fmt.Sprintf("%v", err),
		Location: source.Span(start.Location, p.token.Location),
	}
	p.errors = append(p.errors, errStmt)
	if err := p.next(); err != nil {
		return nil, err
	}
	p.recovery = false
	return errStmt, nil
}

func (p *parser) recoverInvokable() error {
	for {
		switch p.lookahead.Kind {
		case token.EOF:
			// Hit end of file, give up.
			return nil
		case token.Event, token.Function, token.Bool, token.Float, token.Int, token.String, token.Identifier:
			// Next token is the start of a valid invokable.
			return nil
		default:
			if err := p.next(); err != nil {
				return err // An error during recovery just fails.
			}
		}
	}
}

func (p *parser) ParseEvent() (*ast.Event, error) {
	var err error
	node := &ast.Event{
		Keyword: p.token,
	}
	if err := p.tryConsume(token.Event); err != nil {
		return nil, err
	}
	node.Name, err = p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	node.Open, node.Parameters, node.Close, err = p.ParseParameterList()
	if err != nil {
		return nil, err
	}
	end := p.token.Location
	for p.token.Kind == token.Native {
		node.Native = append(node.Native, p.token)
		end = p.token.Location
		if err := p.next(); err != nil {
			return nil, err
		}
	}
	if p.token.Kind == token.Newline && p.lookahead.Kind == token.BraceOpen {
		if err := p.next(); err != nil {
			return nil, err
		}
		node.Comment, err = p.ParseDocComment()
		if err != nil {
			return nil, err
		}
		end = node.Comment.Close.SourceLocation()
	}
	if node.Native != nil {
		node.Location = source.Span(node.Keyword.SourceLocation(), end)
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		return node, nil
	}
	node.Statements, err = p.ParseFunctionStatementBlock(token.EndEvent)
	if err != nil {
		return nil, err
	}
	node.Location = source.Span(node.Keyword.SourceLocation(), p.token.Location)
	node.EndKeyword = p.token
	if err := p.tryConsume(token.EndEvent); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseFunction() (*ast.Function, error) {
	var err error
	node := &ast.Function{}
	from := p.token.Location
	if p.token.Kind != token.Function {
		node.ReturnType, err = p.ParseTypeLiteral()
		if err != nil {
			return nil, err
		}
	}
	node.Keyword = p.token
	if err := p.tryConsume(token.Function); err != nil {
		return nil, err
	}
	node.Name, err = p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	node.Open, node.Parameters, node.Close, err = p.ParseParameterList()
	if err != nil {
		return nil, err
	}
	var end source.Location
	for p.token.Kind == token.Native || p.token.Kind == token.Global {
		if p.token.Kind == token.Native {
			node.Native = append(node.Native, p.token)
		} else {
			node.Global = append(node.Global, p.token)
		}
		end = p.token.Location
		if err := p.next(); err != nil {
			return nil, err
		}
	}
	if p.token.Kind == token.Newline && p.lookahead.Kind == token.BraceOpen {
		if err := p.next(); err != nil {
			return nil, err
		}
		node.Comment, err = p.ParseDocComment()
		if err != nil {
			return nil, err
		}
	}
	if len(node.Native) > 0 {
		node.Location = source.Span(from, end)
		return node, nil
	}
	node.Statements, err = p.ParseFunctionStatementBlock(token.EndFunction)
	if err != nil {
		return nil, err
	}
	node.Location = source.Span(from, p.token.Location)
	node.EndKeyword = p.token
	if err := p.tryConsume(token.EndFunction); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseParameterList() (*ast.Token, []*ast.Parameter, *ast.Token, error) {
	open := p.token
	if err := p.tryConsume(token.ParenthesisOpen); err != nil {
		return nil, nil, nil, err
	}
	var params []*ast.Parameter
	for {
		switch p.token.Kind {
		case token.Comma:
			if err := p.next(); err != nil {
				return nil, nil, nil, err
			}
		case token.ParenthesisClose:
			close := p.token
			if err := p.next(); err != nil {
				return nil, nil, nil, err
			}
			return open, params, close, nil
		default:
			param, err := p.ParseParameter()
			if err != nil {
				return nil, nil, nil, err
			}
			params = append(params, param)
		}
	}
}

func (p *parser) ParseParameter() (*ast.Parameter, error) {
	var err error
	node := &ast.Parameter{}
	node.Type, err = p.ParseTypeLiteral()
	if err != nil {
		return nil, err
	}
	node.Name, err = p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	node.Location = source.Span(node.Type.Location, node.Name.Location)
	if p.token.Kind == token.Assign {
		// Has default.
		node.Operator = p.token
		if err := p.tryConsume(token.Assign); err != nil {
			return nil, err
		}
		node.Value, err = p.ParseLiteral()
		if err != nil {
			return nil, err
		}
		node.Location = source.Span(node.Location, node.Value.SourceLocation())
	}
	return node, nil
}

func (p *parser) ParseFunctionStatementBlock(terminals ...token.Kind) ([]ast.FunctionStatement, error) {
	terms := make(map[token.Kind]struct{})
	for _, t := range terminals {
		terms[t] = struct{}{}
	}
	var stmts []ast.FunctionStatement
	for {
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		if _, ok := terms[p.token.Kind]; ok {
			return stmts, nil
		}
		start := p.token.Location
		stmt, err := p.ParseFunctionStatement()
		if err == nil {
			stmts = append(stmts, stmt)
			continue
		}
		// Error recovery. Attempt to realign to a known statement token and emit an
		// error statement to fill the gap.
		if p.recovery || !p.attemptRecovery {
			// If an error was returned during a recovery operation
			// or we shouldn't even attempt recovery, just propagate it.
			return nil, err
		}
		p.recovery = true
		if err := p.recoverFunctionStatement(); err != nil {
			return nil, err
		}
		errStmt := &ast.ErrorFunctionStatement{
			Message:  fmt.Sprintf("%v", err),
			Location: source.Span(start, p.token.Location),
		}
		p.errors = append(p.errors, errStmt)
		p.recovery = false
		stmts = append(stmts, errStmt)
	}
}

func (p *parser) recoverFunctionStatement() error {
	for {
		switch p.lookahead.Kind {
		case token.EOF:
			// Hit end of file, give up.
			return nil
		case token.Newline:
			// Next token is the start of a valid invokable.
			return nil
		default:
			if err := p.next(); err != nil {
				return err // An error during recovery just fails.
			}
		}
	}
}

func (p *parser) ParseFunctionStatement() (ast.FunctionStatement, error) {
	switch p.token.Kind {
	case token.Return:
		return p.ParseReturn()
	case token.If:
		return p.ParseIf()
	case token.While:
		return p.ParseWhile()
	case token.Bool, token.Int, token.Float, token.String:
		return p.ParseFunctionVariable()
	case token.Identifier:
		switch p.lookahead.Kind {
		case token.Identifier: // p.token is an object type, p.lookahead is a variable name
			return p.ParseFunctionVariable()
		case token.Assign, token.AssignAdd, token.AssignDivide, token.AssignModulo, token.AssignMultiply, token.AssignSubtract:
			return p.ParseAssignment()
		}
	}
	return p.ParseExpression(lowest)
}

func (p *parser) ParseFunctionVariable() (*ast.FunctionVariable, error) {
	var err error
	node := &ast.FunctionVariable{}
	node.Type, err = p.ParseTypeLiteral()
	if err != nil {
		return nil, err
	}
	node.Name, err = p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	end := node.Name.Location
	if p.token.Kind == token.Assign {
		node.Operator = p.token
		if err := p.tryConsume(token.Assign); err != nil {
			return nil, err
		}
		node.Value, err = p.ParseExpression(lowest)
		if err != nil {
			return nil, err
		}
		end = node.Value.SourceLocation()
	}
	node.Location = source.Span(node.Type.Location, end)
	return node, nil
}

func (p *parser) ParseAssignment() (*ast.Assignment, error) {
	start := p.token.Location
	assignee, err := p.ParseExpression(lowest)
	if err != nil {
		return nil, err
	}
	operator := p.token
	if err := p.tryConsume(
		token.Assign,
		token.AssignAdd,
		token.AssignDivide,
		token.AssignModulo,
		token.AssignMultiply,
		token.AssignSubtract); err != nil {
		return nil, err
	}
	expr, err := p.ParseExpression(lowest)
	if err != nil {
		return nil, err
	}
	return &ast.Assignment{
		Assignee: assignee,
		Operator: operator,
		Value:    expr,
		Location: source.Span(start, expr.SourceLocation()),
	}, nil
}

func (p *parser) ParseReturn() (*ast.Return, error) {
	var err error
	node := &ast.Return{
		Keyword:  p.token,
		Location: p.token.Location,
	}
	if err := p.tryConsume(token.Return); err != nil {
		return nil, err
	}
	if p.token.Kind == token.Newline {
		return node, nil
	}
	node.Value, err = p.ParseExpression(lowest)
	if err != nil {
		return nil, err
	}
	node.Location = source.Span(node.Location, node.Value.SourceLocation())
	return node, nil
}

func (p *parser) ParseIf() (*ast.If, error) {
	var err error
	node := &ast.If{
		Keyword: p.token,
	}
	if err := p.tryConsume(token.If); err != nil {
		return nil, err
	}
	node.Condition, err = p.ParseExpression(lowest)
	if err != nil {
		return nil, err
	}
	if err := p.consumeNewlines(); err != nil {
		return nil, err
	}
	node.Statements, err = p.ParseFunctionStatementBlock(token.EndIf, token.Else, token.ElseIf)
	if err != nil {
		return nil, err
	}
	for p.token.Kind == token.ElseIf {
		block := &ast.ElseIf{
			Keyword: p.token,
		}
		if err := p.tryConsume(token.ElseIf); err != nil {
			return nil, err
		}
		block.Condition, err = p.ParseExpression(lowest)
		if err != nil {
			return nil, err
		}
		end := block.Condition.SourceLocation()
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		block.Statements, err = p.ParseFunctionStatementBlock(token.EndIf, token.Else, token.ElseIf)
		if err != nil {
			return nil, err
		}
		if len(block.Statements) > 0 {
			end = block.Statements[len(block.Statements)-1].SourceLocation()
		}
		block.Location = source.Span(block.Keyword.SourceLocation(), end)
		node.ElseIfs = append(node.ElseIfs, block)
	}
	if p.token.Kind == token.Else {
		block := &ast.Else{
			Keyword: p.token,
		}
		if err := p.tryConsume(token.Else); err != nil {
			return nil, err
		}
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		end := block.Keyword.SourceLocation()
		block.Statements, err = p.ParseFunctionStatementBlock(token.EndIf)
		if err != nil {
			return nil, err
		}
		if len(block.Statements) > 0 {
			end = block.Statements[len(block.Statements)-1].SourceLocation()
		}
		node.Location = source.Span(block.Keyword.SourceLocation(), end)
		node.Else = block
	}
	node.Location = source.Span(node.Keyword.SourceLocation(), p.token.Location)
	node.EndKeyword = p.token
	if err := p.tryConsume(token.EndIf); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseWhile() (*ast.While, error) {
	var err error
	node := &ast.While{
		Keyword: p.token,
	}
	if err := p.tryConsume(token.While); err != nil {
		return nil, err
	}
	node.Condition, err = p.ParseExpression(lowest)
	if err != nil {
		return nil, err
	}
	if err := p.consumeNewlines(); err != nil {
		return nil, err
	}
	node.Statements, err = p.ParseFunctionStatementBlock(token.EndWhile)
	if err != nil {
		return nil, err
	}
	node.Location = source.Span(node.Keyword.SourceLocation(), p.token.Location)
	node.EndKeyword = p.token
	if err := p.tryConsume(token.EndWhile); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseProperty() (*ast.Property, error) {
	var err error
	node := &ast.Property{}
	node.Type, err = p.ParseTypeLiteral()
	if err != nil {
		return nil, err
	}
	node.Keyword = p.token
	if err := p.tryConsume(token.Property); err != nil {
		return nil, err
	}
	node.Name, err = p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	end := node.Name.Location
	if p.token.Kind == token.Assign {
		if err := p.tryConsume(token.Assign); err != nil {
			return nil, err
		}
		node.Value, err = p.ParseLiteral()
		if err != nil {
			return nil, err
		}
		end = node.Value.SourceLocation()
	}
	if p.token.Kind == token.Auto {
		node.Auto = p.token
		end = p.token.Location
		if err := p.tryConsume(token.Auto); err != nil {
			return nil, err
		}
	} else if p.token.Kind == token.AutoReadOnly {
		if node.Value == nil {
			return nil, newError(p.token.Location, "expected value to be defined for %s property", token.AutoReadOnly)
		}
		node.AutoReadOnly = p.token
		end = p.token.Location
		if err := p.tryConsume(token.AutoReadOnly); err != nil {
			return nil, err
		}
	}
	if node.Auto != nil || node.AutoReadOnly != nil {
		for p.token.Kind == token.Hidden || p.token.Kind == token.Conditional {
			end = p.token.Location
			if p.token.Kind == token.Hidden {
				node.Hidden = append(node.Hidden, p.token)
			} else {
				node.Conditional = append(node.Conditional, p.token)
			}
			if err := p.tryConsume(token.Hidden, token.Conditional); err != nil {
				return nil, err
			}
		}
		if p.token.Kind == token.Newline && p.lookahead.Kind == token.BraceOpen {
			if err := p.next(); err != nil {
				return nil, err
			}
			node.Comment, err = p.ParseDocComment()
			if err != nil {
				return nil, err
			}
			end = node.Comment.Close.SourceLocation()
		}
		node.Location = source.Span(node.Type.Location, end)
		return node, nil
	}
	// Full Property
	for p.token.Kind == token.Hidden {
		node.Hidden = append(node.Hidden, p.token)
		if err := p.tryConsume(token.Hidden); err != nil {
			return nil, err
		}
	}
	if err := p.tryConsume(token.Newline); err != nil {
		return nil, err
	}
	if p.token.Kind == token.BraceOpen {
		if err := p.next(); err != nil {
			return nil, err
		}
		comment, err := p.ParseDocComment()
		if err != nil {
			return nil, err
		}
		node.Comment = comment
	}
	if err := p.consumeNewlines(); err != nil {
		return nil, err
	}
	first, err := p.ParseFunction()
	if err != nil {
		return nil, err
	}
	if err := p.consumeNewlines(); err != nil {
		return nil, err
	}
	var second *ast.Function
	if p.token.Kind != token.EndProperty {
		second, err = p.ParseFunction()
		if err != nil {
			return nil, err
		}
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
	}
	if first.Name.Normalized == "get" {
		if first.ReturnType == nil {
			return nil, newError(first.Name.Location, "expected '%s' to have a return type of %s, but found none", first.Name.SourceLocation().Text(), node.Type.Location.Text())
		}
		if first.ReturnType.Type != node.Type.Type {
			return nil, newError(first.ReturnType.Location, "expected '%s' to have a return type of %s, but found %s", first.Name.SourceLocation().Text(), node.Type.Location.Text(), first.ReturnType.SourceLocation().Text())
		}
		if len(first.Parameters) != 0 {
			loc := source.Span(first.Parameters[0].Location, first.Parameters[len(first.Parameters)-1].Location)
			return nil, newError(loc, "expected '%s' to have no parameters, but found %d", first.Name.SourceLocation().Text(), len(first.Parameters))
		}
		node.Get = first
	} else if first.Name.Normalized == "set" {
		if first.ReturnType != nil {
			return nil, newError(first.ReturnType.Location, "expected '%s' to have no return type, but found %s", first.Name.SourceLocation().Text(), first.ReturnType.Location.Text())
		}
		if len(first.Parameters) == 0 {
			return nil, newError(first.Name.Location, "expected '%s' to have one parameter, but found none", first.Name.SourceLocation().Text())
		}
		if len(first.Parameters) > 1 {
			loc := source.Span(first.Parameters[0].Location, first.Parameters[len(first.Parameters)-1].Location)
			return nil, newError(loc, "expected '%s' to have one parameter, but found %d", first.Name.SourceLocation().Text(), len(first.Parameters))
		}
		if first.Parameters[0].Type.Type != node.Type.Type {
			return nil, newError(first.ReturnType.Location, "expected '%s' to have a parameter of type %s, but found %s", first.Name.SourceLocation().Text(), node.Type.Location.Text(), first.Parameters[0].Type.Location.Text())
		}
		node.Set = first
	} else {
		return nil, newError(first.SourceLocation(), "expected 'Get' or 'Set' function for property, but found '%s'", first.Name.SourceLocation().Text())
	}
	if second != nil {
		if second.Name.Normalized == "get" {
			if node.Get != nil {
				return nil, newError(second.Location, "expected exactly one 'Get' function, but found two")
			}
			if second.ReturnType == nil {
				return nil, newError(second.Name.Location, "expected '%s' to have a return type of %s, but found none", second.Name.SourceLocation().Text(), node.Type.Location.Text())
			}
			if second.ReturnType.Type != node.Type.Type {
				return nil, newError(second.ReturnType.Location, "expected '%s' to have a return type of %s, but found %s", second.Name.SourceLocation().Text(), node.Type.Location.Text(), second.ReturnType.SourceLocation().Text())
			}
			if len(second.Parameters) != 0 {
				loc := source.Span(second.Parameters[0].Location, second.Parameters[len(second.Parameters)-1].Location)
				return nil, newError(loc, "expected '%s' to have no parameters, but found %d", second.Name.SourceLocation().Text(), len(second.Parameters))
			}
			node.Get = second
		} else if second.Name.Normalized == "set" {
			if node.Set != nil {
				return nil, newError(second.Location, "expected exactly one 'Set' function, but found two")
			}
			if second.ReturnType != nil {
				return nil, newError(second.ReturnType.Location, "expected '%s' to have no return type, but found %s", second.Name.SourceLocation().Text(), second.ReturnType.Location.Text())
			}
			if len(second.Parameters) == 0 {
				return nil, newError(second.Name.Location, "expected '%s' to have one parameter, but found none", second.Name.SourceLocation().Text())
			}
			if len(second.Parameters) > 1 {
				loc := source.Span(second.Parameters[0].Location, second.Parameters[len(second.Parameters)-1].Location)
				return nil, newError(loc, "expected '%s' to have one parameter, but found %d", second.Name.SourceLocation().Text(), len(second.Parameters))
			}
			if second.Parameters[0].Type.Type != node.Type.Type {
				return nil, newError(second.ReturnType.Location, "expected '%s' to have a parameter of type %s, but found %s", second.Name.SourceLocation().Text(), node.Type.Location.Text(), second.Parameters[0].Type.Location.Text())
			}
			node.Set = second
		} else {
			return nil, newError(second.SourceLocation(), "expected 'Get' or 'Set' function for property, but found '%s'", second.Name.SourceLocation().Text())
		}
	}
	node.Location = source.Span(node.Type.Location, p.token.Location)
	node.EndKeyword = p.token
	if err := p.tryConsume(token.EndProperty); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseScriptVariable() (*ast.ScriptVariable, error) {
	var err error
	node := &ast.ScriptVariable{}
	node.Type, err = p.ParseTypeLiteral()
	if err != nil {
		return nil, err
	}
	node.Name, err = p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	end := node.Name.Location
	if p.token.Kind == token.Assign {
		node.Operator = p.token
		if err := p.tryConsume(token.Assign); err != nil {
			return nil, err
		}
		node.Value, err = p.ParseLiteral()
		if err != nil {
			return nil, err
		}
		end = node.Value.SourceLocation()
	}
	for p.token.Kind == token.Conditional {
		end = p.token.Location
		node.Conditional = append(node.Conditional, p.token)
		if err := p.next(); err != nil {
			return nil, err
		}
	}
	node.Location = source.Span(node.Type.Location, end)
	if err := p.tryConsume(token.Newline, token.EOF); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseIdentifier() (*ast.Identifier, error) {
	node := &ast.Identifier{
		Text:       p.token,
		Normalized: string(bytes.ToLower(p.token.Location.Text())),
		Location:   p.token.Location,
	}
	if err := p.tryConsume(token.Identifier, token.Self, token.Parent); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseTypeLiteral() (*ast.TypeLiteral, error) {
	node := &ast.TypeLiteral{
		Text:     p.token,
		Location: p.token.Location,
	}
	switch p.token.Kind {
	case token.Bool:
		node.Type = types.Bool{}
	case token.BoolArray:
		node.Type = types.Array{ElementType: types.Bool{}}
	case token.Int:
		node.Type = types.Int{}
	case token.IntArray:
		node.Type = types.Array{ElementType: types.Int{}}
	case token.Float:
		node.Type = types.Float{}
	case token.FloatArray:
		node.Type = types.Array{ElementType: types.Float{}}
	case token.String:
		node.Type = types.String{}
	case token.StringArray:
		node.Type = types.Array{ElementType: types.String{}}
	case token.Identifier:
		node.Type = types.Object{
			Name: string(bytes.ToLower(p.token.Location.Text())),
		}
	case token.ObjectArray:
		node.Type = types.Array{
			ElementType: types.Object{
				Name: string(bytes.TrimSuffix(bytes.ToLower(p.token.Location.Text()), []byte{'[', ']'})),
			},
		}
	default:
		return nil, unexpectedTokenError(
			p.token,
			token.Bool,
			token.Float,
			token.Int,
			token.String,
			token.Identifier)
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseExpression(precedence int) (ast.Expression, error) {
	prefix := p.prefix[p.token.Kind]
	if prefix == nil {
		want := keys(p.prefix)
		return nil, unexpectedTokenError(p.token, want[0], want[1:]...)
	}
	expr, err := prefix()
	if err != nil {
		return nil, err
	}
	for p.token.Kind != token.Newline && p.token.Kind != token.EOF && precedence < precedenceOf(p.token.Kind) {
		infix := p.infix[p.token.Kind]
		if infix == nil {
			return expr, nil
		}
		expr, err = infix(expr)
		if err != nil {
			return nil, err
		}
	}
	return expr, nil
}

func (p *parser) ParseBinary(left ast.Expression) (*ast.Binary, error) {
	precedence := precedenceOf(p.token.Kind)
	operator := p.token
	if err := p.tryConsume(
		token.LogicalOr,
		token.LogicalAnd,
		token.Equal,
		token.NotEqual,
		token.Greater,
		token.GreaterOrEqual,
		token.Less,
		token.LessOrEqual,
		token.Plus,
		token.Minus,
		token.Divide,
		token.Multiply,
		token.Modulo); err != nil {
		return nil, err
	}
	right, err := p.ParseExpression(precedence)
	if err != nil {
		return nil, err
	}
	return &ast.Binary{
		LeftOperand:  left,
		Operator:     operator,
		RightOperand: right,
		Location:     source.Span(left.SourceLocation(), right.SourceLocation()),
	}, nil
}

func (p *parser) ParseUnary() (ast.Expression, error) {
	if p.token.Kind == token.Minus &&
		(p.lookahead.Kind == token.IntLiteral || p.lookahead.Kind == token.FloatLiteral) {
		return p.ParseLiteral()
	}
	var err error
	node := &ast.Unary{
		Operator: p.token,
	}
	if err := p.tryConsume(token.Minus, token.LogicalNot); err != nil {
		return nil, err
	}
	node.Operand, err = p.ParseExpression(prefix)
	if err != nil {
		return nil, err
	}
	node.Location = source.Span(node.Operator.SourceLocation(), node.Operand.SourceLocation())
	return node, nil
}

func (p *parser) ParseCast(value ast.Expression) (*ast.Cast, error) {
	operator := p.token
	if err := p.tryConsume(token.As); err != nil {
		return nil, err
	}
	typeLiteral, err := p.ParseTypeLiteral()
	if err != nil {
		return nil, err
	}
	return &ast.Cast{
		Value:    value,
		Operator: operator,
		Type:     typeLiteral,
		Location: source.Span(value.SourceLocation(), typeLiteral.Location),
	}, nil
}

func (p *parser) ParseAccess(value ast.Expression) (*ast.Access, error) {
	operator := p.token
	if err := p.tryConsume(token.Dot); err != nil {
		return nil, err
	}
	name, err := p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	return &ast.Access{
		Value:    value,
		Operator: operator,
		Name:     name,
		Location: source.Span(value.SourceLocation(), name.Location),
	}, nil
}

func (p *parser) ParseIndex(array ast.Expression) (*ast.Index, error) {
	open := p.token
	if err := p.tryConsume(token.BracketOpen); err != nil {
		return nil, err
	}
	index, err := p.ParseExpression(lowest)
	if err != nil {
		return nil, err
	}
	close := p.token
	if err := p.tryConsume(token.BracketClose); err != nil {
		return nil, err
	}
	return &ast.Index{
		Value:    array,
		Open:     open,
		Index:    index,
		Close:    close,
		Location: source.Span(array.SourceLocation(), close.Location),
	}, nil
}

func (p *parser) ParseCall(function ast.Expression) (*ast.Call, error) {
	open := p.token
	if err := p.tryConsume(token.ParenthesisOpen); err != nil {
		return nil, err
	}
	args, err := p.ParseArgumentList()
	if err != nil {
		return nil, err
	}
	close := p.token
	if err := p.tryConsume(token.ParenthesisClose); err != nil {
		return nil, err
	}
	return &ast.Call{
		Function:  function,
		Open:      open,
		Arguments: args,
		Close:     close,
		Location:  source.Span(function.SourceLocation(), close.Location),
	}, nil
}

func (p *parser) ParseArgumentList() ([]*ast.Argument, error) {
	var args []*ast.Argument
	for {
		switch p.token.Kind {
		case token.Comma:
			if err := p.next(); err != nil {
				return nil, err
			}
		case token.ParenthesisClose:
			return args, nil
		default:
			arg, err := p.ParseArgument()
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
	}
}

func (p *parser) ParseArgument() (*ast.Argument, error) {
	node := &ast.Argument{}
	if p.token.Kind == token.Identifier && p.lookahead.Kind == token.Assign {
		id, err := p.ParseIdentifier()
		if err != nil {
			return nil, err
		}
		node.Name = id
		node.Operator = p.token
		if err := p.tryConsume(token.Assign); err != nil {
			return nil, err
		}
	}
	value, err := p.ParseExpression(lowest)
	if err != nil {
		return nil, err
	}
	node.Value = value
	if node.Name != nil {
		node.Location = source.Span(node.Name.Location, value.SourceLocation())
	} else {
		node.Location = value.SourceLocation()
	}
	return node, nil
}

func (p *parser) ParseArrayCreation() (*ast.ArrayCreation, error) {
	new := p.token
	if err := p.tryConsume(token.New); err != nil {
		return nil, err
	}
	typeLiteral, err := p.ParseTypeLiteral()
	if err != nil {
		return nil, err
	}
	open := p.token
	if err := p.tryConsume(token.BracketOpen); err != nil {
		return nil, err
	}
	size, err := p.ParseIntLiteral()
	if err != nil {
		return nil, err
	}
	if size.Value < 1 || size.Value > 128 {
		return nil, newError(size.SourceLocation(), "expected array size to be an %s in range [1, 128], but found %d", token.IntLiteral, size.Value)
	}
	close := p.token
	if err := p.tryConsume(token.BracketClose); err != nil {
		return nil, err
	}
	return &ast.ArrayCreation{
		New:      new,
		Type:     typeLiteral,
		Open:     open,
		Size:     size,
		Close:    close,
		Location: source.Span(new.Location, close.Location),
	}, nil
}

func (p *parser) ParseParenthetical() (*ast.Parenthetical, error) {
	var err error
	node := &ast.Parenthetical{
		Open: p.token,
	}
	if err := p.tryConsume(token.ParenthesisOpen); err != nil {
		return nil, err
	}
	node.Value, err = p.ParseExpression(lowest)
	if err != nil {
		return nil, err
	}
	node.Close = p.token
	if err := p.tryConsume(token.ParenthesisClose); err != nil {
		return nil, err
	}
	node.Location = source.Span(node.Open.SourceLocation(), node.Close.SourceLocation())
	return node, nil
}

func (p *parser) ParseLiteral() (ast.Literal, error) {
	switch p.token.Kind {
	case token.Minus:
		// While this overlaps with a unary expression, we lump these together
		// because there are some contexts where a literal is required which can
		// include a sign, but where a Unary is not allowed.
		sign := p.token
		if err := p.next(); err != nil {
			return nil, err
		}
		switch p.token.Kind {
		case token.IntLiteral:
			lit, err := p.ParseIntLiteral()
			if err != nil {
				return nil, err
			}
			lit.Location = source.Span(sign.Location, lit.Location)
			lit.Value = -lit.Value
			return lit, err
		case token.FloatLiteral:
			lit, err := p.ParseFloatLiteral()
			if err != nil {
				return nil, err
			}
			lit.Location = source.Span(sign.Location, lit.Location)
			lit.Value = -lit.Value
			return lit, err
		default:
			return nil, unexpectedTokenError(
				p.token,
				token.IntLiteral,
				token.FloatLiteral)
		}
	case token.True, token.False:
		return p.ParseBoolLiteral()
	case token.IntLiteral:
		return p.ParseIntLiteral()
	case token.FloatLiteral:
		return p.ParseFloatLiteral()
	case token.StringLiteral:
		return p.ParseStringLiteral()
	case token.None:
		return p.ParseNoneLiteral()
	}
	return nil, unexpectedTokenError(
		p.token,
		token.True,
		token.False,
		token.IntLiteral,
		token.FloatLiteral,
		token.StringLiteral,
		token.None)
}

func (p *parser) ParseIntLiteral() (*ast.IntLiteral, error) {
	node := &ast.IntLiteral{
		Text:     p.token,
		Location: p.token.Location,
	}
	if err := p.tryConsume(token.IntLiteral); err != nil {
		return nil, err
	}
	text := strings.ToLower(string(node.Location.Text()))
	val, err := strconv.ParseInt(text, 0, 32)
	if err != nil {
		return nil, newError(node.Location, "failed to parse %q as an integer: %v", text, err)
	}
	node.Value = int(val)
	return node, nil
}

func (p *parser) ParseFloatLiteral() (*ast.FloatLiteral, error) {
	node := &ast.FloatLiteral{
		Text:     p.token,
		Location: p.token.Location,
	}
	if err := p.tryConsume(token.FloatLiteral); err != nil {
		return nil, err
	}
	text := strings.ToLower(string(node.Location.Text()))
	val, err := strconv.ParseFloat(text, 32)
	if err != nil {
		return nil, newError(node.Location, "failed to parse %q as a float: %v", text, err)
	}
	node.Value = float32(val)
	return node, nil
}

func (p *parser) ParseBoolLiteral() (*ast.BoolLiteral, error) {
	node := &ast.BoolLiteral{
		Text:     p.token,
		Value:    p.token.Kind == token.True,
		Location: p.token.Location,
	}
	if err := p.tryConsume(token.True, token.False); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseStringLiteral() (*ast.StringLiteral, error) {
	node := &ast.StringLiteral{
		Text:     p.token,
		Value:    string(p.token.Location.Text()[1 : p.token.Location.Length-1]),
		Location: p.token.Location,
	}
	if err := p.tryConsume(token.StringLiteral); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseNoneLiteral() (*ast.NoneLiteral, error) {
	node := &ast.NoneLiteral{
		Text:     p.token,
		Location: p.token.Location,
	}
	if err := p.tryConsume(token.None); err != nil {
		return nil, err
	}
	return node, nil
}

func registerPrefix[T ast.Expression](p *parser, fn func() (T, error), types ...token.Kind) {
	for _, t := range types {
		p.prefix[t] = func() (ast.Expression, error) { return fn() }
	}
}

func registerInfix[T ast.Expression](p *parser, fn func(ast.Expression) (T, error), types ...token.Kind) {
	for _, t := range types {
		p.infix[t] = func(expr ast.Expression) (ast.Expression, error) { return fn(expr) }
	}
}

func keys[K comparable, V any](data map[K]V) []K {
	keys := make([]K, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}

func tok(t token.Token) *ast.Token {
	return &ast.Token{
		Kind:     t.Kind,
		Location: t.Location,
	}
}
