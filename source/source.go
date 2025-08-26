// Package source provides utilities for referring to source code.
package source

import (
	"fmt"
	"math"
	"slices"
	"unicode/utf8"
)

// File contains information for a source code file.
type File struct {
	path        string
	len         uint32
	content     []byte
	lineOffsets []uint32
}

// NewFile returns a new file or an error if the content is larger than ~4 GiB.
func NewFile(path string, content []byte) (*File, error) {
	if len(content) >= math.MaxUint32 {
		return nil, fmt.Errorf("content exceeds maximum size: %d >= %d", len(content), math.MaxUint32)
	}
	file := &File{
		path:        path,
		len:         uint32(len(content)), // #nosec G115 -- Checked above.
		content:     content,
		lineOffsets: append([]uint32(nil), 0),
	}
	for i := range file.len {
		if file.content[i] == '\n' {
			file.lineOffsets = append(file.lineOffsets, i+1)
		}
	}
	return file, nil
}

// Path returns the path of the file.
func (f *File) Path() string {
	return f.path
}

// Content returns the content of the file.
func (f *File) Content() []byte {
	return f.content
}

// Len returns the number of bytes in the file.
func (f *File) Len() uint32 {
	return f.len
}

// Bytes returns the bytes of content at the given location in this
// file or nil if the location is outside the range of this file.
func (f *File) Bytes(location Location) []byte {
	end := location.ByteOffset + location.Length
	if end > f.len {
		return nil
	}
	return f.content[location.ByteOffset:end]
}

// StartLine returns the 1-indexed line of the inclusive start of the
// location or zero if the location is outside the range of this file.
func (f *File) StartLine(location Location) uint32 {
	if location.ByteOffset >= f.len {
		return 0
	}
	line, exact := slices.BinarySearch(f.lineOffsets, location.ByteOffset)
	if exact {
		return uint32(line + 1) // #nosec G115 -- Checked in NewFile.
	}
	return uint32(line) // #nosec G115 -- Checked in NewFile.
}

// StartColumn returns the 1-indexed column of the inclusive start of the
// location or zero if the location is outside the range of this file.
func (f *File) StartColumn(location Location) uint32 {
	if location.ByteOffset >= f.len {
		return 0
	}
	return uint32(utf8.RuneCount(f.Preamble(location))) + 1 // #nosec G115 -- Checked in NewFile.
}

// EndLine returns the 1-indexed line of the inclusive end of the
// location or zero if the location is outside the range of this file.
func (f *File) EndLine(location Location) uint32 {
	end := max(location.ByteOffset+location.Length-1, 0)
	if end >= f.len {
		return 0
	}
	line, exact := slices.BinarySearch(f.lineOffsets, end)
	if exact {
		return uint32(line + 1) // #nosec G115 -- Checked in NewFile.
	}
	if line == len(f.lineOffsets) {
		return uint32(len(f.lineOffsets)) // #nosec G115 -- Checked in NewFile.
	}
	return uint32(line) // #nosec G115 -- Checked in NewFile.
}

// EndColumn returns the 1-indexed column of the inclusive end of the
// location or zero if the location is outside the range of this file.
func (f *File) EndColumn(location Location) uint32 {
	end := max(location.ByteOffset+location.Length-1, 0)
	if end >= f.len {
		return 0
	}
	return uint32(utf8.RuneCount(f.content[f.lineStart(end) : end+1])) // #nosec G115 -- Checked in NewFile.
}

// Preamble returns the content before a location
// on the same line as the start of the location.
func (f *File) Preamble(location Location) []byte {
	if location.ByteOffset >= f.len {
		return nil
	}
	return f.content[f.lineStart(location.ByteOffset):location.ByteOffset]
}

// Postamble returns the content after a location on the same line as the end of
// the location up to, but not including the trailing newline (and carriage
// return if present).
func (f *File) Postamble(location Location) []byte {
	offset := location.ByteOffset + location.Length
	if offset == f.len {
		return []byte{} // Location is valid, there's just no content left.
	}
	if offset > f.len {
		return nil
	}
	postamble := f.content[offset:f.lineEnd(offset)]
	last := len(postamble) - 1
	if last >= 0 && postamble[last] == '\n' {
		postamble = postamble[:last]
		last--
	}
	if last >= 0 && postamble[last] == '\r' {
		postamble = postamble[:last]
	}
	return postamble
}

// Context returns the content from the start of the line that contains the
// start of the location to the end of the line that contains the end of the
// location up to, but not including the trailing newline (and carriage return
// if present).
func (f *File) Context(location Location) []byte {
	end := location.ByteOffset + location.Length
	if end > f.len || location.ByteOffset >= f.len {
		return nil
	}
	context := f.content[f.lineStart(location.ByteOffset):f.lineEnd(min(end, f.len-1))]
	last := len(context) - 1
	if last >= 0 && context[last] == '\n' {
		context = context[:last]
		last--
	}
	if last >= 0 && context[last] == '\r' {
		context = context[:last]
	}
	return context
}

func (f *File) lineStart(offset uint32) uint32 {
	line, exact := slices.BinarySearch(f.lineOffsets, offset)
	if exact {
		return f.lineOffsets[line]
	}
	return f.lineOffsets[line-1]
}

func (f *File) lineEnd(offset uint32) uint32 {
	line, exact := slices.BinarySearch(f.lineOffsets, offset)
	if exact {
		return f.lineOffsets[line+1]
	}
	if line == len(f.lineOffsets) {
		return f.len
	}
	return f.lineOffsets[line]
}

// Location points to a range of bytes in a source code file.
type Location struct {
	// ByteOffset is the number of bytes from the start of the file for this
	// position.
	ByteOffset uint32
	// Length is the number of bytes in this range.
	Length uint32
}

// String implements [fmt.Stringer].
func (l Location) String() string {
	return fmt.Sprintf("[%d:%d]", l.ByteOffset, l.Length)
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
		ByteOffset: start.ByteOffset,
		Length:     end.ByteOffset - start.ByteOffset + end.Length,
	}
}
