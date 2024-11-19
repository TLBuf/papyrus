// Package parser defines a Papyrus parser.
package parser

import (
	"bytes"
	"strings"

	"github.com/TLBuf/papyrus/pkg/ast"
	"github.com/TLBuf/papyrus/pkg/lexer"
	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/TLBuf/papyrus/pkg/token"
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
	prsr := &parser{
		l:                 lexer.New(file),
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
}

// next advances token and lookahead by one token while skipping loose comment
// tokens. Returns true if parsing should continue, false otherwise.
func (p *parser) next() error {
	p.token = p.lookahead
	t, err := p.l.NextToken()
	if err != nil {
		return newError(err.(lexer.Error).Location, err.(lexer.Error).Message)
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
		return newError(p.token.SourceRange, "expected any of [%s, %s], but found %s", t, strings.Join(strs, ", "), p.token.Type)
	}
	return newError(p.token.SourceRange, "expected %s, but found %s", t, p.token.Type)
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
			File:   p.token.SourceRange.File,
			Length: len(p.token.SourceRange.File.Text),
			Line:   1,
			Column: 1,
		},
	}
	if err := p.ParseScriptHeader(script); err != nil {
		return nil, err
	}
	if p.token.Type == token.DocComment {
		script.Comment = &ast.DocComment{
			Text:        string(p.token.SourceRange.Text()),
			SourceRange: p.token.SourceRange,
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
	switch p.token.Type {
	case token.Import:
		return p.ParseImport()
	case token.Event:
		return p.ParseEvent()
	case token.Auto, token.State:
		return p.ParseState()
	case token.Function:
		return p.ParseFunction(nil)
	case token.Bool, token.Float, token.Int, token.String, token.Identifier:
		typeLiteral, err := p.ParseTypeLiteral()
		if err != nil {
			return nil, err
		}
		switch p.token.Type {
		case token.Identifier:
			return p.ParseScriptVariable(typeLiteral)
		case token.Property:
			return p.ParseProperty(typeLiteral)
		case token.Function:
			return p.ParseFunction(typeLiteral)
		}
	}
	return nil, newError(p.token.SourceRange, "expected Import, Event, State, Function, Property, or Variable, but found %s", p.token.Type)
}

func (p *parser) ParseImport() (*ast.Import, error) {
	start := p.token.SourceRange
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

func (p *parser) ParseState() (*ast.State, error) {
	return nil, newError(p.token.SourceRange, "ParseState unimplemented.")
}

func (p *parser) ParseEvent() (*ast.Event, error) {
	return nil, newError(p.token.SourceRange, "ParseEvent unimplemented.")
}

func (p *parser) ParseFunction(returnType *ast.TypeLiteral) (*ast.Function, error) {
	return nil, newError(p.token.SourceRange, "ParseFunction unimplemented.")
}

func (p *parser) ParseProperty(propertyType *ast.TypeLiteral) (*ast.Property, error) {
	return nil, newError(p.token.SourceRange, "ParseProperty unimplemented.")
}

func (p *parser) ParseScriptVariable(variableType *ast.TypeLiteral) (*ast.ScriptVariable, error) {
	return nil, newError(p.token.SourceRange, "ParseScriptVariable unimplemented.")
}

func (p *parser) ParseIdentifier() (*ast.Identifier, error) {
	rng := p.token.SourceRange
	if err := p.tryConsume(token.Identifier); err != nil {
		return nil, err
	}
	return &ast.Identifier{
		Text:        string(bytes.ToLower(rng.Text())),
		SourceRange: rng,
	}, nil
}

func (p *parser) ParseTypeLiteral() (*ast.TypeLiteral, error) {
	return nil, newError(p.token.SourceRange, "ParseTypeLiteral unimplemented.")
}
