package ast

import "github.com/TLBuf/papyrus/pkg/source"

// DocComment represents a documentation comment.
type DocComment struct {
	// Text is the text of the comment (which may include newlines).
	Text string
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (c *DocComment) Range() source.Range {
	return c.SourceRange
}

var _ Node = (*DocComment)(nil)

// BlockComment represents block comment.
type BlockComment struct {
	// Text is the text of the comment (which may include newlines).
	Text string
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (c *BlockComment) Range() source.Range {
	return c.SourceRange
}

func (*BlockComment) scriptStatement() {}

func (*BlockComment) functionStatement() {}

var _ ScriptStatement = (*BlockComment)(nil)
var _ FunctionStatement = (*BlockComment)(nil)

// LineComment represents line comment.
type LineComment struct {
	// Text is the text of the comment (which will never include a newline).
	Text string
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (c *LineComment) Range() source.Range {
	return c.SourceRange
}

var _ Node = (*LineComment)(nil)

// CommentBlock represents a contiguous group of line comments.
type CommentGroup struct {
	// Comments is the list of contiguous line comments.
	Comments []*LineComment
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (c *CommentGroup) Range() source.Range {
	return c.SourceRange
}

func (*CommentGroup) comment() {}

var _ Node = (*CommentGroup)(nil)

// StatementCommentSet is the set of comments associated with a statement.
type StatementCommentSet struct {
	// Leading is the list of comment groups that lead the statement.
	Leading []*CommentGroup
	// Trailing is the list of comment groups that follow the statement.
	Trailing []*CommentGroup
}

// NodeCommentSet is the set of comments associated with a node.
type NodeCommentSet struct {
	// Suffix is the line comment that follows the node on the same line.
	Suffix *LineComment
	// Blocks is the list of block comments that appear within this node and
	// which are not more closely associated with some child node.
	Blocks []*BlockComment
}
