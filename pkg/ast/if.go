package ast

import "github.com/TLBuf/papyrus/pkg/source"

// If is a statement that evaluates some set of statements if a condition is
// true and potentially a different set of statements if that condition is
// false.
type If struct {
	// Conditional is the main conditional block of the if statement.
	Conditional ConditionalBlock
	// AlternativeConditonals are the ordered set of alternative contional blocks,
	// i.e. ElseIfs.
	AlternativeConditionals []ConditionalBlock
	// Alternative is the list of statements that should be evaluated if the
	// condition and all alternative conditionals are false, i.e. Else.
	Alternative []FunctionStatement
	// Location is the source range of the node.
	Location source.Range
}

// ConditionalBlock is a list of statements that may be conditionally executed.
type ConditionalBlock struct {
	// Condition is the expression that defines the condition to check.
	Condition Expression
	// Statements is the list of statements that should be evaluated if the
	// condition is true.
	Statements []FunctionStatement
}

// Range returns the source range of the node.
func (i *If) Range() source.Range {
	return i.Location
}

func (*If) functionStatement() {}

var _ FunctionStatement = (*If)(nil)
