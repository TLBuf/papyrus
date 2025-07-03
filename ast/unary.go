package ast

import "github.com/TLBuf/papyrus/source"

// Unary is an expression that computes a value from two operands.
type Unary struct {
	InfixTrivia
	// Operator defines the operator token this unary expression uses.
	Operator *Token
	// Operand is the operand.
	Operand Expression
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] assocaited with this node.
func (u *Unary) Trivia() InfixTrivia {
	return u.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (u *Unary) Accept(v Visitor) error {
	return v.VisitUnary(u)
}

// Location returns the source location of the node.
func (u *Unary) Location() source.Location {
	return u.NodeLocation
}

func (*Unary) expression() {}

var _ Expression = (*Unary)(nil)
