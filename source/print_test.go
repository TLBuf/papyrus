package source_test

import (
	"bytes"
	"testing"

	"github.com/TLBuf/papyrus/source"
	"github.com/google/go-cmp/cmp"
)

func TestPrint(t *testing.T) {
	tests := []struct {
		name    string
		snippet source.Snippet
		want    string
	}{
		{
			"single_line_point",
			source.Snippet{
				Start:  source.Indicator{Column: 4},
				End:    source.Indicator{Column: 4},
				Lines:  []source.Line{{Number: 5, Chunks: []source.Chunk{{Text: "1234567890"}}}},
				Width:  20,
				Height: 3,
			},
			`
        ▼
 5 | 1234567890
`,
		},
		{
			"single_line_range",
			source.Snippet{
				Start:  source.Indicator{Column: 3},
				End:    source.Indicator{Column: 6},
				Lines:  []source.Line{{Number: 5, Chunks: []source.Chunk{{Text: "1234567890"}}}},
				Width:  20,
				Height: 3,
			},
			`
       ▼··▼
 5 | 1234567890
`,
		},
		{
			"multi_line_range",
			source.Snippet{
				Start: source.Indicator{Column: 3},
				End:   source.Indicator{Column: 6},
				Lines: []source.Line{
					{Number: 5, Chunks: []source.Chunk{{Text: "1234567890"}}},
					{Chunks: []source.Chunk{{Text: "... 2 lines ..."}}},
					{Number: 8, Chunks: []source.Chunk{{Text: "1234567890"}}},
				},
				Width:  20,
				Height: 3,
			},
			`
       ▼·················
 5 | 1234567890
 - | ... 2 lines ...
 8 | 1234567890
··········▲
`,
		},
		{
			"multi_line_range_number_width",
			source.Snippet{
				Start: source.Indicator{Column: 3},
				End:   source.Indicator{Column: 6},
				Lines: []source.Line{
					{Number: 9, Chunks: []source.Chunk{{Text: "1234567890"}}},
					{Chunks: []source.Chunk{{Text: "... 2 lines ..."}}},
					{Number: 12, Chunks: []source.Chunk{{Text: "1234567890"}}},
				},
				Width:  20,
				Height: 3,
			},
			`
        ▼·················
 09 | 1234567890
 -- | ... 2 lines ...
 12 | 1234567890
···········▲
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			source.Print(&buf, test.snippet)
			got := buf.String()
			if diff := cmp.Diff(test.want[1:], got); diff != "" {
				t.Errorf("Print() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
