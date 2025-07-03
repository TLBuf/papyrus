package ast

import "github.com/TLBuf/papyrus/source"

// While is a statement that evaluates some set of statements repeatedly so long
// as a condition is true.
type While struct {
	LineTrivia
	// StartKeywordLocation is the location of the While keyword that starts
	// the statement.
	StartKeywordLocation source.Location
	// Condition is the expression that defines the condition to check before each
	// iteration.
	Condition Expression
	// Statements is the list of function statements that constitute the body of
	// the while.
	Statements []FunctionStatement
	// EndKeywordLocation is the location of the EndWhile keyword that ends the
	// statement.
	EndKeywordLocation source.Location
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (w *While) Trivia() LineTrivia {
	return w.LineTrivia
}

// Body returns the nodes that comprise the body of this block.
func (w *While) Body() []FunctionStatement {
	return w.Statements
}

// Accept calls the appropriate visitor method for the node.
func (w *While) Accept(v Visitor) error {
	return v.VisitWhile(w)
}

// Location returns the source location of the node.
func (w *While) Location() source.Location {
	return source.Span(w.StartKeywordLocation, w.EndKeywordLocation)
}

func (*While) block() {}

func (*While) statement() {}

func (*While) functionStatement() {}

var _ FunctionStatement = (*While)(nil)
