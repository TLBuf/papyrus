package ast

import (
	"fmt"

	"github.com/TLBuf/papyrus/source"
)

// Variable is a variable definition.
type Variable struct {
	// HasLeadingBlankLine is true if this node was preceded by a blank line.
	HasLeadingBlankLine bool
	// Type is the type literal that defines the type of the variable.
	Type *TypeLiteral
	// Name is the name of the variable.
	Name *Identifier
	// Value is the expression the variable is assigned or nil if there isn't one
	// (and the variable should have the default value for its type).
	//
	// Variables defined at the script level  (i.e. not in a function or event)
	// may only set a [Literal] value.
	Value Expression
	// ConditionalLocations are the locations of the Conditional keywords that
	// mark this variable as conditional (i.e. it can appear in conditions) or
	// empty if this variable is not conditional.
	//
	// Only variables defined at the script level (i.e. not in a function or
	// event) may set this field.
	//
	// Errata: This being multiple values is due to the official Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the variable
	// conditional.
	ConditionalLocations []source.Location
	// OperatorLocation is the location of the assignment operator.
	//
	// This is only valid if Value is not nil.
	OperatorLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// LeadingBlankLine returns true if this node was preceded by a blank line.
func (v *Variable) LeadingBlankLine() bool {
	return v.HasLeadingBlankLine
}

// Accept calls the appropriate visitor method for the node.
func (v *Variable) Accept(r Visitor) error {
	return r.VisitVariable(v)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (v *Variable) Comments() *Comments {
	return v.NodeComments
}

// Location returns the source location of the node.
func (v *Variable) Location() source.Location {
	end := v.Name.Location()
	if v.Value != nil {
		end = v.Value.Location()
	}
	if len(v.ConditionalLocations) > 0 {
		end = v.ConditionalLocations[len(v.ConditionalLocations)-1]
	}
	return source.Span(v.Type.Location(), end)
}

func (v *Variable) String() string {
	return fmt.Sprintf("Variable%s", v.Location())
}

func (*Variable) statement() {}

func (*Variable) scriptStatement() {}

func (*Variable) functionStatement() {}

var (
	_ ScriptStatement   = (*Variable)(nil)
	_ FunctionStatement = (*Variable)(nil)
)
