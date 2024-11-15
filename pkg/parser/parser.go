// Package parser defines a Papyrus parser.
package parser

import (
	"bytes"

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
	if status := prsr.next(); status != statusOK {
		return nil, Error{prsr.issues}
	}
	if status := prsr.next(); status != statusOK {
		return nil, Error{prsr.issues}
	}
	script, status := prsr.ParseScript()
	if status != statusOK {
		return nil, Error{prsr.issues}
	}
	return script, nil
}

type parser struct {
	l *lexer.Lexer

	token     token.Token
	lookahead token.Token

	keepLooseComments bool
	looseComments     []token.Token

	issues []*Issue
}

// next advances token and lookahead by one token while skipping loose comment
// tokens. A non-ok status is returned if the lexer encountered an error while
// reading the input.
func (p *parser) next() status {
	p.token = p.lookahead
	t, err := p.l.NextToken()
	if err != nil {
		p.issue(err.(lexer.Error).Location, err.(lexer.Error).Message)
		return statusFatal
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
	return statusOK
}

// tryConsume advances the token position if the current token matches the given
// token type or records an issue and returns a non-ok status.
func (p *parser) tryConsume(t token.Type) status {
	if p.token.Type != t {
		p.issue(p.token.SourceRange, "expected %s, but found %s", t, p.token.Type)
		return statusError
	}
	return p.next()
}

// consumeLine advances the token position through the next newline reporting
// an issue for any non-newline tokens found.
func (p *parser) consumeLine() status {
	ok := true
	for p.token.Type != token.Newline && p.token.Type != token.EOF {
		p.issue(p.token.SourceRange, "expected %s, but found %s", token.Newline, p.token.Type)
		ok = false
	}
	if s := p.next(); s != statusOK {
		return s
	}
	if !ok {
		return statusError
	}
	return statusOK
}

func (p *parser) ParseScript() (*ast.Script, status) {
	script := &ast.Script{
		SourceRange: source.Range{
			File:   p.token.SourceRange.File,
			Length: len(p.token.SourceRange.File.Text),
			Line:   1,
			Column: 1,
		},
	}
	if status := p.ParseScriptHeader(script); status == statusFatal {
		return nil, status
	}
	return script, statusOK
}

func (p *parser) ParseScriptHeader(script *ast.Script) status {
	if status := p.tryConsume(token.ScriptName); status == statusFatal {
		return status
	}
	var status status
	script.Name, status = p.ParseIdentifier()
	if status == statusFatal {
		return status
	}
	if p.token.Type == token.Extends {
		if status := p.next(); status == statusFatal {
			return status
		}
		script.Extends, status = p.ParseIdentifier()
		if status == statusFatal {
			return status
		}
	}
	for p.token.Type == token.Hidden || p.token.Type == token.Conditional {
		if p.token.Type == token.Hidden {
			script.IsHidden = true
		} else {
			script.IsConditional = true
		}
		if status := p.next(); status == statusFatal {
			return status
		}
	}
	if status := p.consumeLine(); status == statusFatal {
		return status
	}
	return statusOK
}

func (p *parser) ParseIdentifier() (*ast.Identifier, status) {
	rng := p.token.SourceRange
	if status := p.tryConsume(token.Identifier); status != statusOK {
		return nil, status
	}
	return &ast.Identifier{
		Text:        string(bytes.ToLower(rng.Text())),
		SourceRange: rng,
	}, statusOK
}
