package ast

import (
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/token"
)

// UnaryKind is the type of unary operation.
type UnaryKind uint8

const (
	// Negate, '-', flips the sign of the numeric operand.
	Negate = UnaryKind(token.Minus)
	// LogicalNot, '!', evaluates to true if the
	// operand evaluates to false and false otherwise.
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
	InfixTrivia
	// Kind is the kind of unary operation this expression represents.
	Kind UnaryKind
	// OperatorLocation is the location of the unary operator.
	OperatorLocation source.Location
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
