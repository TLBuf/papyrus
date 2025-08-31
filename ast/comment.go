package ast

import (
	"fmt"

	"github.com/TLBuf/papyrus/source"
)

// Comment is a common interface for non-doc comments.
type Comment interface {
	Node

	// Prefix returns true if this node is followed
	// by anything other than a newline token.
	Prefix() bool
	// Suffix returns true if this node is preceded
	// by anything other than a newline token.
	Suffix() bool
	// LeadingBlankLine returns true if this node was preceded by a blank line.
	LeadingBlankLine() bool
	// TrailingBlankLine returns true if this node was followed by a blank line.
	TrailingBlankLine() bool

	comment()
}

// BlockComment represents block comment.
type BlockComment struct {
	// IsPrefix is true if this node is followed
	// by anything other than a newline token.
	IsPrefix bool
	// IsSuffix is true if this node is preceded
	// by anything other than a newline token.
	IsSuffix bool
	// HasLeadingBlankLine is true if this node was preceded by a blank line.
	HasLeadingBlankLine bool
	// HasTrailingBlankLine is true if this node was followed by a blank line.
	HasTrailingBlankLine bool
	// OpenLocation is the location of the opening block comment token.
	OpenLocation source.Location
	// TextLocation is the location of the text of the comment (which may include
	// newlines).
	TextLocation source.Location
	// CloseLocation is the location of the closing block comment token.
	CloseLocation source.Location
}

// Prefix returns true if this node is followed
// by anything other than a newline token.
func (c *BlockComment) Prefix() bool {
	return c.IsPrefix
}

// Suffix returns true if this node is preceded
// by anything other than a newline token.
func (c *BlockComment) Suffix() bool {
	return c.IsSuffix
}

// LeadingBlankLine returns true if this node was preceded by a blank line.
func (c *BlockComment) LeadingBlankLine() bool {
	return c.HasLeadingBlankLine
}

// TrailingBlankLine returns true if this node was followed by a blank line.
func (c *BlockComment) TrailingBlankLine() bool {
	return c.HasTrailingBlankLine
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (c *BlockComment) Accept(v Visitor) error {
	return v.VisitBlockComment(c)
}

// Location returns the source location of the node.
func (c *BlockComment) Location() source.Location {
	return source.Span(c.OpenLocation, c.CloseLocation)
}

func (c *BlockComment) String() string {
	return fmt.Sprintf("BlockComment%s", c.Location())
}

func (*BlockComment) comment() {}

var _ Comment = (*BlockComment)(nil)

// LineComment represents line comment.
type LineComment struct {
	// IsSuffix is true if this node is preceded
	// by anything other than a newline token.
	IsSuffix bool
	// HasLeadingBlankLine is true if this node was preceded by a blank line.
	HasLeadingBlankLine bool
	// HasTrailingBlankLine is true if this node was followed by a blank line.
	HasTrailingBlankLine bool
	// SemicolonLocation is the location of the semicolon that starts the comment.
	SemicolonLocation source.Location
	// TextLocation is the location of the text of the comment (which will never
	// include newlines).
	TextLocation source.Location
}

// Prefix returns false as line comments can never be prefixes.
func (*LineComment) Prefix() bool {
	return false
}

// Suffix returns true if this node is preceded
// by anything other than a newline token.
func (c *LineComment) Suffix() bool {
	return c.IsSuffix
}

// LeadingBlankLine returns true if this node was preceded by a blank line.
func (c *LineComment) LeadingBlankLine() bool {
	return c.HasLeadingBlankLine
}

// TrailingBlankLine returns true if this node was followed by a blank line.
func (c *LineComment) TrailingBlankLine() bool {
	return c.HasTrailingBlankLine
}

// Accept calls the appropriate method on the [Visitor] for the node.
func (c *LineComment) Accept(v Visitor) error {
	return v.VisitLineComment(c)
}

// Location returns the source location of the node.
func (c *LineComment) Location() source.Location {
	return source.Span(c.SemicolonLocation, c.TextLocation)
}

func (c *LineComment) String() string {
	return fmt.Sprintf("LineComment%s", c.Location())
}

func (*LineComment) comment() {}

var _ Comment = (*LineComment)(nil)

// Comments is a set of comments associated with a node on the same line.
type Comments struct {
	// PrefixComments are the loose comments that
	// appear before a node on the same line.
	PrefixComments []Comment
	// SuffixComments are the loose comments that appear after a node on
	// the same line, but which are not associated with another node.
	SuffixComments []Comment
}

// Prefix returns the loose comments that appear before a node.
func (c *Comments) Prefix() []Comment {
	return c.PrefixComments
}

// Suffix retunrs the loose comments that appear after
// a node, but which are not associated with another node.
func (c *Comments) Suffix() []Comment {
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

func (c *Documentation) String() string {
	return fmt.Sprintf("Documentation%s", c.Location())
}

var _ Node = (*Documentation)(nil)
