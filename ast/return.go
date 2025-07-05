package ast

import "github.com/TLBuf/papyrus/source"

// Return is a statement that terminates a function potentially with a value.
type Return struct {
	// HasLeadingBlankLine is true if this node was preceded by a blank line.
	HasLeadingBlankLine bool
	// Value is the expression that defines the value to return or nil if there is
	// none (i.e. the function doesn't return a value).
	Value Expression
	// KeywordLocation is the location of the Return keyword that starts the
	// statement.
	KeywordLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// LeadingBlankLine returns true if this node was preceded by a blank line.
func (r *Return) LeadingBlankLine() bool {
	return r.HasLeadingBlankLine
}

// Accept calls the appropriate visitor method for the node.
func (r *Return) Accept(v Visitor) error {
	return v.VisitReturn(r)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (r *Return) Comments() *Comments {
	return r.NodeComments
}

// Location returns the source location of the node.
func (r *Return) Location() source.Location {
	if r.Value == nil {
		return r.KeywordLocation
	}
	return source.Span(r.KeywordLocation, r.Value.Location())
}

func (*Return) statement() {}

func (*Return) functionStatement() {}

var _ FunctionStatement = (*Return)(nil)
