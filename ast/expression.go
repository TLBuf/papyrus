package ast

// Expression is a common interface for all expression nodes.
type Expression interface {
	Node

	// Comments returns the [InlineComments] associated
	// with this node or nil if there are none.
	Comments() *Comments

	expression()
}

// Literal is a common interface for all expression nodes that describe literal
// values.
type Literal interface {
	Expression

	// Text returns the raw text of the literal.
	Text() []byte

	literal()
}
