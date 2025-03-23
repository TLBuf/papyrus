package ast

import "github.com/TLBuf/papyrus/pkg/source"

// AccessOperator represents the dot operator for performing accesses.
type AccessOperator struct {
	// Location is the source range of the node.
	Location source.Range
}

// Range returns the source range of the node.
func (o *AccessOperator) Range() source.Range {
	return o.Location
}

var _ Node = (*AccessOperator)(nil)

// Access represents a value or function access via an identifier.
type Access struct {
	// Value is the expression that defines the value have something accessed.
	Value Expression
	// Operator is the operator for this access expression.
	Operator *AccessOperator
	// Name is the name of the variable or function being accessed in value.
	Name *Identifier
	// Location is the source range of the node.
	Location source.Range
}

// Range returns the source range of the node.
func (a *Access) Range() source.Range {
	return a.Location
}

func (*Access) expression() {}

var _ Expression = (*Access)(nil)
