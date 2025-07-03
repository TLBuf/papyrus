package ast

import "github.com/TLBuf/papyrus/source"

// ScriptVariable is a variable definition at the script level.
type ScriptVariable struct {
	LineTrivia
	// Type is the type literal that defines the type of the variable.
	Type *TypeLiteral
	// Name is the name of the variable.
	Name *Identifier
	// Operator is the assignment operator or nil if no inital value is assigned.
	Operator *Token
	// Value is the literal the script variable is assigned or nil if there isn't
	// one (and the variable should have the default value for its type).
	Value Literal
	// Conditional are the Conditional tokens that define that this variable is
	// conditional (i.e. it can appear in conditions) or empty if this variable is
	// not conditional.
	//
	// Errata: This being multiple values is due to the offical Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the variable
	// conditional.
	Conditional []*Token
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (s *ScriptVariable) Trivia() LineTrivia {
	return s.LineTrivia
}

// Accept calls the appropriate visitor method for the node.
func (s *ScriptVariable) Accept(v Visitor) error {
	return v.VisitScriptVariable(s)
}

// Location returns the source location of the node.
func (s *ScriptVariable) Location() source.Location {
	return s.NodeLocation
}

func (*ScriptVariable) statement() {}

func (*ScriptVariable) scriptStatement() {}

var _ ScriptStatement = (*ScriptVariable)(nil)

// FunctionVariable is a variable definition within the body of a function (or
// event).
type FunctionVariable struct {
	LineTrivia
	// Type is the type literal that defines the type of the variable.
	Type *TypeLiteral
	// Name is the name of the variable.
	Name *Identifier
	// Operator is the assignment operator or nil if no inital value is assigned.
	Operator *Token
	// Value is the expression the variable is assigned or nil if there isn't one
	// (and the variable should have the default value for its type).
	Value Expression
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (f *FunctionVariable) Trivia() LineTrivia {
	return f.LineTrivia
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (f *FunctionVariable) Accept(v Visitor) error {
	return v.VisitFunctionVariable(f)
}

// Location returns the source location of the node.
func (f *FunctionVariable) Location() source.Location {
	return f.NodeLocation
}

func (*FunctionVariable) statement() {}

func (*FunctionVariable) functionStatement() {}

var _ FunctionStatement = (*FunctionVariable)(nil)
