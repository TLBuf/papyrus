package ast

import "github.com/TLBuf/papyrus/source"

// Cast is an expression that casts a value of some type to another.
type Cast struct {
	InfixTrivia
	// Value is the expression being cast to a new type.
	Value Expression
	// AsLocation is the location of the As operator.
	AsLocation source.Location
	// Type is the type the value is being cast to.
	Type *TypeLiteral
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
	return source.Span(c.Value.Location(), c.Type.Location())
}

func (*Cast) expression() {}

var _ Expression = (*Cast)(nil)
