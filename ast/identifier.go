package ast

import "github.com/TLBuf/papyrus/source"

// Identifier represents an arbitrary identifier.
type Identifier struct {
	// Normalized is the normalized text of the identifier.
	Normalized string
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *InlineComments
}

// Accept calls the appropriate visitor method for the node.
func (i *Identifier) Accept(v Visitor) error {
	return v.VisitIdentifier(i)
}

// Comments returns the [InlineComments] associated
// with this node or nil if there are none.
func (i *Identifier) Comments() *InlineComments {
	return i.NodeComments
}

// Location returns the source location of the node.
func (i *Identifier) Location() source.Location {
	return i.NodeLocation
}

func (*Identifier) expression() {}

var _ Expression = (*Identifier)(nil)
