package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Cast is an expression casts a value of some type to another.
type Call struct {
	Trivia
	// Reciever is the reference to the function being called.
	Reciever Expression
	// Arguments is the list of arguments being passed to the function being
	// called.
	Arguments []*Argument
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (c *Call) SourceLocation() source.Location {
	return c.Location
}

func (*Call) expression() {}

var _ Expression = (*Call)(nil)
