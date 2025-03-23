package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Parameter is a named and typed parameter to an invokable.
type Parameter struct {
	// Type is the type literal that defines the type of the parameter.
	Type *TypeLiteral
	// Name is the name of the parameter.
	Name *Identifier
	// Value is the optional default value of the parameter.
	Value Literal
	// Location is the source range of the node.
	Location source.Range
}

// Range returns the source range of the node.
func (v *Parameter) Range() source.Range {
	return v.Location
}

var _ Node = (*Parameter)(nil)
