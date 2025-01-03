// Package source provides utilities for referring to source code.
package source

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
	// Line is the 1-indexed line of start of the range in the file.
	Line int
	// Column is the 1-indexed column start of the range in the file.
	Column int
}

// Text returns the text this range represents.
func (r Range) Text() []byte {
	return r.File.Text[r.ByteOffset : r.ByteOffset+r.Length]
}

// Span returns a Range that spans two given Ranges.
func Span(start, end Range) Range {
	return Range{
		File:       start.File,
		ByteOffset: start.ByteOffset,
		Length:     end.ByteOffset - start.ByteOffset + end.Length,
		Line:       start.Line,
		Column:     start.Line,
	}
}
