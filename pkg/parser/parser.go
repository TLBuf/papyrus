// Package parser defines a Papyrus parser.
package parser

import (
	"github.com/TLBuf/papyrus/pkg/ast"
	"github.com/TLBuf/papyrus/pkg/lexer"
	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/TLBuf/papyrus/pkg/token"
)

// Parser provides the ability to lex and parse a Papyrus script into an
// [*ast.Script].
type Parser struct {
	l *lexer.Lexer

	token     token.Token
	lookahead token.Token

	keepLooseComments bool
	looseComments     []token.Token
}

// Error defines an error raised by the parser.
type Error struct {
	msg string
	// SourceRange is the source range of the segment of input text that caused an
	// error.
	SourceRange source.Range
}

// Error implments the error interface.
func (e Error) Error() string {
	return e.msg
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

func New(l *lexer.Lexer, opts ...Option) *Parser {
	p := &Parser{l: l}
	for _, opt := range opts {
		opt(p)
	}
	p.next()
	p.next()
	return p
}

func (p *Parser) next() error {
	p.token = p.lookahead
	t, err := p.l.NextToken()
	if err != nil {
		return err
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

func (p *Parser) ParseScript() (*ast.Script, error) {
	return nil, nil
}
