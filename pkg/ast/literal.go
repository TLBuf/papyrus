package ast

import "github.com/TLBuf/papyrus/pkg/source"

// BoolLiteral is a boolean literal (i.e. true or false).
type BoolLiteral struct {
	Trivia
	// Text is the BoolLiteral token.
	Text *Token
	// Value is the parsed value of the string literal.
	Value bool
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *BoolLiteral) Accept(v Visitor) error {
	return v.VisitBoolLiteral(l)
}

// SourceLocation returns the source location of the node.
func (l *BoolLiteral) SourceLocation() source.Location {
	return l.Location
}

func (*BoolLiteral) expression() {}

func (*BoolLiteral) literal() {}

func (*BoolLiteral) functionStatement() {}

var _ Literal = (*BoolLiteral)(nil)

// IntLiteral is an integer literal.
type IntLiteral struct {
	Trivia
	// Text is the IntLiteral token.
	Text *Token
	// Value is the parsed value of the string literal.
	Value int
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *IntLiteral) Accept(v Visitor) error {
	return v.VisitIntLiteral(l)
}

// SourceLocation returns the source location of the node.
func (l *IntLiteral) SourceLocation() source.Location {
	return l.Location
}

func (*IntLiteral) expression() {}

func (*IntLiteral) literal() {}

func (*IntLiteral) functionStatement() {}

var _ Literal = (*IntLiteral)(nil)

// FloatLiteral is a floating-point literal.
type FloatLiteral struct {
	Trivia
	// Text is the FloatLiteral token.
	Text *Token
	// Value is the parsed value of the string literal.
	Value float32
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *FloatLiteral) Accept(v Visitor) error {
	return v.VisitFloatLiteral(l)
}

// SourceLocation returns the source location of the node.
func (l *FloatLiteral) SourceLocation() source.Location {
	return l.Location
}

func (*FloatLiteral) expression() {}

func (*FloatLiteral) literal() {}

func (*FloatLiteral) functionStatement() {}

var _ Literal = (*FloatLiteral)(nil)

// StringLiteral is a string literal.
type StringLiteral struct {
	Trivia
	// Text is the StringLiteral token.
	Text *Token
	// Value is the parsed value of the string literal.
	Value string
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *StringLiteral) Accept(v Visitor) error {
	return v.VisitStringLiteral(l)
}

// SourceLocation returns the source location of the node.
func (l *StringLiteral) SourceLocation() source.Location {
	return l.Location
}

func (*StringLiteral) expression() {}

func (*StringLiteral) literal() {}

func (*StringLiteral) functionStatement() {}

var _ Literal = (*StringLiteral)(nil)

// NoneLiteral is the none literal (i.e. the null object literal).
type NoneLiteral struct {
	Trivia
	// Text is the None token.
	Text *Token
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *NoneLiteral) Accept(v Visitor) error {
	return v.VisitNoneLiteral(l)
}

// SourceLocation returns the source location of the node.
func (l *NoneLiteral) SourceLocation() source.Location {
	return l.Location
}

func (*NoneLiteral) expression() {}

func (*NoneLiteral) literal() {}

func (*NoneLiteral) functionStatement() {}

var _ Literal = (*NoneLiteral)(nil)
