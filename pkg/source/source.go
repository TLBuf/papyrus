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

// Location points to a range of bytes in a source code file.
type Location struct {
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
func (l Location) Text() []byte {
	return l.File.Text[l.ByteOffset : l.ByteOffset+l.Length]
}

// Preamble returns the text on the same line before this range.
func (l Location) Preamble() []byte {
	return l.File.Text[l.ByteOffset-l.PreambleLength : l.ByteOffset]
}

// Postamble returns the text on the same line after this range.
func (l Location) Postamble() []byte {
	return l.File.Text[l.ByteOffset+l.Length : l.ByteOffset+l.Length+l.PostambleLength]
}

// String implements [fmt.Stringer].
func (l Location) String() string {
	return fmt.Sprintf("%s:%d:%d", l.File.Path, l.StartLine, l.StartColumn)
}

// Span returns a Range that spans two given Ranges.
func Span(start, end Location) Location {
	return Location{
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
