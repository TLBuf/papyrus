package ast

import (
	"github.com/TLBuf/papyrus/pkg/source"
)

// Import represents a script import statement.
//
// These are used to bring identifiers available within one script into the scope of another script.
type Import struct {
	Name        *Identifier
	SourceRange source.Range
}

// Range returns the source range of the node.
func (i *Import) Range() source.Range {
	return i.SourceRange
}

func (*Import) scriptStatement() {}

var _ ScriptStatement = (*Import)(nil)
