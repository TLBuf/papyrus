package ast

import "github.com/TLBuf/papyrus/source"

// Function defines a Papyrus function.
type Function struct {
	LineTrivia
	// ReturnType is the type of value this function returns or nil if it doesn't
	// return a value.
	ReturnType *TypeLiteral
	// StartKeywordLocation is the location of the Function keyword that starts
	// the statement.
	StartKeywordLocation source.Location
	// Name is the name of the function.
	Name *Identifier
	// OpenLocation is the location of the opening parenthesis that starts the
	// parameter list.
	OpenLocation source.Location
	// ParameterList is the list of parameters this function defines in order.
	ParameterList []*Parameter
	// CloseLocation is the location of the closing parenthesis that starts the
	// parameter list.
	CloseLocation source.Location
	// GlobalLocations are the locations of the Global keywords that mark this as
	// a global function (i.e. it does not actually run on an object, and has no
	// "Self" variable) or empty if this function in non-global.
	//
	// Errata: This being multiple values is due to the offical Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the function global.
	GlobalLocations []source.Location
	// NativeLocations are the locations of the Native keywords that mark this as
	// a native function or empty if this function in non-native.
	//
	// If non-empty, Statements will be empty and EndKeywordLocation will be
	// invalid.
	//
	// Errata: This being multiple values is due to the offical Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the event native.
	NativeLocations []source.Location
	// Documentation is the documentation comment for this function or nil if
	// there is not one.
	Documentation *Documentation
	// Statements is the list of function statements that constitute the body of
	// the function.
	Statements []FunctionStatement
	// EndKeywordLocation is the location of the EndFunction keyword that ends the
	// statement.
	//
	// This is only valid if NativeLocations is empty.
	EndKeywordLocation source.Location
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Parameters returns the list of parameters defined for this invokable.
func (f *Function) Parameters() []*Parameter {
	return f.ParameterList
}

// Body returns the nodes that comprise the body of this block.
func (f *Function) Body() []FunctionStatement {
	return f.Statements
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (f *Function) Trivia() LineTrivia {
	return f.LineTrivia
}

// Accept calls the appropriate visitor method for the node.
func (f *Function) Accept(v Visitor) error {
	return v.VisitFunction(f)
}

// Location returns the source location of the node.
func (f *Function) Location() source.Location {
	return f.NodeLocation
}

func (*Function) block() {}

func (*Function) statement() {}

func (*Function) scriptStatement() {}

func (*Function) invokable() {}

var _ Invokable = (*Function)(nil)
