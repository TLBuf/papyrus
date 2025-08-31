package ast

import (
	"fmt"

	"github.com/TLBuf/papyrus/source"
)

// Call is an expression the calls a function defined elsewhere.
type Call struct {
	// Function is the reference to the function being called.
	Function Expression
	// Arguments is the list of arguments being passed to the function being
	// called.
	Arguments []*Argument
	// OpenLocation is the location of the opening parenthesis that starts the
	// argument list.
	OpenLocation source.Location
	// CloseLocation is the location of the closing parenthesis that starts the
	// argument list.
	CloseLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate visitor method for the node.
func (c *Call) Accept(v Visitor) error {
	return v.VisitCall(c)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (c *Call) Comments() *Comments {
	return c.NodeComments
}

// Location returns the source location of the node.
func (c *Call) Location() source.Location {
	return source.Span(c.Function.Location(), c.CloseLocation)
}

func (c *Call) String() string {
	return fmt.Sprintf("Call%s", c.Location())
}

func (*Call) expression() {}

var _ Expression = (*Call)(nil)
