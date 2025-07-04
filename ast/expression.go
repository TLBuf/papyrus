package ast

// Expression is a common interface for all expression nodes.
type Expression interface {
	Node

	// Trivia returns the [InfixTrivia] associated with this node.
	Trivia() InfixTrivia
	expression()
}

// Literal is a common interface for all expression nodes that describe literal
// values.
type Literal interface {
	Expression

	literal()
}
