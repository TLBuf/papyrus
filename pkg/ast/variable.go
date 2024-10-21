package ast

import "github.com/TLBuf/papyrus/pkg/source"

// ScriptVariable is a variable definition at the script level.
type ScriptVariable struct {
	// Type is the type literal that defines the type of the variable.
	Type *TypeLiteral
	// Name is the name of the variable.
	Name *Identifier
	// Value is the literal the script variable is assigned or nil if there isn't one
	// (and the variable should have the default value for its type).
	Value Literal
	// IsConditional
	IsConditional bool
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (v *ScriptVariable) Range() source.Range {
	return v.SourceRange
}

func (*ScriptVariable) scriptStatement() {}

var _ ScriptStatement = (*ScriptVariable)(nil)

// FunctionVariable is a variable definition within the body of a function (or event).
type FunctionVariable struct {
	// Type is the type literal that defines the type of the variable.
	Type *TypeLiteral
	// Name is the name of the variable.
	Name *Identifier
	// Value is the expression the variable is assigned or nil if there isn't one
	// (and the variable should have the default value for its type).
	Value Expression
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (v *FunctionVariable) Range() source.Range {
	return v.SourceRange
}

func (*FunctionVariable) functionStatement() {}

var _ FunctionStatement = (*FunctionVariable)(nil)
