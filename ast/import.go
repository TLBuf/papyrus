package ast

import (
	"github.com/TLBuf/papyrus/source"
)

// Import represents a script import statement.
//
// These are used to bring identifiers available within one script into the
// scope of another script.
type Import struct {
	LineTrivia
	// Keyword is the Import keyword token that starts the statement.
	Keyword *Token
	// Name is the name of the script being imported.
	Name *Identifier
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (i *Import) Trivia() LineTrivia {
	return i.LineTrivia
}

// Accept calls the appropriate visitor method for the node.
func (i *Import) Accept(v Visitor) error {
	return v.VisitImport(i)
}

// Location returns the source location of the node.
func (i *Import) Location() source.Location {
	return i.NodeLocation
}

func (*Import) statement() {}

func (*Import) scriptStatement() {}

var _ ScriptStatement = (*Import)(nil)
