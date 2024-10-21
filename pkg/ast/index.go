package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Index represents the access of a specific element in an array.
type Index struct {
	// Value is the expression that defines the array to reference.
	Value Expression
	// OpenOperator is the open bracket.
	OpenOperator *ArrayOpenOperator
	// Index is the expression that defines the index of the element to reference in the array.
	Index Expression
	// CloseOperator is the close bracket.
	CloseOperator *ArrayCloseOperator
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (i *Index) Range() source.Range {
	return i.SourceRange
}

func (*Index) expression() {}

func (*Index) reference() {}

var _ Reference = (*Index)(nil)
