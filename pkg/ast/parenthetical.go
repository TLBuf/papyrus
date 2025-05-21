package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Parenthetical represents a parenthesized expression.
type Parenthetical struct {
	Trivia
	// Open is the open parenthesis token.
	Open *Token
	// Value is the expression within the parentheses.
	Value Expression
	// Close is the close parenthesis token.
	Close *Token
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (p *Parenthetical) Accept(v Visitor) error {
	return v.VisitParenthetical(p)
}

// SourceLocation returns the source location of the node.
func (p *Parenthetical) SourceLocation() source.Location {
	return p.Location
}

func (*Parenthetical) expression() {}

func (*Parenthetical) functionStatement() {}

var _ Expression = (*Parenthetical)(nil)
