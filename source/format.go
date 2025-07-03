package source

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	up      = []byte{'\xE2', '\x96', '\xB2'}
	dot     = []byte{'\xC2', '\xB7'}
	down    = []byte{'\xE2', '\x96', '\xBC'}
	newline = []byte{'\n'}
)

// Format emits the snippet to in text format.
func Format(w io.Writer, snippet Snippet) error {
	if len(snippet.Lines) > 1 {
		return printMultiLine(w, snippet.Lines, snippet.Start.Column, snippet.End.Column, snippet.Width)
	}
	return printSingleLine(w, snippet.Lines[0], snippet.Start.Column, snippet.End.Column)
}

func printSingleLine(w io.Writer, line Line, start, end int) error {
	numWidth := utf8.RuneCountInString(strconv.Itoa(line.Number))
	prefixWidth := numWidth + 4 // " 1 | "
	if _, err := w.Write(bytes.Repeat([]byte{' '}, start+prefixWidth-1)); err != nil {
		return err
	}
	if _, err := w.Write(down); err != nil {
		return err
	}
	if end-start > 1 {
		if _, err := w.Write(bytes.Repeat(dot, end-start-1)); err != nil {
			return err
		}
	}
	if end > start {
		if _, err := w.Write(down); err != nil {
			return err
		}
	}
	if _, err := w.Write(newline); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, " %d | ", line.Number); err != nil {
		return err
	}
	for _, c := range line.Chunks {
		if _, err := io.WriteString(w, c.Text); err != nil {
			return err
		}
	}
	if _, err := w.Write(newline); err != nil {
		return err
	}
	return nil
}

func printMultiLine(w io.Writer, lines []Line, start, end, width int) error {
	numWidth := utf8.RuneCountInString(strconv.Itoa(lines[len(lines)-1].Number))
	prefixWidth := numWidth + 4 // " 1 | "
	if _, err := w.Write(bytes.Repeat([]byte{' '}, start+prefixWidth-1)); err != nil {
		return err
	}
	if _, err := w.Write(down); err != nil {
		return err
	}
	if _, err := w.Write(bytes.Repeat(dot, width-start)); err != nil {
		return err
	}
	if _, err := w.Write(newline); err != nil {
		return err
	}
	for _, line := range lines {
		if line.Number > 0 {
			if _, err := fmt.Fprintf(w, " %0*d | ", numWidth, line.Number); err != nil {
				return err
			}
		} else {
			if _, err := fmt.Fprintf(w, " %s | ", strings.Repeat("-", numWidth)); err != nil {
				return err
			}
		}
		for _, c := range line.Chunks {
			if _, err := io.WriteString(w, c.Text); err != nil {
				return err
			}
		}
		if _, err := w.Write(newline); err != nil {
			return err
		}
	}
	if _, err := w.Write(bytes.Repeat(dot, end-1+prefixWidth)); err != nil {
		return err
	}
	if _, err := w.Write(up); err != nil {
		return err
	}
	if _, err := w.Write(newline); err != nil {
		return err
	}
	return nil
}
