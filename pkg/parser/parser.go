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
	if issue := prsr.next(); issue != nil {
		return nil, Error{[]*Issue{issue}}
	}
	if issue := prsr.next(); issue != nil {
		return nil, Error{[]*Issue{issue}}
	}
	script, issues := prsr.ParseScript()
	if len(issues) > 0 {
		return nil, Error{issues}
	}
	return script, nil
}

type parser struct {
	l *lexer.Lexer

	token     token.Token
	lookahead token.Token

	keepLooseComments bool
	looseComments     []token.Token
}

// next advances token and lookahead by one token while skipping loose comment
// tokens. A non-nil issue is returned if the lexer encountered an error while
// reading the input; such errors are fatal and should halt parsing.
func (p *parser) next() *Issue {
	p.token = p.lookahead
	t, err := p.l.NextToken()
	if err != nil {
		return &Issue{
			Message:  err.(lexer.Error).Message,
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

func (p *parser) ParseScript() (*ast.Script, []*Issue) {
	return nil, nil
}
