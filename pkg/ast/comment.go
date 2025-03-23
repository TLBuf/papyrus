package ast

import "github.com/TLBuf/papyrus/pkg/source"

// DocComment represents a documentation comment.
type DocComment struct {
	// Text is the text of the comment (which may include newlines).
	Text string
	// Location is the source range of the node.
	Location source.Range
}

// Range returns the source range of the node.
func (c *DocComment) Range() source.Range {
	return c.Location
}

var _ Node = (*DocComment)(nil)

// BlockComment represents block comment.
type BlockComment struct {
	// Text is the text of the comment (which may include newlines).
	Text string
	// Location is the source range of the node.
	Location source.Range
}

// Range returns the source range of the node.
func (c *BlockComment) Range() source.Range {
	return c.Location
}

func (*BlockComment) looseComment() {}

var _ LooseComment = (*BlockComment)(nil)

// LineComment represents line comment.
type LineComment struct {
	// Text is the text of the comment (which will never include a newline).
	Text string
	// Location is the source range of the node.
	Location source.Range
}

// Range returns the source range of the node.
func (c *LineComment) Range() source.Range {
	return c.Location
}

func (*LineComment) looseComment() {}

var _ LooseComment = (*LineComment)(nil)
