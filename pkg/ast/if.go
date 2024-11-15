package ast

import "github.com/TLBuf/papyrus/pkg/source"

// If is a statement that evaluates some set of statements if a condition is
// true and potentially a different set of statements if that condition is
// false.
type If struct {
	// Condition is the expression that defines the condition to check.
	Condition Expression
	// Consequence is the list of statements that should be evaluated if the
	// condition is true.
	Consequence []FunctionStatement
	// Alternative is the list of statements that should be evaluated if the
	// condition is false.
	Alternative []FunctionStatement
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (i *If) Range() source.Range {
	return i.SourceRange
}

func (*If) functionStatement() {}

var _ FunctionStatement = (*If)(nil)
