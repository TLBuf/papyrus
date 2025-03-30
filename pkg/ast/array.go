package ast

import "github.com/TLBuf/papyrus/pkg/source"

// NewOperator represents the new operator used to create arrays.
type NewOperator struct {
	Trivia
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (o *NewOperator) SourceLocation() source.Location {
	return o.Location
}

var _ Node = (*NewOperator)(nil)

// ArrayOpenOperator represents the open bracket for an array operation.
type ArrayOpenOperator struct {
	Trivia
	// Location is the source range of the node.
	Location source.Location
}

// Range returns the source location of the node.
func (o *ArrayOpenOperator) Range() source.Location {
	return o.Location
}

var _ Node = (*NewOperator)(nil)

// ArrayCloseOperator represents the close bracket for an array operation.
type ArrayCloseOperator struct {
	Trivia
	// Location is the source range of the node.
	Location source.Location
}

// Range returns the source location of the node.
func (o *ArrayCloseOperator) Range() source.Location {
	return o.Location
}

var _ Node = (*NewOperator)(nil)

// ArrayCreation represents an array creation expression.
type ArrayCreation struct {
	Trivia
	// NewOperator is the new operator.
	NewOperator *NewOperator
	// Type is the type of elements the array can contain.
	Type *TypeLiteral
	// OpenOperator is the open bracket.
	OpenOperator *ArrayOpenOperator
	// Size is the length of the array (between 1 and 128 inclusive).
	Size *IntLiteral
	// CloseOperator is the close bracket.
	CloseOperator *ArrayCloseOperator
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (a *ArrayCreation) SourceLocation() source.Location {
	return a.Location
}

func (*ArrayCreation) expression() {}

var _ Expression = (*ArrayCreation)(nil)
