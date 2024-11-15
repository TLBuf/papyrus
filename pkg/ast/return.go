package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Return is a statement that terminates a function potentially with a value.
type Return struct {
	// Value is the expression that defines the value to return or nil if there is
	// none (i.e. the function doesn't return a value).
	Value Expression
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (r *Return) Range() source.Range {
	return r.SourceRange
}

func (*Return) functionStatement() {}

var _ FunctionStatement = (*Return)(nil)
