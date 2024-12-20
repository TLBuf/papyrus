package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Argument is a named argument for a function call.
type Argument struct {
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
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (a *Argument) Range() source.Range {
	return a.SourceRange
}

var _ Node = (*Argument)(nil)
