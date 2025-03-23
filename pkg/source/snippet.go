package source

import (
	"fmt"
)

const (
	// MinimumSnippetWidth is the absolute minimum width of Snippet's content.
	MinimumSnippetWidth = 20
	// MinimumSnippetHeight is the absolute minimum height of a Snippet's content.
	MinimumSnippetHeight = 3
	// DefaultTabWidth is the number of spaces substituted for tabs.
	DefaultTabWidth = 2
)

// Snippet is a section of source code, formatted to fit in a specified number
// of columns and lines with a range start and end indicated.
type Snippet struct {
	// Start and End indicate a section of the snippet that is of special
	// importance (as snippets typically include a bit of extra context).
	//
	// If there is only one line, these apply to that line, otherwise, Start
	// applies to the first line and End applies to the last.
	Start, End Indicator
	// Lines are the lines of the source code clipped
	Lines []Line
	// Width and Height are the upper bounds on the number of runes wide and lines
	// respectively in the Snippet.
	Width, Height int
}

// Indicator describes a indicator of some piece of content.
type Indicator struct {
	// Column is the 1-indexed column being indicated.
	Column int
}

// Line is a single numbered line of source code.
type Line struct {
	// Number is the 1-indexed line number.
	Number int
	// Chunks are the chunks of text of the line in order.
	Chunks []Chunk
}

// Chunk is a single segment of text that never contains newlines, carriage
// returns, or tabs.
type Chunk struct {
	// Text is the literal text of the chunk.
	Text string
	// IsSource is true if the text of this chunk is source code.
	IsSource bool
}

type snippetOptions struct {
	tabWidth int
}

// SnippetOption
type SnippetOption func(*snippetOptions)

// WithTabWidth returns a [SnippetOption] that overrides the [DefaultTabWidth].
func WithTabWidth(w int) SnippetOption {
	return func(opts *snippetOptions) {
		opts.tabWidth = w
	}
}

// Snippet returns the range formatted to fit in the given `width` and `height`.
//
// An error is returns if `width` is less than [MinimumSnippetWidth] or `height`
// is less than [MinimumSnippetHeight].
func (r Range) Snippet(width, height int, opts ...SnippetOption) (Snippet, error) {
	options := snippetOptions{
		tabWidth: DefaultTabWidth,
	}
	for _, opt := range opts {
		opt(&options)
	}
	if width < MinimumSnippetWidth {
		return Snippet{}, fmt.Errorf("%d is less than minimum snippet width, %d", width, MinimumSnippetWidth)
	}
	if height < MinimumSnippetHeight {
		return Snippet{}, fmt.Errorf("%d is less than minimum snippet height, %d", height, MinimumSnippetHeight)
	}
	var snippet Snippet
	if r.StartLine == r.EndLine {
		snippet = formatSingleLineSnippet(r, width, options.tabWidth)
	} else {
		snippet = formatMultiLineSnippet(r, width, height, options.tabWidth)
	}
	snippet.Width = width
	snippet.Height = height
	return snippet, nil
}

func formatSingleLineSnippet(r Range, width, tabWidth int) Snippet {
	runes := textForSnippet(r)
	chunks, start, end := fitLine(runes, r.StartColumn, r.EndColumn, width, tabWidth)
	return Snippet{
		Start: Indicator{Column: start},
		End:   Indicator{Column: end},
		Lines: []Line{{Number: r.StartLine, Chunks: chunks}},
	}
}

func formatMultiLineSnippet(r Range, width, height, tabWidth int) Snippet {
	text := splitLines(textForSnippet(r))
	first, start, _ := fitLine(text[0], r.StartColumn, 0, width, tabWidth)
	last, end, _ := fitLine(text[len(text)-1], r.EndColumn, 0, width, tabWidth)
	remaining := r.EndLine - r.StartLine - 1
	available := max(0, height-3)
	lines := []Line{{Number: r.StartLine, Chunks: first}}
	if remaining <= available+1 {
		for i := 0; i < remaining; i++ {
			chunks, _, _ := fitLine(text[i+1], 0, 0, width, tabWidth)
			lines = append(lines, Line{Number: r.StartLine + i + 1, Chunks: chunks})
		}
		lines = append(lines, Line{Number: r.EndLine, Chunks: last})
		return Snippet{
			Start: Indicator{Column: start},
			End:   Indicator{Column: end},
			Lines: lines,
		}
	}
	heightA := available/2 + available%2
	heightB := available / 2
	for i := 0; i < heightA; i++ {
		chunks, _, _ := fitLine(text[i+1], 0, 0, width, tabWidth)
		lines = append(lines, Line{Number: r.StartLine + i + 1, Chunks: chunks})
	}
	omitted := remaining - available
	lines = append(lines, Line{Chunks: []Chunk{{Text: fmt.Sprintf("... %d lines ...", omitted)}}})
	for i := 0; i < heightB; i++ {
		chunks, _, _ := fitLine(text[i+omitted+1], 0, 0, width, tabWidth)
		lines = append(lines, Line{Number: r.StartLine + i + omitted + 1, Chunks: chunks})
	}
	lines = append(lines, Line{Number: r.EndLine, Chunks: last})
	return Snippet{
		Start: Indicator{Column: start},
		End:   Indicator{Column: end},
		Lines: lines,
	}
}

func textForSnippet(r Range) []rune {
	return []rune(string(r.File.Text[r.ByteOffset-r.PreambleLength : r.ByteOffset+r.Length+r.PostambleLength]))
}

func splitLines(text []rune) [][]rune {
	var lines [][]rune
	s := -1
	for i, r := range text {
		if r == '\r' || r == '\n' || r == 0 {
			if s >= 0 {
				lines = append(lines, text[s:i])
				s = -1
			}
			continue
		}
		if s < 0 {
			s = i
		}
	}
	if s > 0 {
		lines = append(lines, text[s:])
	}
	return lines
}

func fitLine(text []rune, start, end, width, tabWidth int) ([]Chunk, int, int) {
	text, start, end = substituteTabs(text, start, end, tabWidth)
	if start <= 0 {
		if len(text) < width {
			return []Chunk{{Text: string(text), IsSource: true}}, 0, 0
		}
		return []Chunk{
			{Text: string(text[:width-3]), IsSource: true},
			{Text: "..."},
		}, 0, 0
	}
	if end <= 0 || start == end {
		chunks, start := fitLineOnePoint(text, start, width)
		return chunks, start, start
	}
	return fitLineTwoPoints(text, start, end, width)
}

func fitLineOnePoint(text []rune, column, width int) ([]Chunk, int) {
	length := len(text)
	if length <= width {
		return []Chunk{{Text: string(text), IsSource: true}}, column
	}
	center := width / 2
	if column < center {
		// Column is in the first half, clip the end and send it.
		return []Chunk{
			{Text: string(text[:width-3]), IsSource: true},
			{Text: "..."},
		}, column
	}
	if length-column < center {
		// Column is in the last half, clip the start and send it.
		return []Chunk{
			{Text: "..."},
			{Text: string(text[length-width+3:]), IsSource: true},
		}, column - length + width - 3
	}
	// Pivot around column since it's somewhere in the middle.
	start := column - center - center%2 + 3
	end := column + center - 3
	return []Chunk{
		{Text: "..."},
		{Text: string(text[start:end]), IsSource: true},
		{Text: "..."},
	}, center
}

func fitLineTwoPoints(text []rune, start, end, width int) ([]Chunk, int, int) {
	length := len(text)
	if length <= width {
		return []Chunk{{Text: string(text), IsSource: true}}, start, end
	}
	available := width - 3
	if end < available {
		// Start and end fit, clip the end.
		return []Chunk{
			{Text: string(text[:available]), IsSource: true},
			{Text: "..."},
		}, start, end
	}
	if length-start < available {
		// Start and end are in the last half, clip the start and send it.
		return []Chunk{
			{Text: "..."},
			{Text: string(text[length-available:]), IsSource: true},
		}, start - length + available + 1, end - length + available + 1
	}
	contentWidth := end - start + 1
	available -= 3
	if contentWidth <= available {
		// Both start and end would fit in the, center things and clip each side.
		remaining := available - contentWidth
		a := start - remaining/2 - remaining%2
		b := end + remaining/2 + 1
		return []Chunk{
			{Text: "..."},
			{Text: string(text[a:b]), IsSource: true},
			{Text: "..."},
		}, remaining/2 + remaining%2 + 3, width - remaining/2 - 4
	}
	// Start and end are far enough apart that we have to clip in the middle.
	chunks := make([]Chunk, 0, 5)
	available -= 3
	widthA := available/2 + available%2
	pivotA := 3 + widthA/2 + widthA%2
	widthB := available / 2
	pivotB := length - widthB/2 - widthB%2 - 3
	if start < pivotA {
		chunks = append(chunks, Chunk{Text: string(text[:widthA+3]), IsSource: true})
	} else {
		s := start - widthA/2 + widthA%2
		e := start + widthA/2
		chunks = append(chunks, Chunk{Text: "..."}, Chunk{Text: string(text[s:e]), IsSource: true})
		start = pivotA
	}
	chunks = append(chunks, Chunk{Text: "..."})
	if end > pivotB {
		chunks = append(chunks, Chunk{Text: string(text[length-widthB-3:]), IsSource: true})
		end = width - length + end
	} else {
		s := end - widthB/2 - widthB%2
		e := end + widthB/2
		chunks = append(chunks, Chunk{Text: string(text[s:e]), IsSource: true}, Chunk{Text: "..."})
		end = widthA + widthB/2 + 6
	}
	return chunks, start, end
}

func substituteTabs(text []rune, start, end, tabWidth int) ([]rune, int, int) {
	out := make([]rune, 0, len(text)*2)
	newStart := start - 1
	newEnd := end - 1
	offset := 0
	for i, r := range text {
		if r == '\t' {
			for i := 0; i < tabWidth; i++ {
				out = append(out, ' ')
			}
			offset += tabWidth - 1
		} else {
			out = append(out, r)
		}
		if i == start-1 {
			newStart += offset + 1
		}
		if i == end-1 {
			newEnd += offset + 1
		}
	}
	return out, newStart, newEnd
}
