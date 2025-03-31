package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Event defines a Papyrus event.
//
// Events are like functions that are predefined by the engine.
type Event struct {
	Trivia
	// Keyword is the Event keyword that starts the definition.
	Keyword Token
	// Name is the name of the event.
	Name *Identifier
	// Open is the open parenthesis token that starts the parameter list.
	Open Token
	// Parameters is the list of parameters this event defines in order.
	Parameters []*Parameter
	// Close is the close parenthesis token that ends the parameter list.
	Close Token
	// Native is the Native token that defines that this is a native event or nil
	// if this event in non-native.
	//
	// If non-nil, [Statements] will be empty and [End] will be nil.
	Native Token
	// Comment is the optional documentation comment for this event.
	Comment *DocComment
	// Statements is the list of function statements that constitute the body of
	// the event.
	Statements []FunctionStatement
	// EndKeyword is the EndEvent keyword that ends the definition or nil if the
	// event is native (and thus has no body).
	EndKeyword Token
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (e *Event) SourceLocation() source.Location {
	return e.Location
}

func (*Event) scriptStatement() {}

func (*Event) invokable() {}

var _ Invokable = (*Event)(nil)
