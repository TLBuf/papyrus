// Package parser defines a Papyrus parser.
package parser

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/TLBuf/papyrus/pkg/ast"
	"github.com/TLBuf/papyrus/pkg/lexer"
	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/TLBuf/papyrus/pkg/token"
	"github.com/TLBuf/papyrus/pkg/types"
)

// Parser provides the ability to lex and parse a Papyrus script into an
// [*ast.Script].
type Parser struct {
	keepLooseComments bool
}

type Option func(*Parser)

// WithLooseComments directs the parser on whether or not to retain loose
// comments that may appear (i.e. line and block comments). Doc comments are
// always captured.
func WithLooseComments(keep bool) Option {
	return func(p *Parser) {
		p.keepLooseComments = keep
	}
}

// New returns a [*Parser] that is configured to parser script files.
func New(opts ...Option) *Parser {
	p := &Parser{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// Parser returns the file parsed as an [*ast.Script] or an [Error] if parsing
// encountered one or more issues.
func (p *Parser) Parse(file *source.File) (*ast.Script, error) {
	lex, err := lexer.New(file)
	if err != nil {
		return nil, err
	}
	prsr := &parser{
		l:                 lex,
		keepLooseComments: p.keepLooseComments,
		prefix:            make(map[token.Kind]prefixParser),
		infix:             make(map[token.Kind]infixParser),
	}

	registerPrefix(prsr, prsr.ParseBoolLiteral, token.True, token.False)
	registerPrefix(prsr, prsr.ParseFloatLiteral, token.FloatLiteral)
	registerPrefix(prsr, prsr.ParseIdentifier, token.Identifier)
	registerPrefix(prsr, prsr.ParseIntLiteral, token.IntLiteral)
	registerPrefix(prsr, prsr.ParseNoneLiteral, token.None)
	registerPrefix(prsr, prsr.ParseParenthetical, token.ParenthesisOpen)
	registerPrefix(prsr, prsr.ParseStringLiteral, token.StringLiteral)
	registerPrefix(prsr, prsr.ParseUnary, token.Minus, token.LogicalNot)

	registerInfix(prsr, prsr.ParseAccess, token.Dot)
	registerInfix(prsr, prsr.ParseBinary,
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
	registerInfix(prsr, prsr.ParseCall, token.ParenthesisOpen)
	registerInfix(prsr, prsr.ParseCast, token.As)
	registerInfix(prsr, prsr.ParseIndex, token.BracketOpen)

	if err := prsr.next(); err != nil {
		return nil, err
	}
	if err := prsr.next(); err != nil {
		return nil, err
	}
	return prsr.ParseScript()
}

type parser struct {
	l *lexer.Lexer

	token     token.Token
	lookahead token.Token

	keepLooseComments bool
	looseComments     []ast.LooseComment

	recovery bool
	errors   []ast.Error

	prefix map[token.Kind]prefixParser
	infix  map[token.Kind]infixParser
}

const (
	_ int = iota
	Lowest
	LogicalOr      // ||
	LogicalAnd     // &&
	Comparison     // ==, !=, >, >=, <, <=
	Additive       // +, -
	Multiplicitive // *, /, %
	Prefix         // -x or !y
	Cast           // x As y
	Access         // x.y
	Call           // x(y)
	Index          // x[y]
)

var precedences = map[token.Kind]int{
	token.LogicalOr:       LogicalOr,
	token.LogicalAnd:      LogicalAnd,
	token.Equal:           Comparison,
	token.NotEqual:        Comparison,
	token.Greater:         Comparison,
	token.GreaterOrEqual:  Comparison,
	token.Less:            Comparison,
	token.LessOrEqual:     Comparison,
	token.Plus:            Additive,
	token.Minus:           Additive,
	token.Multiply:        Multiplicitive,
	token.Divide:          Multiplicitive,
	token.Modulo:          Multiplicitive,
	token.As:              Cast,
	token.Dot:             Access,
	token.ParenthesisOpen: Call,
	token.BracketOpen:     Index,
}

func precedenceOf(t token.Kind) int {
	if p, ok := precedences[t]; ok {
		return p
	}
	return Lowest
}

type (
	prefixParser func() (ast.Expression, error)
	infixParser  func(ast.Expression) (ast.Expression, error)
)

// next advances token and lookahead by one token while skipping loose comment
// tokens. Returns true if parsing should continue, false otherwise.
func (p *parser) next() error {
	p.token = p.lookahead
	t, err := p.l.NextToken()
	if err != nil {
		return Error{
			Err:      err,
			Location: err.(lexer.Error).Location,
		}
	}
	p.lookahead = t
	// Consume loose comments immediately so the rest of the
	// parser never has to deal with them directly.
	if p.token.Kind == token.Semicolon {
		tok, err := p.ParseLineComment()
		if err != nil {
			return err
		}
		if p.keepLooseComments {
			p.looseComments = append(p.looseComments, tok)
		}
	}
	if p.token.Kind == token.BlockCommentOpen {
		tok, err := p.ParseBlockComment()
		if err != nil {
			return err
		}
		if p.keepLooseComments {
			p.looseComments = append(p.looseComments, tok)
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
	if len(alts) > 0 {
		return newError(p.token.Location, "expected any of [%s, %s], but found %s", t, tokensTypesToString(alts...), p.token.Kind)
	}
	return newError(p.token.Location, "expected %s, but found %s", t, p.token.Kind)
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
		Semicolon: p.token,
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
	for p.token.Kind != token.EOF {
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		stmt, err := p.ParseScriptStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			node.Statements = append(node.Statements, stmt)
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
	case token.Bool, token.Float, token.Int, token.String, token.Identifier:
		switch p.lookahead.Kind {
		case token.Property:
			stmt, err = p.ParseProperty()
		case token.Function:
			stmt, err = p.ParseFunction()
		case token.Identifier:
			stmt, err = p.ParseScriptVariable()
		default:
			err = fmt.Errorf("expected Import, Event, State, Function, Property, or a variable definition, but found %s", start.Kind)
		}
	default:
		err = fmt.Errorf("expected Import, Event, State, Function, Property, or a variable definition, but found %s", start.Kind)
	}
	if err == nil {
		return stmt, nil
	}
	// Error recovery. Attempt to realign to a known statement token and emit an
	// error statement to fill the gap.
	if p.recovery {
		// If an error was returned during a recovery operation, just propagate it.
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
		err = fmt.Errorf("expected Event or Function, but found %s", start.Kind)
	}
	if err == nil {
		return stmt, nil
	}
	// Error recovery. Attempt to realign to a known statement token and emit an
	// error statement to fill the gap.
	if p.recovery {
		// If an error was returned during a recovery operation, just propagate it.
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
	if p.token.Kind == token.Native {
		node.Native = p.token
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

func (p *parser) ParseParameterList() (ast.Token, []*ast.Parameter, ast.Token, error) {
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
		if p.recovery {
			// If an error was returned during a recovery operation, just propagate it.
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
		if p.lookahead.Kind == token.Identifier { // Object type.
			return p.ParseFunctionVariable()
		}
	}
	return p.ParseAssignment()
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
		node.Value, err = p.ParseLiteral()
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
	assignee, err := p.ParseExpression(Lowest)
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
	expr, err := p.ParseExpression(Lowest)
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
	node.Value, err = p.ParseExpression(Lowest)
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
	node.Condition, err = p.ParseExpression(Lowest)
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
		block.Condition, err = p.ParseExpression(Lowest)
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
	node.Condition, err = p.ParseExpression(Lowest)
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
			return nil, newError(p.token.Location, "expected value to be defined for AutoReadOnly property")
		}
		node.AutoReadOnly = p.token
		end = p.token.Location
		if err := p.tryConsume(token.Auto); err != nil {
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
	if err := p.tryConsume(token.Identifier); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseTypeLiteral() (*ast.TypeLiteral, error) {
	node := &ast.TypeLiteral{
		Text:     p.token,
		Location: p.token.Location,
	}
	var scalar types.Scalar
	switch p.token.Kind {
	case token.Bool:
		scalar = types.Bool{}
	case token.Int:
		scalar = types.Int{}
	case token.Float:
		scalar = types.Float{}
	case token.String:
		scalar = types.String{}
	case token.Identifier:
		scalar = types.Object{
			Name: string(bytes.ToLower(p.token.Location.Text())),
		}
	default:
		return nil, newError(p.token.Location, "expected Bool, Int, Float, String, or an identifier, but found %s", p.token.Kind)
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	if p.token.Kind != token.BracketOpen {
		return node, nil
	}
	node.Open = p.token
	if err := p.tryConsume(token.BracketOpen); err != nil {
		return nil, err
	}
	node.Close = p.token
	if err := p.tryConsume(token.BracketClose); err != nil {
		return nil, err
	}
	return &ast.TypeLiteral{
		Type: types.Array{
			ElementType: scalar,
		},
		Location: source.Span(node.Location, node.Close.SourceLocation()),
	}, nil
}

func (p *parser) ParseExpression(precedence int) (ast.Expression, error) {
	prefix := p.prefix[p.token.Kind]
	if prefix == nil {
		return nil, newError(p.token.Location, "expected any of [%s], but found %s", tokensTypesToString(keys(p.prefix)...), p.token.Kind)
	}
	expr, err := prefix()
	if err != nil {
		return nil, err
	}
	if p.lookahead.Kind != token.Newline && p.lookahead.Kind != token.EOF && precedence < precedenceOf(p.lookahead.Kind) {
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
	node.Operand, err = p.ParseExpression(Prefix)
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
	index, err := p.ParseExpression(Lowest)
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

func (p *parser) ParseCall(reciever ast.Expression) (*ast.Call, error) {
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
		Reciever:  reciever,
		Open:      open,
		Arguments: args,
		Close:     close,
		Location:  source.Span(reciever.SourceLocation(), close.Location),
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
	if p.token.Kind == token.Identifier {
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
	value, err := p.ParseExpression(Lowest)
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
		return nil, newError(size.SourceLocation(), "expected array size to be an IntLiteral in range [1, 128], but found %d", size.Value)
	}
	close := p.token
	if err := p.tryConsume(token.BracketClose); err != nil {
		return nil, err
	}
	return &ast.ArrayCreation{
		NewOperator: new,
		Type:        typeLiteral,
		Open:        open,
		Size:        size,
		Close:       close,
		Location:    source.Span(new.Location, close.Location),
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
	node.Value, err = p.ParseExpression(Lowest)
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
			return nil, fmt.Errorf("expected IntLiteral or FloatLiteral, but found %s", p.token.Kind)
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
	return nil, fmt.Errorf("expected True, False, None, Integer, Float, or String literal, but found %s", p.token.Kind)
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
