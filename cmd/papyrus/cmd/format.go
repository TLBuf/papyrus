package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TLBuf/papyrus/format"
	"github.com/TLBuf/papyrus/parser"
	"github.com/TLBuf/papyrus/source"
	"github.com/spf13/cobra"
)

var (
	formatTabs   bool
	formatUnix   bool
	formatIndent int
)

// Format returns a command that formats one or more Papyrus files.
func Format() *cobra.Command {
	command := &cobra.Command{
		Use:          "format [path...]",
		Short:        "Formats one or more Papyrus files",
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, args []string) error {
			return formatFiles(args...)
		},
	}

	command.Flags().BoolVarP(&formatTabs, "tabs", "t", false, "Controls whether tabs or spaces are used for indentation")
	command.Flags().BoolVarP(
		&formatUnix,
		"unix",
		"u",
		false,
		"Controls whether unix-style (vs Windows) line ending are used",
	)
	command.Flags().IntVarP(
		&formatIndent,
		"indent",
		"i",
		2,
		"Controls the number of spaces used for each indentation level",
	)
	command.MarkFlagsMutuallyExclusive("tabs", "indent")

	return command
}

func formatFiles(paths ...string) error {
	failed := 0
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		if err := formatFile(path); err != nil {
			failed++
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
	}
	if failed > 0 {
		return fmt.Errorf("failed to format %d file(s)", failed)
	}
	return nil
}

func formatFile(path string) error {
	text, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("failed to read %q: %w", path, err)
	}
	file := source.File{
		Path: path,
		Text: text,
	}
	script, err := parser.Parse(file, parser.WithComments(true))
	if err != nil {
		var perr parser.Error
		if !errors.As(err, &perr) {
			return fmt.Errorf("failed to extract a parser.Error from: %w", err)
		}
		snip, serr := perr.Location.Snippet(file, 80, 9)
		if serr != nil {
			return fmt.Errorf("failed to create snippet for parser error: %w: %w", serr, err)
		}
		var sb strings.Builder
		if err := source.Format(&sb, snip); err != nil {
			return fmt.Errorf("failed to format snippet: %w", err)
		}
		return fmt.Errorf("failed to parse: %w\n\n%s %s\n%s", perr.Err, file.Path, perr.Location, sb.String())
	}
	var formatted bytes.Buffer
	if err := format.Format(&formatted, script, format.WithTabs(formatTabs), format.WithUnixLineEndings(formatUnix), format.WithIndentWidth(formatIndent)); err != nil {
		return fmt.Errorf("failed to format %q: %w", path, err)
	}
	if err := os.WriteFile(path, formatted.Bytes(), 0o600); err != nil {
		return fmt.Errorf("failed to write %q: %w", path, err)
	}
	return nil
}
