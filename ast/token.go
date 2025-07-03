package ast

import (
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/token"
)

// Token encodes a single lexical token in the Papyrus language.
//
// Each token has a [Kind] and information about where it is located
// in a source file.
type Token struct {
	// Kind defines the specific kind of token the text represents.
	Kind token.Kind
	// Location describes exactly where in a file the token is.
	NodeLocation source.Location
}

// Accept calls the appropriate visitor method for the node.
func (t *Token) Accept(v Visitor) error {
	return v.VisitToken(t)
}

// Location returns the source location of the node.
func (t *Token) Location() source.Location {
	return t.NodeLocation
}
