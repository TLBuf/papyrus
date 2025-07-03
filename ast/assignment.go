package ast

import "github.com/TLBuf/papyrus/source"

// Assignment is a statement that assigns a new value to a variable (or
// property).
type Assignment struct {
	LineTrivia
	// Assignee is the reference to a variable to assign the value to.
	Assignee Expression
	// Operator defines the operator token this assignment uses.
	Operator *Token
	// Value is the expression that defines the value to use in the assignment.
	Value Expression
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (a *Assignment) Trivia() LineTrivia {
	return a.LineTrivia
}

// Accept calls the appropriate visitor method for the node.
func (a *Assignment) Accept(v Visitor) error {
	return v.VisitAssignment(a)
}

// Location returns the source location of the node.
func (a *Assignment) Location() source.Location {
	return a.NodeLocation
}

func (*Assignment) statement() {}

func (*Assignment) functionStatement() {}

var _ FunctionStatement = (*Assignment)(nil)
