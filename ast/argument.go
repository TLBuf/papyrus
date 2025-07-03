package ast

import "github.com/TLBuf/papyrus/source"

// Argument is a named argument for a function call.
type Argument struct {
	InfixTrivia
	// Name is the name of the parameter for this argument or nil if using
	// positional syntax.
	Name *Identifier
	// Operator is the assignment operator token between the name and value or
	// nil if using positional syntax.
	Operator *Token
	// Value is the expression that defines the value of this argument.
	Value Expression
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] assocaited with this node.
func (a *Argument) Trivia() InfixTrivia {
	return a.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (a *Argument) Accept(v Visitor) error {
	return v.VisitArgument(a)
}

// Location returns the source location of the node.
func (a *Argument) Location() source.Location {
	return a.NodeLocation
}

var _ Node = (*Argument)(nil)
