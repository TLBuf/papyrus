package ast

import "github.com/TLBuf/papyrus/source"

// Parenthetical represents a parenthesized expression.
type Parenthetical struct {
	InfixTrivia
	// Open is the open parenthesis token.
	Open *Token
	// Value is the expression within the parentheses.
	Value Expression
	// Close is the close parenthesis token.
	Close *Token
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
