package ast

import (
	"github.com/TLBuf/papyrus/pkg/source"
)

// Import represents a script import statement.
//
// These are used to bring identifiers available within one script into the
// scope of another script.
type Import struct {
	Trivia
	// Name is the name of the script being imported.
	Name *Identifier
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (i *Import) SourceLocation() source.Location {
	return i.Location
}

func (*Import) scriptStatement() {}

var _ ScriptStatement = (*Import)(nil)
