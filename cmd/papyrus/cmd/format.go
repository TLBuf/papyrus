package cmd

import (
	"bytes"
	"fmt"
	"os"
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

func Format() *cobra.Command {
	command := &cobra.Command{
		Use:          "format [path...]",
		Short:        "Formats one or more Papyrus files",
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return formatFiles(args...)
		},
	}

	command.Flags().BoolVarP(&formatTabs, "tabs", "t", false, "Controls whether tabs or spaces are used for indentation")
	command.Flags().BoolVarP(&formatUnix, "unix", "u", false, "Controls whether unix-style (vs Windows) line ending are used")
	command.Flags().IntVarP(&formatIndent, "indent", "i", 2, "Controls the number of spaces used for each indentation level")
	command.MarkFlagsMutuallyExclusive("tabs", "indent")

	return command
}

func formatFiles(paths ...string) error {
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		if err := formatFile(path); err != nil {
			return err
		}
	}
	return nil
}

func formatFile(path string) error {
	text, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read %q: %w", path, err)
	}
	file := &source.File{
		Path: path,
		Text: text,
	}
	script, err := parser.Parse(file)
	if err != nil {
		return fmt.Errorf("failed to parse %q: %w", path, err)
	}
	var formatted bytes.Buffer
	if err := format.Format(&formatted, script, format.WithTabs(formatTabs), format.WithUnixLineEndings(formatUnix), format.WithIndentWidth(formatIndent)); err != nil {
		return fmt.Errorf("failed to format %q: %w", path, err)
	}
	if err := os.WriteFile(path, formatted.Bytes(), 0600); err != nil {
		return fmt.Errorf("failed to write %q: %w", path, err)
	}
	return nil
}
