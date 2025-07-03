package ast

import "github.com/TLBuf/papyrus/source"

// Call is an expression the calls a function defined elsewhere.
type Call struct {
	InfixTrivia
	// Function is the reference to the function being called.
	Function Expression
	// Open is the open parenthesis token.
	Open *Token
	// Arguments is the list of arguments being passed to the function being
	// called.
	Arguments []*Argument
	// Close is the close parenthesis token.
	Close *Token
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] assocaited with this node.
func (c *Call) Trivia() InfixTrivia {
	return c.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (c *Call) Accept(v Visitor) error {
	return v.VisitCall(c)
}

// Location returns the source location of the node.
func (c *Call) Location() source.Location {
	return c.NodeLocation
}

func (*Call) expression() {}

var _ Expression = (*Call)(nil)
