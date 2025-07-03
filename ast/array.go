package ast

import "github.com/TLBuf/papyrus/source"

// ArrayCreation is an expression that creates a new array of a fixed length.
type ArrayCreation struct {
	InfixTrivia
	// NewLocation is the location of the new operator.
	NewLocation source.Location
	// Type is the type of elements the array can contain.
	Type *TypeLiteral
	// OpenLocation is the location of the opening bracket.
	OpenLocation source.Location
	// Size is the length of the array (between 1 and 128 inclusive).
	Size *IntLiteral
	// CloseLocation is the location of the closing bracket.
	CloseLocation source.Location
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] assocaited with this node.
func (a *ArrayCreation) Trivia() InfixTrivia {
	return a.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (a *ArrayCreation) Accept(v Visitor) error {
	return v.VisitArrayCreation(a)
}

// Location returns the source location of the node.
func (a *ArrayCreation) Location() source.Location {
	return a.NodeLocation
}

func (*ArrayCreation) expression() {}

var _ Expression = (*ArrayCreation)(nil)
