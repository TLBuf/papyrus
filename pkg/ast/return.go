package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Return is a statement that terminates a function potentially with a value.
type Return struct {
	Trivia
	// Keyword is the Return keyword token.
	Keyword *Token
	// Value is the expression that defines the value to return or nil if there is
	// none (i.e. the function doesn't return a value).
	Value Expression
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (r *Return) Accept(v Visitor) error {
	return v.VisitReturn(r)
}

// SourceLocation returns the source location of the node.
func (r *Return) SourceLocation() source.Location {
	return r.Location
}

func (*Return) functionStatement() {}

var _ FunctionStatement = (*Return)(nil)
