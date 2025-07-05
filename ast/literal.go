package ast

import "github.com/TLBuf/papyrus/source"

// BoolLiteral is a boolean literal (i.e. true or false).
type BoolLiteral struct {
	// Value is the parsed value of the string literal.
	Value bool
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate visitor method for the node.
func (l *BoolLiteral) Accept(v Visitor) error {
	return v.VisitBoolLiteral(l)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (l *BoolLiteral) Comments() *Comments {
	return l.NodeComments
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
	// Value is the parsed value of the string literal.
	Value int
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *IntLiteral) Accept(v Visitor) error {
	return v.VisitIntLiteral(l)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (l *IntLiteral) Comments() *Comments {
	return l.NodeComments
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
	// Value is the parsed value of the string literal.
	Value float32
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *FloatLiteral) Accept(v Visitor) error {
	return v.VisitFloatLiteral(l)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (l *FloatLiteral) Comments() *Comments {
	return l.NodeComments
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
	// Value is the parsed value of the string literal.
	Value string
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *StringLiteral) Accept(v Visitor) error {
	return v.VisitStringLiteral(l)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (l *StringLiteral) Comments() *Comments {
	return l.NodeComments
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
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
	// NodeComments are the comments on before and/or after a node on the
	// same line or nil if the node has no comments associated with it.
	NodeComments *Comments
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (l *NoneLiteral) Accept(v Visitor) error {
	return v.VisitNoneLiteral(l)
}

// Comments returns the [Comments] associated
// with this node or nil if there are none.
func (l *NoneLiteral) Comments() *Comments {
	return l.NodeComments
}

// Location returns the source location of the node.
func (l *NoneLiteral) Location() source.Location {
	return l.NodeLocation
}

func (*NoneLiteral) expression() {}

func (*NoneLiteral) literal() {}

var _ Literal = (*NoneLiteral)(nil)
