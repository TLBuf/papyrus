package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Cast is an expression that casts a value of some type to another.
type Cast struct {
	Trivia
	// Value is the expression being cast to a new type.
	Value Expression
	// Operator is the As operator token.
	Operator Token
	// Type is the type the value is being cast to.
	Type *TypeLiteral
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (c *Cast) SourceLocation() source.Location {
	return c.Location
}

func (*Cast) expression() {}

var _ Expression = (*Cast)(nil)
