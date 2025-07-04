package ast

import "github.com/TLBuf/papyrus/source"

// Comment is a common interface for non-doc comments.
type Comment interface {
	Node
	comment()
}

// BlockComment represents block comment.
type BlockComment struct {
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
	// OpenLocation is the location of the opening block comment token.
	OpenLocation source.Location
	// TextLocation is the location of the text of the comment (which may include
	// newlines).
	TextLocation source.Location
	// CloseLocation is the location of the closing block comment token.
	CloseLocation source.Location
}

// PrecedingBlankLine returns true if this node was preceded by a blank line.
func (c *BlockComment) PrecedingBlankLine() bool {
	return c.HasPrecedingBlankLine
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (c *BlockComment) Accept(v Visitor) error {
	return v.VisitBlockComment(c)
}

// Location returns the source location of the node.
func (c *BlockComment) Location() source.Location {
	return source.Span(c.OpenLocation, c.CloseLocation)
}

func (*BlockComment) comment() {}

var _ Comment = (*BlockComment)(nil)

// CommentBlock represents a block of one or more line comments.
//
// Not to be confused with a block comment, this construct consists of one or
// more standard line comments that appear on their own lines and one after the
// other without intervening blank lines.
type CommentBlock struct {
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
	// Comments are the line comments in this block in the order they appear.
	Elements []*LineComment
}

// PrecedingBlankLine returns true if this node was preceded by a blank line.
func (c *CommentBlock) PrecedingBlankLine() bool {
	return c.HasPrecedingBlankLine
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (c *CommentBlock) Accept(v Visitor) error {
	return v.VisitCommentBlock(c)
}

// Location returns the source location of the node.
func (c *CommentBlock) Location() source.Location {
	if len(c.Elements) == 1 {
		return c.Elements[0].Location()
	}
	return source.Span(c.Elements[0].Location(), c.Elements[len(c.Elements)-1].Location())
}

func (*CommentBlock) comment() {}

var _ Comment = (*CommentBlock)(nil)

// LineComment represents line comment.
type LineComment struct {
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
	// SemicolonLocation is the location of the semicolon that starts the comment.
	SemicolonLocation source.Location
	// TextLocation is the location of the text of the comment (which will never
	// include newlines).
	TextLocation source.Location
}

// PrecedingBlankLine returns true if this node was preceded by a blank line.
func (c *LineComment) PrecedingBlankLine() bool {
	return c.HasPrecedingBlankLine
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (c *LineComment) Accept(v Visitor) error {
	return v.VisitLineComment(c)
}

// Location returns the source location of the node.
func (c *LineComment) Location() source.Location {
	return source.Span(c.SemicolonLocation, c.TextLocation)
}

func (*LineComment) comment() {}

var _ Comment = (*LineComment)(nil)

// CrosslineComments is a set of comments associated with a node on the
// surrounding lines.
type CrosslineComments struct {
	// LeadingComments are the loose comments that appear on lines before a node.
	LeadingComments []Comment
	// TrailingComments are the loose comments that appear on lines
	// after a node, but which are not associated with another node.
	TrailingComments []Comment
}

// Leading returns the loose comments that appear before a node.
func (c *CrosslineComments) Leading() []Comment {
	if c == nil {
		return nil
	}
	return c.LeadingComments
}

// Trailing retunrs the loose comments that appear after a node, but which are
// not associated with another node.
func (c *CrosslineComments) Trailing() []Comment {
	if c == nil {
		return nil
	}
	return c.LeadingComments
}

// InlineComments is a set of comments associated with a node on the same line.
type InlineComments struct {
	// PrefixComments are the loose comments that
	// appear before a node on the same line.
	PrefixComments []Comment
	// SuffixComments are the loose comments that appear after a node on
	// the same line, but which are not associated with another node.
	SuffixComments []Comment
}

// Prefix returns the loose comments that appear before a node.
func (c *InlineComments) Prefix() []Comment {
	if c == nil {
		return nil
	}
	return c.PrefixComments
}

// Suffix retunrs the loose comments that appear after
// a node, but which are not associated with another node.
func (c *InlineComments) Suffix() []Comment {
	if c == nil {
		return nil
	}
	return c.SuffixComments
}

// Documentation represents a documentation comment.
type Documentation struct {
	// OpenLocation is the location of the opening brace.
	OpenLocation source.Location
	// TextLocation is the location of the text of the comment (which may include
	// newlines).
	TextLocation source.Location
	// CloseLocation is the location of the closing brace.
	CloseLocation source.Location
}

// Accept calls the appropriate visitor method for the node.
func (c *Documentation) Accept(v Visitor) error {
	return v.VisitDocumentation(c)
}

// Location returns the source location of the node.
func (c *Documentation) Location() source.Location {
	return source.Span(c.OpenLocation, c.CloseLocation)
}

var _ Node = (*Documentation)(nil)
