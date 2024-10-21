package ast

import "github.com/TLBuf/papyrus/pkg/source"

// AssignmentOperatorKind defines the various types of assignment operations.
type AssignmentOperatorKind int

const (
	// UnknownAssignmentOperatorKind is the default (and invalid) assignment operator.
	UnknownAssignmentOperatorKind AssignmentOperatorKind = iota
	// Assign is the basic assignment opertator, '='.
	//
	// The variable is updated to the value of the expression.
	Assign
	// AssignAdd is the assign with addition operator, '+='.
	//
	// The varable is updated to the value of the variable added to the value of the expression.
	AssignAdd
	// AssignSubtract is the assign with subtraction operator, '-='.
	//
	// The varable is updated to the value of the variable less the value of the expression.
	AssignSubtract
	// AssignMultiply is the assign with multiplication operator, '*='.
	//
	// The varable is updated to the value of the variable multiplied by the value of the expression.
	AssignMultiply
	// AssignDivide is the assign with division operator, '/='.
	//
	// The varable is updated to the value of the variable divided by the value of the expression.
	AssignDivide
	// AssignModulo is the assign with modulus operator, '%='.
	//
	// The varable is updated to the value of the variable modulo the value of the expression.
	AssignModulo
)

func (o AssignmentOperatorKind) String() string {
	name, ok := AssignmentOperatorKindNames[o]
	if ok {
		return name
	}
	return "<unknown>"
}

var AssignmentOperatorKindNames = map[AssignmentOperatorKind]string{
	Assign:         "=",
	AssignAdd:      "+=",
	AssignSubtract: "-=",
	AssignMultiply: "*=",
	AssignDivide:   "/=",
	AssignModulo:   "%=",
}

// AssignmentOperator represents an assignment operator.
type AssignmentOperator struct {
	// Kind is the type of assignment operator.
	Kind AssignmentOperatorKind
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (o *AssignmentOperator) Range() source.Range {
	return o.SourceRange
}

var _ Node = (*AssignmentOperator)(nil)

// Assignment is a statement that assigns a new value to a variable (or property).
type Assignment struct {
	// Assignee is the reference to a variable to assign the value to.
	Assignee Reference
	// Operator defines the operator this assignment uses.
	Operator *AssignmentOperator
	// Value is the expression that defines the value to use in the assignment.
	Value Expression
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (a *Assignment) Range() source.Range {
	return a.SourceRange
}

func (*Assignment) functionStatement() {}

var _ FunctionStatement = (*Assignment)(nil)
