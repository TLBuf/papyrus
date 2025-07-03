package ast

import "github.com/TLBuf/papyrus/source"

// Assignment is a statement that assigns a new value to a variable (or
// property).
type Assignment struct {
	Trivia
	// Assignee is the reference to a variable to assign the value to.
	Assignee Expression
	// Operator defines the operator token this assignment uses.
	Operator *Token
	// Value is the expression that defines the value to use in the assignment.
	Value Expression
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate visitor method for the node.
func (a *Assignment) Accept(v Visitor) error {
	return v.VisitAssignment(a)
}

// SourceLocation returns the source location of the node.
func (a *Assignment) SourceLocation() source.Location {
	return a.Location
}

func (*Assignment) statement() {}

func (*Assignment) functionStatement() {}

var _ FunctionStatement = (*Assignment)(nil)
