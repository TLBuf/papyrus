package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Parameter is a named and typed parameter to an invokable.
type Parameter struct {
	Trivia
	// Type is the type literal that defines the type of the parameter.
	Type *TypeLiteral
	// Name is the name of the parameter.
	Name *Identifier
	// Operator is the assignment operator or nil if no default value is defined.
	Operator *Token
	// Value is the optional default value of the parameter or nil if no default
	// value is defined.
	Value Literal
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (p *Parameter) Accept(v Visitor) error {
	return v.VisitParameter(p)
}

// SourceLocation returns the source location of the node.
func (v *Parameter) SourceLocation() source.Location {
	return v.Location
}

var _ Node = (*Parameter)(nil)
