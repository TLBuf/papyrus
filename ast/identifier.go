package ast

import "github.com/TLBuf/papyrus/source"

// Identifier represents an arbitrary identifier.
type Identifier struct {
	InfixTrivia
	// Normalized is the normalized text of the identifier.
	Normalized string
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] assocaited with this node.
func (i *Identifier) Trivia() InfixTrivia {
	return i.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (i *Identifier) Accept(v Visitor) error {
	return v.VisitIdentifier(i)
}

// Location returns the source location of the node.
func (i *Identifier) Location() source.Location {
	return i.NodeLocation
}

func (*Identifier) expression() {}

var _ Expression = (*Identifier)(nil)
