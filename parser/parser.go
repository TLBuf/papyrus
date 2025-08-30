// Package parser defines a Papyrus parser.
package parser

import (
	"slices"
	"strings"

	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/issue"
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

// Parse returns the file parsed as an [*ast.Script] or false if parsing failed.
// If this returns false, the log is guarnteed to contain at least one issue.
func Parse(file *source.File, log *issue.Log, opts ...Option) (script *ast.Script, ok bool) {
	lex := lexer.New(file, log)
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

	registerPrefix(p, p.ParseExpressionArrayCreation, token.New)
	registerPrefix(p, p.ParseExpressionBoolLiteral, token.True, token.False)
	registerPrefix(p, p.ParseExpressionFloatLiteral, token.FloatLiteral)
	registerPrefix(p, p.ParseExpressionIdentifier, token.Identifier, token.Self, token.Parent, token.Length)
	registerPrefix(p, p.ParseExpressionIntLiteral, token.IntLiteral)
	registerPrefix(p, p.ParseExpressionNoneLiteral, token.None)
	registerPrefix(p, p.ParseExpressionParenthetical, token.ParenthesisOpen)
	registerPrefix(p, p.ParseExpressionStringLiteral, token.StringLiteral)
	registerPrefix(p, p.ParseExpressionUnary, token.Minus, token.LogicalNot)

	registerInfix(p, p.ParseExpressionAccess, token.Dot)
	registerInfix(p, p.ParseExpressionCall, token.ParenthesisOpen)
	registerInfix(p, p.ParseExpressionCast, token.As)
	registerInfix(p, p.ParseExpressionIndex, token.BracketOpen)
	registerInfix(p,
		p.ParseExpressionBinary,
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

	defer func() {
		if r := recover(); r != nil {
			// If we failed at this level, recovery
			// failed, so don't return a broken Script.
			log.Append(r.(*issue.Issue))
			script = nil
			ok = false
		}
	}()

	p.advance()
	p.advance()

	script = p.ParseScript()
	if log.HasInternal() || log.HasError() {
		return script, false
	}
	if p.keepComments {
		p.attachLooseComments(script, p.inlineComments)
	}
	return script, true
}

type parser struct {
	log  *issue.Log
	file *source.File
	lex  *lexer.Lexer

	token     token.Token
	lookahead token.Token

	blankLine          bool
	keepComments       bool
	inlineComments     []ast.Comment
	standaloneComments []ast.Comment

	fatal           bool
	attemptRecovery bool
	recovery        bool

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
	prefixParser func() ast.Expression
	infixParser  func(ast.Expression) ast.Expression
)

// tryConsume advances the token position if the current token matches the given
// token type or returns an error, consuming any comments along the way.
func (p *parser) tryConsume(def *issue.Definition, t token.Kind, alts ...token.Kind) {
	if p.token.Kind == t || slices.Contains(alts, p.token.Kind) {
		p.consume()
		return
	}
	p.failWithDetail(def, p.token.Location, "Encountered %s %s token.", p.token.Kind.Article(), p.token.Kind)
}

// consumeExpected advances the token position if the current token matches the
// given token type while skipping loose comment tokens or raises an internal
// issue. This is used for cases where the parser should be guarnteed to always
// find a match even if the input is malformed.
func (p *parser) consumeExpected(t token.Kind, alts ...token.Kind) {
	if p.token.Kind == t || slices.Contains(alts, p.token.Kind) {
		p.consume()
		return
	}
	p.fatal = true // Force unexpected parser state to skip recovery.
	p.unexpectedToken(intenalInvalidState, p.token, t, alts...)
}

// consume advances token and lookahead by one
// token while skipping loose comment tokens.
func (p *parser) consume() {
	suffix := p.token.Kind != token.Newline
	p.advance()
	if p.token.Kind == token.Illegal {
		return
	}
	// Consume loose comments immediately so the rest of the
	// parser never has to deal with them directly.
	p.consumeComments(suffix)
}

// advanceExpected advances the token position if the current token matches the
// given token type or raises an internal issue. This is used for cases where
// the parser should be guarnteed to always find a match even if the input is
// malformed.
func (p *parser) advanceExpected(t token.Kind, alts ...token.Kind) {
	if p.token.Kind == t {
		p.advance()
		return
	}
	for _, alt := range alts {
		if p.token.Kind == alt {
			p.advance()
			return
		}
	}
	p.fatal = true // Force unexpected parser state to skip recovery.
	p.unexpectedToken(intenalInvalidState, p.token, t, alts...)
}

// advance advances token and lookahead by one.
func (p *parser) advance() {
	newline := p.token.Kind == token.Newline
	p.token = p.lookahead
	tok, ok := p.lex.Next()
	if !ok {
		// Force lexer errors to skip recovery.
		p.fatal = true
		panic(p.log.Last())
	}
	p.lookahead = tok
	if !p.blankLine && newline && p.token.Kind == token.Newline {
		p.blankLine = true
	}
}

func (p *parser) consumeComments(suffix bool) {
	for p.token.Kind == token.Semicolon || p.token.Kind == token.BlockCommentOpen {
		var comment ast.Comment
		switch p.token.Kind {
		case token.Semicolon:
			comment = p.ParseLineComment(suffix)
		case token.BlockCommentOpen:
			comment = p.ParseBlockComment(suffix)
		}
		if p.keepComments && comment != nil {
			if !comment.Prefix() && !comment.Suffix() {
				p.standaloneComments = append(p.standaloneComments, comment)
				continue
			}
			p.inlineComments = append(p.inlineComments, comment)
			return
		}
		if p.token.Kind == token.Newline &&
			(p.lookahead.Kind == token.Semicolon || p.lookahead.Kind == token.BlockCommentOpen) {
			p.advance()
		}
	}
}

func (p *parser) commentStatement() *ast.CommentStatement {
	stmt := &ast.CommentStatement{
		Elements: make([]ast.Comment, 0, len(p.standaloneComments)),
	}
	stmt.Elements = append(stmt.Elements, p.standaloneComments...)
	p.standaloneComments = p.standaloneComments[:0]
	return stmt
}

func (p *parser) unexpectedToken(def *issue.Definition, got token.Token, want token.Kind, alts ...token.Kind) {
	if len(alts) > 0 {
		var detail strings.Builder
		_, _ = detail.WriteString("Expected ")
		_, _ = detail.WriteString(want.Article())
		_, _ = detail.WriteRune(' ')
		_, _ = detail.WriteString(want.String())
		for _, alt := range alts[:len(alts)-1] {
			_, _ = detail.WriteString(", ")
			_, _ = detail.WriteString(alt.String())
		}
		_, _ = detail.WriteString("or ")
		_, _ = detail.WriteString(alts[len(alts)-1].String())
		_, _ = detail.WriteString(" token, but encountered ")
		_, _ = detail.WriteString(got.Kind.Article())
		_, _ = detail.WriteRune(' ')
		_, _ = detail.WriteString(got.String())
		_, _ = detail.WriteString(" token.")
		p.failWithDetail(def, got.Location, "%v", detail)
	}
	p.failWithDetail(
		def,
		got.Location,
		"Expected %s token, but encountered %s %s token.",
		want,
		got.Kind.Article(),
		got.Kind,
	)
}

// consumeNewlines advances the token position through the as many newlines as
// possible until a non-newline token is found.
func (p *parser) consumeNewlines() {
	for p.token.Kind == token.Newline {
		p.consume()
	}
}

func (p *parser) hasLeadingBlankLine() bool {
	l := p.blankLine
	p.blankLine = false
	return l
}

func (p *parser) ParseDocComment() *ast.Documentation {
	node := &ast.Documentation{
		OpenLocation: p.token.Location,
	}
	p.advanceExpected(token.BraceOpen)
	node.TextLocation = p.token.Location
	p.advanceExpected(token.Comment)
	node.CloseLocation = p.token.Location
	p.advanceExpected(token.BraceClose)
	return node
}

func (p *parser) ParseBlockComment(suffix bool) *ast.BlockComment {
	node := &ast.BlockComment{
		IsPrefix:            true,
		IsSuffix:            suffix,
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		OpenLocation:        p.token.Location,
	}
	p.advanceExpected(token.BlockCommentOpen)
	node.TextLocation = p.token.Location
	p.advanceExpected(token.Comment)
	node.CloseLocation = p.token.Location
	p.advanceExpected(token.BlockCommentClose)
	if p.token.Kind == token.Newline || p.token.Kind == token.EOF {
		node.IsPrefix = false
		if p.lookahead.Kind == token.Newline || p.lookahead.Kind == token.EOF {
			node.HasTrailingBlankLine = true
		}
	}
	return node
}

func (p *parser) ParseLineComment(suffix bool) *ast.LineComment {
	node := &ast.LineComment{
		IsSuffix:            suffix,
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		SemicolonLocation:   p.token.Location,
	}
	p.advanceExpected(token.Semicolon)
	node.TextLocation = p.token.Location
	p.advanceExpected(token.Comment)
	if (p.token.Kind == token.Newline || p.token.Kind == token.EOF) &&
		(p.lookahead.Kind == token.Newline || p.lookahead.Kind == token.EOF) {
		node.HasTrailingBlankLine = true
	}
	return node
}

func (p *parser) ParseScript() *ast.Script {
	node := &ast.Script{
		File:         p.file,
		NodeLocation: source.NewLocation(0, p.file.Len()),
	}
	for p.token.Kind != token.ScriptName {
		p.consumeNewlines()
		p.consumeComments(false)
		if len(p.standaloneComments) > 0 {
			node.HeaderComments = append(node.HeaderComments, p.standaloneComments...)
			p.standaloneComments = p.standaloneComments[:0]
		}
	}
	node.KeywordLocation = p.token.Location
	p.tryConsume(errorExpectedScriptName, token.ScriptName)
	node.Name = p.ParseIdentifier(errorExpectedScriptNameIdent)
	if p.token.Kind == token.Extends {
		node.ExtendsLocation = p.token.Location
		p.consume()
		node.Parent = p.ParseIdentifier(errorExpectedExtendsIdent)
	}
	for p.token.Kind == token.Hidden || p.token.Kind == token.Conditional {
		if p.token.Kind == token.Hidden {
			node.HiddenLocations = append(node.HiddenLocations, p.token.Location)
		} else {
			node.ConditionalLocations = append(node.ConditionalLocations, p.token.Location)
		}
		p.consume()
	}
	p.consumeNewlines()
	if p.token.Kind == token.BraceOpen {
		node.Documentation = p.ParseDocComment()
	}
	for p.token.Kind != token.EOF {
		p.consumeNewlines()
		if len(p.standaloneComments) > 0 {
			node.Statements = append(node.Statements, p.commentStatement())
			continue
		}
		if p.token.Kind == token.EOF {
			break
		}
		stmt := p.ParseScriptStatement()
		if stmt != nil {
			node.Statements = append(node.Statements, stmt)
		}
	}
	if len(p.standaloneComments) > 0 {
		node.Statements = append(node.Statements, p.commentStatement())
	}
	return node
}

func (p *parser) ParseScriptStatement() (stmt ast.ScriptStatement) {
	start := p.token
	if p.attemptRecovery {
		// Error recovery. Attempt to realign to a known statement token and emit an
		// error statement to fill the gap.
		defer func() {
			if r := recover(); r != nil {
				if p.fatal || p.recovery || r.(*issue.Issue).Definition().Severity() == issue.Internal {
					// If a recovery fails or it's an internal issue, just propagate.
					panic(r)
				}
				p.recovery = true
				p.recoverScriptStatement()
				stmt = &ast.ErrorStatement{
					Issue:        r.(*issue.Issue),
					NodeLocation: source.Span(start.Location, p.token.Location),
				}
				p.consume()
				p.recovery = false
			}
		}()
	}

	switch p.token.Kind {
	case token.Import:
		return p.ParseImport()
	case token.Event:
		return p.ParseEvent()
	case token.Auto, token.State:
		return p.ParseState()
	case token.Function:
		return p.ParseFunction(nil)
	case token.Bool,
		token.Float,
		token.Int,
		token.String,
		token.Identifier:
		typeLiteral := p.ParseTypeLiteral(intenalInvalidState)
		switch p.token.Kind {
		case token.Property:
			return p.ParseProperty(typeLiteral)
		case token.Function:
			return p.ParseFunction(typeLiteral)
		case token.Identifier:
			return p.ParseScriptVariable(typeLiteral)
		default:
			p.unexpectedToken(
				errorExpectedScriptStatementKeyword,
				p.token,
				token.Property,
				token.Function,
				token.Identifier)
		}
	default:
		p.unexpectedToken(
			errorExpectedScriptStatement,
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
	return nil
}

func (p *parser) recoverScriptStatement() {
	for {
		switch p.lookahead.Kind {
		case token.EOF:
			// Hit end of file, give up.
			return
		case token.Import,
			token.Event,
			token.Auto,
			token.State,
			token.Function:
			// Next token is the start of a script statement.
		case token.Bool,
			token.Float,
			token.Int,
			token.String,
			token.Identifier:
			if p.token.Kind == token.Newline {
				return // Next token is likely the start of a script statement.
			}
		}
		p.consume()
	}
}

func (p *parser) ParseImport() *ast.Import {
	node := &ast.Import{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		KeywordLocation:     p.token.Location,
	}
	p.consumeExpected(token.Import)
	node.Name = p.ParseIdentifier(errorExpectedImportIdent)
	p.tryConsume(errorExpectedImportEnd, token.Newline, token.EOF)
	return node
}

func (p *parser) ParseState() ast.ScriptStatement {
	node := &ast.State{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
	}
	start := p.token.Location
	if p.token.Kind == token.Auto {
		node.IsAuto = true
		node.AutoLocation = p.token.Location
		p.consume()
	}
	node.StartKeywordLocation = p.token.Location
	p.tryConsume(errorExpectedAutoStateKeyword, token.State)
	node.Name = p.ParseIdentifier(errorExpectedStateIdent)
	for p.token.Kind != token.EndState {
		if p.token.Kind == token.EOF {
			// State was never closed, proactively create an error statement.
			loc := source.Span(start, p.token.Location)
			err := issue.New(errorUnclosedState, p.file, loc)
			p.log.Append(err)
			stmt := &ast.ErrorStatement{
				Issue:        err,
				NodeLocation: loc,
			}
			return stmt
		}
		p.consumeNewlines()
		if len(p.standaloneComments) > 0 {
			node.Invokables = append(node.Invokables, p.commentStatement())
		}
		if p.token.Kind == token.EndState {
			break
		}
		node.Invokables = append(node.Invokables, p.ParseInvokable())
	}
	node.EndKeywordLocation = p.token.Location
	p.consumeExpected(token.EndState)
	p.tryConsume(errorStateEnd, token.Newline, token.EOF)
	return node
}

func (p *parser) ParseInvokable() (stmt ast.Invokable) {
	start := p.token
	if p.attemptRecovery {
		// Error recovery. Attempt to realign to a known statement token and emit an
		// error statement to fill the gap.
		defer func() {
			if r := recover(); r != nil {
				if p.fatal || p.recovery || r.(*issue.Issue).Definition().Severity() == issue.Internal {
					// If a recovery fails or it's an internal issue, just propagate.
					panic(r)
				}
				p.recovery = true
				p.recoverInvokable()
				stmt = &ast.ErrorStatement{
					Issue:        r.(*issue.Issue),
					NodeLocation: source.Span(start.Location, p.token.Location),
				}
				p.consume()
				p.recovery = false
			}
		}()
	}

	switch p.token.Kind {
	case token.Event:
		return p.ParseEvent()
	case token.Function:
		return p.ParseFunction(nil)
	case token.Bool,
		token.Float,
		token.Int,
		token.String,
		token.Identifier:
		typeLiteral := p.ParseTypeLiteral(intenalInvalidState)
		switch p.token.Kind {
		case token.Function:
			return p.ParseFunction(typeLiteral)
		default:
			p.unexpectedToken(
				errorExpectedStateStatementKeyword,
				p.token,
				token.Function)
		}
	default:
		p.unexpectedToken(
			errorExpectedStateStatement,
			p.token,
			token.Event,
			token.Function,
			token.Bool,
			token.Float,
			token.Int,
			token.String,
			token.Identifier)
	}
	return nil
}

func (p *parser) recoverInvokable() {
	for {
		switch p.lookahead.Kind {
		case token.EOF:
			// Hit end of file, give up.
			return
		case token.EndState:
			// Hit end of state, give up.
			return
		case token.Event, token.Function:
			// Next token is the start of a statement.
			return
		case token.Bool,
			token.Float,
			token.Int,
			token.String,
			token.Identifier:
			if p.token.Kind == token.Newline {
				return // Next token is likely the start of a valid statement.
			}
		}
		p.consume()
	}
}

func (p *parser) ParseEvent() *ast.Event {
	node := &ast.Event{
		HasLeadingBlankLine:  p.hasLeadingBlankLine(),
		StartKeywordLocation: p.token.Location,
	}
	p.consumeExpected(token.Event)
	node.Name = p.ParseIdentifier(errorExpectedEventIdent)
	node.OpenLocation = p.token.Location
	p.tryConsume(errorExpectedEventOpenParen, token.ParenthesisOpen)
	node.ParameterList = p.ParseParameterList()
	node.CloseLocation = p.token.Location
	p.consumeExpected(token.ParenthesisClose)
	for p.token.Kind == token.Native {
		node.NativeLocations = append(node.NativeLocations, p.token.Location)
		p.consume()
	}
	if p.token.Kind == token.Newline {
		p.consumeNewlines()
		if p.token.Kind == token.BraceOpen {
			node.Documentation = p.ParseDocComment()
		}
	}
	if len(node.NativeLocations) > 0 {
		p.consumeNewlines()
		return node
	}
	node.Statements = p.ParseFunctionStatementBlock(node.StartKeywordLocation, errorUnclosedEvent, token.EndEvent)
	node.EndKeywordLocation = p.token.Location
	p.consumeExpected(token.EndEvent)
	p.tryConsume(errorEventEnd, token.Newline, token.EOF)
	return node
}

func (p *parser) ParseFunction(returnType *ast.TypeLiteral) *ast.Function {
	node := &ast.Function{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		ReturnType:          returnType, // May be nil.
	}
	node.StartKeywordLocation = p.token.Location
	p.consumeExpected(token.Function)
	node.Name = p.ParseIdentifier(errorExpectedFunctionIdent)
	node.OpenLocation = p.token.Location
	p.tryConsume(errorExpectedFunctionOpenParen, token.ParenthesisOpen)
	node.ParameterList = p.ParseParameterList()
	node.CloseLocation = p.token.Location
	p.consumeExpected(token.ParenthesisClose)
	for p.token.Kind == token.Native || p.token.Kind == token.Global {
		if p.token.Kind == token.Native {
			node.NativeLocations = append(node.NativeLocations, p.token.Location)
		} else {
			node.GlobalLocations = append(node.GlobalLocations, p.token.Location)
		}
		p.consume()
	}
	if p.token.Kind == token.Newline {
		p.consumeNewlines()
		if p.token.Kind == token.BraceOpen {
			node.Documentation = p.ParseDocComment()
		}
	}
	if len(node.NativeLocations) > 0 {
		return node
	}
	node.Statements = p.ParseFunctionStatementBlock(node.StartKeywordLocation, errorUnclosedFunction, token.EndFunction)
	node.EndKeywordLocation = p.token.Location
	p.consumeExpected(token.EndFunction)
	p.tryConsume(errorFunctionEnd, token.Newline, token.EOF)
	return node
}

func (p *parser) ParseParameterList() []*ast.Parameter {
	var params []*ast.Parameter
	for {
		switch p.token.Kind {
		case token.EOF:
			p.fail(errorUnclosedParamListEOF, p.token.Location)
		case token.Newline:
			p.fail(errorUnclosedParamListNewline, p.token.Location)
		case token.Comma:
			p.consume()
		case token.ParenthesisClose:
			return params
		default:
			params = append(params, p.ParseParameter())
		}
	}
}

func (p *parser) ParseParameter() *ast.Parameter {
	node := &ast.Parameter{}
	p.ParseTypeLiteral(errorExpectedParamTypeLiteral)
	node.Name = p.ParseIdentifier(errorExpectedParamIdent)
	if p.token.Kind == token.Assign {
		// Has default.
		node.OperatorLocation = p.token.Location
		p.advanceExpected(token.Assign)
		node.DefaultValue = p.ParseLiteral(errorExpectedParamLiteral)
	}
	return node
}

func (p *parser) ParseFunctionStatementBlock(
	start source.Location,
	unclosed *issue.Definition,
	terminals ...token.Kind,
) []ast.FunctionStatement {
	var stmts []ast.FunctionStatement
	for {
		if p.token.Kind == token.EOF {
			p.fail(unclosed, source.Span(start, p.token.Location))
		}
		p.consumeNewlines()
		if len(p.standaloneComments) > 0 {
			stmts = append(stmts, p.commentStatement())
		}
		if slices.Contains(terminals, p.token.Kind) {
			return stmts
		}
		stmts = append(stmts, p.ParseFunctionStatement())
	}
}

func (p *parser) ParseFunctionStatement() (stmt ast.FunctionStatement) {
	start := p.token
	if p.attemptRecovery {
		// Error recovery. Attempt to realign to a known statement token and emit an
		// error statement to fill the gap.
		defer func() {
			if r := recover(); r != nil {
				if p.fatal || p.recovery || r.(*issue.Issue).Definition().Severity() == issue.Internal {
					// If a recovery fails or it's an internal issue, just propagate.
					panic(r)
				}
				p.recovery = true
				p.recoverFunctionStatement()
				stmt = &ast.ErrorStatement{
					Issue:        r.(*issue.Issue),
					NodeLocation: source.Span(start.Location, p.token.Location),
				}
				p.recovery = false
			}
		}()
	}
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
	expr := p.ParseExpression(errorExpectedFunctionStatementExpr, lowest)
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
	}
}

func (p *parser) recoverFunctionStatement() {
	for {
		switch p.lookahead.Kind {
		case token.EOF:
			// Hit end of file, give up.
			return
		case token.Newline:
			// Next token is the start of a statement.
			return
		default:
			p.consume()
		}
	}
}

func (p *parser) ParseFunctionVariable() *ast.Variable {
	node := &ast.Variable{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
	}
	node.Type = p.ParseTypeLiteral(intenalInvalidState)
	node.Name = p.ParseIdentifier(errorExpectedFunctionVariableIdent)
	if p.token.Kind == token.Assign {
		node.OperatorLocation = p.token.Location
		p.consume()
		node.Value = p.ParseExpression(errorExpectedFunctionVariableExpr, lowest)
	}
	return node
}

func (p *parser) ParseAssignment(assignee ast.Expression) *ast.Assignment {
	if assignee == nil {
		assignee = p.ParseExpression(errorExpectedAssignmentAssigneeExpr, lowest)
	}
	node := &ast.Assignment{
		Kind:                ast.AssignmentKind(p.token.Kind),
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		Assignee:            assignee,
		OperatorLocation:    p.token.Location,
	}
	p.consumeExpected(
		token.Assign,
		token.AssignAdd,
		token.AssignDivide,
		token.AssignModulo,
		token.AssignMultiply,
		token.AssignSubtract)
	p.ParseExpression(errorExpectedAssignmentValueExpr, lowest)
	return node
}

func (p *parser) ParseReturn() *ast.Return {
	node := &ast.Return{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		KeywordLocation:     p.token.Location,
	}
	p.consumeExpected(token.Return)
	if p.token.Kind == token.Newline || p.token.Kind == token.EOF {
		return node
	}
	node.Value = p.ParseExpression(errorExpectedReturnExpr, lowest)
	return node
}

func (p *parser) ParseIf() *ast.If {
	node := &ast.If{
		HasLeadingBlankLine:  p.hasLeadingBlankLine(),
		StartKeywordLocation: p.token.Location,
	}
	p.consumeExpected(token.If)
	node.Condition = p.ParseExpression(errorExpectedIfConditionExpr, lowest)
	node.Statements = p.ParseFunctionStatementBlock(
		node.StartKeywordLocation,
		errorUnclosedIf,
		token.EndIf,
		token.Else,
		token.ElseIf,
	)
	for p.token.Kind == token.ElseIf {
		block := &ast.ElseIf{
			HasLeadingBlankLine: p.hasLeadingBlankLine(),
			KeywordLocation:     p.token.Location,
		}
		p.consume()
		block.Condition = p.ParseExpression(errorExpectedElseIfConditionExpr, lowest)
		block.Statements = p.ParseFunctionStatementBlock(
			block.KeywordLocation,
			errorUnclosedElseIf,
			token.EndIf,
			token.Else,
			token.ElseIf,
		)
		node.ElseIfs = append(node.ElseIfs, block)
	}
	if p.token.Kind == token.Else {
		block := &ast.Else{
			HasLeadingBlankLine: p.hasLeadingBlankLine(),
			KeywordLocation:     p.token.Location,
		}
		p.consume()
		block.Statements = p.ParseFunctionStatementBlock(block.KeywordLocation, errorUnclosedElse, token.EndIf)
		node.Else = block
	}
	node.EndKeywordLocation = p.token.Location
	p.consumeExpected(token.EndIf)
	return node
}

func (p *parser) ParseWhile() *ast.While {
	node := &ast.While{
		HasLeadingBlankLine:  p.hasLeadingBlankLine(),
		StartKeywordLocation: p.token.Location,
	}
	p.consumeExpected(token.While)
	node.Condition = p.ParseExpression(errorExpectedWhileExpr, lowest)
	node.Statements = p.ParseFunctionStatementBlock(node.StartKeywordLocation, errorUnclosedWhile, token.EndWhile)
	node.EndKeywordLocation = p.token.Location
	p.consumeExpected(token.EndWhile)
	return node
}

func (p *parser) ParseProperty(typeLiteral *ast.TypeLiteral) *ast.Property {
	start := p.token.Location
	if typeLiteral != nil {
		start = typeLiteral.Location()
	}
	node := &ast.Property{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		Type:                typeLiteral,
	}
	node.StartKeywordLocation = p.token.Location
	p.consumeExpected(token.Property)
	node.Name = p.ParseIdentifier(errorExpectedPropertyIdent)
	if p.token.Kind == token.Assign {
		node.OperatorLocation = p.token.Location
		p.consume()
		node.Value = p.ParseLiteral(errorExpectedPropertyLiteral)
	}
	switch p.token.Kind {
	case token.Auto:
		node.Kind = ast.Auto
		node.AutoLocation = p.token.Location
		p.consume()
	case token.AutoReadOnly:
		if node.Value == nil {
			p.fail(errorExpectedPropertyReadOnlyValue, source.Span(start, p.token.Location))
		}
		node.Kind = ast.AutoReadOnly
		node.AutoLocation = p.token.Location
		p.consume()
	}
	if node.Kind == ast.Auto || node.Kind == ast.AutoReadOnly {
		for p.token.Kind == token.Hidden || p.token.Kind == token.Conditional {
			if p.token.Kind == token.Hidden {
				node.HiddenLocations = append(node.HiddenLocations, p.token.Location)
			} else {
				node.ConditionalLocations = append(node.ConditionalLocations, p.token.Location)
			}
			p.consume()
		}
		if p.token.Kind == token.Newline {
			p.consumeNewlines()
			if p.token.Kind == token.BraceOpen {
				node.Documentation = p.ParseDocComment()
			}
		}
		return node
	}
	// Full Property
	for p.token.Kind == token.Hidden {
		node.HiddenLocations = append(node.HiddenLocations, p.token.Location)
		p.consume()
	}
	p.consumeNewlines()
	if p.token.Kind == token.BraceOpen {
		p.advance()
		node.Documentation = p.ParseDocComment()
	}
	p.consumeNewlines()
	var returnType *ast.TypeLiteral
	if p.token.Kind != token.Function {
		returnType = p.ParseTypeLiteral(errorExpectedFullPropertyStatement)
		if p.token.Kind != token.Function {
			p.fail(errorExpectedFullPropertyKeywordType, p.token.Location)
		}
	}
	first := p.ParseFunction(returnType)
	p.consumeNewlines()
	var second *ast.Function
	if p.token.Kind == token.EOF {
		p.fail(errorUnclosedFullProperty, source.Span(start, p.token.Location))
	}
	if p.token.Kind != token.EndProperty {
		var returnType *ast.TypeLiteral
		if p.token.Kind != token.Function {
			returnType = p.ParseTypeLiteral(errorExpectedFullPropertyStatement)
			if p.token.Kind != token.Function {
				p.fail(errorExpectedFullPropertyKeywordType, p.token.Location)
			}
		}
		second = p.ParseFunction(returnType)
		p.consumeNewlines()
	}
	if p.token.Kind == token.EOF {
		p.fail(errorUnclosedFullProperty, source.Span(start, p.token.Location))
	}
	switch {
	case strings.EqualFold(first.Name.Text, "get"):
		if first.ReturnType == nil {
			p.fail(errorExpectedFullPropertyGetReturnType, first.Location())
		}
		if len(first.ParameterList) != 0 {
			p.fail(errorExpectedFullPropertyGetParams, first.Location())
		}
		node.Get = first
	case strings.EqualFold(first.Name.Text, "set"):
		if first.ReturnType != nil {
			p.fail(errorExpectedFullPropertySetReturnType, first.Location())
		}
		if len(first.ParameterList) != 1 {
			p.fail(errorExpectedFullPropertySetParams, first.Location())
		}
		node.Set = first
	default:
		p.fail(errorExpectedFullPropertyGetOrSet, first.Name.Location())
	}
	switch {
	case second == nil:
		// Do nothing, property is either read only or write only.
	case strings.EqualFold(second.Name.Text, "get"):
		if node.Get != nil {
			panic(
				issue.New(
					errorExpectedFullPropertyGetDuplicate,
					p.file,
					second.Location(),
				).AppendRelated(
					p.file,
					node.Get.Location(),
					"'Get' already defined.",
				),
			)
		}
		if second.ReturnType == nil {
			p.fail(errorExpectedFullPropertyGetReturnType, second.Location())
		}
		if len(second.ParameterList) != 0 {
			p.fail(errorExpectedFullPropertyGetParams, second.Location())
		}
		node.Get = second
	case strings.EqualFold(second.Name.Text, "set"):
		if node.Set != nil {
			panic(
				issue.New(
					errorExpectedFullPropertySetDuplicate,
					p.file,
					second.Location(),
				).AppendRelated(
					p.file,
					node.Set.Location(),
					"'Set' already defined.",
				),
			)
		}
		if second.ReturnType != nil {
			p.fail(errorExpectedFullPropertySetReturnType, second.Location())
		}
		if len(first.ParameterList) != 1 {
			p.fail(errorExpectedFullPropertySetParams, second.Location())
		}
		node.Set = second
	default:
		p.fail(errorExpectedFullPropertyGetOrSet, second.Name.Location())
	}
	node.EndKeywordLocation = p.token.Location
	p.tryConsume(errorFullPropertyExtra, token.EndProperty)
	p.tryConsume(errorFullPropertyEnd, token.Newline, token.EOF)
	return node
}

func (p *parser) ParseScriptVariable(typeLiteral *ast.TypeLiteral) *ast.Variable {
	node := &ast.Variable{
		HasLeadingBlankLine: p.hasLeadingBlankLine(),
		Type:                typeLiteral,
	}
	node.Name = p.ParseIdentifier(errorExpectedScriptVariableIdent)
	if p.token.Kind == token.Assign {
		node.OperatorLocation = p.token.Location
		p.consume()
		node.Value = p.ParseLiteral(errorExpectedScriptVariableLiteral)
	}
	for p.token.Kind == token.Conditional {
		node.ConditionalLocations = append(node.ConditionalLocations, p.token.Location)
		p.consume()
	}
	p.tryConsume(errorExpectedScriptVariableEnd, token.Newline, token.EOF)
	return node
}

func (p *parser) ParseIdentifier(def *issue.Definition) *ast.Identifier {
	node := &ast.Identifier{
		Text:         string(p.token.Text),
		NodeLocation: p.token.Location,
	}
	p.tryConsume(def, token.Identifier, token.Self, token.Parent, token.Length)
	return node
}

func (p *parser) ParseTypeLiteral(def *issue.Definition) *ast.TypeLiteral {
	node := &ast.TypeLiteral{
		Name: &ast.Identifier{
			Text:         string(p.token.Text),
			NodeLocation: p.token.Location,
		},
	}
	p.tryConsume(def, token.Identifier, token.Bool, token.Int, token.Float, token.String)
	if p.token.Kind == token.ArrayType {
		node.BracketLocation = p.token.Location
		node.IsArray = true
		p.consume()
	}
	return node
}

func (p *parser) ParseExpression(def *issue.Definition, precedence int) ast.Expression {
	prefix := p.prefix[p.token.Kind]
	if prefix == nil {
		want := keys(p.prefix)
		p.unexpectedToken(def, p.token, want[0], want[1:]...)
		return nil // Unreachable, unexpectedToken panics.
	}
	expr := prefix()
	for p.token.Kind != token.Newline && p.token.Kind != token.EOF && precedence < precedenceOf(p.token.Kind) {
		infix := p.infix[p.token.Kind]
		if infix == nil {
			return expr
		}
		expr = infix(expr)
	}
	return expr
}

func (p *parser) ParseExpressionIdentifier() *ast.Identifier {
	node := &ast.Identifier{
		Text:         string(p.token.Text),
		NodeLocation: p.token.Location,
	}
	p.advanceExpected(token.Identifier, token.Self, token.Parent, token.Length)
	return node
}

func (p *parser) ParseExpressionBinary(left ast.Expression) *ast.Binary {
	precedence := precedenceOf(p.token.Kind)
	node := &ast.Binary{
		Kind:             ast.BinaryKind(p.token.Kind),
		LeftOperand:      left,
		OperatorLocation: p.token.Location,
	}
	p.consumeExpected(
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
	node.RightOperand = p.ParseExpression(errorExpectedBinaryExpr, precedence)
	return node
}

func (p *parser) ParseExpressionUnary() ast.Expression {
	if p.token.Kind == token.Minus {
		if p.lookahead.Kind == token.IntLiteral {
			return p.ParseExpressionIntLiteral()
		}
		if p.lookahead.Kind == token.FloatLiteral {
			return p.ParseExpressionFloatLiteral()
		}
	}
	node := &ast.Unary{
		Kind:             ast.UnaryKind(p.token.Kind),
		OperatorLocation: p.token.Location,
	}
	p.consumeExpected(token.Minus, token.LogicalNot)
	node.Operand = p.ParseExpression(errorExpectedUnaryExpr, prefix)
	return node
}

func (p *parser) ParseExpressionCast(value ast.Expression) *ast.Cast {
	node := &ast.Cast{
		Value:      value,
		AsLocation: p.token.Location,
	}
	p.consumeExpected(token.As)
	node.Type = p.ParseTypeLiteral(errorExpectedCastTypeLiteral)
	return node
}

func (p *parser) ParseExpressionAccess(value ast.Expression) *ast.Access {
	node := &ast.Access{
		Value:       value,
		DotLocation: p.token.Location,
	}
	p.consumeExpected(token.Dot)
	node.Name = p.ParseIdentifier(errorExpectedAccessIdent)
	return node
}

func (p *parser) ParseExpressionIndex(array ast.Expression) *ast.Index {
	node := &ast.Index{
		Value:        array,
		OpenLocation: p.token.Location,
	}
	p.consumeExpected(token.BracketOpen)
	node.Index = p.ParseExpression(errorExpectedIndexExpr, lowest)
	node.CloseLocation = p.token.Location
	p.tryConsume(errorExpectedIndexCloseBracket, token.BracketClose)
	return node
}

func (p *parser) ParseExpressionCall(function ast.Expression) *ast.Call {
	node := &ast.Call{
		Function:     function,
		OpenLocation: p.token.Location,
	}
	p.consumeExpected(token.ParenthesisOpen)
	node.Arguments = p.ParseArgumentList()
	node.CloseLocation = p.token.Location
	p.consumeExpected(token.ParenthesisClose)
	return node
}

func (p *parser) ParseArgumentList() []*ast.Argument {
	var args []*ast.Argument
	for {
		switch p.token.Kind {
		case token.EOF:
			p.fail(errorUnclosedArgListEOF, p.token.Location)
		case token.Newline:
			p.fail(errorUnclosedArgListNewline, p.token.Location)
		case token.Comma:
			p.consume()
		case token.ParenthesisClose:
			return args
		default:
			args = append(args, p.ParseArgument())
		}
	}
}

func (p *parser) ParseArgument() *ast.Argument {
	node := &ast.Argument{}
	if p.token.Kind == token.Identifier && p.lookahead.Kind == token.Assign {
		node.Name = p.ParseIdentifier(intenalInvalidState)
		node.OperatorLocation = p.token.Location
		p.consumeExpected(token.Assign)
	}
	node.Value = p.ParseExpression(errorExpectedArgExpr, lowest)
	return node
}

func (p *parser) ParseExpressionArrayCreation() *ast.ArrayCreation {
	node := &ast.ArrayCreation{
		NewLocation: p.token.Location,
	}
	p.advanceExpected(token.New)
	node.Type = p.ParseTypeLiteral(errorExpectedArrayCreationTypeLiteral)
	node.OpenLocation = p.token.Location
	p.tryConsume(errorExpectedArrayCreationOpenBracket, token.BracketOpen)
	if p.token.Kind != token.IntLiteral {
		p.failWithDetail(errorExpectedArrayCreationInt, p.token.Location, "Encountered a %s token.", p.token.Kind)
	}
	node.Size = p.ParseExpressionIntLiteral()
	node.CloseLocation = p.token.Location
	p.tryConsume(errorExpectedArrayCreationCloseBracket, token.BracketClose)
	return node
}

func (p *parser) ParseExpressionParenthetical() *ast.Parenthetical {
	node := &ast.Parenthetical{
		OpenLocation: p.token.Location,
	}
	p.advanceExpected(token.ParenthesisOpen)
	p.ParseExpression(errorExpectedParenExpr, lowest)
	node.CloseLocation = p.token.Location
	p.tryConsume(errorExpectedParenClose, token.ParenthesisClose)
	return node
}

func (p *parser) ParseLiteral(def *issue.Definition) ast.Literal {
	switch p.token.Kind {
	case token.True, token.False:
		p.ParseExpressionBoolLiteral()
	case token.IntLiteral:
		p.ParseExpressionIntLiteral()
	case token.FloatLiteral:
		p.ParseExpressionFloatLiteral()
	case token.StringLiteral:
		p.ParseExpressionStringLiteral()
	case token.None:
		p.ParseExpressionNoneLiteral()
	}
	p.unexpectedToken(def,
		p.token,
		token.True,
		token.False,
		token.IntLiteral,
		token.FloatLiteral,
		token.StringLiteral,
		token.None)
	return nil // Unreachable, unexpectedToken panics.
}

func (p *parser) ParseExpressionIntLiteral() *ast.IntLiteral {
	node := &ast.IntLiteral{
		RawText:      p.file.Bytes(p.token.Location),
		NodeLocation: p.token.Location,
	}
	p.consumeExpected(token.IntLiteral)
	return node
}

func (p *parser) ParseExpressionFloatLiteral() *ast.FloatLiteral {
	node := &ast.FloatLiteral{
		RawText:      p.file.Bytes(p.token.Location),
		NodeLocation: p.token.Location,
	}
	p.consumeExpected(token.FloatLiteral)
	return node
}

func (p *parser) ParseExpressionBoolLiteral() *ast.BoolLiteral {
	node := &ast.BoolLiteral{
		RawText:      p.file.Bytes(p.token.Location),
		NodeLocation: p.token.Location,
	}
	p.consumeExpected(token.StringLiteral, token.False)
	return node
}

func (p *parser) ParseExpressionStringLiteral() *ast.StringLiteral {
	node := &ast.StringLiteral{
		RawText:      p.file.Bytes(p.token.Location),
		NodeLocation: p.token.Location,
	}
	p.consumeExpected(token.StringLiteral)
	return node
}

func (p *parser) ParseExpressionNoneLiteral() *ast.NoneLiteral {
	node := &ast.NoneLiteral{
		RawText:      p.file.Bytes(p.token.Location),
		NodeLocation: p.token.Location,
	}
	p.consumeExpected(token.None)
	return node
}

func (p *parser) fail(def *issue.Definition, loc source.Location) {
	panic(issue.New(def, p.file, loc))
}

func (p *parser) failWithDetail(def *issue.Definition, loc source.Location, msg string, args ...any) {
	panic(issue.New(def, p.file, loc).WithDetail(msg, args...))
}

func registerPrefix[T ast.Expression](p *parser, fn func() T, kinds ...token.Kind) {
	for _, t := range kinds {
		p.prefix[t] = func() ast.Expression { return fn() }
	}
}

func registerInfix[T ast.Expression](p *parser, fn func(ast.Expression) T, kinds ...token.Kind) {
	for _, t := range kinds {
		p.infix[t] = func(expr ast.Expression) ast.Expression { return fn(expr) }
	}
}

func keys[K comparable, V any](data map[K]V) []K {
	keys := make([]K, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	return keys
}
