package source_test

import (
	"strings"
	"testing"

	"github.com/TLBuf/papyrus/source"
	"github.com/google/go-cmp/cmp"
)

func TestSnippet(t *testing.T) {
	tests := []struct {
		name     string
		file     *source.File
		location source.Location
		want     source.Snippet
	}{
		{
			"point_single_line_fits",
			file("1234567890\r\n"),
			source.Location{
				ByteOffset: 2,
				Length:     1,
			},
			source.Snippet{
				Start: source.Indicator{Column: 3},
				End:   source.Indicator{Column: 3},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "1234567890", IsSource: true},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"point_single_line_tabs",
			file("123\t4567890\r\n"),
			source.Location{
				ByteOffset: 2,
				Length:     3,
			},
			source.Snippet{
				Start: source.Indicator{Column: 3},
				End:   source.Indicator{Column: 6},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "123  4567890", IsSource: true},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"point_single_line_first_half",
			file(strings.Repeat("12345678901234567890123456789012345678\r\n", 6)),
			source.Location{
				ByteOffset: 42,
				Length:     1,
			},
			source.Snippet{
				Start: source.Indicator{Column: 3},
				End:   source.Indicator{Column: 3},
				Lines: []source.Line{{Number: 2, Chunks: []source.Chunk{
					{Text: "12345678901234567", IsSource: true},
					{Text: "..."},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"point_single_line_second_half",
			file(strings.Repeat("12345678901234567890123456789012345678\r\n", 6)),
			source.Location{
				ByteOffset: 75,
				Length:     1,
			},
			source.Snippet{
				Start: source.Indicator{Column: 18},
				End:   source.Indicator{Column: 18},
				Lines: []source.Line{{Number: 2, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "23456789012345678", IsSource: true},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"point_single_line_middle",
			file(strings.Repeat("12345678901234567890123456789012345678\r\n", 6)),
			source.Location{
				ByteOffset: 60,
				Length:     1,
			},
			source.Snippet{
				Start: source.Indicator{Column: 10},
				End:   source.Indicator{Column: 10},
				Lines: []source.Line{{Number: 2, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "56789012345678", IsSource: true},
					{Text: "..."},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"range_single_line_fits",
			file("1234567890\r\n"),
			source.Location{
				ByteOffset: 2,
				Length:     5,
			},
			source.Snippet{
				Start: source.Indicator{Column: 3},
				End:   source.Indicator{Column: 7},
				Lines: []source.Line{{Number: 1, Chunks: []source.Chunk{
					{Text: "1234567890", IsSource: true},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"range_single_line_first_half",
			file(strings.Repeat("12345678901234567890123456789012345678\r\n", 6)),
			source.Location{
				ByteOffset: 42,
				Length:     5,
			},
			source.Snippet{
				Start: source.Indicator{Column: 3},
				End:   source.Indicator{Column: 7},
				Lines: []source.Line{{Number: 2, Chunks: []source.Chunk{
					{Text: "12345678901234567", IsSource: true},
					{Text: "..."},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"range_single_line_second_half",
			file(strings.Repeat("12345678901234567890123456789012345678\r\n", 6)),
			source.Location{
				ByteOffset: 71,
				Length:     5,
			},
			source.Snippet{
				Start: source.Indicator{Column: 11},
				End:   source.Indicator{Column: 15},
				Lines: []source.Line{{Number: 2, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "23456789012345678", IsSource: true},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"range_single_line_middle",
			file(strings.Repeat("12345678901234567890123456789012345678\r\n", 6)),
			source.Location{
				ByteOffset: 59,
				Length:     3,
			},
			source.Snippet{
				Start: source.Indicator{Column: 9},
				End:   source.Indicator{Column: 11},
				Lines: []source.Line{{Number: 2, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "56789012345678", IsSource: true},
					{Text: "..."},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"range_single_line_middle_and_end",
			file(strings.Repeat("12345678901234567890123456789012345678\r\n", 6)),
			source.Location{
				ByteOffset: 44,
				Length:     18,
			},
			source.Snippet{
				Start: source.Indicator{Column: 5},
				End:   source.Indicator{Column: 14},
				Lines: []source.Line{{Number: 2, Chunks: []source.Chunk{
					{Text: "123456789", IsSource: true},
					{Text: "..."},
					{Text: "01234", IsSource: true},
					{Text: "..."},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"range_single_line_middle_and_start",
			file(strings.Repeat("12345678901234567890123456789012345678\r\n", 6)),
			source.Location{
				ByteOffset: 56,
				Length:     18,
			},
			source.Snippet{
				Start: source.Indicator{Column: 6},
				End:   source.Indicator{Column: 16},
				Lines: []source.Line{{Number: 2, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "567890", IsSource: true},
					{Text: "..."},
					{Text: "12345678", IsSource: true},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"range_single_line_start_middle_end",
			file(strings.Repeat("12345678901234567890123456789012345678\r\n", 6)),
			source.Location{
				ByteOffset: 50,
				Length:     18,
			},
			source.Snippet{
				Start: source.Indicator{Column: 6},
				End:   source.Indicator{Column: 14},
				Lines: []source.Line{{Number: 2, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "901234", IsSource: true},
					{Text: "..."},
					{Text: "67890", IsSource: true},
					{Text: "..."},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"range_single_line_start_middle_end",
			file(strings.Repeat("12345678901234567890123456789012345678\r\n", 6)),
			source.Location{
				ByteOffset: 56,
				Length:     7,
			},
			source.Snippet{
				Start: source.Indicator{Column: 7},
				End:   source.Indicator{Column: 13},
				Lines: []source.Line{{Number: 2, Chunks: []source.Chunk{
					{Text: "..."},
					{Text: "45678901234567", IsSource: true},
					{Text: "..."},
				}}},
				Width:  20,
				Height: 4,
			},
		},
		{
			"range_multi_line",
			file(strings.Repeat("12345678901234567890123456789012345678\r\n", 6)),
			source.Location{
				ByteOffset: 2,
				Length:     201,
			},
			source.Snippet{
				Start: source.Indicator{Column: 3},
				End:   source.Indicator{Column: 3},
				Lines: []source.Line{
					{Number: 1, Chunks: []source.Chunk{{Text: "12345678901234567", IsSource: true}, {Text: "..."}}},
					{Number: 2, Chunks: []source.Chunk{{Text: "12345678901234567", IsSource: true}, {Text: "..."}}},
					{Chunks: []source.Chunk{{Text: "... 3 lines ..."}}},
					{Number: 6, Chunks: []source.Chunk{{Text: "12345678901234567", IsSource: true}, {Text: "..."}}},
				},
				Width:  20,
				Height: 4,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.location.Snippet(test.file, 20, 4)
			if err != nil {
				t.Fatalf("Snippet() returned an unexpected error: %v", err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("Snippet() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
