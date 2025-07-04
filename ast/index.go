package ast

import "github.com/TLBuf/papyrus/source"

// Index represents the access of a specific element in an array.
type Index struct {
	InfixTrivia

	// Value is the expression that defines the array to reference.
	Value Expression
	// OpenLocation is the location of the opening bracket.
	OpenLocation source.Location
	// Index is the expression that defines the index of the element to reference
	// in the array.
	Index Expression
	// CloseLocation is the location of the closing bracket.
	CloseLocation source.Location
}

// Trivia returns the [InfixTrivia] associated with this node.
func (i *Index) Trivia() InfixTrivia {
	return i.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (i *Index) Accept(v Visitor) error {
	return v.VisitIndex(i)
}

// Location returns the source location of the node.
func (i *Index) Location() source.Location {
	return source.Span(i.Value.Location(), i.CloseLocation)
}

func (*Index) expression() {}

var _ Expression = (*Index)(nil)
