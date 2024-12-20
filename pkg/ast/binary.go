package ast

import "github.com/TLBuf/papyrus/pkg/source"

// BinaryOperatorKind defines the various types of binary operations.
type BinaryOperatorKind int

const (
	// UnknownBinaryOperatorKind is the default (and invalid) binary operator.
	UnknownBinaryOperatorKind BinaryOperatorKind = iota
	// LogicalOr is the logical OR operator, '||'.
	LogicalOr
	// LogicalAnd is the logical AND operator, '&&'.
	LogicalAnd
	// Equal is the equality operator, '=='.
	Equal
	// NotEqual is the inequality operator, '!='.
	NotEqual
	// Greater is the greater than releational operator, '>'.
	Greater
	// GreaterOrEqual is the greater than or equal to releational operator, '>='.
	GreaterOrEqual
	// Less is the less than releational operator, '<'.
	Less
	// LessOrEqual is the less than or equal to releational operator, '<='.
	LessOrEqual
	// Add is the addition operator, '+'.
	Add
	// Subtract is the subtraction operator, '-'.
	Subtract
	// Multiply is the multiplication operator, '*'.
	Multiply
	// Divide is the division operator, '/'.
	Divide
	// Modulo is the modulus operator, '%'.
	Modulo
)

func (o BinaryOperatorKind) String() string {
	name, ok := BinaryOperatorKindNames[o]
	if ok {
		return name
	}
	return "<unknown>"
}

var BinaryOperatorKindNames = map[BinaryOperatorKind]string{
	LogicalOr:      "||",
	LogicalAnd:     "&&",
	Equal:          "==",
	NotEqual:       "!=",
	Greater:        ">",
	GreaterOrEqual: ">=",
	Less:           "<",
	LessOrEqual:    "<=",
	Add:            "+",
	Subtract:       "-",
	Multiply:       "*",
	Divide:         "/",
	Modulo:         "%",
}

// BinaryOperator represents a binary operator.
type BinaryOperator struct {
	// Kind is the type of binary operator.
	Kind BinaryOperatorKind
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (o *BinaryOperator) Range() source.Range {
	return o.SourceRange
}

var _ Node = (*BinaryOperator)(nil)

// Binary is an expression that computes a value from two operands.
type Binary struct {
	// LeftOperand is the operand on the left of the operator.
	LeftOperand Expression
	// Operator defines the operator this binary expression uses.
	Operator *BinaryOperator
	// RightOperand is the operand on the right of the operator.
	RightOperand Expression
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (b *Binary) Range() source.Range {
	return b.SourceRange
}

func (*Binary) expression() {}

var _ Expression = (*Binary)(nil)
