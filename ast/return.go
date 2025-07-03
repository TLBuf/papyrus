package ast

import "github.com/TLBuf/papyrus/source"

// Return is a statement that terminates a function potentially with a value.
type Return struct {
	LineTrivia
	// Keyword is the Return keyword token.
	Keyword *Token
	// Value is the expression that defines the value to return or nil if there is
	// none (i.e. the function doesn't return a value).
	Value Expression
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (r *Return) Trivia() LineTrivia {
	return r.LineTrivia
}

// Accept calls the appropriate visitor method for the node.
func (r *Return) Accept(v Visitor) error {
	return v.VisitReturn(r)
}

// Location returns the source location of the node.
func (r *Return) Location() source.Location {
	return r.NodeLocation
}

func (*Return) statement() {}

func (*Return) functionStatement() {}

var _ FunctionStatement = (*Return)(nil)
