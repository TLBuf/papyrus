package ast

import "github.com/TLBuf/papyrus/pkg/source"

// While is a statement that evaluates some set of statements repeatedly so long
// as a condition is true.
type While struct {
	Trivia
	// Keyword is the While keyword that starts the statement.
	Keyword Token
	// Condition is the expression that defines the condition to check before each
	// iteration.
	Condition Expression
	// Statements is the list of function statements that constitute the body of
	// the while.
	Statements []FunctionStatement
	// Keyword is the EndWhile keyword that ends the statement.
	EndKeyword Token
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (w *While) SourceLocation() source.Location {
	return w.Location
}

func (*While) functionStatement() {}

var _ FunctionStatement = (*While)(nil)
