package ast

import (
	"fmt"

	"github.com/TLBuf/papyrus/source"
)

// Statement is a common interface for all statement nodes.
type Statement interface {
	Node

	// LeadingBlankLine returns true if this node was preceded by a blank line.
	LeadingBlankLine() bool

	// Comments returns the [Comments] associated
	// with this node or nil if there are none.
	Comments() *Comments

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
	// HasLeadingBlankLine is true if this node was preceded by a blank line.
	HasLeadingBlankLine bool
	// Expression is the expression that makes up the statement.
	Expression Expression
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// LeadingBlankLine returns true if this node was preceded by a blank line.
func (s *ExpressionStatement) LeadingBlankLine() bool {
	return s.HasLeadingBlankLine
}

// Accept calls the appropriate visitor method for the node.
func (s *ExpressionStatement) Accept(v Visitor) error {
	return v.VisitExpressionStatement(s)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (s *ExpressionStatement) Comments() *Comments {
	return s.NodeComments
}

// Location returns the source location of the node.
func (s *ExpressionStatement) Location() source.Location {
	return s.Expression.Location()
}

func (s *ExpressionStatement) String() string {
	return fmt.Sprintf("ExpressionStatement%s", s.Location())
}

func (*ExpressionStatement) statement() {}

func (*ExpressionStatement) functionStatement() {}

var _ FunctionStatement = (*ExpressionStatement)(nil)

// CommentStatement is a special statement (both script
// and function) that is comprised entirely of comments.
type CommentStatement struct {
	// Elements are the comments that make up this statement.
	Elements []Comment
}

// LeadingBlankLine returns true if this node was preceded by a blank line.
func (s *CommentStatement) LeadingBlankLine() bool {
	return s.Elements[0].LeadingBlankLine()
}

// Accept calls the appropriate visitor method for the node.
func (s *CommentStatement) Accept(v Visitor) error {
	return v.VisitCommentStatement(s)
}

// Parameters implements the [Invokable] interface and always returns nil.
func (*CommentStatement) Parameters() []*Parameter {
	return nil
}

// Comments implements the [Statement] interface and always returns nil.
func (*CommentStatement) Comments() *Comments {
	return nil
}

// Location returns the source location of the node.
func (s *CommentStatement) Location() source.Location {
	if len(s.Elements) == 1 {
		return s.Elements[0].Location()
	}
	return source.Span(s.Elements[0].Location(), s.Elements[len(s.Elements)-1].Location())
}

func (s *CommentStatement) String() string {
	return fmt.Sprintf("CommentStatement%s", s.Location())
}

func (*CommentStatement) statement() {}

func (*CommentStatement) invokable() {}

func (*CommentStatement) functionStatement() {}

func (*CommentStatement) scriptStatement() {}

var (
	_ FunctionStatement = (*CommentStatement)(nil)
	_ ScriptStatement   = (*CommentStatement)(nil)
	_ Invokable         = (*CommentStatement)(nil)
)
