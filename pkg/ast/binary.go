package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Binary is an expression that computes a value from two operands.
type Binary struct {
	Trivia
	// LeftOperand is the operand on the left of the operator.
	LeftOperand Expression
	// Operator defines the operator token this binary expression uses.
	Operator *Token
	// RightOperand is the operand on the right of the operator.
	RightOperand Expression
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (b *Binary) Accept(v Visitor) error {
	return v.VisitBinary(b)
}

// SourceLocation returns the source location of the node.
func (b *Binary) SourceLocation() source.Location {
	return b.Location
}

func (*Binary) expression() {}

func (*Binary) functionStatement() {}

var _ Expression = (*Binary)(nil)
