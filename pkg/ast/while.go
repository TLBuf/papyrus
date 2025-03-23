package ast

import "github.com/TLBuf/papyrus/pkg/source"

// While is a statement that evaluates some set of statements repeatedly so long
// as a condition is true.
type While struct {
	// Condition is the expression that defines the condition to check before each
	// iteration.
	Condition Expression
	// Statements is the list of function statements that constitute the body of
	// the while.
	Statements []FunctionStatement
	// Location is the source range of the node.
	Location source.Range
}

// Range returns the source range of the node.
func (w *While) Range() source.Range {
	return w.Location
}

func (*While) functionStatement() {}

var _ FunctionStatement = (*While)(nil)
