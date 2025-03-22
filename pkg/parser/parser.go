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
	}
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
	looseComments     []token.Token

	recovery bool
	errors   []ast.Error
}

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
	if p.token.Type == token.LineComment || p.token.Type == token.BlockComment {
		if p.keepLooseComments {
			p.looseComments = append(p.looseComments, p.token)
		}
		return p.next()
	}
	return nil
}

// tryConsume advances the token position if the current token matches the given
// token type or returns an error.
func (p *parser) tryConsume(t token.Type, alts ...token.Type) error {
	if p.token.Type == t {
		return p.next()
	}
	for _, t := range alts {
		if p.token.Type == t {
			return p.next()
		}
	}
	if len(alts) > 0 {
		strs := make([]string, len(alts))
		for i, alt := range alts {
			strs[i] = alt.String()
		}
		return newError(p.token.Location, "expected any of [%s, %s], but found %s", t, strings.Join(strs, ", "), p.token.Type)
	}
	return newError(p.token.Location, "expected %s, but found %s", t, p.token.Type)
}

// consumeNewlines advances the token position through the as many newlines as
// possible until a non-newline token is found.
func (p *parser) consumeNewlines() error {
	for p.token.Type == token.Newline {
		if err := p.next(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) ParseScript() (*ast.Script, error) {
	script := &ast.Script{
		SourceRange: source.Range{
			File:        p.token.Location.File,
			Length:      len(p.token.Location.File.Text),
			StartLine:   1,
			StartColumn: 1,
		},
	}
	if err := p.ParseScriptHeader(script); err != nil {
		return nil, err
	}
	if p.token.Type == token.DocComment {
		script.Comment = &ast.DocComment{
			Text:        string(p.token.Location.Text()),
			SourceRange: p.token.Location,
		}
		if err := p.next(); err != nil {
			return nil, err
		}
	}
	for p.token.Type != token.EOF {
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		stmt, err := p.ParseScriptStatement()
		if err != nil {
			return nil, err
		}
		if stmt != nil {
			script.Statements = append(script.Statements, stmt)
		}
	}
	return script, nil
}

func (p *parser) ParseScriptHeader(script *ast.Script) error {
	if err := p.tryConsume(token.ScriptName); err != nil {
		return err
	}
	var err error
	if script.Name, err = p.ParseIdentifier(); err != nil {
		return err
	}
	if p.token.Type == token.Extends {
		if err := p.next(); err != nil {
			return err
		}
		if script.Extends, err = p.ParseIdentifier(); err != nil {
			return err
		}
	}
	for p.token.Type == token.Hidden || p.token.Type == token.Conditional {
		if p.token.Type == token.Hidden {
			script.IsHidden = true
		} else {
			script.IsConditional = true
		}
		if err := p.next(); err != nil {
			return err
		}
	}
	return p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseScriptStatement() (ast.ScriptStatement, error) {
	start := p.token
	var stmt ast.ScriptStatement
	var err error
	switch p.token.Type {
	case token.Import:
		stmt, err = p.ParseImport()
	case token.Event:
		stmt, err = p.ParseEvent()
	case token.Auto, token.State:
		stmt, err = p.ParseState()
	case token.Function:
		stmt, err = p.ParseFunction(nil)
	case token.Bool, token.Float, token.Int, token.String, token.Identifier:
		var typeLiteral *ast.TypeLiteral
		typeLiteral, err = p.ParseTypeLiteral()
		if err != nil {
			return nil, err
		}
		switch p.token.Type {
		case token.Identifier:
			stmt, err = p.ParseScriptVariable(typeLiteral)
		case token.Property:
			stmt, err = p.ParseProperty(typeLiteral)
		case token.Function:
			stmt, err = p.ParseFunction(typeLiteral)
		}
	default:
		err = fmt.Errorf("expected Import, Event, State, Function, Property, or a variable definition, but found %s", start.Type)
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
		Message:     fmt.Sprintf("%v", err),
		SourceRange: source.Span(start.Location, p.token.Location),
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
		switch p.lookahead.Type {
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
	start := p.token.Location
	if err := p.next(); err != nil {
		return nil, err
	}
	ident, err := p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	node := &ast.Import{
		Name:        ident,
		SourceRange: source.Span(start, ident.SourceRange),
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseState() (ast.ScriptStatement, error) {
	start := p.token.Location
	isAuto := p.token.Type == token.Auto
	if isAuto {
		if err := p.next(); err != nil {
			return nil, err
		}
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	name, err := p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	node := &ast.State{
		Name:   name,
		IsAuto: isAuto,
	}
	for p.token.Type != token.EndState {
		if p.token.Type == token.EOF {
			// State was never closed, proactively create a
			errStmt := &ast.ErrorScriptStatement{
				Message:     fmt.Sprintf("hit end of file while parsing state %q, did you forget %s?", name.SourceRange.Text(), token.EndState),
				SourceRange: source.Span(start, p.token.Location),
			}
			p.errors = append(p.errors, errStmt)
			return errStmt, nil
		}
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		if p.token.Type == token.EndState {
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
	node.SourceRange = source.Span(start, p.token.Location)
	if err := p.next(); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseInvokable() (ast.Invokable, error) {
	start := p.token
	var stmt ast.Invokable
	var err error
	switch p.token.Type {
	case token.Event:
		stmt, err = p.ParseEvent()
	case token.Function:
		stmt, err = p.ParseFunction(nil)
	case token.Bool, token.Float, token.Int, token.String, token.Identifier:
		var typeLiteral *ast.TypeLiteral
		typeLiteral, err = p.ParseTypeLiteral()
		if err != nil {
			return nil, err
		}
		switch p.token.Type {
		case token.Function:
			stmt, err = p.ParseFunction(typeLiteral)
		}
	default:
		err = fmt.Errorf("expected Event or Function, but found %s", start.Type)
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
		Message:     fmt.Sprintf("%v", err),
		SourceRange: source.Span(start.Location, p.token.Location),
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
		switch p.lookahead.Type {
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
	start := p.token.Location
	if err := p.next(); err != nil {
		return nil, err
	}
	name, err := p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	node := &ast.Event{
		Name: name,
	}
	params, err := p.ParseParameterList()
	if err != nil {
		return nil, err
	}
	node.Parameters = params
	if p.token.Type == token.Native {
		if err := p.next(); err != nil {
			return nil, err
		}
		node.IsNative = true
		node.SourceRange = source.Span(start, p.token.Location)
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		return node, err
	}
	stmts, err := p.ParseFunctionStatementBlock(token.EndEvent)
	if err != nil {
		return nil, err
	}
	node.Statements = stmts
	node.SourceRange = source.Span(start, p.token.Location)
	if err := p.next(); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseFunction(returnType *ast.TypeLiteral) (*ast.Function, error) {
	start := p.token.Location
	if returnType != nil {
		start = returnType.SourceRange
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	name, err := p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	node := &ast.Function{
		Name:       name,
		ReturnType: returnType,
	}
	params, err := p.ParseParameterList()
	if err != nil {
		return nil, err
	}
	node.Parameters = params
	var end source.Range
	for p.token.Type == token.Native || p.token.Type == token.Global {
		if p.token.Type == token.Native {
			node.IsNative = true
		} else {
			node.IsGlobal = true
		}
		end = p.token.Location
		if err := p.next(); err != nil {
			return nil, err
		}
	}
	if err := p.consumeNewlines(); err != nil {
		return nil, err
	}
	if p.token.Type == token.DocComment {
		node.Comment = &ast.DocComment{
			Text:        string(p.token.Location.Text()[1 : p.token.Location.Length-1]),
			SourceRange: p.token.Location,
		}
		end = p.token.Location
		if err := p.next(); err != nil {
			return nil, err
		}
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
	}
	if node.IsNative {
		node.SourceRange = source.Span(start, end)
		return node, nil
	}
	stmts, err := p.ParseFunctionStatementBlock(token.EndFunction)
	if err != nil {
		return nil, err
	}
	node.Statements = stmts
	node.SourceRange = source.Span(start, p.token.Location)
	if err := p.next(); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseParameterList() ([]*ast.Parameter, error) {
	if err := p.tryConsume(token.LParen); err != nil {
		return nil, err
	}
	var params []*ast.Parameter
	for {
		switch p.token.Type {
		case token.Comma:
			if err := p.next(); err != nil {
				return nil, err
			}
		case token.RParen:
			if err := p.next(); err != nil {
				return nil, err
			}
			return params, nil
		default:
			param, err := p.ParseParameter()
			if err != nil {
				return nil, err
			}
			params = append(params, param)
		}
	}
}

func (p *parser) ParseParameter() (*ast.Parameter, error) {
	start := p.token.Location
	typeLiteral, err := p.ParseTypeLiteral()
	if err != nil {
		return nil, err
	}
	name, err := p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	node := &ast.Parameter{
		Type: typeLiteral,
		Name: name,
	}
	node.SourceRange = source.Span(start, name.SourceRange)
	if p.token.Type == token.Assign {
		// Has default.
		if err := p.next(); err != nil {
			return nil, err
		}
		literal, err := p.ParseLiteral()
		if err != nil {
			return nil, err
		}
		node.Value = literal
		node.SourceRange = source.Span(start, literal.Range())
	}
	return node, nil
}

func (p *parser) ParseFunctionStatementBlock(terminals ...token.Type) ([]ast.FunctionStatement, error) {
	terms := make(map[token.Type]struct{})
	for _, t := range terminals {
		terms[t] = struct{}{}
	}
	var stmts []ast.FunctionStatement
	for {
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		if _, ok := terms[p.token.Type]; ok {
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
			Message:     fmt.Sprintf("%v", err),
			SourceRange: source.Span(start, p.token.Location),
		}
		p.errors = append(p.errors, errStmt)
		p.recovery = false
		stmts = append(stmts, errStmt)
	}
}

func (p *parser) recoverFunctionStatement() error {
	for {
		switch p.lookahead.Type {
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
	return nil, fmt.Errorf("ParseFunctionStatement unimplmented")
}

func (p *parser) ParseFunctionVariable() (*ast.FunctionVariable, error) {
	start := p.token.Location
	typeLiteral, err := p.ParseTypeLiteral()
	if err != nil {
		return nil, err
	}
	name, err := p.ParseIdentifier()
	if err != nil {
		return nil, err
	}
	if err := p.tryConsume(token.Assign); err != nil {
		return nil, err
	}
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	return &ast.FunctionVariable{
		Type:        typeLiteral,
		Name:        name,
		Value:       expr,
		SourceRange: source.Span(start, expr.Range()),
	}, nil
}

func (p *parser) ParseReturn() (*ast.Return, error) {
	start := p.token.Location
	if err := p.tryConsume(token.Return); err != nil {
		return nil, err
	}
	if p.token.Type == token.Newline {
		return &ast.Return{
			SourceRange: start,
		}, nil
	}
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	return &ast.Return{
		Value:       expr,
		SourceRange: source.Span(start, expr.Range()),
	}, nil
}

func (p *parser) ParseIf() (*ast.If, error) {
	start := p.token.Location
	if err := p.tryConsume(token.If); err != nil {
		return nil, err
	}
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	if err := p.consumeNewlines(); err != nil {
		return nil, err
	}
	stmts, err := p.ParseFunctionStatementBlock(token.EndIf, token.Else, token.ElseIf)
	if err != nil {
		return nil, err
	}
	node := &ast.If{
		Conditional: ast.ConditionalBlock{
			Condition:  expr,
			Statements: stmts,
		},
	}
	for {
		if p.token.Type != token.ElseIf {
			break
		}
		if err := p.next(); err != nil {
			return nil, err
		}
		expr, err := p.ParseExpression()
		if err != nil {
			return nil, err
		}
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		stmts, err := p.ParseFunctionStatementBlock(token.EndIf, token.Else, token.ElseIf)
		if err != nil {
			return nil, err
		}
		block := ast.ConditionalBlock{
			Condition:  expr,
			Statements: stmts,
		}
		node.AlternativeConditionals = append(node.AlternativeConditionals, block)
	}
	if p.token.Type == token.Else {
		if err := p.next(); err != nil {
			return nil, err
		}
		if err := p.consumeNewlines(); err != nil {
			return nil, err
		}
		stmts, err := p.ParseFunctionStatementBlock(token.EndIf)
		if err != nil {
			return nil, err
		}
		node.Alternative = stmts
	}
	node.SourceRange = source.Span(start, p.token.Location)
	if err := p.tryConsume(token.EndIf); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseWhile() (*ast.While, error) {
	start := p.token.Location
	if err := p.tryConsume(token.If); err != nil {
		return nil, err
	}
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	if err := p.consumeNewlines(); err != nil {
		return nil, err
	}
	stmts, err := p.ParseFunctionStatementBlock(token.EndWhile)
	if err != nil {
		return nil, err
	}
	node := &ast.While{
		Condition:   expr,
		Statements:  stmts,
		SourceRange: source.Span(start, p.token.Location),
	}
	if err := p.tryConsume(token.EndWhile); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseParenthetical() (*ast.Parenthetical, error) {
	start := p.token.Location
	if err := p.tryConsume(token.LParen); err != nil {
		return nil, err
	}
	expr, err := p.ParseExpression()
	if err != nil {
		return nil, err
	}
	node := &ast.Parenthetical{
		Value:       expr,
		SourceRange: source.Span(start, p.token.Location),
	}
	if err := p.tryConsume(token.RParen); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseProperty(propertyType *ast.TypeLiteral) (*ast.Property, error) {
	return nil, newError(p.token.Location, "ParseProperty unimplemented.")
}

func (p *parser) ParseScriptVariable(variableType *ast.TypeLiteral) (*ast.ScriptVariable, error) {
	return nil, newError(p.token.Location, "ParseScriptVariable unimplemented.")
}

func (p *parser) ParseIdentifier() (*ast.Identifier, error) {
	rng := p.token.Location
	if err := p.tryConsume(token.Identifier); err != nil {
		return nil, err
	}
	return &ast.Identifier{
		Text:        string(bytes.ToLower(rng.Text())),
		SourceRange: rng,
	}, nil
}

func (p *parser) ParseTypeLiteral() (*ast.TypeLiteral, error) {
	start := p.token.Location
	var scalar types.Scalar
	switch p.token.Type {
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
		return nil, newError(p.token.Location, "expected Bool, Int, Float, String, or an identifier, but found %s", p.token.Type)
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	if p.token.Type != token.LBracket {
		return &ast.TypeLiteral{
			Type:        scalar,
			SourceRange: start,
		}, nil
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	end := p.token.Location
	if err := p.tryConsume(token.RBracket); err != nil {
		return nil, err
	}
	return &ast.TypeLiteral{
		Type: types.Array{
			ElementType: scalar,
		},
		SourceRange: source.Span(start, end),
	}, nil
}

func (p *parser) ParseExpression() (ast.Expression, error) {
	return nil, newError(p.token.Location, "ParseExpression unimplemented.")
}

func (p *parser) ParseLiteral() (ast.Literal, error) {
	switch p.token.Type {
	case token.Subtract:
		sign := p.token
		if err := p.next(); err != nil {
			return nil, err
		}
		return p.ParseNumber(&sign)
	case token.True, token.False:
		return &ast.BoolLiteral{
			Value:       p.token.Type == token.True,
			SourceRange: p.token.Location,
		}, nil
	case token.IntLiteral, token.FloatLiteral:
		return p.ParseNumber(nil)
	case token.StringLiteral:
		return &ast.StringLiteral{
			Value:       string(p.token.Location.Text()[1 : p.token.Location.Length-1]),
			SourceRange: p.token.Location,
		}, nil
	case token.None:
		return &ast.NoneLiteral{
			SourceRange: p.token.Location,
		}, nil
	}
	return nil, fmt.Errorf("expected True, False, None, Integer, Float, or String literal, but found %s", p.token.Type)
}

func (p *parser) ParseNumber(sign *token.Token) (ast.Literal, error) {
	switch p.token.Type {
	case token.IntLiteral:
		tok := p.token
		if err := p.tryConsume(token.IntLiteral); err != nil {
			return nil, err
		}
		text := strings.ToLower(string(tok.Location.Text()))
		val, err := strconv.ParseInt(text, 0, 32)
		if err != nil {
			return nil, newError(tok.Location, "failed to parse %q as an integer: %v", text, err)
		}
		srcRange := tok.Location
		if sign != nil {
			srcRange = source.Span(sign.Location, tok.Location)
			val = -val
		}
		return &ast.IntLiteral{
			Value:       int(val),
			SourceRange: srcRange,
		}, nil
	case token.FloatLiteral:
		tok := p.token
		if err := p.tryConsume(token.FloatLiteral); err != nil {
			return nil, err
		}
		text := strings.ToLower(string(tok.Location.Text()))
		val, err := strconv.ParseFloat(text, 32)
		if err != nil {
			return nil, newError(tok.Location, "failed to parse %q as a float: %v", text, err)
		}
		srcRange := tok.Location
		if sign != nil {
			srcRange = source.Span(sign.Location, tok.Location)
			val = -val
		}
		return &ast.FloatLiteral{
			Value:       float32(val),
			SourceRange: srcRange,
		}, nil
	}
	return nil, fmt.Errorf("expected Int or Float literal, but found %s", p.token.Type)
}
