package ast

import "github.com/TLBuf/papyrus/source"

// Statement is a common interface for all statement nodes.
type Statement interface {
	Node

	// PrecedingBlankLine returns true if this node was preceded by a blank line.
	PrecedingBlankLine() bool

	// Comments returns the [CrosslineComments] associated
	// with this node or nil if there are none.
	Comments() *CrosslineComments

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
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
	// Expression is the expression that makes up the statement.
	Expression Expression
	// NodeComments are the comments on lines before and/or after a
	// node or nil if the node has no comments associated with it.
	NodeComments *CrosslineComments
}

// PrecedingBlankLine returns true if this node was preceded by a blank line.
func (s *ExpressionStatement) PrecedingBlankLine() bool {
	return s.HasPrecedingBlankLine
}

// Accept calls the appropriate visitor method for the node.
func (s *ExpressionStatement) Accept(v Visitor) error {
	return v.VisitExpressionStatement(s)
}

// Comments returns the [CrosslineComments] associated
// with this node or nil if there are none.
func (s *ExpressionStatement) Comments() *CrosslineComments {
	return s.NodeComments
}

// Location returns the source location of the node.
func (s *ExpressionStatement) Location() source.Location {
	return s.Expression.Location()
}

func (*ExpressionStatement) statement() {}

func (*ExpressionStatement) functionStatement() {}

var _ FunctionStatement = (*ExpressionStatement)(nil)
