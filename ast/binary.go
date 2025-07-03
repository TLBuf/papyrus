package ast

import "github.com/TLBuf/papyrus/source"

// Binary is an expression that computes a value from two operands.
type Binary struct {
	InfixTrivia
	// LeftOperand is the operand on the left of the operator.
	LeftOperand Expression
	// Operator defines the operator token this binary expression uses.
	Operator *Token
	// RightOperand is the operand on the right of the operator.
	RightOperand Expression
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] assocaited with this node.
func (b *Binary) Trivia() InfixTrivia {
	return b.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (b *Binary) Accept(v Visitor) error {
	return v.VisitBinary(b)
}

// Location returns the source location of the node.
func (b *Binary) Location() source.Location {
	return b.NodeLocation
}

func (*Binary) expression() {}

var _ Expression = (*Binary)(nil)
