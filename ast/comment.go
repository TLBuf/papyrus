package ast

import "github.com/TLBuf/papyrus/source"

// Comment is a common interface for non-doc comments.
type Comment interface {
	Node
	comment()
}

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
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (c *BlockComment) Accept(v Visitor) error {
	return v.VisitBlockComment(c)
}

// Location returns the source location of the node.
func (c *BlockComment) Location() source.Location {
	return c.NodeLocation
}

func (*BlockComment) comment() {}

var _ Comment = (*BlockComment)(nil)

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
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (c *LineComment) Accept(v Visitor) error {
	return v.VisitLineComment(c)
}

// Location returns the source location of the node.
func (c *LineComment) Location() source.Location {
	return c.NodeLocation
}

func (*LineComment) comment() {}

var _ Comment = (*LineComment)(nil)

// Comments is a set of comments associated with a node.
type Comments struct {
	// LeadingComments are the loose comments that appear before a node.
	LeadingComments []Comment
	// TrailingComments are the loose comments that appear after a node, but which
	// are not assocaited with another node.
	TrailingComments []Comment
}

// Leading returns the loose comments that appear before a node.
func (c *Comments) Leading() []Comment {
	if c == nil {
		return nil
	}
	return c.LeadingComments
}

// Trailing retunrs the loose comments that appear after a node, but which are
// not assocaited with another node.
func (c *Comments) Trailing() []Comment {
	if c == nil {
		return nil
	}
	return c.LeadingComments
}

// Documentation represents a documentation comment.
type Documentation struct {
	// Open is the opening brace token for the documentation.
	//
	// This is always of kind [token.BraceOpen].
	Open *Token
	// Text is the token for the text of the documentation (which may include
	// newlines).
	//
	// This is always of kind [token.Comment].
	Text *Token
	// Close is the closing brace token for the docdocumentation.
	//
	// This is always of kind [token.BraceClose].
	Close *Token
	// NodeLocation is the source location of the node.
	NodeLocation source.Location
}

// Accept calls the appropriate visitor method for the node.
func (c *Documentation) Accept(v Visitor) error {
	return v.VisitDocumentation(c)
}

// Location returns the source location of the node.
func (c *Documentation) Location() source.Location {
	return c.NodeLocation
}

var _ Node = (*Documentation)(nil)
