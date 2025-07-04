package ast

import "github.com/TLBuf/papyrus/source"

// Argument is a named argument for a function call.
type Argument struct {
	InfixTrivia

	// Name is the name of the parameter for this argument or nil if using
	// positional syntax.
	Name *Identifier
	// OperatorLocation is the location of the assignment operator.
	//
	// This is only valid if Name is not nil.
	OperatorLocation source.Location
	// Value is the expression that defines the value of this argument.
	Value Expression
}

// Trivia returns the [InfixTrivia] associated with this node.
func (a *Argument) Trivia() InfixTrivia {
	return a.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (a *Argument) Accept(v Visitor) error {
	return v.VisitArgument(a)
}

// Location returns the source location of the node.
func (a *Argument) Location() source.Location {
	if a.Name == nil {
		return a.Value.Location()
	}
	return source.Span(a.Name.Location(), a.Value.Location())
}

var _ Node = (*Argument)(nil)
