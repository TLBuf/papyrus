package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Identifier represents an arbitrary identifier.
type Identifier struct {
	// Text is the normalized text of the identifier.
	Text string
	// Location is the source range of the node.
	Location source.Range
}

// Range returns the source range of the node.
func (i *Identifier) Range() source.Range {
	return i.Location
}

func (*Identifier) expression() {}

func (*Identifier) reference() {}

var _ Reference = (*Identifier)(nil)
