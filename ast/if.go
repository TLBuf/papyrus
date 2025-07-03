package ast

import "github.com/TLBuf/papyrus/source"

// If is a statement that evaluates some set of statements if a condition is
// true and potentially a different set of statements if that condition is
// false.
type If struct {
	LineTrivia
	// StartKeywordLocation is the location of the If keyword that starts the
	// statement.
	StartKeywordLocation source.Location
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
	// EndKeywordLocation is the location of the EndIf keyword that ends the
	// statement.
	EndKeywordLocation source.Location
}

// Body returns the nodes that comprise the body of this block.
func (i *If) Body() []FunctionStatement {
	return i.Statements
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (i *If) Trivia() LineTrivia {
	return i.LineTrivia
}

// Accept calls the appropriate visitor method for the node.
func (i *If) Accept(v Visitor) error {
	return v.VisitIf(i)
}

// Location returns the source location of the node.
func (i *If) Location() source.Location {
	return source.Span(i.StartKeywordLocation, i.EndKeywordLocation)
}

func (*If) block() {}

func (*If) statement() {}

func (*If) functionStatement() {}

var _ FunctionStatement = (*If)(nil)

// ElseIf is a list of statements that may be executed if a condition is true
// and all previous conditions evaluate to false.
type ElseIf struct {
	LineTrivia
	// KeywordLocation is the location of the ElseIf keyword that starts the
	// block.
	KeywordLocation source.Location
	// Condition is the expression that defines the condition to check.
	Condition Expression
	// Statements is the list of statements that should be evaluated if the
	// condition is true.
	Statements []FunctionStatement
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (e *ElseIf) Trivia() LineTrivia {
	return e.LineTrivia
}

// Body returns the nodes that comprise the body of this block.
func (e *ElseIf) Body() []FunctionStatement {
	return e.Statements
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (e *ElseIf) Accept(v Visitor) error {
	return v.VisitElseIf(e)
}

// Location returns the source location of the node.
func (e *ElseIf) Location() source.Location {
	return source.Span(e.KeywordLocation, e.Statements[len(e.Statements)-1].Location())
}

func (*ElseIf) block() {}

var _ Node = (*ElseIf)(nil)

// Else is a list of statements that may be executed if all previous conditions
// evaluate to false.
type Else struct {
	LineTrivia
	// KeywordLocation is the location of the Else keyword that starts the block.
	KeywordLocation source.Location
	// Statements is the list of statements that should be evaluated.
	Statements []FunctionStatement
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (e *Else) Trivia() LineTrivia {
	return e.LineTrivia
}

// Body returns the nodes that comprise the body of this block.
func (e *Else) Body() []FunctionStatement {
	return e.Statements
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (e *Else) Accept(v Visitor) error {
	return v.VisitElse(e)
}

// Location returns the source location of the node.
func (e *Else) Location() source.Location {
	return source.Span(e.KeywordLocation, e.Statements[len(e.Statements)-1].Location())
}

func (*Else) block() {}

var _ Node = (*Else)(nil)
