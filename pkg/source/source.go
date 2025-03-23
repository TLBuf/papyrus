// Package source provides utilities for referring to source code.
package source

import "fmt"

// File contains information for a source code file.
type File struct {
	// The path of the file.
	Path string
	// The full text of the file.
	Text []byte
}

// Range points to a range of bytes in a source code file.
type Range struct {
	// File is the file that contains the range.
	File *File
	// ByteOffset is the number of bytes from the start of the file for this
	// position.
	ByteOffset int
	// Length is the number of bytes in this range.
	Length int
	// StartLine is the 1-indexed line of the inclusive start of the range.
	StartLine int
	// StartColumn is the 1-indexed column of the inclusive start of the range.
	StartColumn int
	// EndLine is the 1-indexed line of the inclusive end of the range.
	EndLine int
	// EndColumn is the 1-indexed column of the inclusive end of the range.
	EndColumn int
	// PreambleLength is the number of bytes before the start of the range on the
	// same line as StartLine.
	PreambleLength int
	// PostambleLength is the number of bytes after the end of the range on the
	// same line as EndLine.
	PostambleLength int
}

// Text returns the text this range represents.
func (r Range) Text() []byte {
	return r.File.Text[r.ByteOffset : r.ByteOffset+r.Length]
}

// Preamble returns the text on the same line before this range.
func (r Range) Preamble() []byte {
	return r.File.Text[r.ByteOffset-r.PreambleLength : r.ByteOffset]
}

// Postamble returns the text on the same line after this range.
func (r Range) Postamble() []byte {
	return r.File.Text[r.ByteOffset+r.Length : r.ByteOffset+r.Length+r.PostambleLength]
}

// String implements fmt.Stringer.
func (r Range) String() string {
	return fmt.Sprintf("%s:%d:%d", r.File.Path, r.StartLine, r.StartColumn)
}

// Span returns a Range that spans two given Ranges.
func Span(start, end Range) Range {
	return Range{
		File:            start.File,
		ByteOffset:      start.ByteOffset,
		Length:          end.ByteOffset - start.ByteOffset + end.Length,
		StartLine:       start.StartLine,
		StartColumn:     start.StartColumn,
		EndLine:         end.EndLine,
		EndColumn:       end.EndColumn,
		PreambleLength:  start.PreambleLength,
		PostambleLength: end.PostambleLength,
	}
}
