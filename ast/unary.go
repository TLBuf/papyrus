package ast

import "github.com/TLBuf/papyrus/source"

// Unary is an expression that computes a value from two operands.
type Unary struct {
	InfixTrivia
	// Operator defines the operator token this unary expression uses.
	Operator *Token
	// Operand is the operand.
	Operand Expression
	// Location is the source range of the node.
	Location source.Location
}

// Trivia returns the [InfixTrivia] assocaited with this node.
func (u *Unary) Trivia() InfixTrivia {
	return u.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (u *Unary) Accept(v Visitor) error {
	return v.VisitUnary(u)
}

// SourceLocation returns the source location of the node.
func (u *Unary) SourceLocation() source.Location {
	return u.Location
}

func (*Unary) expression() {}

var _ Expression = (*Unary)(nil)
