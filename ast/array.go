package ast

import "github.com/TLBuf/papyrus/source"

// ArrayCreation is an expression that creates a new array of a fixed length.
type ArrayCreation struct {
	// Type is the type of elements the array can contain.
	Type *TypeLiteral
	// Size is the length of the array (between 1 and 128 inclusive).
	Size *IntLiteral
	// NewLocation is the location of the new operator.
	NewLocation source.Location
	// OpenLocation is the location of the opening bracket.
	OpenLocation source.Location
	// CloseLocation is the location of the closing bracket.
	CloseLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate visitor method for the node.
func (a *ArrayCreation) Accept(v Visitor) error {
	return v.VisitArrayCreation(a)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (a *ArrayCreation) Comments() *Comments {
	return a.NodeComments
}

// Location returns the source location of the node.
func (a *ArrayCreation) Location() source.Location {
	return source.Span(a.NewLocation, a.CloseLocation)
}

func (*ArrayCreation) expression() {}

var _ Expression = (*ArrayCreation)(nil)
