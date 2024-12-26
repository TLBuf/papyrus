package ast

import "github.com/TLBuf/papyrus/pkg/source"

// ScriptStatement is a common interface for all script statement nodes.
type ErrorScriptStatement struct {
	// Message is a human-readable message describing the error encountered.
	Message string
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (e *ErrorScriptStatement) Range() source.Range {
	return e.SourceRange
}

// Message returns a human-readable message describing the error encountered.
func (e *ErrorScriptStatement) ErrorMessage() string {
	return e.Message
}

func (*ErrorScriptStatement) scriptStatement() {}

// FunctionStatement is a common interface for all function (and event)
// statement nodes.
type ErrorFunctionStatement struct {
	// Message is a human-readable message describing the error encountered.
	Message string
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (e *ErrorFunctionStatement) Range() source.Range {
	return e.SourceRange
}

// Message returns a human-readable message describing the error encountered.
func (e *ErrorFunctionStatement) ErrorMessage() string {
	return e.Message
}

func (*ErrorFunctionStatement) functionStatement() {}

// Expression is a common interface for all expression nodes.
type ErrorExpression struct {
	// Message is a human-readable message describing the error encountered.
	Message string
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (e *ErrorExpression) Range() source.Range {
	return e.SourceRange
}

// Message returns a human-readable message describing the error encountered.
func (e *ErrorExpression) ErrorMessage() string {
	return e.Message
}

func (*ErrorExpression) expression() {}
