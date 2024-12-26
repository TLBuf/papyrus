// Package ast defines the Papyrus AST.
package ast

import (
	"github.com/TLBuf/papyrus/pkg/source"
)

// Node is a common interfface for all AST nodes.
type Node interface {
	// Range returns the source range of the node.
	Range() source.Range
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
	Node
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

// Reference is a common interface for references to values.
type Reference interface {
	Expression
	reference()
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
