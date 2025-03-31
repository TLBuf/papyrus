package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Identifier represents an arbitrary identifier.
type Identifier struct {
	Trivia
	// Text is the Identifier token.
	Text Token
	// Normalized is the normalized text of the identifier.
	Normalized string
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (i *Identifier) SourceLocation() source.Location {
	return i.Location
}

func (*Identifier) expression() {}

var _ Expression = (*Identifier)(nil)
