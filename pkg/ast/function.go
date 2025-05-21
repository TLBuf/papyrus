package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Function defines a Papyrus function.
type Function struct {
	Trivia
	// ReturnType is the type of value this function returns or nil if it doesn't
	// return a value.
	ReturnType *TypeLiteral
	// Keyword is the Function keyword that starts the definition.
	Keyword *Token
	// Name is the name of the function.
	Name *Identifier
	// Open is the open parenthesis token that starts the parameter list.
	Open *Token
	// Parameters is the list of parameters this function defines in order.
	Parameters []*Parameter
	// Close is the close parenthesis token that ends the parameter list.
	Close *Token
	// IsGlobal defines whether this function is considered global (i.e. it does
	// not actually run on an object, and has no "Self" variable).

	// Global are the Global tokens that define that this is a global function
	// (i.e. it does not actually run on an object, and has no "Self" variable) or
	// empty if this function in non-global.
	//
	// Errata: This being multiple values is due to the offical Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the function global.
	Global []*Token
	// Native are the Native tokens that define that this is a native function or
	// empty if this function in non-native.
	//
	// If non-empty, [Statements] will be empty and [EndKeyword] will be nil.
	//
	// Errata: This being multiple values is due to the offical Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the function native.
	Native []*Token
	// Comment is the optional documentation comment for this function.
	Comment *DocComment
	// Statements is the list of function statements that constitute the body of
	// the function.
	Statements []FunctionStatement
	// EndKeyword is the EndFunction keyword that ends the definition or nil if
	// the function is native (and thus has no body).
	EndKeyword *Token
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (f *Function) Accept(v Visitor) error {
	return v.VisitFunction(f)
}

// SourceLocation returns the source location of the node.
func (f *Function) SourceLocation() source.Location {
	return f.Location
}

func (*Function) scriptStatement() {}

func (*Function) invokable() {}

var _ Invokable = (*Function)(nil)
