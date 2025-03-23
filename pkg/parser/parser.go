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
		prefix:            make(map[token.Type]prefixParser),
		infix:             make(map[token.Type]infixParser),
	}

	registerPrefix(prsr, prsr.ParseBoolLiteral, token.True, token.False)
	registerPrefix(prsr, prsr.ParseFloatLiteral, token.FloatLiteral)
	registerPrefix(prsr, prsr.ParseIdentifier, token.Identifier)
	registerPrefix(prsr, prsr.ParseIntLiteral, token.IntLiteral)
	registerPrefix(prsr, prsr.ParseNoneLiteral, token.None)
	registerPrefix(prsr, prsr.ParseParenthetical, token.LParen)
	registerPrefix(prsr, prsr.ParseStringLiteral, token.StringLiteral)
	registerPrefix(prsr, prsr.ParseUnary, token.Subtract, token.LogicalNot)

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
		token.Add,
		token.Subtract,
		token.Divide,
		token.Multiply,
		token.Modulo)
	registerInfix(prsr, prsr.ParseCall, token.LParen)
	registerInfix(prsr, prsr.ParseCast, token.As)
	registerInfix(prsr, prsr.ParseIndex, token.LBracket)

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

	prefix map[token.Type]prefixParser
	infix  map[token.Type]infixParser
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

var precedences = map[token.Type]int{
	token.LogicalOr:      LogicalOr,
	token.LogicalAnd:     LogicalAnd,
	token.Equal:          Comparison,
	token.NotEqual:       Comparison,
	token.Greater:        Comparison,
	token.GreaterOrEqual: Comparison,
	token.Less:           Comparison,
	token.LessOrEqual:    Comparison,
	token.Add:            Additive,
	token.Subtract:       Additive,
	token.Multiply:       Multiplicitive,
	token.Divide:         Multiplicitive,
	token.Modulo:         Multiplicitive,
	token.As:             Cast,
	token.Dot:            Access,
	token.LParen:         Call,
	token.LBracket:       Index,
}

func precedenceOf(t token.Type) int {
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
		return newError(p.token.Location, "expected any of [%s, %s], but found %s", t, tokensTypesToString(alts...), p.token.Type)
	}
	return newError(p.token.Location, "expected %s, but found %s", t, p.token.Type)
}

func tokensTypesToString(types ...token.Type) string {
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
	for p.token.Type == token.Newline {
		if err := p.next(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) ParseScript() (*ast.Script, error) {
	script := &ast.Script{
		Location: source.Range{
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
			Text:     string(p.token.Location.Text()),
			Location: p.token.Location,
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
		Name:     ident,
		Location: source.Span(start, ident.Location),
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
				Message:  fmt.Sprintf("hit end of file while parsing state %q, did you forget %s?", name.Location.Text(), token.EndState),
				Location: source.Span(start, p.token.Location),
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
	node.Location = source.Span(start, p.token.Location)
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
		node.Location = source.Span(start, p.token.Location)
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
	node.Location = source.Span(start, p.token.Location)
	if err := p.next(); err != nil {
		return nil, err
	}
	return node, p.tryConsume(token.Newline, token.EOF)
}

func (p *parser) ParseFunction(returnType *ast.TypeLiteral) (*ast.Function, error) {
	start := p.token.Location
	if returnType != nil {
		start = returnType.Location
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
			Text:     string(p.token.Location.Text()[1 : p.token.Location.Length-1]),
			Location: p.token.Location,
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
		node.Location = source.Span(start, end)
		return node, nil
	}
	stmts, err := p.ParseFunctionStatementBlock(token.EndFunction)
	if err != nil {
		return nil, err
	}
	node.Statements = stmts
	node.Location = source.Span(start, p.token.Location)
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
	node.Location = source.Span(start, name.Location)
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
		node.Location = source.Span(start, literal.Range())
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
	switch p.token.Type {
	case token.Return:
		return p.ParseReturn()
	case token.If:
		return p.ParseIf()
	case token.While:
		return p.ParseWhile()
	case token.Bool, token.Int, token.Float, token.String:
		return p.ParseFunctionVariable()
	case token.Identifier:
		if p.lookahead.Type == token.Identifier { // Object type.
			return p.ParseFunctionVariable()
		}
	}
	return p.ParseAssignment()
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
	expr, err := p.ParseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	return &ast.FunctionVariable{
		Type:     typeLiteral,
		Name:     name,
		Value:    expr,
		Location: source.Span(start, expr.Range()),
	}, nil
}

func (p *parser) ParseAssignment() (*ast.Assignment, error) {
	start := p.token.Location
	assignee, err := p.ParseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	operator, err := p.ParseAssignmentOperator()
	if err != nil {
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
		Location: source.Span(start, expr.Range()),
	}, nil
}

func (p *parser) ParseAssignmentOperator() (*ast.AssignmentOperator, error) {
	operator := &ast.AssignmentOperator{
		Location: p.token.Location,
	}
	switch p.token.Type {
	case token.Assign:
		operator.Kind = ast.Assign
	case token.AssignAdd:
		operator.Kind = ast.AssignAdd
	case token.AssignDivide:
		operator.Kind = ast.AssignDivide
	case token.AssignModulo:
		operator.Kind = ast.AssignModulo
	case token.AssignMultiply:
		operator.Kind = ast.AssignMultiply
	case token.AssignSubtract:
		operator.Kind = ast.AssignSubtract
	default:
		types := tokensTypesToString(
			token.Assign,
			token.AssignAdd,
			token.AssignDivide,
			token.AssignModulo,
			token.AssignMultiply,
			token.AssignSubtract)
		return nil, newError(p.token.Location, "expected any of [%s], but found %s", types, p.token.Type)
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	return operator, nil
}

func (p *parser) ParseReturn() (*ast.Return, error) {
	start := p.token.Location
	if err := p.tryConsume(token.Return); err != nil {
		return nil, err
	}
	if p.token.Type == token.Newline {
		return &ast.Return{
			Location: start,
		}, nil
	}
	expr, err := p.ParseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	return &ast.Return{
		Value:    expr,
		Location: source.Span(start, expr.Range()),
	}, nil
}

func (p *parser) ParseIf() (*ast.If, error) {
	start := p.token.Location
	if err := p.tryConsume(token.If); err != nil {
		return nil, err
	}
	expr, err := p.ParseExpression(Lowest)
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
		expr, err := p.ParseExpression(Lowest)
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
	node.Location = source.Span(start, p.token.Location)
	if err := p.tryConsume(token.EndIf); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseWhile() (*ast.While, error) {
	start := p.token.Location
	if err := p.tryConsume(token.While); err != nil {
		return nil, err
	}
	expr, err := p.ParseExpression(Lowest)
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
		Condition:  expr,
		Statements: stmts,
		Location:   source.Span(start, p.token.Location),
	}
	if err := p.tryConsume(token.EndWhile); err != nil {
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
		Text:     string(bytes.ToLower(rng.Text())),
		Location: rng,
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
			Type:     scalar,
			Location: start,
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
		Location: source.Span(start, end),
	}, nil
}

func (p *parser) ParseExpression(precedence int) (ast.Expression, error) {
	prefix := p.prefix[p.token.Type]
	if prefix == nil {
		return nil, newError(p.token.Location, "expected any of [%s], but found %s", tokensTypesToString(keys(p.prefix)...), p.token.Type)
	}
	expr, err := prefix()
	if err != nil {
		return nil, err
	}
	if p.lookahead.Type != token.Newline && p.lookahead.Type != token.EOF && precedence < precedenceOf(p.lookahead.Type) {
		infix := p.infix[p.token.Type]
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
	precedence := precedenceOf(p.token.Type)
	operator, err := p.ParseBinaryOperator()
	if err != nil {
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
		Location:     source.Span(left.Range(), right.Range()),
	}, nil
}

func (p *parser) ParseBinaryOperator() (*ast.BinaryOperator, error) {
	operator := &ast.BinaryOperator{Location: p.token.Location}
	switch p.token.Type {
	case token.LogicalOr:
		operator.Kind = ast.LogicalOr
	case token.LogicalAnd:
		operator.Kind = ast.LogicalAnd
	case token.Equal:
		operator.Kind = ast.Equal
	case token.NotEqual:
		operator.Kind = ast.NotEqual
	case token.Greater:
		operator.Kind = ast.Greater
	case token.GreaterOrEqual:
		operator.Kind = ast.GreaterOrEqual
	case token.Less:
		operator.Kind = ast.Less
	case token.LessOrEqual:
		operator.Kind = ast.LessOrEqual
	case token.Add:
		operator.Kind = ast.Add
	case token.Subtract:
		operator.Kind = ast.Subtract
	case token.Multiply:
		operator.Kind = ast.Multiply
	case token.Divide:
		operator.Kind = ast.Divide
	case token.Modulo:
		operator.Kind = ast.Modulo
	default:
		types := tokensTypesToString(
			token.LogicalOr,
			token.LogicalAnd,
			token.Equal,
			token.NotEqual,
			token.Greater,
			token.GreaterOrEqual,
			token.Less,
			token.LessOrEqual,
			token.Add,
			token.Subtract,
			token.Divide,
			token.Multiply,
			token.Modulo,
		)
		return nil, newError(p.token.Location, "expected any of [%s], but found %s", types, p.token.Type)
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	return operator, nil
}

func (p *parser) ParseUnary() (ast.Expression, error) {
	if p.token.Type == token.Subtract &&
		(p.lookahead.Type == token.IntLiteral || p.lookahead.Type == token.FloatLiteral) {
		return p.ParseLiteral()
	}
	operator, err := p.ParseUnaryOperator()
	if err != nil {
		return nil, err
	}
	expr, err := p.ParseExpression(Prefix)
	if err != nil {
		return nil, err
	}
	return &ast.Unary{
		Operator: operator,
		Operand:  expr,
		Location: source.Span(operator.Location, expr.Range()),
	}, nil
}

func (p *parser) ParseUnaryOperator() (*ast.UnaryOperator, error) {
	operator := &ast.UnaryOperator{
		Location: p.token.Location,
	}
	switch p.token.Type {
	case token.Subtract:
		operator.Kind = ast.Negate
	case token.LogicalNot:
		operator.Kind = ast.LogicalNot
	default:
		types := tokensTypesToString(token.Subtract, token.LogicalNot)
		return nil, newError(p.token.Location, "expected any of [%s], but found %s", types, p.token.Type)
	}
	if err := p.next(); err != nil {
		return nil, err
	}
	return operator, nil
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
		Operator: &ast.AsOperator{Location: operator.Location},
		Type:     typeLiteral,
		Location: source.Span(value.Range(), typeLiteral.Location),
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
		Operator: &ast.AccessOperator{Location: operator.Location},
		Name:     name,
		Location: source.Span(value.Range(), name.Location),
	}, nil
}

func (p *parser) ParseAccessOperator() (*ast.AccessOperator, error) {
	operator := &ast.AccessOperator{
		Location: p.token.Location,
	}
	if err := p.tryConsume(token.Dot); err != nil {
		return nil, err
	}
	return operator, nil
}

func (p *parser) ParseIndex(array ast.Expression) (*ast.Index, error) {
	open, err := p.ParseArrayOpenOperator()
	if err != nil {
		return nil, err
	}
	index, err := p.ParseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	close, err := p.ParseArrayCloseOperator()
	if err != nil {
		return nil, err
	}
	return &ast.Index{
		Value:         array,
		OpenOperator:  open,
		Index:         index,
		CloseOperator: close,
		Location:      source.Span(array.Range(), close.Location),
	}, nil
}

func (p *parser) ParseCall(reciever ast.Expression) (*ast.Call, error) {
	if err := p.tryConsume(token.LParen); err != nil {
		return nil, err
	}
	args, err := p.ParseArgumentList()
	if err != nil {
		return nil, err
	}
	end := p.token.Location
	if err := p.tryConsume(token.RParen); err != nil {
		return nil, err
	}
	return &ast.Call{
		Reciever:  reciever,
		Arguments: args,
		Location:  source.Span(reciever.Range(), end),
	}, nil
}

func (p *parser) ParseArgumentList() ([]*ast.Argument, error) {
	var args []*ast.Argument
	for {
		switch p.token.Type {
		case token.Comma:
			if err := p.next(); err != nil {
				return nil, err
			}
		case token.RParen:
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
	if p.token.Type == token.Identifier {
		id, err := p.ParseIdentifier()
		if err != nil {
			return nil, err
		}
		assign := p.token
		if err := p.tryConsume(token.Assign); err != nil {
			return nil, err
		}
		node.Name = id
		node.Operator = &ast.AssignmentOperator{
			Kind:     ast.Assign,
			Location: assign.Location,
		}
	}
	value, err := p.ParseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	node.Value = value
	if node.Name != nil {
		node.Location = source.Span(node.Name.Location, value.Range())
	} else {
		node.Location = value.Range()
	}
	return node, nil
}

func (p *parser) ParseArrayCreation() (*ast.ArrayCreation, error) {
	new, err := p.ParseNewOperator()
	if err != nil {
		return nil, err
	}
	typeLiteral, err := p.ParseTypeLiteral()
	if err != nil {
		return nil, err
	}
	open, err := p.ParseArrayOpenOperator()
	if err != nil {
		return nil, err
	}
	size, err := p.ParseIntLiteral()
	if err != nil {
		return nil, err
	}
	if size.Value < 1 || size.Value > 128 {
		return nil, newError(size.Range(), "expected array size to be an IntLiteral in range [1, 128], but found %d", size.Value)
	}
	close, err := p.ParseArrayCloseOperator()
	if err != nil {
		return nil, err
	}
	return &ast.ArrayCreation{
		NewOperator:   new,
		Type:          typeLiteral,
		OpenOperator:  open,
		Size:          size,
		CloseOperator: close,
		Location:      source.Span(new.Location, close.Location),
	}, nil
}

func (p *parser) ParseNewOperator() (*ast.NewOperator, error) {
	operator := &ast.NewOperator{
		Location: p.token.Location,
	}
	if err := p.tryConsume(token.New); err != nil {
		return nil, err
	}
	return operator, nil
}

func (p *parser) ParseArrayOpenOperator() (*ast.ArrayOpenOperator, error) {
	operator := &ast.ArrayOpenOperator{
		Location: p.token.Location,
	}
	if err := p.tryConsume(token.LBracket); err != nil {
		return nil, err
	}
	return operator, nil
}

func (p *parser) ParseArrayCloseOperator() (*ast.ArrayCloseOperator, error) {
	operator := &ast.ArrayCloseOperator{
		Location: p.token.Location,
	}
	if err := p.tryConsume(token.RBracket); err != nil {
		return nil, err
	}
	return operator, nil
}

func (p *parser) ParseParenthetical() (*ast.Parenthetical, error) {
	start := p.token.Location
	if err := p.tryConsume(token.LParen); err != nil {
		return nil, err
	}
	expr, err := p.ParseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	node := &ast.Parenthetical{
		Value:    expr,
		Location: source.Span(start, p.token.Location),
	}
	if err := p.tryConsume(token.RParen); err != nil {
		return nil, err
	}
	return node, nil
}

func (p *parser) ParseLiteral() (ast.Literal, error) {
	switch p.token.Type {
	case token.Subtract:
		// While this overlaps with a unary expression, we lump these together
		// because there are some contexts where a literal is required which can
		// include a sign, but where a Unary is not allowed.
		sign := p.token
		if err := p.next(); err != nil {
			return nil, err
		}
		switch p.token.Type {
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
			return nil, fmt.Errorf("expected IntLiteral or FloatLiteral, but found %s", p.token.Type)
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
	return nil, fmt.Errorf("expected True, False, None, Integer, Float, or String literal, but found %s", p.token.Type)
}

func (p *parser) ParseIntLiteral() (*ast.IntLiteral, error) {
	tok := p.token
	if err := p.tryConsume(token.IntLiteral); err != nil {
		return nil, err
	}
	text := strings.ToLower(string(tok.Location.Text()))
	val, err := strconv.ParseInt(text, 0, 32)
	if err != nil {
		return nil, newError(tok.Location, "failed to parse %q as an integer: %v", text, err)
	}
	return &ast.IntLiteral{
		Value:    int(val),
		Location: tok.Location,
	}, nil
}

func (p *parser) ParseFloatLiteral() (*ast.FloatLiteral, error) {
	tok := p.token
	if err := p.tryConsume(token.FloatLiteral); err != nil {
		return nil, err
	}
	text := strings.ToLower(string(tok.Location.Text()))
	val, err := strconv.ParseFloat(text, 32)
	if err != nil {
		return nil, newError(tok.Location, "failed to parse %q as a float: %v", text, err)
	}
	return &ast.FloatLiteral{
		Value:    float32(val),
		Location: tok.Location,
	}, nil
}

func (p *parser) ParseBoolLiteral() (*ast.BoolLiteral, error) {
	tok := p.token
	if err := p.tryConsume(token.True, token.False); err != nil {
		return nil, err
	}
	return &ast.BoolLiteral{
		Value:    tok.Type == token.True,
		Location: tok.Location,
	}, nil
}

func (p *parser) ParseStringLiteral() (*ast.StringLiteral, error) {
	tok := p.token
	if err := p.tryConsume(token.StringLiteral); err != nil {
		return nil, err
	}
	return &ast.StringLiteral{
		Value:    string(tok.Location.Text()[1 : tok.Location.Length-1]),
		Location: tok.Location,
	}, nil
}

func (p *parser) ParseNoneLiteral() (*ast.NoneLiteral, error) {
	tok := p.token
	if err := p.tryConsume(token.None); err != nil {
		return nil, err
	}
	return &ast.NoneLiteral{
		Location: tok.Location,
	}, nil
}

func registerPrefix[T ast.Expression](p *parser, fn func() (T, error), types ...token.Type) {
	for _, t := range types {
		p.prefix[t] = func() (ast.Expression, error) { return fn() }
	}
}

func registerInfix[T ast.Expression](p *parser, fn func(ast.Expression) (T, error), types ...token.Type) {
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
