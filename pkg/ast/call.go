package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Call is an expression the calls a function defined elsewhere.
type Call struct {
	Trivia
	// Reciever is the reference to the function being called.
	Reciever Expression
	// Open is the open parenthesis token.
	Open *Token
	// Arguments is the list of arguments being passed to the function being
	// called.
	Arguments []*Argument
	// Close is the close parenthesis token.
	Close *Token
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (c *Call) Accept(v Visitor) error {
	return v.VisitCall(c)
}

// SourceLocation returns the source location of the node.
func (c *Call) SourceLocation() source.Location {
	return c.Location
}

func (*Call) expression() {}

func (*Call) functionStatement() {}

var _ Expression = (*Call)(nil)
