// Package ast defines the Papyrus AST.
package ast

import (
	"github.com/TLBuf/papyrus/source"
)

// Node is a common interfface for all AST nodes.
type Node interface {
	// Accept calls the appropriate visitor method for the node.
	Accept(Visitor) error
	// SourceLocation returns the source location of the node.
	SourceLocation() source.Location
}

// ScriptStatement is a common interface for all script statement nodes.
type ScriptStatement interface {
	Node
	scriptStatement()
}

// FunctionStatement is a common interface for all function (and event)
// statement nodes.
type FunctionStatement interface {
	Node
	functionStatement()
}

// Expression is a common interface for all expression nodes.
type Expression interface {
	FunctionStatement
	expression()
}

// Literal is a common interface for all expression nodes that describe literal
// values.
type Literal interface {
	Expression
	literal()
}

// Invokable is a common interface for statements that define invokable entities
// (i.e. functions and events).
type Invokable interface {
	ScriptStatement
	invokable()
}

// LooseComment is a common interface for loose comments (i.e. non-doc
// comments).
type LooseComment interface {
	Node
	looseComment()
}

// Error is a common interface for error nodes that are produced when invalid
// input is encountered.
type Error interface {
	Node
	// Message returns a human-readable message describing the error encountered.
	ErrorMessage() string
}

// Trivia contains supplemental information that has no semantic meaning, but
// which humans may find useful (i.e. comments).
type Trivia struct {
	// LeadingComments are the loose comments that appear on the lines immediately
	// before the node.
	LeadingComments []LooseComment
	// PrefixComments are the loose comments that appear before the node, but on
	// the same line as the node.
	PrefixComments []LooseComment
	// SuffixComments are the loose comments that appear after the node, but on
	// the same line as the node.
	SuffixComments []LooseComment
	// TrailingComments are the loose comments that appear on the lines immedately
	// after the node, but which are not assocaited with another node.
	TrailingComments []LooseComment
}
