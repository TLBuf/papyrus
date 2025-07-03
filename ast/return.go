package ast

import "github.com/TLBuf/papyrus/source"

// Return is a statement that terminates a function potentially with a value.
type Return struct {
	LineTrivia
	// KeywordLocation is the location of the Return keyword that starts the
	// statement.
	KeywordLocation source.Location
	// Value is the expression that defines the value to return or nil if there is
	// none (i.e. the function doesn't return a value).
	Value Expression
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
	if r.Value == nil {
		return r.KeywordLocation
	}
	return source.Span(r.KeywordLocation, r.Value.Location())
}

func (*Return) statement() {}

func (*Return) functionStatement() {}

var _ FunctionStatement = (*Return)(nil)
