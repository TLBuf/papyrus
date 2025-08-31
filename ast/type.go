package ast

import (
	"fmt"

	"github.com/TLBuf/papyrus/source"
)

// TypeLiteral represents a literal type name in source.
type TypeLiteral struct {
	// IsArray is true if this is an array type.
	IsArray bool
	// Name is the scalar type of the literal, either
	// this is a primitive type or the name of an object.
	Name *Identifier
	// BracketLocation is the location of both brackets defining this type as an
	// array type.
	//
	// This is only valid if IsArray is true.
	BracketLocation source.Location
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
	if !t.IsArray {
		return t.Name.Location()
	}
	return source.Span(t.Name.Location(), t.BracketLocation)
}

func (t *TypeLiteral) String() string {
	return fmt.Sprintf("TypeLiteral%s", t.Location())
}

var _ Node = (*TypeLiteral)(nil)
