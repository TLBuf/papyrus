package ast

import (
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/token"
)

// BinaryKind is the type of binary operation.
type BinaryKind uint8

const (
	// Add is the operation that adds the left and right operands together.
	Add = BinaryKind(token.Plus)
	// Divide is the operation that  divides
	// the left operand by the right operand.
	Divide = BinaryKind(token.Divide)
	// Equal is the operation that evaluates to true
	// if the two operands are equal and false otherwise.
	Equal = BinaryKind(token.Equal)
	// Greater is the operation that evaluates to true if the left
	// operand is greater than the right operand and false otherwise.
	Greater = BinaryKind(token.Greater)
	// GreaterOrEqual is the operation that  evaluates to true if the left
	// operand is greater than or equal to the right operand and false otherwise.
	GreaterOrEqual = BinaryKind(token.GreaterOrEqual)
	// Less is the operation that evaluates to true if the left
	// operand is less than the right operand and false otherwise.
	Less = BinaryKind(token.Less)
	// LessOrEqual is the operation that evaluates to true if the left operand
	// is less than or equal to the right operand and false otherwise.
	LessOrEqual = BinaryKind(token.LessOrEqual)
	// LogicalAnd is the operation that evaluates to true if both the left
	// operand and the right operand evaluate to true and false otherwise.
	LogicalAnd = BinaryKind(token.LogicalAnd)
	// LogicalOr is the operation that evaluates to true if either the left
	// operand or the right operand evaluate to true and false otherwise.
	LogicalOr = BinaryKind(token.LogicalOr)
	// Subtract is the operation that subtracts
	// the right operand from the left operand.
	Subtract = BinaryKind(token.Minus)
	// Modulo is the operation that evaluates to the remainder from dividing
	// the left operand by the right operand using integer division.
	Modulo = BinaryKind(token.Modulo)
	// Multiply is the operation that multiplies
	// the left and right operands together.
	Multiply = BinaryKind(token.Multiply)
	// NotEqual is the operation that evaluates to true if
	// the two operands are not equal and false otherwise.
	NotEqual = BinaryKind(token.NotEqual)
)

// Symbol returns the string representation of this
// kind or an empty string if it is invalid.
func (k BinaryKind) Symbol() string {
	switch k {
	case Add:
		return token.Plus.Symbol()
	case Divide:
		return token.Divide.Symbol()
	case Equal:
		return token.Equal.Symbol()
	case Greater:
		return token.Greater.Symbol()
	case GreaterOrEqual:
		return token.GreaterOrEqual.Symbol()
	case Less:
		return token.Less.Symbol()
	case LessOrEqual:
		return token.LessOrEqual.Symbol()
	case LogicalAnd:
		return token.LogicalAnd.Symbol()
	case LogicalOr:
		return token.LogicalOr.Symbol()
	case Subtract:
		return token.Minus.Symbol()
	case Modulo:
		return token.Modulo.Symbol()
	case Multiply:
		return token.Multiply.Symbol()
	case NotEqual:
		return token.NotEqual.Symbol()
	default:
		return ""
	}
}

// Binary is an expression that computes a value from two operands.
type Binary struct {
	// Kind is the kind of binary operation this expression represents.
	Kind BinaryKind
	// LeftOperand is the operand on the left of the operator.
	LeftOperand Expression
	// RightOperand is the operand on the right of the operator.
	RightOperand Expression
	// OperatorLocation is the location of the binary operator.
	OperatorLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate visitor method for the node.
func (b *Binary) Accept(v Visitor) error {
	return v.VisitBinary(b)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (b *Binary) Comments() *Comments {
	return b.NodeComments
}

// Location returns the source location of the node.
func (b *Binary) Location() source.Location {
	return source.Span(b.LeftOperand.Location(), b.RightOperand.Location())
}

func (*Binary) expression() {}

var _ Expression = (*Binary)(nil)
