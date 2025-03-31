package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Unary is an expression that computes a value from two operands.
type Unary struct {
	Trivia
	// Operator defines the operator token this unary expression uses.
	Operator Token
	// Operand is the operand.
	Operand Expression
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (u *Unary) SourceLocation() source.Location {
	return u.Location
}

func (*Unary) expression() {}

var _ Expression = (*Unary)(nil)
