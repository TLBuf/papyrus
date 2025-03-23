package source_test

import (
	"strings"
	"testing"

	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/google/go-cmp/cmp"
)

var file = &source.File{
	Text: []byte(strings.Repeat("12345678901234567890123456789012345678\r\n", 6)),
}

func TestSnippet(t *testing.T) {
	tests := []struct {
		name     string
		location source.Range
		want     source.Snippet
	}{
		{
			"point_single_line_fits",
			source.Range{
				File:            &source.File{Text: []byte("1234567890\r\n")},
				ByteOffset:      2,
				Length:          1,
				StartLine:       1,
				StartColumn:     3,
				EndLine:         1,
				EndColumn:       3,
				PreambleLength:  2,
				PostambleLength: 7,
			},
			source.Snippet{
				Start: source.Indicator{Column: 3},
				End:   source.Indicator{Column: 3},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "1234567890", IsSource: true},
				}}},
			},
		},
		{
			"point_single_line_first_half",
			source.Range{
				File:            file,
				ByteOffset:      42,
				Length:          1,
				StartLine:       1,
				StartColumn:     3,
				EndLine:         1,
				EndColumn:       3,
				PreambleLength:  2,
				PostambleLength: 35,
			},
			source.Snippet{
				Start: source.Indicator{Column: 3},
				End:   source.Indicator{Column: 3},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "12345678901234567", IsSource: true},
					{Text: "..."},
				}}},
			},
		},
		{
			"point_single_line_second_half",
			source.Range{
				File:            file,
				ByteOffset:      75,
				Length:          1,
				StartLine:       1,
				StartColumn:     36,
				EndLine:         1,
				EndColumn:       36,
				PreambleLength:  35,
				PostambleLength: 2,
			},
			source.Snippet{
				Start: source.Indicator{Column: 15},
				End:   source.Indicator{Column: 15},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "23456789012345678", IsSource: true},
				}}},
			},
		},
		{
			"point_single_line_middle",
			source.Range{
				File:            file,
				ByteOffset:      60,
				Length:          1,
				StartLine:       1,
				StartColumn:     21,
				EndLine:         1,
				EndColumn:       21,
				PreambleLength:  20,
				PostambleLength: 17,
			},
			source.Snippet{
				Start: source.Indicator{Column: 10},
				End:   source.Indicator{Column: 10},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "56789012345678", IsSource: true},
					{Text: "..."},
				}}},
			},
		},
		{
			"range_single_line_fits",
			source.Range{
				File:            &source.File{Text: []byte("1234567890\r\n")},
				ByteOffset:      2,
				Length:          5,
				StartLine:       1,
				StartColumn:     3,
				EndLine:         1,
				EndColumn:       7,
				PreambleLength:  2,
				PostambleLength: 3,
			},
			source.Snippet{
				Start: source.Indicator{Column: 3},
				End:   source.Indicator{Column: 7},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "1234567890", IsSource: true},
				}}},
			},
		},
		{
			"range_single_line_first_half",
			source.Range{
				File:            file,
				ByteOffset:      42,
				Length:          5,
				StartLine:       1,
				StartColumn:     3,
				EndLine:         1,
				EndColumn:       7,
				PreambleLength:  2,
				PostambleLength: 31,
			},
			source.Snippet{
				Start: source.Indicator{Column: 3},
				End:   source.Indicator{Column: 7},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "12345678901234567", IsSource: true},
					{Text: "..."},
				}}},
			},
		},
		{
			"range_single_line_second_half",
			source.Range{
				File:            file,
				ByteOffset:      71,
				Length:          5,
				StartLine:       1,
				StartColumn:     32,
				EndLine:         1,
				EndColumn:       36,
				PreambleLength:  32,
				PostambleLength: 2,
			},
			source.Snippet{
				Start: source.Indicator{Column: 11},
				End:   source.Indicator{Column: 15},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "23456789012345678", IsSource: true},
				}}},
			},
		},
		{
			"range_single_line_middle",
			source.Range{
				File:            file,
				ByteOffset:      59,
				Length:          3,
				StartLine:       1,
				StartColumn:     20,
				EndLine:         1,
				EndColumn:       22,
				PreambleLength:  19,
				PostambleLength: 16,
			},
			source.Snippet{
				Start: source.Indicator{Column: 9},
				End:   source.Indicator{Column: 11},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "56789012345678", IsSource: true},
					{Text: "..."},
				}}},
			},
		},
		{
			"range_single_line_middle_and_end",
			source.Range{
				File:            file,
				ByteOffset:      44,
				Length:          18,
				StartLine:       1,
				StartColumn:     5,
				EndLine:         1,
				EndColumn:       22,
				PreambleLength:  4,
				PostambleLength: 16,
			},
			source.Snippet{
				Start: source.Indicator{Column: 5},
				End:   source.Indicator{Column: 14},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "123456789", IsSource: true},
					{Text: "..."},
					{Text: "01234", IsSource: true},
					{Text: "..."},
				}}},
			},
		},
		{
			"range_single_line_middle_and_start",
			source.Range{
				File:            file,
				ByteOffset:      56,
				Length:          18,
				StartLine:       1,
				StartColumn:     17,
				EndLine:         1,
				EndColumn:       34,
				PreambleLength:  16,
				PostambleLength: 4,
			},
			source.Snippet{
				Start: source.Indicator{Column: 6},
				End:   source.Indicator{Column: 16},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "567890", IsSource: true},
					{Text: "..."},
					{Text: "12345678", IsSource: true},
				}}},
			},
		}, //          |                |
		{ // 12345678901234567890123456789012345678
			"range_single_line_start_middle_end",
			source.Range{
				File:            file,
				ByteOffset:      50,
				Length:          18,
				StartLine:       1,
				StartColumn:     11,
				EndLine:         1,
				EndColumn:       28,
				PreambleLength:  10,
				PostambleLength: 10,
			},
			source.Snippet{
				Start: source.Indicator{Column: 6},
				End:   source.Indicator{Column: 14},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "901234", IsSource: true},
					{Text: "..."},
					{Text: "67890", IsSource: true},
					{Text: "..."},
				}}},
			},
		},
		{
			"range_single_line_start_middle_end",
			source.Range{
				File:            file,
				ByteOffset:      56,
				Length:          7,
				StartLine:       1,
				StartColumn:     17,
				EndLine:         1,
				EndColumn:       23,
				PreambleLength:  16,
				PostambleLength: 15,
			},
			source.Snippet{
				Start: source.Indicator{Column: 7},
				End:   source.Indicator{Column: 13},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "45678901234567", IsSource: true},
					{Text: "..."},
				}}},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.location.Snippet(20, 4)
			if err != nil {
				t.Fatalf("Snippet() returned an unexpected error: %v", err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Snippet() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
