package ast

import (
	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/TLBuf/papyrus/pkg/token"
)

// Token encodes a single lexical token in the Papyrus language.
//
// Each token has a [Kind] and information about where it is located
// in a source file.
type Token struct {
	// Kind defines the specific kind of token the text represents.
	Kind token.Kind
	// Location describes exactly where in a file the token is.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (t *Token) Accept(v Visitor) error {
	return v.VisitToken(t)
}

// SourceLocation returns the source location of the node.
func (t *Token) SourceLocation() source.Location {
	return t.Location
}
