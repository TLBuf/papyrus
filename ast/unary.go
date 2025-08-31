package ast

import (
	"fmt"

	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/token"
)

// UnaryKind is the type of unary operation.
type UnaryKind uint8

const (
	// Negate is the operation that flips the sign of the numeric operand.
	Negate = UnaryKind(token.Minus)
	// LogicalNot is the operation that evaluates to true if
	// the operand evaluates to false and false otherwise.
	LogicalNot = UnaryKind(token.LogicalNot)
)

// Symbol returns the string representation of this
// kind or an empty string if it is invalid.
func (k UnaryKind) Symbol() string {
	switch k {
	case Negate:
		return token.Minus.Symbol()
	case LogicalNot:
		return token.LogicalNot.Symbol()
	default:
		return ""
	}
}

// Unary is an expression that computes a value from two operands.
type Unary struct {
	// Kind is the kind of unary operation this expression represents.
	Kind UnaryKind
	// Operand is the operand.
	Operand Expression
	// OperatorLocation is the location of the unary operator.
	OperatorLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate visitor method for the node.
func (u *Unary) Accept(v Visitor) error {
	return v.VisitUnary(u)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (u *Unary) Comments() *Comments {
	return u.NodeComments
}

// Location returns the source location of the node.
func (u *Unary) Location() source.Location {
	return source.Span(u.OperatorLocation, u.Operand.Location())
}

func (u *Unary) String() string {
	return fmt.Sprintf("Unary%s", u.Location())
}

func (*Unary) expression() {}

var _ Expression = (*Unary)(nil)
