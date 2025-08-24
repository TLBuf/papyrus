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
	// ByteOffset is the number of bytes from the start of the file for this
	// position.
	ByteOffset uint32
	// Length is the number of bytes in this range.
	Length uint32
	// StartLine is the 1-indexed line of the inclusive start of the range.
	StartLine uint32
	// StartColumn is the 1-indexed column of the inclusive start of the range.
	StartColumn uint32
	// EndLine is the 1-indexed line of the inclusive end of the range.
	EndLine uint32
	// EndColumn is the 1-indexed column of the inclusive end of the range.
	EndColumn uint32
	// PreambleLength is the number of bytes before the start of the range on the
	// same line as StartLine.
	PreambleLength uint32
	// PostambleLength is the number of bytes after the end of the range on the
	// same line as EndLine.
	PostambleLength uint32
}

// Text returns the text this range represents.
func (l Location) Text(file *File) []byte {
	ln := len(file.Text)
	return file.Text[min(int(l.ByteOffset), ln):min(int(l.ByteOffset+l.Length), ln)]
}

// Preamble returns the text on the same line before this range.
func (l Location) Preamble(file *File) []byte {
	ln := len(file.Text)
	return file.Text[min(int(l.ByteOffset-l.PreambleLength), ln):min(int(l.ByteOffset), ln)]
}

// Postamble returns the text on the same line after this range.
func (l Location) Postamble(file *File) []byte {
	ln := len(file.Text)
	return file.Text[min(int(l.ByteOffset+l.Length), ln):min(int(l.ByteOffset+l.Length+l.PostambleLength), ln)]
}

// String implements [fmt.Stringer].
func (l Location) String() string {
	return fmt.Sprintf("[%d:%d]", l.StartLine, l.StartColumn)
}

// Compare returns 0 if this location has the same byte offset as the given
// location, a negative number if this location has a smaller byte offset, or a
// positive number of this location has a larger byte offset.
func (l Location) Compare(o Location) int {
	return int(l.ByteOffset) - int(o.ByteOffset)
}

// Span returns a Range that spans two given Ranges.
func Span(start, end Location) Location {
	if end.ByteOffset < start.ByteOffset {
		panic(fmt.Sprintf("end before start: %d < %d", end.ByteOffset, start.ByteOffset))
	}
	return Location{
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
