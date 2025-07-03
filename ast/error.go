package ast

import "github.com/TLBuf/papyrus/source"

// ScriptStatement is a common interface for all script statement nodes.
type ErrorScriptStatement struct {
	Trivia
	// Message is a human-readable message describing the error encountered.
	Message string
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate visitor method for the node.
func (e *ErrorScriptStatement) Accept(v Visitor) error {
	return v.VisitErrorScriptStatement(e)
}

// SourceLocation returns the source location of the node.
func (e *ErrorScriptStatement) SourceLocation() source.Location {
	return e.Location
}

// Message returns a human-readable message describing the error encountered.
func (e *ErrorScriptStatement) ErrorMessage() string {
	return e.Message
}

func (*ErrorScriptStatement) statement() {}

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

// Accept calls the appropriate visitor method for the node.
func (e *ErrorFunctionStatement) Accept(v Visitor) error {
	return v.VisitErrorFunctionStatement(e)
}

// SourceLocation returns the source location of the node.
func (e *ErrorFunctionStatement) SourceLocation() source.Location {
	return e.Location
}

// Message returns a human-readable message describing the error encountered.
func (e *ErrorFunctionStatement) ErrorMessage() string {
	return e.Message
}

func (*ErrorFunctionStatement) statement() {}

func (*ErrorFunctionStatement) functionStatement() {}
