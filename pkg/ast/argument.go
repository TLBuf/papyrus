package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Argument is a named argument for a function call.
type Argument struct {
	Trivia
	// Name is the name of the parameter for this argument or nil if using
	// positional syntax.
	Name *Identifier
	// Operator is the assignment operator token between the name and value or
	// nil if using positional syntax.
	Operator *Token
	// Value is the expression that defines the value of this argument.
	Value Expression
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (a *Argument) Accept(v Visitor) error {
	return v.VisitArgument(a)
}

// SourceLocation returns the source location of the node.
func (a *Argument) SourceLocation() source.Location {
	return a.Location
}

var _ Node = (*Argument)(nil)
