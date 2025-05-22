package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Index represents the access of a specific element in an array.
type Index struct {
	Trivia
	// Value is the expression that defines the array to reference.
	Value Expression
	// Open is the open bracket token.
	Open *Token
	// Index is the expression that defines the index of the element to reference
	// in the array.
	Index Expression
	// Close is the close bracket token.
	Close *Token
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (i *Index) SourceLocation() source.Location {
	return i.Location
}

func (*Index) expression() {}

func (*Index) functionStatement() {}

var _ Expression = (*Index)(nil)
