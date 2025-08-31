package ast

import (
	"fmt"

	"github.com/TLBuf/papyrus/source"
)

// Import represents a script import statement.
//
// These are used to bring identifiers available within one script into the
// scope of another script.
type Import struct {
	// HasLeadingBlankLine is true if this node was preceded by a blank line.
	HasLeadingBlankLine bool
	// Name is the name of the script being imported.
	Name *Identifier
	// KeywordLocation is the location of the Import keyword that starts the
	// statement.
	KeywordLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// LeadingBlankLine returns true if this node was preceded by a blank line.
func (i *Import) LeadingBlankLine() bool {
	return i.HasLeadingBlankLine
}

// Accept calls the appropriate visitor method for the node.
func (i *Import) Accept(v Visitor) error {
	return v.VisitImport(i)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (i *Import) Comments() *Comments {
	return i.NodeComments
}

// Location returns the source location of the node.
func (i *Import) Location() source.Location {
	return source.Span(i.KeywordLocation, i.Name.Location())
}

func (i *Import) String() string {
	return fmt.Sprintf("Import%s", i.Location())
}

func (*Import) statement() {}

func (*Import) scriptStatement() {}

var _ ScriptStatement = (*Import)(nil)
