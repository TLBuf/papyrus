package ast

import (
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/types"
)

// TypeLiteral represents a literal type name in source.
type TypeLiteral struct {
	// Type is the type the literal represents.
	Type types.Type
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate visitor method for the node.
func (t *TypeLiteral) Accept(v Visitor) error {
	return v.VisitTypeLiteral(t)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (t *TypeLiteral) Comments() *Comments {
	return t.NodeComments
}

// Location returns the source location of the node.
func (t *TypeLiteral) Location() source.Location {
	return t.NodeLocation
}

var _ Node = (*TypeLiteral)(nil)
