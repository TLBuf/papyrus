// Package ast defines the Papyrus AST.
package ast

import "github.com/TLBuf/papyrus/source"

// Script represents a single Papyrus script file.
type Script struct {
	// File is the source file for this script.
	File *source.File
	// Name is the name of script.
	Name *Identifier
	// Parent is the name of the script this one extends from or nil if this
	// script doesn't extend another.
	Parent *Identifier
	// Documentation is the documentation comment for this script or nil if there
	// is not one.
	Documentation *Documentation
	// Statements is the list of statements that constitute the body of the
	// script.
	Statements []ScriptStatement
	// HiddenLocations are the locations of the Hidden keywords that mark this
	// script as hidden (i.e. it doesn't appear in the editor) or empty if this
	// script is not hidden.
	//
	// Errata: This being multiple values is due to the official Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the script hidden.
	HiddenLocations []source.Location
	// ConditionalLocations are the locations of the Conditional keywords that
	// mark this script as conditional (i.e. conditional properties it defines can
	// appear in conditions) or empty if this script is not conditional.
	//
	// Errata: This being multiple values is due to the official Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the script conditional.
	ConditionalLocations []source.Location
	// KeywordLocation is the location of the ScriptName keyword that starts the
	// script.
	KeywordLocation source.Location
	// ExtendsLocation is the location of the Extends keyword that indicates this
	// script extends another.
	//
	// This is only valid if Parent is not nil.
	ExtendsLocation source.Location
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
	// HeaderComments are the standalone comments
	// that appear before the script header.
	HeaderComments []Comment
}

// Accept calls the appropriate visitor method for the node.
func (s *Script) Accept(v Visitor) error {
	return v.VisitScript(s)
}

// Location returns the source location of the node.
func (s *Script) Location() source.Location {
	return s.NodeLocation
}

var _ Node = (*Script)(nil)
