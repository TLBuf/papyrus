package ast

import "github.com/TLBuf/papyrus/source"

// Cast is an expression that casts a value of some type to another.
type Cast struct {
	// Value is the expression being cast to a new type.
	Value Expression
	// Type is the type the value is being cast to.
	Type *TypeLiteral
	// AsLocation is the location of the As operator.
	AsLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *InlineComments
}

// Accept calls the appropriate visitor method for the node.
func (c *Cast) Accept(v Visitor) error {
	return v.VisitCast(c)
}

// Comments returns the [InlineComments] associated
// with this node or nil if there are none.
func (c *Cast) Comments() *InlineComments {
	return c.NodeComments
}

// Location returns the source location of the node.
func (c *Cast) Location() source.Location {
	return source.Span(c.Value.Location(), c.Type.Location())
}

func (*Cast) expression() {}

var _ Expression = (*Cast)(nil)
