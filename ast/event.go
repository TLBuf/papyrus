package ast

import "github.com/TLBuf/papyrus/source"

// Event defines a Papyrus event.
//
// Events are like functions that are predefined by the engine.
type Event struct {
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
	// Name is the name of the event.
	Name *Identifier
	// ParameterList is the list of parameters this event defines in order.
	ParameterList []*Parameter
	// Documentation is the documentation comment for this event or nil if there
	// is not one.
	Documentation *Documentation
	// Statements is the list of function statements that constitute the body of
	// the event.
	Statements []FunctionStatement
	// NativeLocations are the locations of the Native keywords that mark this as
	// a native event or empty if this event in non-native.
	//
	// If non-empty, Statements will be empty and EndKeywordLocation will be
	// invalid.
	//
	// Errata: This being multiple values is due to the official Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the event native.
	NativeLocations []source.Location
	// StartKeywordLocation is the location of the Event keyword that starts the
	// statement.
	StartKeywordLocation source.Location
	// OpenLocation is the location of the opening parenthesis that starts the
	// parameter list.
	OpenLocation source.Location
	// CloseLocation is the location of the closing parenthesis that starts the
	// parameter list.
	CloseLocation source.Location
	// EndKeywordLocation is the location of the EndEvent keyword that ends the
	// statement.
	//
	// This is only valid if NativeLocations is empty.
	EndKeywordLocation source.Location
	// NodeComments are the comments on lines before and/or after a
	// node or nil if the node has no comments associated with it.
	NodeComments *CrosslineComments
}

// Parameters returns the list of parameters defined for this invokable.
func (e *Event) Parameters() []*Parameter {
	return e.ParameterList
}

// Body returns the nodes that comprise the body of this block.
func (e *Event) Body() []FunctionStatement {
	return e.Statements
}

// PrecedingBlankLine returns true if this node was preceded by a blank line.
func (e *Event) PrecedingBlankLine() bool {
	return e.HasPrecedingBlankLine
}

// Accept calls the appropriate visitor method for the node.
func (e *Event) Accept(v Visitor) error {
	return v.VisitEvent(e)
}

// Comments returns the [CrosslineComments] associated
// with this node or nil if there are none.
func (e *Event) Comments() *CrosslineComments {
	return e.NodeComments
}

// Location returns the source location of the node.
func (e *Event) Location() source.Location {
	if len(e.NativeLocations) == 0 {
		return source.Span(e.StartKeywordLocation, e.EndKeywordLocation)
	}
	end := e.NativeLocations[len(e.NativeLocations)-1]
	if e.Documentation != nil {
		end = e.Documentation.Location()
	}
	return source.Span(e.StartKeywordLocation, end)
}

func (*Event) block() {}

func (*Event) statement() {}

func (*Event) scriptStatement() {}

func (*Event) invokable() {}

var _ Invokable = (*Event)(nil)
