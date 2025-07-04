package ast

import "github.com/TLBuf/papyrus/source"

// If is a statement that evaluates some set of statements if a condition is
// true and potentially a different set of statements if that condition is
// false.
type If struct {
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
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
	// StartKeywordLocation is the location of the If keyword that starts the
	// statement.
	StartKeywordLocation source.Location
	// EndKeywordLocation is the location of the EndIf keyword that ends the
	// statement.
	EndKeywordLocation source.Location
	// NodeComments are the comments on lines before and/or after a
	// node or nil if the node has no comments associated with it.
	NodeComments *CrosslineComments
}

// Body returns the nodes that comprise the body of this block.
func (i *If) Body() []FunctionStatement {
	return i.Statements
}

// PrecedingBlankLine returns true if this node was preceded by a blank line.
func (i *If) PrecedingBlankLine() bool {
	return i.HasPrecedingBlankLine
}

// Accept calls the appropriate visitor method for the node.
func (i *If) Accept(v Visitor) error {
	return v.VisitIf(i)
}

// Comments returns the [CrosslineComments] associated
// with this node or nil if there are none.
func (i *If) Comments() *CrosslineComments {
	return i.NodeComments
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
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
	// Condition is the expression that defines the condition to check.
	Condition Expression
	// Statements is the list of statements that should be evaluated if the
	// condition is true.
	Statements []FunctionStatement
	// KeywordLocation is the location of the ElseIf keyword that starts the
	// block.
	KeywordLocation source.Location
	// NodeComments are the comments on lines before and/or after a
	// node or nil if the node has no comments associated with it.
	NodeComments *CrosslineComments
}

// Body returns the nodes that comprise the body of this block.
func (e *ElseIf) Body() []FunctionStatement {
	return e.Statements
}

// PrecedingBlankLine returns true if this node was preceded by a blank line.
func (e *ElseIf) PrecedingBlankLine() bool {
	return e.HasPrecedingBlankLine
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (e *ElseIf) Accept(v Visitor) error {
	return v.VisitElseIf(e)
}

// Comments returns the [CrosslineComments] associated
// with this node or nil if there are none.
func (e *ElseIf) Comments() *CrosslineComments {
	return e.NodeComments
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
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
	// Statements is the list of statements that should be evaluated.
	Statements []FunctionStatement
	// KeywordLocation is the location of the Else keyword that starts the block.
	KeywordLocation source.Location
	// NodeComments are the comments on lines before and/or after a
	// node or nil if the node has no comments associated with it.
	NodeComments *CrosslineComments
}

// Body returns the nodes that comprise the body of this block.
func (e *Else) Body() []FunctionStatement {
	return e.Statements
}

// PrecedingBlankLine returns true if this node was preceded by a blank line.
func (e *Else) PrecedingBlankLine() bool {
	return e.HasPrecedingBlankLine
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (e *Else) Accept(v Visitor) error {
	return v.VisitElse(e)
}

// Comments returns the [CrosslineComments] associated
// with this node or nil if there are none.
func (e *Else) Comments() *CrosslineComments {
	return e.NodeComments
}

// Location returns the source location of the node.
func (e *Else) Location() source.Location {
	return source.Span(e.KeywordLocation, e.Statements[len(e.Statements)-1].Location())
}

func (*Else) block() {}

var _ Node = (*Else)(nil)
