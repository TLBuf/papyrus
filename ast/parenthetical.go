package ast

import "github.com/TLBuf/papyrus/source"

// Parenthetical represents a parenthesized expression.
type Parenthetical struct {
	InfixTrivia
	// OpenLocation is the location of the opening parenthesis.
	OpenLocation source.Location
	// Value is the expression within the parentheses.
	Value Expression
	// CloseLocation is the location of the closing parenthesis.
	CloseLocation source.Location
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Trivia returns the [InfixTrivia] assocaited with this node.
func (p *Parenthetical) Trivia() InfixTrivia {
	return p.InfixTrivia
}

// Accept calls the appropriate visitor method for the node.
func (p *Parenthetical) Accept(v Visitor) error {
	return v.VisitParenthetical(p)
}

// Location returns the source location of the node.
func (p *Parenthetical) Location() source.Location {
	return p.NodeLocation
}

func (*Parenthetical) expression() {}

var _ Expression = (*Parenthetical)(nil)
