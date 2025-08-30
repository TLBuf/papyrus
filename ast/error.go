package ast

import (
	"github.com/TLBuf/papyrus/issue"
	"github.com/TLBuf/papyrus/source"
)

// ErrorStatement is a statement that failed to parse.
type ErrorStatement struct {
	// Issue is the issue describing the error encountered.
	Issue *issue.Issue
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Parameters implements the [Invokable] interface and always returns nil.
func (*ErrorStatement) Parameters() []*Parameter {
	return nil
}

// LeadingBlankLine implements the [Statement]
// interface and always returns false.
func (*ErrorStatement) LeadingBlankLine() bool {
	return false
}

// Accept calls the appropriate visitor method for the node.
func (e *ErrorStatement) Accept(v Visitor) error {
	return v.VisitErrorStatement(e)
}

// Comments implements the [Statement] interface and always returns nil.
func (*ErrorStatement) Comments() *Comments {
	return nil
}

// Location returns the source location of the node.
func (e *ErrorStatement) Location() source.Location {
	return e.NodeLocation
}

func (*ErrorStatement) statement() {}

func (*ErrorStatement) scriptStatement() {}

func (*ErrorStatement) functionStatement() {}

func (*ErrorStatement) invokable() {}
