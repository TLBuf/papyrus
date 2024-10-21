package ast

import "github.com/TLBuf/papyrus/pkg/source"

// AsOperator represents the as operator used to cast values.
type AsOperator struct {
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (o *AsOperator) Range() source.Range {
	return o.SourceRange
}

var _ Node = (*AsOperator)(nil)

// Cast is an expression casts a value of some type to another.
type Cast struct {
	// Value is the expression being cast to a new type.
	Value Expression
	// Operator is the as operator.
	Operator *AsOperator
	// Type is the type the value is being cast to.
	Type *TypeLiteral
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (c *Cast) Range() source.Range {
	return c.SourceRange
}

func (*Cast) expression() {}

var _ Expression = (*Cast)(nil)
