package ast

import "github.com/TLBuf/papyrus/source"

// Event defines a Papyrus event.
//
// Events are like functions that are predefined by the engine.
type Event struct {
	LineTrivia
	// Keyword is the Event keyword that starts the definition.
	Keyword *Token
	// Name is the name of the event.
	Name *Identifier
	// Open is the open parenthesis token that starts the parameter list.
	Open *Token
	// ParameterList is the list of parameters this event defines in order.
	ParameterList []*Parameter
	// Close is the close parenthesis token that ends the parameter list.
	Close *Token
	// Native are the Native tokens that define that this is a native event or
	// empty if this event in non-native.
	//
	// If non-empty, [Statements] will be empty and [EndKeyword] will be nil.
	//
	// Errata: This being multiple values is due to the offical Papyrus parser
	// accepting any number of flag tokens. They are all included here for
	// completeness, but only one is required to consider the event native.
	Native []*Token
	// Comment is the optional documentation comment for this event.
	Comment *Documentation
	// Statements is the list of function statements that constitute the body of
	// the event.
	Statements []FunctionStatement
	// EndKeyword is the EndEvent keyword that ends the definition or nil if the
	// event is native (and thus has no body).
	EndKeyword *Token
	// Location is the source range of the node.
	Location source.Location
}

// Parameters returns the list of parameters defined for this invokable.
func (e *Event) Parameters() []*Parameter {
	return e.ParameterList
}

// Body returns the nodes that comprise the body of this block.
func (e *Event) Body() []FunctionStatement {
	return e.Statements
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (e *Event) Trivia() LineTrivia {
	return e.LineTrivia
}

// Accept calls the appropriate visitor method for the node.
func (e *Event) Accept(v Visitor) error {
	return v.VisitEvent(e)
}

// SourceLocation returns the source location of the node.
func (e *Event) SourceLocation() source.Location {
	return e.Location
}

func (*Event) block() {}

func (*Event) statement() {}

func (*Event) scriptStatement() {}

func (*Event) invokable() {}

var _ Invokable = (*Event)(nil)
