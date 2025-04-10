package ast

import "github.com/TLBuf/papyrus/pkg/source"

// ScriptStatement is a common interface for all script statement nodes.
type ErrorScriptStatement struct {
	Trivia
	// Message is a human-readable message describing the error encountered.
	Message string
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (e *ErrorScriptStatement) SourceLocation() source.Location {
	return e.Location
}

// Message returns a human-readable message describing the error encountered.
func (e *ErrorScriptStatement) ErrorMessage() string {
	return e.Message
}

func (*ErrorScriptStatement) scriptStatement() {}

func (*ErrorScriptStatement) invokable() {}

// FunctionStatement is a common interface for all function (and event)
// statement nodes.
type ErrorFunctionStatement struct {
	Trivia
	// Message is a human-readable message describing the error encountered.
	Message string
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (e *ErrorFunctionStatement) SourceLocation() source.Location {
	return e.Location
}

// Message returns a human-readable message describing the error encountered.
func (e *ErrorFunctionStatement) ErrorMessage() string {
	return e.Message
}

func (*ErrorFunctionStatement) functionStatement() {}

// Expression is a common interface for all expression nodes.
type ErrorExpression struct {
	Trivia
	// Message is a human-readable message describing the error encountered.
	Message string
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (e *ErrorExpression) SourceLocation() source.Location {
	return e.Location
}

// Message returns a human-readable message describing the error encountered.
func (e *ErrorExpression) ErrorMessage() string {
	return e.Message
}

func (*ErrorExpression) expression() {}
