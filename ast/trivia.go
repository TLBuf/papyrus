package ast

// InfixTrivia contains supplemental infix information for a node that has no
// semantic meaning, but which humans may find useful (i.e. comments).
type InfixTrivia struct {
	// Comments are the comments on before and/or after a node on the same line or
	// nil if the node has no comments associated with it.
	Comments *Comments
}

// LineTrivia contains supplemental information for an entire line (generally a
// statement) that has no semantic meaning, but which humans may find useful
// (i.e. comments).
type LineTrivia struct {
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
	// Comments are the comments on lines before and/or after a node or nil if the
	// node has no comments associated with it.
	Comments *Comments
}
