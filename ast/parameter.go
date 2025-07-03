package ast

import "github.com/TLBuf/papyrus/source"

// Parameter is a named and typed parameter to an invokable.
type Parameter struct {
	InfixTrivia
	// Type is the type literal that defines the type of the parameter.
	Type *TypeLiteral
	// Name is the name of the parameter.
	Name *Identifier
	// Operator is the assignment operator or nil if no default value is defined.
	Operator *Token
	// Value is the optional default value of the parameter or nil if no default
	// value is defined.
	Value Literal
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Accept calls the appropriate visitor method for the node.
func (p *Parameter) Accept(v Visitor) error {
	return v.VisitParameter(p)
}

// Location returns the source location of the node.
func (p *Parameter) Location() source.Location {
	return p.NodeLocation
}

var _ Node = (*Parameter)(nil)
