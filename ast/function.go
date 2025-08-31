package ast

import "github.com/TLBuf/papyrus/source"

// Function defines a Papyrus function.
type Function struct {
	// HasLeadingBlankLine is true if this node was preceded by a blank line.
	HasLeadingBlankLine bool
	// ReturnType is the type of value this function returns or nil if it doesn't
	// return a value.
	ReturnType *TypeLiteral
	// Name is the name of the function.
	Name *Identifier
	// ParameterList is the list of parameters this function defines in order.
	ParameterList []*Parameter
	// Documentation is the documentation comment for this function or nil if
	// there is not one.
	Documentation *Documentation
	// Statements is the list of function statements that constitute the body of
	// the function.
	Statements []FunctionStatement
	// GlobalLocations are the locations of the Global keywords that mark this as
	// a global function (i.e. it does not actually run on an object, and has no
	// "Self" variable) or empty if this function in non-global.
	//
	// Errata: This being multiple values is due to the official Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the function global.
	GlobalLocations []source.Location
	// NativeLocations are the locations of the Native keywords that mark this as
	// a native function or empty if this function in non-native.
	//
	// If non-empty, Statements will be empty and EndKeywordLocation will be
	// invalid.
	//
	// Errata: This being multiple values is due to the official Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the event native.
	NativeLocations []source.Location
	// StartKeywordLocation is the location of the Function keyword that starts
	// the statement.
	StartKeywordLocation source.Location
	// OpenLocation is the location of the opening parenthesis that starts the
	// parameter list.
	OpenLocation source.Location
	// CloseLocation is the location of the closing parenthesis that starts the
	// parameter list.
	CloseLocation source.Location
	// EndKeywordLocation is the location of the EndFunction keyword that ends the
	// statement.
	//
	// This is only valid if NativeLocations is empty.
	EndKeywordLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Parameters returns the list of parameters defined for this invokable.
func (f *Function) Parameters() []*Parameter {
	return f.ParameterList
}

// Body returns the nodes that comprise the body of this block.
func (f *Function) Body() []FunctionStatement {
	return f.Statements
}

// LeadingBlankLine returns true if this node was preceded by a blank line.
func (f *Function) LeadingBlankLine() bool {
	return f.HasLeadingBlankLine
}

// Accept calls the appropriate visitor method for the node.
func (f *Function) Accept(v Visitor) error {
	return v.VisitFunction(f)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (f *Function) Comments() *Comments {
	return f.NodeComments
}

// Location returns the source location of the node.
func (f *Function) Location() source.Location {
	start := f.StartKeywordLocation
	if f.ReturnType != nil {
		start = f.ReturnType.Location()
	}
	if f.EndKeywordLocation != source.UnknownLocation {
		return source.Span(start, f.EndKeywordLocation)
	}
	if f.Documentation != nil {
		return source.Span(start, f.Documentation.Location())
	}
	end := f.CloseLocation
	if len(f.NativeLocations) > 0 {
		end = f.NativeLocations[len(f.NativeLocations)-1]
	}
	if len(f.GlobalLocations) > 0 {
		last := f.GlobalLocations[len(f.GlobalLocations)-1]
		if last.Start() > end.Start() {
			end = last
		}
	}
	return source.Span(start, end)
}

func (*Function) block() {}

func (*Function) statement() {}

func (*Function) scriptStatement() {}

func (*Function) invokable() {}

var _ Invokable = (*Function)(nil)
