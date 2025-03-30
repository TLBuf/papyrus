package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Parenthetical represents a parenthesized expression.
type Parenthetical struct {
	Trivia
	// Value is the expression within the parentheses.
	Value Expression
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (p *Parenthetical) SourceLocation() source.Location {
	return p.Location
}

func (*Parenthetical) expression() {}

var _ Expression = (*Parenthetical)(nil)
