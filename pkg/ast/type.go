package ast

import (
	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/TLBuf/papyrus/pkg/types"
)

// TypeLiteral represents a literal type name in source.
type TypeLiteral struct {
	// Type is the type the literal represents.
	Type types.Type
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (t *TypeLiteral) Range() source.Range {
	return t.SourceRange
}

var _ Node = (*TypeLiteral)(nil)
