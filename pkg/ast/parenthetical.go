package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Parenthetical represents a parenthesized expression.
type Parenthetical struct {
	// Value is the expression within the parentheses.
	Value Expression
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (p *Parenthetical) Range() source.Range {
	return p.SourceRange
}

func (*Parenthetical) expression() {}

var _ Expression = (*Parenthetical)(nil)
