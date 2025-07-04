package ast

import "github.com/TLBuf/papyrus/source"

// BoolLiteral is a boolean literal (i.e. true or false).
type BoolLiteral struct {
	InfixTrivia

	// Value is the parsed value of the string literal.
	Value bool
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] associated with this node.
func (l *BoolLiteral) Trivia() InfixTrivia {
	return l.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (l *BoolLiteral) Accept(v Visitor) error {
	return v.VisitBoolLiteral(l)
}

// Location returns the source location of the node.
func (l *BoolLiteral) Location() source.Location {
	return l.NodeLocation
}

func (*BoolLiteral) expression() {}

func (*BoolLiteral) literal() {}

var _ Literal = (*BoolLiteral)(nil)

// IntLiteral is an integer literal.
type IntLiteral struct {
	InfixTrivia

	// Value is the parsed value of the string literal.
	Value int
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] associated with this node.
func (l *IntLiteral) Trivia() InfixTrivia {
	return l.InfixTrivia
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *IntLiteral) Accept(v Visitor) error {
	return v.VisitIntLiteral(l)
}

// Location returns the source location of the node.
func (l *IntLiteral) Location() source.Location {
	return l.NodeLocation
}

func (*IntLiteral) expression() {}

func (*IntLiteral) literal() {}

var _ Literal = (*IntLiteral)(nil)

// FloatLiteral is a floating-point literal.
type FloatLiteral struct {
	InfixTrivia

	// Value is the parsed value of the string literal.
	Value float32
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] associated with this node.
func (l *FloatLiteral) Trivia() InfixTrivia {
	return l.InfixTrivia
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *FloatLiteral) Accept(v Visitor) error {
	return v.VisitFloatLiteral(l)
}

// Location returns the source location of the node.
func (l *FloatLiteral) Location() source.Location {
	return l.NodeLocation
}

func (*FloatLiteral) expression() {}

func (*FloatLiteral) literal() {}

var _ Literal = (*FloatLiteral)(nil)

// StringLiteral is a string literal.
type StringLiteral struct {
	InfixTrivia

	// Value is the parsed value of the string literal.
	Value string
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] associated with this node.
func (l *StringLiteral) Trivia() InfixTrivia {
	return l.InfixTrivia
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *StringLiteral) Accept(v Visitor) error {
	return v.VisitStringLiteral(l)
}

// Location returns the source location of the node.
func (l *StringLiteral) Location() source.Location {
	return l.NodeLocation
}

func (*StringLiteral) expression() {}

func (*StringLiteral) literal() {}

var _ Literal = (*StringLiteral)(nil)

// NoneLiteral is the none literal (i.e. the null object literal).
type NoneLiteral struct {
	InfixTrivia

	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] associated with this node.
func (l *NoneLiteral) Trivia() InfixTrivia {
	return l.InfixTrivia
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *NoneLiteral) Accept(v Visitor) error {
	return v.VisitNoneLiteral(l)
}

// Location returns the source location of the node.
func (l *NoneLiteral) Location() source.Location {
	return l.NodeLocation
}

func (*NoneLiteral) expression() {}

func (*NoneLiteral) literal() {}

var _ Literal = (*NoneLiteral)(nil)
