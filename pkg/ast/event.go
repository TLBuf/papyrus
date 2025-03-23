package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Event defines a Papyrus event.
//
// Events are like functions that are predefined by the engine.
type Event struct {
	// Name is the name of the event.
	Name *Identifier
	// Parameters is the list of parameters this event defines in order.
	Parameters []*Parameter
	// IsNative defines whether this is a native event. If true, this event will
	// have no statements.
	IsNative bool
	// Comment is the optional documentation comment for this event.
	Comment *DocComment
	// Statements is the list of function statements that constitute the body of
	// the event.
	Statements []FunctionStatement
	// Location is the source range of the node.
	Location source.Range
}

// Range returns the source range of the node.
func (e *Event) Range() source.Range {
	return e.Location
}

func (*Event) scriptStatement() {}

func (*Event) invokable() {}

var _ Invokable = (*Event)(nil)
