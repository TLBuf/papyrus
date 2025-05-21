// Package ast defines the Papyrus AST.
package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Script represents a single Papyrus script file.
type Script struct {
	Trivia
	// Keyword is the ScriptName keyword token that starts the script.
	Keyword *Token
	// Name is the name of script.
	Name *Identifier
	// Extends is the Extends token that indicates this script extends another.
	//
	// If this is non-nil, [Parent] will also be non-nil (and vice versa).
	Extends *Token
	// Parent is the name of the script this one extends from or nil if this
	// script doesn't extend another.
	Parent *Identifier
	// Hidden are the Hidden tokens that define that this script is hidden (i.e.
	// it doesn't appear in the editor) or empty if this script is not hidden.
	//
	// Errata: This being multiple values is due to the offical Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the script hidden.
	Hidden []*Token
	// Conditional are the Conditional tokens that define that this script is
	// conditional (i.e. conditional properties it defines can appear in
	// conditions) or empty if this script is not conditional.
	//
	// Errata: This being multiple values is due to the offical Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the script conditional.
	Conditional []*Token
	// Comment is the documentation comment for this script.
	Comment *DocComment
	// Statements is the list of statements that constitute the body of the
	// script.
	Statements []ScriptStatement
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (s *Script) Accept(v Visitor) error {
	return v.VisitScript(s)
}

// SourceLocation returns the source location of the node.
func (s *Script) SourceLocation() source.Location {
	return s.Location
}

var _ Node = (*Script)(nil)
