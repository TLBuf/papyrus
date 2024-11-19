package ast

import (
	"github.com/TLBuf/papyrus/pkg/source"
)

// Import represents a script import statement.
//
// These are used to bring identifiers available within one script into the
// scope of another script.
type Import struct {
	// Name is the name of the script being imported.
	Name *Identifier
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (i *Import) Range() source.Range {
	return i.SourceRange
}

func (*Import) scriptStatement() {}

var _ ScriptStatement = (*Import)(nil)
