package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Property defines a script property.
//
// Properties are like variables but which can be accessed in the editor and
// referenced by the engine.
type Property struct {
	// Name is the name of the event.
	Name *Identifier
	// Type is the type of this property.
	Type *TypeLiteral
	// Parameters is the list of parameters this event defines in order.
	Parameters []Parameter
	// IsHidden defines whether this is a hidden property (i.e. it doesn't appear
	// in the editor).
	IsHidden bool
	// IsConditional defines whether this is a conditional property (i.e. it can
	// referenced in conditions).
	IsConditional bool
	// IsAuto defines whether this property uses the auto syntax (i.e. it has not
	// get or set function definitions).
	IsAuto bool
	// IsReadOnly defines whether this property is marked read-only.
	IsReadOnly bool
	// Comment is the optional documentation comment for this event.
	Comment *DocComment
	// Value is the literal that defines the initial value of the property. This
	// is nil if IsAuto is false.
	Value Literal
	// Get is the get function for this property or nil if undefined or IsAuto is
	// true.
	//
	// If IsAuto is false, either Get or Set (or both) will be non-nil.
	//
	// This function is never global or native, has no parameters, and returns the
	// same type as this property's type.
	Get *Function
	// Set is the set function for this property or nil if undefined or IsAuto is
	// true.
	//
	// If IsAuto is false, either Get or Set (or both) will be non-nil.
	//
	// This function is never global or native, has one parameter that is the same
	// type as this property's type, and returns nothing.
	Set *Function
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (p *Property) Range() source.Range {
	return p.SourceRange
}

func (*Property) scriptStatement() {}

var _ ScriptStatement = (*Property)(nil)
