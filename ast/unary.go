package ast

import "github.com/TLBuf/papyrus/source"

// Unary is an expression that computes a value from two operands.
type Unary struct {
	Trivia
	// Operator defines the operator token this unary expression uses.
	Operator *Token
	// Operand is the operand.
	Operand Expression
	// Location is the source range of the node.
	Location source.Location
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
