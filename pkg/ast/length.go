package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Length represents a length access of an array value.
type Length struct {
	// Value represents the array value having its length taken.
	Value Expression
	// Operator is the access operator for this length expression.
	AccessOperator *AccessOperator
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (l *Length) Range() source.Range {
	return l.SourceRange
}

func (*Length) expression() {}

var _ Expression = (*Length)(nil)
