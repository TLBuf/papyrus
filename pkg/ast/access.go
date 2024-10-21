package ast

import "github.com/TLBuf/papyrus/pkg/source"

// AccessOperator represents the dot operator for performing accesses.
type AccessOperator struct {
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (o *AccessOperator) Range() source.Range {
	return o.SourceRange
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
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (a *Access) Range() source.Range {
	return a.SourceRange
}

func (*Access) expression() {}

func (*Access) reference() {}

var _ Reference = (*Access)(nil)
