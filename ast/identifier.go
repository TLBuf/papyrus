package ast

import (
	"fmt"

	"github.com/TLBuf/papyrus/source"
)

// Identifier represents an arbitrary identifier.
type Identifier struct {
	// Text is the text of the identifier.
	Text string
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate visitor method for the node.
func (i *Identifier) Accept(v Visitor) error {
	return v.VisitIdentifier(i)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (i *Identifier) Comments() *Comments {
	return i.NodeComments
}

// Location returns the source location of the node.
func (i *Identifier) Location() source.Location {
	return i.NodeLocation
}

func (i *Identifier) String() string {
	return fmt.Sprintf("Identifier%s", i.Location())
}

func (*Identifier) expression() {}

var _ Expression = (*Identifier)(nil)
