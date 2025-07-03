package ast

import (
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/token"
)

// AssignmentKind is the type of assignment operation.
type AssignmentKind uint8

const (
	// Assign, '=', assigns the value to the assignee.
	Assign = AssignmentKind(token.Assign)
	// AssignAdd, '+=', adds the assignee to the
	// value and assigns the result to assignee.
	AssignAdd = AssignmentKind(token.AssignAdd)
	// AssignDivide, '/=', divides the assignee by
	// the value and assigns the result to assignee.
	AssignDivide = AssignmentKind(token.AssignDivide)
	// AssignDivide, '%=', assigns the remainder from dividing the
	// assignee by the value using integer division to the assignee.
	AssignModulo = AssignmentKind(token.AssignModulo)
	// AssignDivide, '/=', multiplies the assignee by
	// the value and assigns the result to assignee.
	AssignMultiply = AssignmentKind(token.AssignMultiply)
	// AssignAdd, '+=', subtracts the value from the
	// assignee and assigns the result to assignee.
)

// Symbol returns the string representation of this
// kind or an empty string if it is invalid.
func (k AssignmentKind) Symbol() string {
	switch k {
	case Assign:
		return token.Assign.Symbol()
	case AssignAdd:
		return token.AssignAdd.Symbol()
	case AssignDivide:
		return token.AssignDivide.Symbol()
	case AssignModulo:
		return token.AssignModulo.Symbol()
	case AssignMultiply:
		return token.AssignMultiply.Symbol()
	default:
		return ""
	}
}

// Assignment is a statement that assigns a new value to a variable (or
// property).
type Assignment struct {
	LineTrivia
	// Kind is the kind of assignment this expression represents.
	Kind AssignmentKind
	// Assignee is the reference to a variable to assign the value to.
	Assignee Expression
	// OperatorLocation is the location of the assignment operator.
	OperatorLocation source.Location
	// Value is the expression that defines the value to use in the assignment.
	Value Expression
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
	return source.Span(a.Assignee.Location(), a.Value.Location())
}

func (*Assignment) statement() {}

func (*Assignment) functionStatement() {}

var _ FunctionStatement = (*Assignment)(nil)
