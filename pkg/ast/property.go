package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Property defines a script property.
//
// Properties are like variables but which can be accessed in the editor and
// referenced by the engine.
type Property struct {
	Trivia
	// Type is the type of this property.
	Type *TypeLiteral
	// Keyword is the Property keyword that starts the definition.
	Keyword *Token
	// Name is the name of the property.
	Name *Identifier
	// Operator is the assign operator that defines an initial value or nil if
	// the property uses the type's default value.
	//
	// If this is non-nil, [Value] will be non-nil (and vice versa). If [Auto] and
	// [AutoReadOnly] are nil, this must be nil.
	Operator *Token
	// Value is the literal that defines the initial value of the property.
	//
	// If [Auto] and [AutoReadOnly] are nil, this must be nil.
	Value Literal
	// Auto is the Auto token that defines that this is an auto property or nil
	// if this property is a read-only auto property or full property.
	//
	// If non-nil, [Get], [Set], and [AutoReadOnly] will be nil.
	Auto *Token
	// AutoReadOnly is the AutoReadOnly token that defines that this is a
	// read-only auto property or nil if this property is an auto property or full
	// property.
	//
	// If non-nil, [Get], [Set], and [Auto] will be nil. If non-nil, [Operator]
	// and [Value] must also be non-nil.
	AutoReadOnly *Token
	// Hidden are the Hidden tokens that define that this property is hidden (i.e.
	// it doesn't appear in the editor) or empty if this property is not hidden.
	//
	// Errata: This being multiple values is due to the offical Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the property hidden.
	Hidden []*Token
	// Conditional are the Conditional tokens that define that this property is
	// conditional (i.e. it can appear in conditions) or empty if this property is
	// not conditional.
	//
	// Errata: This being multiple values is due to the offical Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the property
	// conditional.
	Conditional []*Token
	// Comment is the optional documentation comment for this property.
	Comment *DocComment
	// Get is the get function for this property or nil if undefined.
	//
	// If [Auto] and [AutoReadOnly] are nil, either Get or Set (or both) will be
	// non-nil.
	//
	// This function is never global or native, has no parameters, and returns the
	// same type as this property's type.
	Get *Function
	// Set is the set function for this property or nil if undefined.
	//
	// If [Auto] and [AutoReadOnly] are nil, either Get or Set (or both) will be
	// non-nil.
	//
	// This function is never global or native, has one parameter that is the same
	// type as this property's type, and returns nothing.
	Set *Function
	// EndKeyword is the EndProperty keyword that ends the definition or nil if
	// the property is Auto or AutoReadOnly (and thus has no Get or Set
	// functions).
	EndKeyword *Token
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (p *Property) Accept(v Visitor) error {
	return v.VisitProperty(p)
}

// SourceLocation returns the source location of the node.
func (p *Property) SourceLocation() source.Location {
	return p.Location
}

func (*Property) scriptStatement() {}

var _ ScriptStatement = (*Property)(nil)
