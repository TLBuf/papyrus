package ast

import "github.com/TLBuf/papyrus/source"

// Parameter is a named and typed parameter to an invokable.
type Parameter struct {
	// Type is the type literal that defines the type of the parameter.
	Type *TypeLiteral
	// Name is the name of the parameter.
	Name *Identifier
	// DefaultValue is the optional default value of the parameter or nil if no
	// default value is defined.
	DefaultValue Literal
	// OperatorLocation is the location of the assignment operator.
	//
	// This is only valid if DefaultValue is not nil.
	OperatorLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *InlineComments
}

// Accept calls the appropriate visitor method for the node.
func (p *Parameter) Accept(v Visitor) error {
	return v.VisitParameter(p)
}

// Comments returns the [InlineComments] associated
// with this node or nil if there are none.
func (p *Parameter) Comments() *InlineComments {
	return p.NodeComments
}

// Location returns the source location of the node.
func (p *Parameter) Location() source.Location {
	if p.DefaultValue == nil {
		return source.Span(p.Type.Location(), p.Name.Location())
	}
	return source.Span(p.Type.Location(), p.DefaultValue.Location())
}

var _ Node = (*Parameter)(nil)
