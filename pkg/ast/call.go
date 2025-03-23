package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Cast is an expression casts a value of some type to another.
type Call struct {
	// Function is the reference to the function being called.
	Function *Reference
	// Arguments is the list of arguments being passed to the function being
	// called.
	Arguments []*Argument
	// Location is the source range of the node.
	Location source.Range
}

// Range returns the source range of the node.
func (c *Call) Range() source.Range {
	return c.Location
}

func (*Call) expression() {}

var _ Expression = (*Call)(nil)
