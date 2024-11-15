// Package ast defines the Papyrus AST.
package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Script represents a single Papyrus script file.
type Script struct {
	// Name is the name of script.
	Name *Identifier
	// Extends is the name of the script this one extends from or nil if this
	// script doesn't extend another.
	Extends *Identifier
	// Comment is the documentation comment for this script.
	Comment *DocComment
	// IsHidden defines whether this is a hidden script (i.e. it doesn't appear in
	// the editor).
	IsHidden bool
	// IsConditional defines whether this is a conditional script (i.e. its
	// properties can referenced in conditions).
	IsConditional bool
	// Statements is the list of statements that constitute the body of the
	// script.
	Statements []ScriptStatement
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (s *Script) Range() source.Range {
	return s.SourceRange
}

var _ Node = (*Script)(nil)
