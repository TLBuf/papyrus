package ast

import "github.com/TLBuf/papyrus/source"

// Parenthetical represents a parenthesized expression.
type Parenthetical struct {
	// Value is the expression within the parentheses.
	Value Expression
	// OpenLocation is the location of the opening parenthesis.
	OpenLocation source.Location
	// CloseLocation is the location of the closing parenthesis.
	CloseLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate visitor method for the node.
func (p *Parenthetical) Accept(v Visitor) error {
	return v.VisitParenthetical(p)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (p *Parenthetical) Comments() *Comments {
	return p.NodeComments
}

// Location returns the source location of the node.
func (p *Parenthetical) Location() source.Location {
	return source.Span(p.OpenLocation, p.CloseLocation)
}

func (*Parenthetical) expression() {}

var _ Expression = (*Parenthetical)(nil)
