package ast

import (
	"fmt"

	"github.com/TLBuf/papyrus/source"
)

// Argument is a named argument for a function call.
type Argument struct {
	// Name is the name of the parameter for this argument or nil if using
	// positional syntax.
	Name *Identifier
	// Value is the expression that defines the value of this argument.
	Value Expression
	// OperatorLocation is the location of the assignment operator.
	//
	// This is only valid if Name is not nil.
	OperatorLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate visitor method for the node.
func (a *Argument) Accept(v Visitor) error {
	return v.VisitArgument(a)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (a *Argument) Comments() *Comments {
	return a.NodeComments
}

// Location returns the source location of the node.
func (a *Argument) Location() source.Location {
	if a.Name == nil {
		return a.Value.Location()
	}
	return source.Span(a.Name.Location(), a.Value.Location())
}

func (a *Argument) String() string {
	return fmt.Sprintf("Argument%s", a.Location())
}

var _ Node = (*Argument)(nil)
