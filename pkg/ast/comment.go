package ast

import "github.com/TLBuf/papyrus/pkg/source"

// DocComment represents a documentation comment.
type DocComment struct {
	// Open is the opening brace token for the doc comment.
	//
	// This is always of kind [token.BraceOpen].
	Open *Token
	// Text is the token for the text of the comment (which may include newlines).
	//
	// This is always of kind [token.Comment].
	Text *Token
	// Close is the closing brace token for the doc comment.
	//
	// This is always of kind [token.BraceClose].
	Close *Token
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (c *DocComment) SourceLocation() source.Location {
	return c.Location
}

var _ Node = (*DocComment)(nil)

// BlockComment represents block comment.
type BlockComment struct {
	// Open is the opening brace token for the block comment.
	//
	// This is always of kind [token.BlockCommentOpen].
	Open *Token
	// Text is the token for the text of the comment (which may include newlines).
	//
	// This is always of kind [token.Comment].
	Text *Token
	// Close is the closing brace token for the block comment.
	//
	// This is always of kind [token.BlockCommentClose].
	Close *Token
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (c *BlockComment) Accept(v Visitor) error {
	return v.VisitBlockComment(c)
}

// SourceLocation returns the source location of the node.
func (c *BlockComment) SourceLocation() source.Location {
	return c.Location
}

func (*BlockComment) looseComment() {}

var _ LooseComment = (*BlockComment)(nil)

// LineComment represents line comment.
type LineComment struct {
	// Open is the semicolon that starts the comment.
	//
	// This is always of kind [token.Open].
	Open *Token
	// Text is the token for the text of the comment (which will never include a newline).
	//
	// This is always of kind [token.Comment].
	Text *Token
	// Location is the source range of the node.
	Location source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (c *LineComment) Accept(v Visitor) error {
	return v.VisitLineComment(c)
}

// SourceLocation returns the source location of the node.
func (c *LineComment) SourceLocation() source.Location {
	return c.Location
}

func (*LineComment) looseComment() {}

var _ LooseComment = (*LineComment)(nil)
