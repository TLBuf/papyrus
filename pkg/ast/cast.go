package ast

import "github.com/TLBuf/papyrus/pkg/source"

// AsOperator represents the as operator used to cast values.
type AsOperator struct {
	Trivia
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (o *AsOperator) SourceLocation() source.Location {
	return o.Location
}

var _ Node = (*AsOperator)(nil)

// Cast is an expression casts a value of some type to another.
type Cast struct {
	Trivia
	// Value is the expression being cast to a new type.
	Value Expression
	// Operator is the as operator.
	Operator *AsOperator
	// Type is the type the value is being cast to.
	Type *TypeLiteral
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (c *Cast) SourceLocation() source.Location {
	return c.Location
}

func (*Cast) expression() {}

var _ Expression = (*Cast)(nil)
