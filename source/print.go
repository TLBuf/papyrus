package source

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	up   = "\xE2\x96\xB2"
	dot  = "\xC2\xB7"
	down = "\xE2\x96\xBC"
)

// Print emits the snippet to in text format.
func Print(w io.Writer, snippet Snippet) {
	if len(snippet.Lines) > 1 {
		printMultiLine(w, snippet.Lines, snippet.Start.Column, snippet.End.Column, snippet.Width)
	} else {
		printSingleLine(w, snippet.Lines[0], snippet.Start.Column, snippet.End.Column)
	}
}

func printSingleLine(w io.Writer, line Line, start, end int) {
	numWidth := utf8.RuneCountInString(strconv.Itoa(line.Number))
	prefixWidth := numWidth + 4 // " 1 | "
	fmt.Fprintf(w, "%s%s", strings.Repeat(" ", start-1+prefixWidth), down)
	if end-start > 1 {
		fmt.Fprintf(w, "%s", strings.Repeat(dot, end-start-1))
	}
	if end > start {
		fmt.Fprint(w, down)
	}
	fmt.Fprint(w, "\n")
	fmt.Fprintf(w, " %d | ", line.Number)
	for _, c := range line.Chunks {
		fmt.Fprint(w, c.Text)
	}
	fmt.Fprint(w, "\n")
}

func printMultiLine(w io.Writer, lines []Line, start, end, width int) {
	numWidth := utf8.RuneCountInString(strconv.Itoa(lines[len(lines)-1].Number))
	prefixWidth := numWidth + 4 // " 1 | "
	fmt.Fprintf(w, "%s%s%s\n", strings.Repeat(" ", start-1+prefixWidth), down, strings.Repeat(dot, width-start))
	for _, line := range lines {
		if line.Number > 0 {
			fmt.Fprintf(w, " %0*d | ", numWidth, line.Number)
		} else {
			fmt.Fprintf(w, " %s | ", strings.Repeat("-", numWidth))
		}
		for _, c := range line.Chunks {
			fmt.Fprint(w, c.Text)
		}
		fmt.Fprint(w, "\n")
	}
	fmt.Fprintf(w, "%s%s\n", strings.Repeat(dot, end-1+prefixWidth), up)
}
