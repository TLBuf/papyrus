// Package source provides utilities for referring to source code.
package source

// File contains information for a source code file.
type File struct {
	// The path of the file.
	Path string
	// The full text of the file.
	Text []byte
}

// Position points to a single rune in a source code file.
type Position struct {
	// ByteOffset is the number of bytes from the start of the file for this position.
	ByteOffset int
	// Line is the 1-indexed line of the position in the file.
	Line int
	// Column is the 1-indexed column of the position in the file.
	Column int
}

// Range points to a range of runes in a source code file.
type Range struct {
	// File is the file that contains the range.
	File *File
	// Start is the inclusive position of the start of the range.
	Start Position
	// Start is the inclusive position of the end of the range.
	End Position
}

// Text returns the text this range represents.
func (r Range) Text() []byte {
	return r.File.Text[r.Start.ByteOffset : r.End.ByteOffset+1]
}
