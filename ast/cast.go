package ast

import "github.com/TLBuf/papyrus/source"

// Cast is an expression that casts a value of some type to another.
type Cast struct {
	InfixTrivia
	// Value is the expression being cast to a new type.
	Value Expression
	// Operator is the As operator token.
	Operator *Token
	// Type is the type the value is being cast to.
	Type *TypeLiteral
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] assocaited with this node.
func (c *Cast) Trivia() InfixTrivia {
	return c.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (c *Cast) Accept(v Visitor) error {
	return v.VisitCast(c)
}

// Location returns the source location of the node.
func (c *Cast) Location() source.Location {
	return c.NodeLocation
}

func (*Cast) expression() {}

var _ Expression = (*Cast)(nil)
