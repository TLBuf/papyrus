package ast

import "github.com/TLBuf/papyrus/pkg/source"

// ArrayCreation is an expression that creates a new array of a fixed length.
type ArrayCreation struct {
	Trivia
	// NewOperator is the new operator token.
	NewOperator *Token
	// Type is the type of elements the array can contain.
	Type *TypeLiteral
	// OpenOperator is the open bracket token.
	Open *Token
	// Size is the length of the array (between 1 and 128 inclusive).
	Size *IntLiteral
	// CloseOperator is the close bracket token
	Close *Token
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate visitor method for the node.
func (a *ArrayCreation) Accept(v Visitor) error {
	return v.VisitArrayCreation(a)
}

// SourceLocation returns the source location of the node.
func (a *ArrayCreation) SourceLocation() source.Location {
	return a.Location
}

func (*ArrayCreation) expression() {}

func (*ArrayCreation) functionStatement() {}

var _ Expression = (*ArrayCreation)(nil)
