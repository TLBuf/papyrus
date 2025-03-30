package ast

import "github.com/TLBuf/papyrus/pkg/source"

// AccessOperator represents the dot operator for performing accesses.
type AccessOperator struct {
	Trivia
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (o *AccessOperator) SourceLocation() source.Location {
	return o.Location
}

var _ Node = (*AccessOperator)(nil)

// Access represents a value or function access via an identifier.
type Access struct {
	Trivia
	// Value is the expression that defines the value have something accessed.
	Value Expression
	// Operator is the operator for this access expression.
	Operator *AccessOperator
	// Name is the name of the variable or function being accessed in value.
	Name *Identifier
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (a *Access) SourceLocation() source.Location {
	return a.Location
}

func (*Access) expression() {}

var _ Expression = (*Access)(nil)
