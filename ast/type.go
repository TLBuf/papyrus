package ast

import (
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/types"
)

// TypeLiteral represents a literal type name in source.
type TypeLiteral struct {
	InfixTrivia

	// Type is the type the literal represents.
	Type types.Type
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] associated with this node.
func (t *TypeLiteral) Trivia() InfixTrivia {
	return t.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (t *TypeLiteral) Accept(v Visitor) error {
	return v.VisitTypeLiteral(t)
}

// Location returns the source location of the node.
func (t *TypeLiteral) Location() source.Location {
	return t.NodeLocation
}

var _ Node = (*TypeLiteral)(nil)
