package ast

import "github.com/TLBuf/papyrus/source"

// Access is an expression that reference a value or function that belongs to
// some scope.
type Access struct {
	InfixTrivia

	// Value is the expression that defines the value have something accessed.
	Value Expression
	// DotLocation is the location of the dot operator.
	DotLocation source.Location
	// Name is the name of the variable or function being accessed in value.
	Name *Identifier
}

// Trivia returns the [InfixTrivia] associated with this node.
func (a *Access) Trivia() InfixTrivia {
	return a.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (a *Access) Accept(v Visitor) error {
	return v.VisitAccess(a)
}

// Location returns the source location of the node.
func (a *Access) Location() source.Location {
	return source.Span(a.Value.Location(), a.Name.Location())
}

func (*Access) expression() {}

var _ Expression = (*Access)(nil)
