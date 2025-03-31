package ast

import "github.com/TLBuf/papyrus/pkg/source"

// Access is an expression that reference a value or function that belongs to
// some scope.
type Access struct {
	Trivia
	// Value is the expression that defines the value have something accessed.
	Value Expression
	// Operator is the dot operator token for this access.
	Operator Token
	// Name is the name of the variable or function being accessed in value.
	Name *Identifier
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (a *Access) SourceLocation() source.Location {
	return a.Location
}

func (*Access) expression() {}

var _ Expression = (*Access)(nil)
