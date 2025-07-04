package ast

import "github.com/TLBuf/papyrus/source"

// Index represents the access of a specific element in an array.
type Index struct {
	// Value is the expression that defines the array to reference.
	Value Expression
	// Index is the expression that defines the index of the element to reference
	// in the array.
	Index Expression
	// OpenLocation is the location of the opening bracket.
	OpenLocation source.Location
	// CloseLocation is the location of the closing bracket.
	CloseLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *InlineComments
}

// Accept calls the appropriate visitor method for the node.
func (i *Index) Accept(v Visitor) error {
	return v.VisitIndex(i)
}

// Comments returns the [InlineComments] associated
// with this node or nil if there are none.
func (i *Index) Comments() *InlineComments {
	return i.NodeComments
}

// Location returns the source location of the node.
func (i *Index) Location() source.Location {
	return source.Span(i.Value.Location(), i.CloseLocation)
}

func (*Index) expression() {}

var _ Expression = (*Index)(nil)
