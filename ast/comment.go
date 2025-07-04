package ast

import "github.com/TLBuf/papyrus/source"

// Comment is a common interface for non-doc comments.
type Comment interface {
	Node
	comment()
}

// BlockComment represents block comment.
type BlockComment struct {
	LineTrivia

	// OpenLocation is the location of the opening block comment token.
	OpenLocation source.Location
	// TextLocation is the location of the text of the comment (which may include
	// newlines).
	TextLocation source.Location
	// CloseLocation is the location of the closing block comment token.
	CloseLocation source.Location
}

// Trivia returns the [LineTrivia] associated with this node.
func (c *BlockComment) Trivia() LineTrivia {
	return c.LineTrivia
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
	LineTrivia

	// Comments are the line comments in this block in the order they appear.
	Elements []*LineComment
}

// Trivia returns the [LineTrivia] associated with this node.
func (c *CommentBlock) Trivia() LineTrivia {
	return c.LineTrivia
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
	LineTrivia

	// SemicolonLocation is the location of the semicolon that starts the comment.
	SemicolonLocation source.Location
	// TextLocation is the location of the text of the comment (which will never
	// include newlines).
	TextLocation source.Location
}

// Trivia returns the [LineTrivia] associated with this node.
func (c *LineComment) Trivia() LineTrivia {
	return c.LineTrivia
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

// Comments is a set of comments associated with a node.
type Comments struct {
	// LeadingComments are the loose comments that appear before a node.
	LeadingComments []Comment
	// TrailingComments are the loose comments that appear after a node, but which
	// are not associated with another node.
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
// not associated with another node.
func (c *Comments) Trailing() []Comment {
	if c == nil {
		return nil
	}
	return c.LeadingComments
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
