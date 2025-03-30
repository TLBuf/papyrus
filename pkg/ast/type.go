package ast

import (
	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/TLBuf/papyrus/pkg/types"
)

// TypeLiteral represents a literal type name in source.
type TypeLiteral struct {
	Trivia
	// Type is the type the literal represents.
	Type types.Type
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (t *TypeLiteral) SourceLocation() source.Location {
	return t.Location
}

var _ Node = (*TypeLiteral)(nil)
