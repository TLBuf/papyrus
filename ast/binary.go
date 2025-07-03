package ast

import (
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/token"
)

// BinaryKind is the type of binary operation.
type BinaryKind uint8

const (
	// Add, '+', adds the left and right operands together.
	Add = BinaryKind(token.Plus)
	// Divide, '/', divides the left operand by the right operand.
	Divide = BinaryKind(token.Divide)
	// Equal, '==', evaluates to true if the two
	// operands are equal and false otherwise.
	Equal = BinaryKind(token.Equal)
	// Greater, '>', evaluates to true if the left operand is
	// greater than the right operand and false otherwise.
	Greater = BinaryKind(token.Greater)
	// GreaterOrEqual, '>=', evaluates to true if the left operand is
	// greater than or equal to the right operand and false otherwise.
	GreaterOrEqual = BinaryKind(token.GreaterOrEqual)
	// Less, '<', evaluates to true if the left operand is
	// less than the right operand and false otherwise.
	Less = BinaryKind(token.Less)
	// LessOrEqual, '<=', evaluates to true if the left operand is
	// less than or equal to the right operand and false otherwise.
	LessOrEqual = BinaryKind(token.LessOrEqual)
	// LogicalAnd, '&&', evaluates to true if both the left operand
	// and the right operand evaluate to true and false otherwise.
	LogicalAnd = BinaryKind(token.LogicalAnd)
	// LogicalOr, '||', evaluates to true if either the left operand
	// or the right operand evaluate to true and false otherwise.
	LogicalOr = BinaryKind(token.LogicalOr)
	// Subtract, '-', subtracts the right operand from the left operand.
	Subtract = BinaryKind(token.Minus)
	// Modulo, '%', evaluates to the remainder from dividing the
	// left operand by the right operand using integer division.
	Modulo = BinaryKind(token.Modulo)
	// Multiply, '*', multiplies the left and right operands together.
	Multiply = BinaryKind(token.Multiply)
	// NotEqual, '!=', evaluates to true if the two
	// operands are not equal and false otherwise.
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
	InfixTrivia
	// Kind is the kind of binary operation this expression represents.
	Kind BinaryKind
	// LeftOperand is the operand on the left of the operator.
	LeftOperand Expression
	// OperatorLocation is the location of the binary operator.
	OperatorLocation source.Location
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
