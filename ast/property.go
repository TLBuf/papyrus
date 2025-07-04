package ast

import (
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/token"
)

// PropertyKind is the type of property.
type PropertyKind uint8

const (
	// Full is a full property definition that defines either a get or set
	// function or both.
	Full = PropertyKind(0)
	// Auto is a mutable property that has implicitly defined get and set
	// functions.
	Auto = PropertyKind(token.Auto)
	// AutoReadOnly is an immutable property that has an implicit get function.
	AutoReadOnly = PropertyKind(token.AutoReadOnly)
)

// Property defines a script property.
//
// Properties are like variables but which can be accessed in the editor and
// referenced by the engine.
type Property struct {
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
	// Kind is the kind of property this statement represents.
	Kind PropertyKind
	// Type is the type of this property.
	Type *TypeLiteral
	// Name is the name of the property.
	Name *Identifier
	// Value is the literal that defines the initial value of the property.
	//
	// If Kind is [Full] this is nil.
	Value Literal
	// Documentation is the documentation comment for this property or nil if
	// there is not one.
	Documentation *Documentation
	// Get is the get function for this property or nil if undefined.
	//
	// If Kind is [Full] are nil, either Get or Set (or both) will be non-nil.
	//
	// This function is never global or native, has no parameters, and returns the
	// same type as this property's type.
	Get *Function
	// Set is the set function for this property or nil if undefined.
	//
	// If Kind is [Full] are nil, either Get or Set (or both) will be non-nil.
	//
	// This function is never global or native, has one parameter that is the same
	// type as this property's type, and returns nothing.
	Set *Function
	// HiddenLocations are the locations of the Hidden keywords that mark this
	// property as hidden (i.e. it doesn't appear in the editor) or empty if this
	// property is not hidden.
	//
	// Errata: This being multiple values is due to the official Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the property hidden.
	HiddenLocations []source.Location
	// ConditionalLocations are the locations of the Conditional keywords that
	// mark this property as conditional (i.e. it can appear in conditions) or
	// empty if this property is not conditional.
	//
	// Errata: This being multiple values is due to the official Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the property
	// conditional.
	ConditionalLocations []source.Location
	// StartKeywordLocation is the location of the Property keyword that starts
	// the statement.
	StartKeywordLocation source.Location
	// OperatorLocation is the location of the assignment operator that defines an
	// initial value.
	//
	// This is only valid if Value is not nil.
	OperatorLocation source.Location
	// AutoLocation is the Auto/AutoReadOnly keyword that defines that this is an
	// auto (or read-only) property
	//
	// This is only valid if Kind is not [Auto] or [AutoReadOnly].
	AutoLocation source.Location
	// EndKeywordLocation is the location of the EndProperty keyword that ends the
	// statement.
	//
	// This is only valid if Kind is [Full].
	EndKeywordLocation source.Location
	// NodeComments are the comments on lines before and/or after a
	// node or nil if the node has no comments associated with it.
	NodeComments *CrosslineComments
}

// PrecedingBlankLine returns true if this node was preceded by a blank line.
func (p *Property) PrecedingBlankLine() bool {
	return p.HasPrecedingBlankLine
}

// Accept calls the appropriate visitor method for the node.
func (p *Property) Accept(v Visitor) error {
	return v.VisitProperty(p)
}

// Comments returns the [CrosslineComments] associated
// with this node or nil if there are none.
func (p *Property) Comments() *CrosslineComments {
	return p.NodeComments
}

// Location returns the source location of the node.
func (p *Property) Location() source.Location {
	if p.Kind == Full {
		return source.Span(p.Type.Location(), p.EndKeywordLocation)
	}
	end := p.AutoLocation
	if len(p.HiddenLocations) > 0 {
		end = p.HiddenLocations[len(p.HiddenLocations)-1]
	}
	if len(p.ConditionalLocations) > 0 {
		last := p.ConditionalLocations[len(p.ConditionalLocations)-1]
		if last.ByteOffset > end.ByteOffset {
			end = last
		}
	}
	return source.Span(p.Type.Location(), end)
}

func (*Property) statement() {}

func (*Property) scriptStatement() {}

var _ ScriptStatement = (*Property)(nil)
