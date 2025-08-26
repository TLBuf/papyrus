package source_test

import (
	"testing"

	"github.com/TLBuf/papyrus/source"
	"github.com/google/go-cmp/cmp"
)

func TestStartLine(t *testing.T) {
	tests := []struct {
		name     string
		file     *source.File
		location source.Location
		want     uint32
	}{
		{
			name:     "EmptyFile",
			file:     file(""),
			location: location(0, 1),
			want:     0,
		}, {
			name:     "FirstByte",
			file:     file("01\n34\n67\n"),
			location: location(0, 1),
			want:     1,
		}, {
			name:     "LastByte",
			file:     file("01\n34\n67\n"),
			location: location(8, 1),
			want:     3,
		}, {
			name:     "LineStart",
			file:     file("01\n34\n67\n"),
			location: location(3, 1),
			want:     2,
		}, {
			name:     "LineEnd",
			file:     file("01\n34\n67\n"),
			location: location(5, 1),
			want:     2,
		}, {
			name:     "PastEnd",
			file:     file("01\n34\n67\n"),
			location: location(10, 1),
			want:     0,
		}, {
			name:     "LastByteNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(7, 1),
			want:     3,
		}, {
			name:     "PastEndNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(10, 1),
			want:     0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.file.StartLine(test.location)
			if got != test.want {
				t.Errorf("StartLine(%v) = %d, want %d", test.location, got, test.want)
			}
		})
	}
}

func TestEndLine(t *testing.T) {
	tests := []struct {
		name     string
		file     *source.File
		location source.Location
		want     uint32
	}{
		{
			name:     "EmptyFile",
			file:     file(""),
			location: location(0, 1),
			want:     0,
		}, {
			name:     "FirstByte",
			file:     file("01\n34\n67\n"),
			location: location(0, 1),
			want:     1,
		}, {
			name:     "LastByte",
			file:     file("01\n34\n67\n"),
			location: location(8, 1),
			want:     3,
		}, {
			name:     "LineStart",
			file:     file("01\n34\n67\n"),
			location: location(3, 1),
			want:     2,
		}, {
			name:     "LineEnd",
			file:     file("01\n34\n67\n"),
			location: location(5, 1),
			want:     2,
		}, {
			name:     "PastEnd",
			file:     file("01\n34\n67\n"),
			location: location(10, 1),
			want:     0,
		}, {
			name:     "LastByteNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(7, 1),
			want:     3,
		}, {
			name:     "PastEndNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(10, 1),
			want:     0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.file.EndLine(test.location)
			if got != test.want {
				t.Errorf("EndLine(%v) = %d, want %d", test.location, got, test.want)
			}
		})
	}
}

func TestStartColumn(t *testing.T) {
	tests := []struct {
		name     string
		file     *source.File
		location source.Location
		want     uint32
	}{
		{
			name:     "EmptyFile",
			file:     file(""),
			location: location(0, 1),
			want:     0,
		}, {
			name:     "FirstByte",
			file:     file("01\n34\n67\n"),
			location: location(0, 1),
			want:     1,
		}, {
			name:     "LastByte",
			file:     file("01\n34\n67\n"),
			location: location(8, 1),
			want:     3,
		}, {
			name:     "LineStart",
			file:     file("01\n34\n67\n"),
			location: location(3, 1),
			want:     1,
		}, {
			name:     "LineEnd",
			file:     file("01\n34\n67\n"),
			location: location(4, 1),
			want:     2,
		}, {
			name:     "PastEnd",
			file:     file("01\n34\n67\n"),
			location: location(10, 1),
			want:     0,
		}, {
			name:     "LastByteNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(7, 1),
			want:     2,
		}, {
			name:     "PastEndNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(10, 1),
			want:     0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.file.StartColumn(test.location)
			if got != test.want {
				t.Errorf("StartColumn(%v) = %d, want %d", test.location, got, test.want)
			}
		})
	}
}

func TestEndColumn(t *testing.T) {
	tests := []struct {
		name     string
		file     *source.File
		location source.Location
		want     uint32
	}{
		{
			name:     "EmptyFile",
			file:     file(""),
			location: location(0, 1),
			want:     0,
		}, {
			name:     "FirstByte",
			file:     file("01\n34\n67\n"),
			location: location(0, 1),
			want:     1,
		}, {
			name:     "LastByte",
			file:     file("01\n34\n67\n"),
			location: location(8, 1),
			want:     3,
		}, {
			name:     "LineStart",
			file:     file("01\n34\n67\n"),
			location: location(3, 1),
			want:     1,
		}, {
			name:     "LineEnd",
			file:     file("01\n34\n67\n"),
			location: location(4, 1),
			want:     2,
		}, {
			name:     "PastEnd",
			file:     file("01\n34\n67\n"),
			location: location(10, 1),
			want:     0,
		}, {
			name:     "LastByteNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(7, 1),
			want:     2,
		}, {
			name:     "PastEndNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(10, 1),
			want:     0,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.file.EndColumn(test.location)
			if got != test.want {
				t.Errorf("EndColumn(%v) = %d, want %d", test.location, got, test.want)
			}
		})
	}
}

func TestPremable(t *testing.T) {
	tests := []struct {
		name     string
		file     *source.File
		location source.Location
		want     []byte
	}{
		{
			name:     "EmptyFile",
			file:     file(""),
			location: location(0, 1),
			want:     nil,
		}, {
			name:     "FirstByte",
			file:     file("01\n34\n67\n"),
			location: location(0, 1),
			want:     []byte(""),
		}, {
			name:     "LastByte",
			file:     file("01\n34\n67\n"),
			location: location(8, 1),
			want:     []byte("67"),
		}, {
			name:     "LineStart",
			file:     file("01\n34\n67\n"),
			location: location(3, 1),
			want:     []byte(""),
		}, {
			name:     "LineEnd",
			file:     file("01\n34\n67\n"),
			location: location(4, 1),
			want:     []byte("3"),
		}, {
			name:     "PastEnd",
			file:     file("01\n34\n67\n"),
			location: location(10, 1),
			want:     nil,
		}, {
			name:     "LastByteNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(7, 1),
			want:     []byte("6"),
		}, {
			name:     "PastEndNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(10, 1),
			want:     nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.file.Preamble(test.location)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Preamble(%v) mismatch (-want +got):\n%s", test.location, diff)
			}
		})
	}
}

func TestPostamble(t *testing.T) {
	tests := []struct {
		name     string
		file     *source.File
		location source.Location
		want     []byte
	}{
		{
			name:     "EmptyFile",
			file:     file(""),
			location: location(0, 1),
			want:     nil,
		}, {
			name:     "FirstByte",
			file:     file("01\n34\n67\n"),
			location: location(0, 1),
			want:     []byte("1"),
		}, {
			name:     "LastByte",
			file:     file("01\n34\n67\n"),
			location: location(8, 1),
			want:     []byte(""),
		}, {
			name:     "LineStart",
			file:     file("01\n34\n67\n"),
			location: location(3, 1),
			want:     []byte("4"),
		}, {
			name:     "LineEnd",
			file:     file("01\n34\n67\n"),
			location: location(4, 1),
			want:     []byte(""),
		}, {
			name:     "PastEnd",
			file:     file("01\n34\n67\n"),
			location: location(10, 1),
			want:     nil,
		}, {
			name:     "LastByteNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(7, 1),
			want:     []byte(""),
		}, {
			name:     "PastEndNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(10, 1),
			want:     nil,
		}, {
			name:     "CarriageReturn",
			file:     file("01\r\n45\r\n89"),
			location: location(4, 1),
			want:     []byte("5"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.file.Postamble(test.location)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Postamble(%v) mismatch (-want +got):\n%s", test.location, diff)
			}
		})
	}
}

func TestContext(t *testing.T) {
	tests := []struct {
		name     string
		file     *source.File
		location source.Location
		want     []byte
	}{
		{
			name:     "EmptyFile",
			file:     file(""),
			location: location(0, 1),
			want:     nil,
		}, {
			name:     "FirstByte",
			file:     file("01\n34\n67\n"),
			location: location(0, 1),
			want:     []byte("01"),
		}, {
			name:     "LastByte",
			file:     file("01\n34\n67\n"),
			location: location(8, 1),
			want:     []byte("67"),
		}, {
			name:     "LineStart",
			file:     file("01\n34\n67\n"),
			location: location(3, 1),
			want:     []byte("34"),
		}, {
			name:     "LineEnd",
			file:     file("01\n34\n67\n"),
			location: location(4, 1),
			want:     []byte("34"),
		}, {
			name:     "PastEnd",
			file:     file("01\n34\n67\n"),
			location: location(10, 1),
			want:     nil,
		}, {
			name:     "LastByteNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(7, 1),
			want:     []byte("67"),
		}, {
			name:     "PastEndNoTrailingNewline",
			file:     file("01\n34\n67"),
			location: location(10, 1),
			want:     nil,
		}, {
			name:     "CarriageReturn",
			file:     file("01\r\n45\r\n89"),
			location: location(4, 1),
			want:     []byte("45"),
		}, {
			name:     "CrossLine",
			file:     file("01\r\n45\r\n89"),
			location: location(1, 4),
			want:     []byte("01\r\n45"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.file.Context(test.location)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Context(%v) mismatch (-want +got):\n%s", test.location, diff)
			}
		})
	}
}

func file(content string) *source.File {
	f, _ := source.NewFile("test.psc", []byte(content))
	return f
}

func location(offset, length uint32) source.Location {
	return source.Location{
		ByteOffset: offset,
		Length:     length,
	}
}
