package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Function defines a Papyrus function.
type Function struct {
	Trivia
	// Name is the name of the function.
	Name *Identifier
	// ReturnType is the type of value this function returns or nil if it doesn't
	// return a value.
	ReturnType *TypeLiteral
	// Parameters is the list of parameters this function defines in order.
	Parameters []*Parameter
	// IsGlobal defines whether this function is considered global (i.e. it does
	// not actually run on an object, and has no "Self" variable).
	IsGlobal bool
	// IsNative defines whether this is a native function. If true, this function
	// will have no statements.
	IsNative bool
	// Comment is the optional documentation comment for this function.
	Comment *DocComment
	// Statements is the list of function statements that constitute the body of
	// the function.
	Statements []FunctionStatement
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (f *Function) SourceLocation() source.Location {
	return f.Location
}

func (*Function) scriptStatement() {}

func (*Function) invokable() {}

var _ Invokable = (*Function)(nil)
