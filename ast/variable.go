package ast

import "github.com/TLBuf/papyrus/source"

// ScriptVariable is a variable definition at the script level.
type ScriptVariable struct {
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
	// Type is the type literal that defines the type of the variable.
	Type *TypeLiteral
	// Name is the name of the variable.
	Name *Identifier
	// Value is the literal the script variable is assigned or nil if there isn't
	// one (and the variable should have the default value for its type).
	Value Literal
	// ConditionalLocations are the locations of the Conditional keywords that
	// mark this variable as conditional (i.e. it can appear in conditions) or
	// empty if this variable is not conditional.
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
	// NodeComments are the comments on lines before and/or after a
	// node or nil if the node has no comments associated with it.
	NodeComments *CrosslineComments
}

// PrecedingBlankLine returns true if this node was preceded by a blank line.
func (s *ScriptVariable) PrecedingBlankLine() bool {
	return s.HasPrecedingBlankLine
}

// Accept calls the appropriate visitor method for the node.
func (s *ScriptVariable) Accept(v Visitor) error {
	return v.VisitScriptVariable(s)
}

// Comments returns the [CrosslineComments] associated
// with this node or nil if there are none.
func (s *ScriptVariable) Comments() *CrosslineComments {
	return s.NodeComments
}

// Location returns the source location of the node.
func (s *ScriptVariable) Location() source.Location {
	end := s.Name.Location()
	if s.Value != nil {
		end = s.Value.Location()
	}
	if len(s.ConditionalLocations) > 0 {
		end = s.ConditionalLocations[len(s.ConditionalLocations)-1]
	}
	return source.Span(s.Type.Location(), end)
}

func (*ScriptVariable) statement() {}

func (*ScriptVariable) scriptStatement() {}

var _ ScriptStatement = (*ScriptVariable)(nil)

// FunctionVariable is a variable definition within the body of a function (or
// event).
type FunctionVariable struct {
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
	// Type is the type literal that defines the type of the variable.
	Type *TypeLiteral
	// Name is the name of the variable.
	Name *Identifier
	// Value is the expression the variable is assigned or nil if there isn't one
	// (and the variable should have the default value for its type).
	Value Expression
	// OperatorLocation is the location of the assignment operator.
	//
	// This is only valid if Value is not nil.
	OperatorLocation source.Location
	// NodeComments are the comments on lines before and/or after a
	// node or nil if the node has no comments associated with it.
	NodeComments *CrosslineComments
}

// PrecedingBlankLine returns true if this node was preceded by a blank line.
func (f *FunctionVariable) PrecedingBlankLine() bool {
	return f.HasPrecedingBlankLine
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (f *FunctionVariable) Accept(v Visitor) error {
	return v.VisitFunctionVariable(f)
}

// Comments returns the [CrosslineComments] associated
// with this node or nil if there are none.
func (f *FunctionVariable) Comments() *CrosslineComments {
	return f.NodeComments
}

// Location returns the source location of the node.
func (f *FunctionVariable) Location() source.Location {
	if f.Value != nil {
		return source.Span(f.Type.Location(), f.Value.Location())
	}
	return source.Span(f.Type.Location(), f.Name.Location())
}

func (*FunctionVariable) statement() {}

func (*FunctionVariable) functionStatement() {}

var _ FunctionStatement = (*FunctionVariable)(nil)
