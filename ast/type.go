package ast

import (
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/types"
)

// TypeLiteral represents a literal type name in source.
type TypeLiteral struct {
	InfixTrivia
	// Text is the type literal token.
	Text *Token
	// Open is the open bracket token that identifies an array type.
	//
	// When this is non-nil, [Close] must also be non-nil.
	Open *Token
	// Close is the close bracket token that identifies an array type.
	//
	// When this is non-nil, [Open] must also be non-nil.
	Close *Token
	// Type is the type the literal represents.
	Type types.Type
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] assocaited with this node.
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
