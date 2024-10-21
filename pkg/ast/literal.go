package ast

import "github.com/TLBuf/papyrus/pkg/source"

// BoolLiteral is a boolean literal (i.e. true or false).
type BoolLiteral struct {
	// Value is the parsed value of the string literal.
	Value bool
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (l *BoolLiteral) Range() source.Range {
	return l.SourceRange
}

func (*BoolLiteral) expression() {}

func (*BoolLiteral) literal() {}

var _ Literal = (*BoolLiteral)(nil)

// IntLiteral is an integer literal.
type IntLiteral struct {
	// Value is the parsed value of the string literal.
	Value int
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (l *IntLiteral) Range() source.Range {
	return l.SourceRange
}

func (*IntLiteral) expression() {}

func (*IntLiteral) literal() {}

var _ Literal = (*IntLiteral)(nil)

// FloatLiteral is a floating-point literal.
type FloatLiteral struct {
	// Value is the parsed value of the string literal.
	Value float32
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (l *FloatLiteral) Range() source.Range {
	return l.SourceRange
}

func (*FloatLiteral) expression() {}

func (*FloatLiteral) literal() {}

var _ Literal = (*FloatLiteral)(nil)

// StringLiteral is a string literal.
type StringLiteral struct {
	// Value is the parsed value of the string literal.
	Value string
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (l *StringLiteral) Range() source.Range {
	return l.SourceRange
}

func (*StringLiteral) expression() {}

func (*StringLiteral) literal() {}

var _ Literal = (*StringLiteral)(nil)

// NoneLiteral is the none literal (i.e. the null object literal).
type NoneLiteral struct {
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (l *NoneLiteral) Range() source.Range {
	return l.SourceRange
}

func (*NoneLiteral) expression() {}

func (*NoneLiteral) literal() {}

var _ Literal = (*NoneLiteral)(nil)
