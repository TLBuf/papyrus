package ast

import "github.com/TLBuf/papyrus/pkg/source"

// UnaryOperatorKind defines the various types of unary operations.
type UnaryOperatorKind int

const (
	// UnknownUnaryOperatorKind is the default (and invalid) unary operator.
	UnknownUnaryOperatorKind UnaryOperatorKind = iota
	// Negate is the negation operator, '-'.
	Negate
	// LogicalNot is the logical NOT operator, '!'.
	LogicalNot
)

func (o UnaryOperatorKind) String() string {
	name, ok := UnaryOperatorKindNames[o]
	if ok {
		return name
	}
	return "<unknown>"
}

var UnaryOperatorKindNames = map[UnaryOperatorKind]string{
	Negate:     "-",
	LogicalNot: "!",
}

// UnaryOperator represents a unary operator.
type UnaryOperator struct {
	Trivia
	// Kind is the type of unary operator.
	Kind UnaryOperatorKind
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (o *UnaryOperator) SourceLocation() source.Location {
	return o.Location
}

var _ Node = (*UnaryOperator)(nil)

// Unary is an expression that computes a value from two operands.
type Unary struct {
	Trivia
	// Operator defines the operator this unary expression uses.
	Operator *UnaryOperator
	// Operand is the operand.
	Operand Expression
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (u *Unary) SourceLocation() source.Location {
	return u.Location
}

func (*Unary) expression() {}

var _ Expression = (*Unary)(nil)
