package ast

import "github.com/TLBuf/papyrus/source"

// ErrorStatement is a statement that failed to parse.
type ErrorStatement struct {
	LineTrivia
	// Message is a human-readable message describing the error encountered.
	Message string
	// Location is the source range of the node.
	Location source.Location
}

// Message returns a human-readable message describing the error encountered.
func (e *ErrorStatement) ErrorMessage() string {
	return e.Message
}

// Parameters implements the [Invokable] interface and always returns nil.
func (*ErrorStatement) Parameters() []*Parameter {
	return nil
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (e *ErrorStatement) Trivia() LineTrivia {
	return e.LineTrivia
}

// Accept calls the appropriate visitor method for the node.
func (e *ErrorStatement) Accept(v Visitor) error {
	return v.VisitErrorStatement(e)
}

// SourceLocation returns the source location of the node.
func (e *ErrorStatement) SourceLocation() source.Location {
	return e.Location
}

func (*ErrorStatement) statement() {}

func (*ErrorStatement) scriptStatement() {}

func (*ErrorStatement) functionStatement() {}

func (*ErrorStatement) invokable() {}
