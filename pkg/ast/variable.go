package ast

import "github.com/TLBuf/papyrus/pkg/source"

// ScriptVariable is a variable definition at the script level.
type ScriptVariable struct {
	Trivia
	// Type is the type literal that defines the type of the variable.
	Type *TypeLiteral
	// Name is the name of the variable.
	Name *Identifier
	// Operator is the assignment operator or nil if no inital value is assigned.
	Operator Token
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
	Conditional []Token
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (v *ScriptVariable) SourceLocation() source.Location {
	return v.Location
}

func (*ScriptVariable) scriptStatement() {}

var _ ScriptStatement = (*ScriptVariable)(nil)

// FunctionVariable is a variable definition within the body of a function (or
// event).
type FunctionVariable struct {
	Trivia
	// Type is the type literal that defines the type of the variable.
	Type *TypeLiteral
	// Name is the name of the variable.
	Name *Identifier
	// Operator is the assignment operator or nil if no inital value is assigned.
	Operator Token
	// Value is the expression the variable is assigned or nil if there isn't one
	// (and the variable should have the default value for its type).
	Value Expression
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (v *FunctionVariable) SourceLocation() source.Location {
	return v.Location
}

func (*FunctionVariable) functionStatement() {}

var _ FunctionStatement = (*FunctionVariable)(nil)
