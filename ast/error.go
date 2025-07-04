package ast

import "github.com/TLBuf/papyrus/source"

// ErrorStatement is a statement that failed to parse.
type ErrorStatement struct {
	LineTrivia

	// ErrorMessage is a human-readable message describing the error encountered.
	ErrorMessage string
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Message returns a human-readable message describing the error encountered.
func (e *ErrorStatement) Message() string {
	return e.ErrorMessage
}

// Parameters implements the [Invokable] interface and always returns nil.
func (*ErrorStatement) Parameters() []*Parameter {
	return nil
}

// Trivia returns the [LineTrivia] associated with this node.
func (e *ErrorStatement) Trivia() LineTrivia {
	return e.LineTrivia
}

// Accept calls the appropriate visitor method for the node.
func (e *ErrorStatement) Accept(v Visitor) error {
	return v.VisitErrorStatement(e)
}

// Location returns the source location of the node.
func (e *ErrorStatement) Location() source.Location {
	return e.NodeLocation
}

func (*ErrorStatement) statement() {}

func (*ErrorStatement) scriptStatement() {}

func (*ErrorStatement) functionStatement() {}

func (*ErrorStatement) invokable() {}
