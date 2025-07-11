// Package parser defines a Papyrus parser.
package parser

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/lexer"
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/token"
)

// Option defines an option to configure how parsing is performed.
type Option interface{ apply(*parser) }

type option func(*parser)

// apply implements the [Option] interface.
func (o option) apply(p *parser) {
	o(p)
}

// WithComments controls block and line (i.e. loose) comment processing.
//
// If enabled, loose comments will be attached to the appropriate nodes and/or
// appear as [ast.CommentStatement] nodes. This is only required when the nodes
// may need to be written back out as source, e.g. when formatting.
func WithComments(enabled bool) Option {
	return option(func(p *parser) {
		p.keepComments = enabled
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

// Parse returns the file parsed as an [*ast.Script] or
// an [Error] if parsing encountered one or more issues.
func Parse(file source.File, opts ...Option) (*ast.Script, error) {
	lex, err := lexer.New(file)
	if err != nil {
		var lerr lexer.Error
		if !errors.As(err, &lerr) {
			return nil, Error{
				Err: fmt.Errorf("failed to initialize lexer: failed to extract a lexer.Error from: %w", err),
				Location: source.Location{
					ByteOffset:  0,
					Length:      1,
					StartLine:   1,
					StartColumn: 1,
					EndLine:     1,
					EndColumn:   1,
				},
			}
		}
		return nil, Error{
			Err:      fmt.Errorf("failed to initialize lexer: %w", err),
			Location: lerr.Location,
		}
	}
	p := &parser{
		file:            file,
		lex:             lex,
		keepComments:    false,
		attemptRecovery: false,
		prefix:          make(map[token.Kind]prefixParser),
		infix:           make(map[token.Kind]infixParser),
	}
	for _, opt := range opts {
		opt.apply(p)
	}

	registerPrefix(p, p.ParseArrayCreation, token.New)
	registerPrefix(p, p.ParseBoolLiteral, token.True, token.False)
	registerPrefix(p, p.ParseFloatLiteral, token.FloatLiteral)
	registerPrefix(p, p.ParseIdentifier, token.Identifier, token.Self, token.Parent, token.Length)
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

	if err := p.advance(); err != nil {
		return nil, err
	}
	if err := p.advance(); err != nil {
		return nil, err
	}
	script, err := p.ParseScript()
	if err != nil {
		return nil, err
	}
	if p.keepComments {
		if err := attachLooseComments(script, p.inlineComments); err != nil {
			return nil, err
		}
	}

	return script, nil
}

type parser struct {
	file source.File
	lex  *lexer.Lexer

	token     token.Token
	lookahead token.Token

	blankLine          bool
	keepComments       bool
	inlineComments     []ast.Comment
	standaloneComments []ast.Comment

	attemptRecovery bool
	recovery        bool
	errors          []*ast.ErrorStatement

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

// tryConsume advances the token position if the current token matches the given
// token type or returns an error, consuming any comments along the way.
func (p *parser) tryConsume(t token.Kind, alts ...token.Kind) error {
	if p.token.Kind == t {
		return p.consume()
	}
	if slices.Contains(alts, p.token.Kind) {
		return p.consume()
	}
	return unexpectedTokenError(p.token, t, alts...)
}

// consume advances token and lookahead by one
// token while skipping loose comment tokens.
func (p *parser) consume() (err error) {
	suffix := p.token.Kind != token.Newline
	if err := p.advance(); err != nil {
		return err
	}
	if p.token.Kind == token.Illegal {
		return nil
	}
	// Consume loose comments immediately so the rest of the
	// parser never has to deal with them directly.
	return p.consumeComments(suffix)
}

// tryAdvance advances the token position if the current token matches the given
// token type or returns an error, but does not consume comments.
func (p *parser) tryAdvance(t token.Kind) error {
	if p.token.Kind == t {
		return p.advance()
	}
	return unexpectedTokenError(p.token, t)
}

// advance advances token and lookahead by one.
func (p *parser) advance() (err error) {
	newline := p.token.Kind == token.Newline
	p.token = p.lookahead
	p.lookahead, err = p.lex.Next()
	if err != nil {
		var lerr lexer.Error
		if !errors.As(err, &lerr) {
			return Error{
				Err:      fmt.Errorf("failed to extract a lexer.Error from: %w", err),
				Location: p.token.Location,
			}
		}
		return Error{
			Err:      err,
			Location: lerr.Location,
		}
	}
	if !p.blankLine && newline && p.token.Kind == token.Newline {
		p.blankLine = true
	}
	return nil
}

func (p *parser) consumeComments(suffix bool) (err error) {
	for p.token.Kind == token.Semicolon || p.token.Kind == token.BlockCommentOpen {
		var comment ast.Comment
		switch p.token.Kind {
		case token.Semicolon:
			if comment, err = p.ParseLineComment(suffix); err != nil {
				return err
			}
		case token.BlockCommentOpen:
			if comment, err = p.ParseBlockComment(suffix); err != nil {
				return err
			}
		}
		if p.keepComments && comment != nil {
			if !comment.Prefix() && !comment.Suffix() {
				p.standaloneComments = append(p.standaloneComments, comment)
				continue
			}
			p.inlineComments = append(p.inlineComments, comment)
			return nil
		}
		if p.token.Kind == token.Newline &&
			(p.lookahead.Kind == token.Semicolon || p.lookahead.Kind == token.BlockCommentOpen) {
			if err := p.advance(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *parser) commentStatement() *ast.CommentStatement {
	stmt := &ast.CommentStatement{
		Elements: make([]ast.Comment, 0, len(p.standaloneComments)),
	}
	stmt.Elements = append(stmt.Elements, p.standaloneComments...)
	p.standaloneComments = p.standaloneComments[:0]
	return stmt
}

func unexpectedTokenError(got token.Token, want token.Kind, alts ...token.Kind) error {
	if len(alts) > 0 {
		return newError(
			got.Location,
			"expected any of [%s, %s], but found: %s",
			want,
			tokensTypesToString(alts...),
			got.Kind,
		)
	}
	return newError(got.Location, "expected '%s', but found: %s", want, got.Kind)
}

func tokensTypesToString(kinds ...token.Kind) string {
	if len(kinds) == 0 {
		return ""
	}
	if len(kinds) == 1 {
		return kinds[0].String()
	}
	strs := make([]string, len(kinds))
	for i, t := range kinds {
		strs[i] = t.String()
	}
	return strings.Join(strs, ", ")
}

// consumeNewlines advances the token position through the as many newlines as
// possible until a non-newline token is found.
func (p *parser) consumeNewlines() error {
	for p.token.Kind == token.Newline {
		if err := p.consume(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) hasLeadingBlankLine() bool {
	l := p.blankLine
	p.blankLine = false
	return l
}

func (p *parser) ParseDocComment() (*ast.Documentation, error) {
	node := &ast.Documentation{
		OpenLocation: p.token.Location,
	}
	if err := p.tryConsume(token.BraceOpen); err != nil {
		return nil, err
	}
	node.TextLocation = p.token.Location
	if err := p.tryConsume(token.Comment); err != nil {
		return nil, err
	}
	node.CloseLocation = p.token.Location
	if err := p.tryConsume(token.BraceClose); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseBlockComment(suffix bool) (*ast.BlockComment, error) {
	node := &ast.BlockComment{
		IsPrefix:            true,
		IsSuffix:            suffix,
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		OpenLocation:        p.token.Location,
	}
	if err := p.tryAdvance(token.BlockCommentOpen); err != nil {
		return nil, err
	}
	node.TextLocation = p.token.Location
	if err := p.tryAdvance(token.Comment); err != nil {
		return nil, err
	}
	node.CloseLocation = p.token.Location
	if err := p.tryAdvance(token.BlockCommentClose); err != nil {
		return nil, err
	}
	if p.token.Kind == token.Newline || p.token.Kind == token.EOF {
		node.IsPrefix = false
		if p.lookahead.Kind == token.Newline || p.lookahead.Kind == token.EOF {
			node.HasTrailingBlankLine = true
		}
	}
	return node, nil
}

func (p *parser) ParseLineComment(suffix bool) (*ast.LineComment, error) {
	node := &ast.LineComment{
		IsSuffix:            suffix,
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		SemicolonLocation:   p.token.Location,
	}
	if err := p.tryAdvance(token.Semicolon); err != nil {
		return nil, err
	}
	node.TextLocation = p.token.Location
	if err := p.tryAdvance(token.Comment); err != nil {
		return nil, err
	}
	if (p.token.Kind == token.Newline || p.token.Kind == token.EOF) &&
		(p.lookahead.Kind == token.Newline || p.lookahead.Kind == token.EOF) {
		node.HasTrailingBlankLine = true
	}
	return node, nil
}

func (p *parser) ParseScript() (*ast.Script, error) {
	var err error
	node := &ast.Script{
		File: p.file,
		NodeLocation: source.Location{
			Length:      uint32(len(p.file.Text)), // #nosec G115 -- Checked at start of parser.Parse via lexer.New.
			StartLine:   1,
			StartColumn: 1,
		},
	}
	for p.token.Kind != token.ScriptName {
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		if err := p.consumeComments(false); err != nil {
			return nil, err
		}
		if len(p.standaloneComments) > 0 {
			node.HeaderComments = append(node.HeaderComments, p.standaloneComments...)
			p.standaloneComments = p.standaloneComments[:0]
		}
	}
	node.KeywordLocation = p.token.Location
	if err := p.tryConsume(token.ScriptName); err != nil {
		return nil, err
	}
	if node.Name, err = p.ParseIdentifier(); err != nil {
		return nil, err
	}
	if p.token.Kind == token.Extends {
		node.ExtendsLocation = p.token.Location
		if err := p.consume(); err != nil {
			return nil, err
		}
		if node.Parent, err = p.ParseIdentifier(); err != nil {
			return nil, err
		}
	}
	for p.token.Kind == token.Hidden || p.token.Kind == token.Conditional {
		if p.token.Kind == token.Hidden {
			node.HiddenLocations = append(node.HiddenLocations, p.token.Location)
		} else {
			node.ConditionalLocations = append(node.ConditionalLocations, p.token.Location)
		}
		if err := p.consume(); err != nil {
			return nil, err
		}
	}
	if err := p.consumeNewlines(); err != nil {
		return nil, err
	}
	if p.token.Kind == token.BraceOpen {
		if node.Documentation, err = p.ParseDocComment(); err != nil {
			return nil, err
		}
	}
	for p.token.Kind != token.EOF {
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		if len(p.standaloneComments) > 0 {
			node.Statements = append(node.Statements, p.commentStatement())
			continue
		}
		if p.token.Kind == token.EOF {
			break
		}
		stmt, err := p.ParseScriptStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			node.Statements = append(node.Statements, stmt)
		}
	}
	if len(p.standaloneComments) > 0 {
		node.Statements = append(node.Statements, p.commentStatement())
	}
	return node, nil
}

func (p *parser) ParseScriptStatement() (stmt ast.ScriptStatement, err error) {
	start := p.token
	switch p.token.Kind {
	case token.Import:
		stmt, err = p.ParseImport()
	case token.Event:
		stmt, err = p.ParseEvent()
	case token.Auto, token.State:
		stmt, err = p.ParseState()
	case token.Function:
		stmt, err = p.ParseFunction(nil)
	case token.Bool,
		token.Float,
		token.Int,
		token.String,
		token.Identifier:
		var typeLiteral *ast.TypeLiteral
		if typeLiteral, err = p.ParseTypeLiteral(); err != nil {
			return nil, err
		}
		switch p.token.Kind {
		case token.Property:
			stmt, err = p.ParseProperty(typeLiteral)
		case token.Function:
			stmt, err = p.ParseFunction(typeLiteral)
		case token.Identifier:
			stmt, err = p.ParseScriptVariable(typeLiteral)
		default:
			err = unexpectedTokenError(
				p.token,
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
	errStmt := &ast.ErrorStatement{
		ErrorMessage: fmt.Sprintf("%v", err),
		NodeLocation: source.Span(start.Location, p.token.Location),
	}
	p.errors = append(p.errors, errStmt)
	if err := p.consume(); err != nil {
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
		case token.Import,
			token.Event,
			token.Auto,
			token.State,
			token.Function,
			token.Bool,
			token.Float,
			token.Int,
			token.String,
			token.Identifier:
			// Next token is the start of a valid statement.
			return nil
		default:
			if err := p.consume(); err != nil {
				return err // An error during recovery just fails.
			}
		}
	}
}

func (p *parser) ParseImport() (*ast.Import, error) {
	var err error
	node := &ast.Import{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		KeywordLocation:     p.token.Location,
	}
	if err := p.tryConsume(token.Import); err != nil {
		return nil, err
	}
	if node.Name, err = p.ParseIdentifier(); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseState() (ast.ScriptStatement, error) {
	var err error
	node := &ast.State{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
	}
	start := p.token.Location
	if p.token.Kind == token.Auto {
		node.IsAuto = true
		node.AutoLocation = p.token.Location
		if err := p.consume(); err != nil {
			return nil, err
		}
	}
	node.StartKeywordLocation = p.token.Location
	if err := p.tryConsume(token.State); err != nil {
		return nil, err
	}
	if node.Name, err = p.ParseIdentifier(); err != nil {
		return nil, err
	}
	for p.token.Kind != token.EndState {
		if p.token.Kind == token.EOF {
			// State was never closed, proactively create an error statement.
			errStmt := &ast.ErrorStatement{
				ErrorMessage: fmt.Sprintf(
					"hit end of file while parsing state %q, did you forget %s?",
					node.Name.NodeLocation.Text(p.file),
					token.EndState,
				),
				NodeLocation: source.Span(start, p.token.Location),
			}
			p.errors = append(p.errors, errStmt)
			return errStmt, nil
		}
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		if len(p.standaloneComments) > 0 {
			node.Invokables = append(node.Invokables, p.commentStatement())
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
	node.EndKeywordLocation = p.token.Location
	if err := p.tryConsume(token.EndState); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseInvokable() (stmt ast.Invokable, err error) {
	start := p.token
	switch p.token.Kind {
	case token.Event:
		stmt, err = p.ParseEvent()
	case token.Function:
		stmt, err = p.ParseFunction(nil)
	case token.Bool, token.Float, token.Int, token.String, token.Identifier:
		var typeLiteral *ast.TypeLiteral
		if typeLiteral, err = p.ParseTypeLiteral(); err != nil {
			return nil, err
		}
		stmt, err = p.ParseFunction(typeLiteral)
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
	errStmt := &ast.ErrorStatement{
		ErrorMessage: fmt.Sprintf("%v", err),
		NodeLocation: source.Span(start.Location, p.token.Location),
	}
	p.errors = append(p.errors, errStmt)
	if err := p.consume(); err != nil {
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
			if err := p.consume(); err != nil {
				return err // An error during recovery just fails.
			}
		}
	}
}

func (p *parser) ParseEvent() (*ast.Event, error) {
	var err error
	node := &ast.Event{
		HasLeadingBlankLine:  p.hasLeadingBlankLine(),
		StartKeywordLocation: p.token.Location,
	}
	if err := p.tryConsume(token.Event); err != nil {
		return nil, err
	}
	if node.Name, err = p.ParseIdentifier(); err != nil {
		return nil, err
	}
	node.OpenLocation = p.token.Location
	if err := p.tryConsume(token.ParenthesisOpen); err != nil {
		return nil, err
	}
	if node.ParameterList, err = p.ParseParameterList(); err != nil {
		return nil, err
	}
	node.CloseLocation = p.token.Location
	if err := p.tryConsume(token.ParenthesisClose); err != nil {
		return nil, err
	}
	for p.token.Kind == token.Native {
		node.NativeLocations = append(node.NativeLocations, p.token.Location)
		if err := p.consume(); err != nil {
			return nil, err
		}
	}
	if p.token.Kind == token.Newline {
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		if p.token.Kind == token.BraceOpen {
			if node.Documentation, err = p.ParseDocComment(); err != nil {
				return nil, err
			}
		}
	}
	if len(node.NativeLocations) > 0 {
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		return node, nil
	}
	if node.Statements, err = p.ParseFunctionStatementBlock(token.EndEvent); err != nil {
		return nil, err
	}
	node.EndKeywordLocation = p.token.Location
	if err := p.tryConsume(token.EndEvent); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseFunction(returnType *ast.TypeLiteral) (*ast.Function, error) {
	var err error
	node := &ast.Function{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		ReturnType:          returnType, // May be nil.
	}
	node.StartKeywordLocation = p.token.Location
	if err := p.tryConsume(token.Function); err != nil {
		return nil, err
	}
	if node.Name, err = p.ParseIdentifier(); err != nil {
		return nil, err
	}
	node.OpenLocation = p.token.Location
	if err := p.tryConsume(token.ParenthesisOpen); err != nil {
		return nil, err
	}
	if node.ParameterList, err = p.ParseParameterList(); err != nil {
		return nil, err
	}
	node.CloseLocation = p.token.Location
	if err := p.tryConsume(token.ParenthesisClose); err != nil {
		return nil, err
	}
	for p.token.Kind == token.Native || p.token.Kind == token.Global {
		if p.token.Kind == token.Native {
			node.NativeLocations = append(node.NativeLocations, p.token.Location)
		} else {
			node.GlobalLocations = append(node.GlobalLocations, p.token.Location)
		}
		if err := p.consume(); err != nil {
			return nil, err
		}
	}
	if p.token.Kind == token.Newline {
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		if p.token.Kind == token.BraceOpen {
			if node.Documentation, err = p.ParseDocComment(); err != nil {
				return nil, err
			}
		}
	}
	if len(node.NativeLocations) > 0 {
		return node, nil
	}
	if node.Statements, err = p.ParseFunctionStatementBlock(token.EndFunction); err != nil {
		return nil, err
	}
	node.EndKeywordLocation = p.token.Location
	if err := p.tryConsume(token.EndFunction); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseParameterList() (params []*ast.Parameter, err error) {
	for {
		switch p.token.Kind {
		case token.Comma:
			if err := p.consume(); err != nil {
				return params, err
			}
		case token.ParenthesisClose:
			return params, nil
		default:
			param, err := p.ParseParameter()
			if err != nil {
				return params, err
			}
			params = append(params, param)
		}
	}
}

func (p *parser) ParseParameter() (*ast.Parameter, error) {
	var err error
	node := &ast.Parameter{}
	if node.Type, err = p.ParseTypeLiteral(); err != nil {
		return nil, err
	}
	if node.Name, err = p.ParseIdentifier(); err != nil {
		return nil, err
	}
	if p.token.Kind == token.Assign {
		// Has default.
		node.OperatorLocation = p.token.Location
		if err := p.tryConsume(token.Assign); err != nil {
			return nil, err
		}
		if node.DefaultValue, err = p.ParseLiteral(); err != nil {
			return nil, err
		}
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
		if len(p.standaloneComments) > 0 {
			stmts = append(stmts, p.commentStatement())
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
		errStmt := &ast.ErrorStatement{
			ErrorMessage: fmt.Sprintf("%v", err),
			NodeLocation: source.Span(start, p.token.Location),
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
			if err := p.consume(); err != nil {
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
	case token.Bool,
		token.Int,
		token.Float,
		token.String:
		return p.ParseFunctionVariable()
	case token.Identifier:
		switch p.lookahead.Kind {
		case token.Identifier, token.ArrayType:
			return p.ParseFunctionVariable()
		case token.Assign,
			token.AssignAdd,
			token.AssignDivide,
			token.AssignModulo,
			token.AssignMultiply,
			token.AssignSubtract:
			return p.ParseAssignment(nil)
		}
	}
	expr, err := p.ParseExpression(lowest)
	if err != nil {
		return nil, err
	}
	switch p.token.Kind {
	case token.Assign,
		token.AssignAdd,
		token.AssignDivide,
		token.AssignModulo,
		token.AssignMultiply,
		token.AssignSubtract:
		return p.ParseAssignment(expr)
	}
	return &ast.ExpressionStatement{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		Expression:          expr,
	}, nil
}

func (p *parser) ParseFunctionVariable() (*ast.Variable, error) {
	var err error
	node := &ast.Variable{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
	}
	if node.Type, err = p.ParseTypeLiteral(); err != nil {
		return nil, err
	}
	if node.Name, err = p.ParseIdentifier(); err != nil {
		return nil, err
	}
	if p.token.Kind == token.Assign {
		node.OperatorLocation = p.token.Location
		if err := p.tryConsume(token.Assign); err != nil {
			return nil, err
		}
		if node.Value, err = p.ParseExpression(lowest); err != nil {
			return nil, err
		}
	}
	return node, nil
}

func (p *parser) ParseAssignment(assignee ast.Expression) (node *ast.Assignment, err error) {
	if assignee == nil {
		var err error
		if assignee, err = p.ParseExpression(lowest); err != nil {
			return nil, err
		}
	}
	node = &ast.Assignment{
		Kind:                ast.AssignmentKind(p.token.Kind),
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		Assignee:            assignee,
		OperatorLocation:    p.token.Location,
	}
	if err := p.tryConsume(
		token.Assign,
		token.AssignAdd,
		token.AssignDivide,
		token.AssignModulo,
		token.AssignMultiply,
		token.AssignSubtract); err != nil {
		return nil, err
	}
	if node.Value, err = p.ParseExpression(lowest); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseReturn() (*ast.Return, error) {
	var err error
	node := &ast.Return{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		KeywordLocation:     p.token.Location,
	}
	if err := p.tryConsume(token.Return); err != nil {
		return nil, err
	}
	if p.token.Kind == token.Newline {
		return node, nil
	}
	if node.Value, err = p.ParseExpression(lowest); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseIf() (*ast.If, error) {
	var err error
	node := &ast.If{
		HasLeadingBlankLine:  p.hasLeadingBlankLine(),
		StartKeywordLocation: p.token.Location,
	}
	if err := p.tryConsume(token.If); err != nil {
		return nil, err
	}
	if node.Condition, err = p.ParseExpression(lowest); err != nil {
		return nil, err
	}
	if node.Statements, err = p.ParseFunctionStatementBlock(token.EndIf, token.Else, token.ElseIf); err != nil {
		return nil, err
	}
	for p.token.Kind == token.ElseIf {
		block := &ast.ElseIf{
			HasLeadingBlankLine: p.hasLeadingBlankLine(),
			KeywordLocation:     p.token.Location,
		}
		if err := p.tryConsume(token.ElseIf); err != nil {
			return nil, err
		}
		if block.Condition, err = p.ParseExpression(lowest); err != nil {
			return nil, err
		}
		if block.Statements, err = p.ParseFunctionStatementBlock(token.EndIf, token.Else, token.ElseIf); err != nil {
			return nil, err
		}
		node.ElseIfs = append(node.ElseIfs, block)
	}
	if p.token.Kind == token.Else {
		block := &ast.Else{
			HasLeadingBlankLine: p.hasLeadingBlankLine(),
			KeywordLocation:     p.token.Location,
		}
		if err := p.tryConsume(token.Else); err != nil {
			return nil, err
		}
		if block.Statements, err = p.ParseFunctionStatementBlock(token.EndIf); err != nil {
			return nil, err
		}
		node.Else = block
	}
	node.EndKeywordLocation = p.token.Location
	if err := p.tryConsume(token.EndIf); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseWhile() (*ast.While, error) {
	var err error
	node := &ast.While{
		HasLeadingBlankLine:  p.hasLeadingBlankLine(),
		StartKeywordLocation: p.token.Location,
	}
	if err := p.tryConsume(token.While); err != nil {
		return nil, err
	}
	if node.Condition, err = p.ParseExpression(lowest); err != nil {
		return nil, err
	}
	if node.Statements, err = p.ParseFunctionStatementBlock(token.EndWhile); err != nil {
		return nil, err
	}
	node.EndKeywordLocation = p.token.Location
	if err := p.tryConsume(token.EndWhile); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseProperty(typeLiteral *ast.TypeLiteral) (*ast.Property, error) {
	var err error
	node := &ast.Property{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		Type:                typeLiteral,
	}
	node.StartKeywordLocation = p.token.Location
	if err := p.tryConsume(token.Property); err != nil {
		return nil, err
	}
	if node.Name, err = p.ParseIdentifier(); err != nil {
		return nil, err
	}
	if p.token.Kind == token.Assign {
		node.OperatorLocation = p.token.Location
		if err := p.tryConsume(token.Assign); err != nil {
			return nil, err
		}
		if node.Value, err = p.ParseLiteral(); err != nil {
			return nil, err
		}
	}
	switch p.token.Kind {
	case token.Auto:
		node.Kind = ast.Auto
		node.AutoLocation = p.token.Location
		if err := p.tryConsume(token.Auto); err != nil {
			return nil, err
		}
	case token.AutoReadOnly:
		if node.Value == nil {
			return nil, newError(p.token.Location, "expected value to be defined for %s property", token.AutoReadOnly)
		}
		node.Kind = ast.AutoReadOnly
		node.AutoLocation = p.token.Location
		if err := p.tryConsume(token.AutoReadOnly); err != nil {
			return nil, err
		}
	}
	if node.Kind == ast.Auto || node.Kind == ast.AutoReadOnly {
		for p.token.Kind == token.Hidden || p.token.Kind == token.Conditional {
			if p.token.Kind == token.Hidden {
				node.HiddenLocations = append(node.HiddenLocations, p.token.Location)
			} else {
				node.ConditionalLocations = append(node.ConditionalLocations, p.token.Location)
			}
			if err := p.tryConsume(token.Hidden, token.Conditional); err != nil {
				return nil, err
			}
		}
		if p.token.Kind == token.Newline {
			if err := p.consumeNewlines(); err != nil {
				return nil, err
			}
			if p.token.Kind == token.BraceOpen {
				if node.Documentation, err = p.ParseDocComment(); err != nil {
					return nil, err
				}
			}
		}
		return node, nil
	}
	// Full Property
	for p.token.Kind == token.Hidden {
		node.HiddenLocations = append(node.HiddenLocations, p.token.Location)
		if err := p.tryConsume(token.Hidden); err != nil {
			return nil, err
		}
	}
	if err := p.consumeNewlines(); err != nil {
		return nil, err
	}
	if p.token.Kind == token.BraceOpen {
		if err := p.consume(); err != nil {
			return nil, err
		}
		comment, err := p.ParseDocComment()
		if err != nil {
			return nil, err
		}
		node.Documentation = comment
	}
	if err := p.consumeNewlines(); err != nil {
		return nil, err
	}
	var returnType *ast.TypeLiteral
	if p.token.Kind != token.Function {
		if returnType, err = p.ParseTypeLiteral(); err != nil {
			return nil, err
		}
	}
	first, err := p.ParseFunction(returnType)
	if err != nil {
		return nil, err
	}
	if err := p.consumeNewlines(); err != nil {
		return nil, err
	}
	var second *ast.Function
	if p.token.Kind != token.EndProperty {
		var returnType *ast.TypeLiteral
		if p.token.Kind != token.Function {
			if returnType, err = p.ParseTypeLiteral(); err != nil {
				return nil, err
			}
		}
		second, err = p.ParseFunction(returnType)
		if err != nil {
			return nil, err
		}
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
	}
	switch first.Name.Normalized {
	case "get":
		if first.ReturnType == nil {
			return nil, newError(
				first.Name.Location(),
				"expected '%s' to have a return type of %s, but found none",
				first.Name.Location().Text(p.file),
				node.Type.Location().Text(p.file),
			)
		}
		if len(first.ParameterList) != 0 {
			loc := source.Span(first.ParameterList[0].Location(), first.ParameterList[len(first.ParameterList)-1].Location())
			return nil, newError(
				loc,
				"expected '%s' to have no parameters, but found %d",
				first.Name.Location().Text(p.file),
				len(first.ParameterList),
			)
		}
		node.Get = first
	case "set":
		if first.ReturnType != nil {
			return nil, newError(
				first.ReturnType.Location(),
				"expected '%s' to have no return type, but found %s",
				first.Name.Location().Text(p.file),
				first.ReturnType.Location().Text(p.file),
			)
		}
		if len(first.ParameterList) == 0 {
			return nil, newError(
				first.Name.Location(),
				"expected '%s' to have one parameter, but found none",
				first.Name.Location().Text(p.file),
			)
		}
		if len(first.ParameterList) > 1 {
			loc := source.Span(first.ParameterList[0].Location(), first.ParameterList[len(first.ParameterList)-1].Location())
			return nil, newError(
				loc,
				"expected '%s' to have one parameter, but found %d",
				first.Name.Location().Text(p.file),
				len(first.ParameterList),
			)
		}
		node.Set = first
	default:
		return nil, newError(
			first.Location(),
			"expected 'Get' or 'Set' function for property, but found '%s'",
			first.Name.Location().Text(p.file),
		)
	}
	if second != nil {
		switch second.Name.Normalized {
		case "get":
			if node.Get != nil {
				return nil, newError(second.Location(), "expected exactly one 'Get' function, but found two")
			}
			if second.ReturnType == nil {
				return nil, newError(
					second.Name.Location(),
					"expected '%s' to have a return type of %s, but found none",
					second.Name.Location().Text(p.file),
					node.Type.Location().Text(p.file),
				)
			}
			if len(second.ParameterList) != 0 {
				loc := source.Span(
					second.ParameterList[0].Location(),
					second.ParameterList[len(second.ParameterList)-1].Location(),
				)
				return nil, newError(
					loc,
					"expected '%s' to have no parameters, but found %d",
					second.Name.Location().Text(p.file),
					len(second.ParameterList),
				)
			}
			node.Get = second
		case "set":
			if node.Set != nil {
				return nil, newError(second.Location(), "expected exactly one 'Set' function, but found two")
			}
			if second.ReturnType != nil {
				return nil, newError(
					second.ReturnType.Location(),
					"expected '%s' to have no return type, but found %s",
					second.Name.Location().Text(p.file),
					second.ReturnType.Location().Text(p.file),
				)
			}
			if len(second.ParameterList) == 0 {
				return nil, newError(
					second.Name.Location(),
					"expected '%s' to have one parameter, but found none",
					second.Name.Location().Text(p.file),
				)
			}
			if len(second.ParameterList) > 1 {
				loc := source.Span(
					second.ParameterList[0].Location(),
					second.ParameterList[len(second.ParameterList)-1].Location(),
				)
				return nil, newError(
					loc,
					"expected '%s' to have one parameter, but found %d",
					second.Name.Location().Text(p.file),
					len(second.ParameterList),
				)
			}
			node.Set = second
		default:
			return nil, newError(
				second.Location(),
				"expected 'Get' or 'Set' function for property, but found '%s'",
				second.Name.Location().Text(p.file),
			)
		}
	}
	node.EndKeywordLocation = p.token.Location
	if err := p.tryConsume(token.EndProperty); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseScriptVariable(typeLiteral *ast.TypeLiteral) (*ast.Variable, error) {
	var err error
	node := &ast.Variable{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		Type:                typeLiteral,
	}
	if node.Name, err = p.ParseIdentifier(); err != nil {
		return nil, err
	}
	if p.token.Kind == token.Assign {
		node.OperatorLocation = p.token.Location
		if err := p.tryConsume(token.Assign); err != nil {
			return nil, err
		}
		if node.Value, err = p.ParseLiteral(); err != nil {
			return nil, err
		}
	}
	for p.token.Kind == token.Conditional {
		node.ConditionalLocations = append(node.ConditionalLocations, p.token.Location)
		if err := p.consume(); err != nil {
			return nil, err
		}
	}
	if err := p.tryConsume(token.Newline, token.EOF); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseIdentifier() (*ast.Identifier, error) {
	node := &ast.Identifier{
		Normalized:   string(bytes.ToLower(p.token.Text)),
		NodeLocation: p.token.Location,
	}
	if err := p.tryConsume(token.Identifier, token.Self, token.Parent, token.Length); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseTypeLiteral() (node *ast.TypeLiteral, err error) {
	node = &ast.TypeLiteral{
		Name: &ast.Identifier{
			Normalized:   string(bytes.ToLower(p.token.Text)),
			NodeLocation: p.token.Location,
		},
	}
	if err := p.tryConsume(token.Identifier, token.Bool, token.Int, token.Float, token.String); err != nil {
		return nil, err
	}
	if p.token.Kind == token.ArrayType {
		node.BracketLocation = p.token.Location
		node.IsArray = true
		if err := p.consume(); err != nil {
			return nil, err
		}
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

func (p *parser) ParseBinary(left ast.Expression) (node *ast.Binary, err error) {
	precedence := precedenceOf(p.token.Kind)
	node = &ast.Binary{
		Kind:             ast.BinaryKind(p.token.Kind),
		LeftOperand:      left,
		OperatorLocation: p.token.Location,
	}
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
	if node.RightOperand, err = p.ParseExpression(precedence); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseUnary() (ast.Expression, error) {
	if p.token.Kind == token.Minus &&
		(p.lookahead.Kind == token.IntLiteral || p.lookahead.Kind == token.FloatLiteral) {
		return p.ParseLiteral()
	}
	var err error
	node := &ast.Unary{
		Kind:             ast.UnaryKind(p.token.Kind),
		OperatorLocation: p.token.Location,
	}
	if err := p.tryConsume(token.Minus, token.LogicalNot); err != nil {
		return nil, err
	}
	if node.Operand, err = p.ParseExpression(prefix); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseCast(value ast.Expression) (node *ast.Cast, err error) {
	node = &ast.Cast{
		Value:      value,
		AsLocation: p.token.Location,
	}
	if err := p.tryConsume(token.As); err != nil {
		return nil, err
	}
	if node.Type, err = p.ParseTypeLiteral(); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseAccess(value ast.Expression) (node *ast.Access, err error) {
	node = &ast.Access{
		Value:       value,
		DotLocation: p.token.Location,
	}
	if err := p.tryConsume(token.Dot); err != nil {
		return nil, err
	}
	if node.Name, err = p.ParseIdentifier(); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseIndex(array ast.Expression) (node *ast.Index, err error) {
	node = &ast.Index{
		Value:        array,
		OpenLocation: p.token.Location,
	}
	if err := p.tryConsume(token.BracketOpen); err != nil {
		return nil, err
	}
	if node.Index, err = p.ParseExpression(lowest); err != nil {
		return nil, err
	}
	node.CloseLocation = p.token.Location
	if err := p.tryConsume(token.BracketClose); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseCall(function ast.Expression) (node *ast.Call, err error) {
	node = &ast.Call{
		Function:     function,
		OpenLocation: p.token.Location,
	}
	if err := p.tryConsume(token.ParenthesisOpen); err != nil {
		return nil, err
	}
	if node.Arguments, err = p.ParseArgumentList(); err != nil {
		return nil, err
	}
	node.CloseLocation = p.token.Location
	if err := p.tryConsume(token.ParenthesisClose); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseArgumentList() ([]*ast.Argument, error) {
	var args []*ast.Argument
	for {
		switch p.token.Kind {
		case token.Comma:
			if err := p.consume(); err != nil {
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
		node.OperatorLocation = p.token.Location
		if err := p.tryConsume(token.Assign); err != nil {
			return nil, err
		}
	}
	value, err := p.ParseExpression(lowest)
	if err != nil {
		return nil, err
	}
	node.Value = value
	return node, nil
}

func (p *parser) ParseArrayCreation() (node *ast.ArrayCreation, err error) {
	node = &ast.ArrayCreation{
		NewLocation: p.token.Location,
	}
	if err := p.tryConsume(token.New); err != nil {
		return nil, err
	}
	if node.Type, err = p.ParseTypeLiteral(); err != nil {
		return nil, err
	}
	node.OpenLocation = p.token.Location
	if err := p.tryConsume(token.BracketOpen); err != nil {
		return nil, err
	}
	if node.Size, err = p.ParseIntLiteral(); err != nil {
		return nil, err
	}
	node.CloseLocation = p.token.Location
	if err := p.tryConsume(token.BracketClose); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseParenthetical() (*ast.Parenthetical, error) {
	var err error
	node := &ast.Parenthetical{
		OpenLocation: p.token.Location,
	}
	if err := p.tryConsume(token.ParenthesisOpen); err != nil {
		return nil, err
	}
	if node.Value, err = p.ParseExpression(lowest); err != nil {
		return nil, err
	}
	node.CloseLocation = p.token.Location
	if err := p.tryConsume(token.ParenthesisClose); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseLiteral() (ast.Literal, error) {
	switch p.token.Kind {
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
		RawText:      p.token.Location.Text(p.file),
		NodeLocation: p.token.Location,
	}
	if err := p.tryConsume(token.IntLiteral); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseFloatLiteral() (*ast.FloatLiteral, error) {
	node := &ast.FloatLiteral{
		RawText:      p.token.Location.Text(p.file),
		NodeLocation: p.token.Location,
	}
	if err := p.tryConsume(token.FloatLiteral); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseBoolLiteral() (*ast.BoolLiteral, error) {
	node := &ast.BoolLiteral{
		RawText:      p.token.Location.Text(p.file),
		NodeLocation: p.token.Location,
	}
	if err := p.tryConsume(token.True, token.False); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseStringLiteral() (*ast.StringLiteral, error) {
	node := &ast.StringLiteral{
		RawText:      p.token.Location.Text(p.file),
		NodeLocation: p.token.Location,
	}
	if err := p.tryConsume(token.StringLiteral); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseNoneLiteral() (*ast.NoneLiteral, error) {
	node := &ast.NoneLiteral{
		RawText:      p.token.Location.Text(p.file),
		NodeLocation: p.token.Location,
	}
	if err := p.tryConsume(token.None); err != nil {
		return nil, err
	}
	return node, nil
}

func registerPrefix[T ast.Expression](p *parser, fn func() (T, error), kinds ...token.Kind) {
	for _, t := range kinds {
		p.prefix[t] = func() (ast.Expression, error) { return fn() }
	}
}

func registerInfix[T ast.Expression](p *parser, fn func(ast.Expression) (T, error), kinds ...token.Kind) {
	for _, t := range kinds {
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
