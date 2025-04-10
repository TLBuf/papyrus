package ast

import "github.com/TLBuf/papyrus/pkg/source"

// BoolLiteral is a boolean literal (i.e. true or false).
type BoolLiteral struct {
	Trivia
	// Text is the BoolLiteral token.
	Text Token
	// Value is the parsed value of the string literal.
	Value bool
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (l *BoolLiteral) SourceLocation() source.Location {
	return l.Location
}

func (*BoolLiteral) expression() {}

func (*BoolLiteral) literal() {}

var _ Literal = (*BoolLiteral)(nil)

// IntLiteral is an integer literal.
type IntLiteral struct {
	Trivia
	// Text is the IntLiteral token.
	Text Token
	// Value is the parsed value of the string literal.
	Value int
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (l *IntLiteral) SourceLocation() source.Location {
	return l.Location
}

func (*IntLiteral) expression() {}

func (*IntLiteral) literal() {}

var _ Literal = (*IntLiteral)(nil)

// FloatLiteral is a floating-point literal.
type FloatLiteral struct {
	Trivia
	// Text is the FloatLiteral token.
	Text Token
	// Value is the parsed value of the string literal.
	Value float32
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (l *FloatLiteral) SourceLocation() source.Location {
	return l.Location
}

func (*FloatLiteral) expression() {}

func (*FloatLiteral) literal() {}

var _ Literal = (*FloatLiteral)(nil)

// StringLiteral is a string literal.
type StringLiteral struct {
	Trivia
	// Text is the StringLiteral token.
	Text Token
	// Value is the parsed value of the string literal.
	Value string
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (l *StringLiteral) SourceLocation() source.Location {
	return l.Location
}

func (*StringLiteral) expression() {}

func (*StringLiteral) literal() {}

var _ Literal = (*StringLiteral)(nil)

// NoneLiteral is the none literal (i.e. the null object literal).
type NoneLiteral struct {
	Trivia
	// Text is the None token.
	Text Token
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (l *NoneLiteral) SourceLocation() source.Location {
	return l.Location
}

func (*NoneLiteral) expression() {}

func (*NoneLiteral) literal() {}

var _ Literal = (*NoneLiteral)(nil)
