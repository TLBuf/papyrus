package ast

import "github.com/TLBuf/papyrus/pkg/source"

// NewOperator represents the new operator used to create arrays.
type NewOperator struct {
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (o *NewOperator) Range() source.Range {
	return o.SourceRange
}

var _ Node = (*NewOperator)(nil)

// ArrayOpenOperator represents the open bracket for an array operation.
type ArrayOpenOperator struct {
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (o *ArrayOpenOperator) Range() source.Range {
	return o.SourceRange
}

var _ Node = (*NewOperator)(nil)

// ArrayCloseOperator represents the close bracket for an array operation.
type ArrayCloseOperator struct {
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (o *ArrayCloseOperator) Range() source.Range {
	return o.SourceRange
}

var _ Node = (*NewOperator)(nil)

// ArrayCreation represents an array creation expression.
type ArrayCreation struct {
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
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (a *ArrayCreation) Range() source.Range {
	return a.SourceRange
}

func (*ArrayCreation) expression() {}

var _ Expression = (*ArrayCreation)(nil)
