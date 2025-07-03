package ast

import "github.com/TLBuf/papyrus/source"

// Statement is a common interface for all statement nodes.
type Statement interface {
	Node
	// Trivia returns the [LineTrivia] assocaited with this node.
	Trivia() LineTrivia
	statement()
}

// ScriptStatement is a common interface for all script statement nodes.
type ScriptStatement interface {
	Statement
	scriptStatement()
}

// Invokable is a common interface for statements that define invokable entities
// (i.e. functions and events).
type Invokable interface {
	ScriptStatement
	// Parameters returns the list of parameters defined for this invokable.
	Parameters() []*Parameter
	invokable()
}

// FunctionStatement is a common interface for all function (and event)
// statement nodes.
type FunctionStatement interface {
	Statement
	functionStatement()
}

// ExpressionStatement is a special function
// statement that is just an expression.
type ExpressionStatement struct {
	LineTrivia
	// Expression is the expression that makes up the statement.
	Expression Expression
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (s *ExpressionStatement) Trivia() LineTrivia {
	return s.LineTrivia
}

// Accept calls the appropriate visitor method for the node.
func (s *ExpressionStatement) Accept(v Visitor) error {
	return v.VisitExpressionStatement(s)
}

// Location returns the source location of the node.
func (s *ExpressionStatement) Location() source.Location {
	return s.Expression.Location()
}

func (*ExpressionStatement) statement() {}

func (*ExpressionStatement) functionStatement() {}

var _ FunctionStatement = (*ExpressionStatement)(nil)
