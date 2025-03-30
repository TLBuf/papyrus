package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Argument is a named argument for a function call.
type Argument struct {
	Trivia
	// Name is the name of the parameter for this argument or nil if using
	// positional syntax.
	Name *Identifier
	// Operator is the assignment operator betweent he name and value or nil if
	// using positional syntax.
	//
	// When present the kind is always [Assign].
	Operator *AssignmentOperator
	// Value is the expression that defines the value of this argument.
	Value Expression
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (a *Argument) SourceLocation() source.Location {
	return a.Location
}

var _ Node = (*Argument)(nil)
