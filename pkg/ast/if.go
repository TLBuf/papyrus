package ast

import "github.com/TLBuf/papyrus/pkg/source"

// If is a statement that evaluates some set of statements if a condition is
// true and potentially a different set of statements if that condition is
// false.
type If struct {
	Trivia
	// Keyword is the If keyword that starts the statement.
	Keyword *Token
	// Condition is the expression that defines the first condition to check.
	Condition Expression
	// Statements is the list of statements that should be evaluated if the first
	// condition is true.
	Statements []FunctionStatement
	// ElseIfs are the ordered list of ElseIf blocks (or empty if there are none).
	ElseIfs []*ElseIf
	// Else is the block that should be executed if the first [Condition] and all
	// [ElseIf] conditions evaluate to false or nil if there is no else block.
	Else *Else
	// EndKeyword is the EndIf keyword that ends the statement.
	EndKeyword *Token
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate visitor method for the node.
func (i *If) Accept(v Visitor) error {
	return v.VisitIf(i)
}

// SourceLocation returns the source location of the node.
func (i *If) SourceLocation() source.Location {
	return i.Location
}

func (*If) functionStatement() {}

var _ FunctionStatement = (*If)(nil)

// ElseIf is a list of statements that may be executed if a condition is true
// and all previous conditions evaluate to false.
type ElseIf struct {
	Trivia
	// Keyword is the ElseIf keyword that starts the block.
	Keyword *Token
	// Condition is the expression that defines the condition to check.
	Condition Expression
	// Statements is the list of statements that should be evaluated if the
	// condition is true.
	Statements []FunctionStatement
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (e *ElseIf) Accept(v Visitor) error {
	return v.VisitElseIf(e)
}

// SourceLocation returns the source location of the node.
func (e *ElseIf) SourceLocation() source.Location {
	return e.Location
}

var _ Node = (*ElseIf)(nil)

// Else is a list of statements that may be executed if all previous conditions
// evaluate to false.
type Else struct {
	Trivia
	// Keyword is the Else keyword that starts the block.
	Keyword *Token
	// Statements is the list of statements that should be evaluated.
	Statements []FunctionStatement
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (e *Else) Accept(v Visitor) error {
	return v.VisitElse(e)
}

// SourceLocation returns the source location of the node.
func (e *Else) SourceLocation() source.Location {
	return e.Location
}

var _ Node = (*Else)(nil)
