package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Index represents the access of a specific element in an array.
type Index struct {
	Trivia
	// Value is the expression that defines the array to reference.
	Value Expression
	// OpenOperator is the open bracket.
	OpenOperator *ArrayOpenOperator
	// Index is the expression that defines the index of the element to reference
	// in the array.
	Index Expression
	// CloseOperator is the close bracket.
	CloseOperator *ArrayCloseOperator
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (i *Index) SourceLocation() source.Location {
	return i.Location
}

func (*Index) expression() {}

var _ Expression = (*Index)(nil)
