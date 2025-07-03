// Package ast defines the Papyrus AST.
package ast

import (
	"github.com/TLBuf/papyrus/source"
)

// Node is a common interfface for all AST nodes.
type Node interface {
	// Accept calls the appropriate visitor method for the node.
	Accept(Visitor) error
	// Location returns the source location of the node.
	Location() source.Location
}

// Block is a common interface for nodes that form blocks of statements.
type Block interface {
	Node
	// Body returns the nodes that comprise the body of this block.
	Body() []FunctionStatement
	block()
}
