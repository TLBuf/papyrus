package ast

import "github.com/TLBuf/papyrus/source"

// Access is an expression that reference a value or function that belongs to
// some scope.
type Access struct {
	// Value is the expression that defines the value have something accessed.
	Value Expression
	// Name is the name of the variable or function being accessed in value.
	Name *Identifier
	// DotLocation is the location of the dot operator.
	DotLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate visitor method for the node.
func (a *Access) Accept(v Visitor) error {
	return v.VisitAccess(a)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (a *Access) Comments() *Comments {
	return a.NodeComments
}

// Location returns the source location of the node.
func (a *Access) Location() source.Location {
	return source.Span(a.Value.Location(), a.Name.Location())
}

func (*Access) expression() {}

var _ Expression = (*Access)(nil)
