package ast

import (
	"fmt"

	"github.com/TLBuf/papyrus/source"
)

// While is a statement that evaluates some set of statements repeatedly so long
// as a condition is true.
type While struct {
	// HasLeadingBlankLine is true if this node was preceded by a blank line.
	HasLeadingBlankLine bool
	// Condition is the expression that defines the condition to check before each
	// iteration.
	Condition Expression
	// Statements is the list of function statements that constitute the body of
	// the while.
	Statements []FunctionStatement
	// StartKeywordLocation is the location of the While keyword that starts
	// the statement.
	StartKeywordLocation source.Location
	// EndKeywordLocation is the location of the EndWhile keyword that ends the
	// statement.
	EndKeywordLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Body returns the nodes that comprise the body of this block.
func (w *While) Body() []FunctionStatement {
	return w.Statements
}

// LeadingBlankLine returns true if this node was preceded by a blank line.
func (w *While) LeadingBlankLine() bool {
	return w.HasLeadingBlankLine
}

// Accept calls the appropriate visitor method for the node.
func (w *While) Accept(v Visitor) error {
	return v.VisitWhile(w)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (w *While) Comments() *Comments {
	return w.NodeComments
}

// Location returns the source location of the node.
func (w *While) Location() source.Location {
	return source.Span(w.StartKeywordLocation, w.EndKeywordLocation)
}

func (w *While) String() string {
	return fmt.Sprintf("While%s", w.Location())
}

func (*While) block() {}

func (*While) statement() {}

func (*While) functionStatement() {}

var _ FunctionStatement = (*While)(nil)
