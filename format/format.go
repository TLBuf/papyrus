// Package format provides utilities for writing formatted Papyrus code.
package format

import (
	"bytes"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/token"
	"github.com/TLBuf/papyrus/types"
)

const (
	// DefaultIndentWidth is the default number of spaces
	// used per indentation level whentabs are not enabled.
	DefaultIndentWidth = 2
	// DefaultUseTabs is whether or not the
	// formatter uses tabs for indentations.
	DefaultUseTabs = false
	// DefaultUnixLineEndings is whether or not
	// UNIX-style line endings should be used.
	DefaultUnixLineEndings = false
)

// Option defines a format option.
type Option func(f *formatter) error

// IndentWidth returns an [Option] that sets the number of spaces used
// for each indentation level when spaces are used for intendation.
func WithIndentWidth(width int) Option {
	return func(f *formatter) error {
		if width < 0 {
			return fmt.Errorf("cannot set negative indent width: %d", width)
		}
		f.indentWidth = width
		return nil
	}
}

// WithTabs controls whether tabs (versus spaces) are used for indentation.
func WithTabs(tabs bool) Option {
	return func(f *formatter) error {
		f.useTabs = tabs
		return nil
	}
}

// WithUnixLineEndings controls whether line endings should be UNIX-style
// (versus Windows-style with carriage returns).
func WithUnixLineEndings(unix bool) Option {
	return func(f *formatter) error {
		f.unixLineEndings = unix
		return nil
	}
}

// WithKeywords controls what text the formatter uses for keywords.
//
// Any empty fields are ignored and the default value is used.
func WithKeywords(overrides Keywords) Option {
	return func(f *formatter) (err error) {
		f.keywords, err = keywords(overrides)
		return err
	}
}

// Format writes the formatted script.
func Format(w io.Writer, script *ast.Script, opts ...Option) error {
	f := &formatter{
		out:             w,
		indentWidth:     DefaultIndentWidth,
		useTabs:         DefaultUseTabs,
		unixLineEndings: DefaultUnixLineEndings,
		keywords:        defaultKeywords,
	}
	for _, opt := range opts {
		if err := opt(f); err != nil {
			return fmt.Errorf("invalid option: %w", err)
		}
	}
	return f.VisitScript(script)
}

type formatter struct {
	out             io.Writer
	indentWidth     int
	useTabs         bool
	unixLineEndings bool
	keywords        Keywords
	level           int
}

func (f *formatter) VisitAccess(node *ast.Access) error {
	if err := node.Value.Accept(f); err != nil {
		return fmt.Errorf("failed to format Value: %w", err)
	}
	if err := node.Operator.Accept(f); err != nil {
		return fmt.Errorf("failed to format Operator: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format Name: %w", err)
	}
	return nil
}

func (f *formatter) VisitArgument(node *ast.Argument) error {
	if node.Name != nil {
		if err := node.Name.Accept(f); err != nil {
			return fmt.Errorf("failed to format Name: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Operator.Accept(f); err != nil {
			return fmt.Errorf("failed to format Operator: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
	}
	if err := node.Value.Accept(f); err != nil {
		return fmt.Errorf("failed to format Value: %w", err)
	}
	return nil
}

func (f *formatter) VisitArrayCreation(node *ast.ArrayCreation) error {
	if err := node.New.Accept(f); err != nil {
		return fmt.Errorf("failed to format New: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Type.Accept(f); err != nil {
		return fmt.Errorf("failed to format Operator: %w", err)
	}
	if err := node.Open.Accept(f); err != nil {
		return fmt.Errorf("failed to format Open: %w", err)
	}
	if err := node.Size.Accept(f); err != nil {
		return fmt.Errorf("failed to format Size: %w", err)
	}
	if err := node.Close.Accept(f); err != nil {
		return fmt.Errorf("failed to format Close: %w", err)
	}
	return nil
}

func (f *formatter) VisitAssignment(node *ast.Assignment) error {
	if err := node.Assignee.Accept(f); err != nil {
		return fmt.Errorf("failed to format Assignee: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Operator.Accept(f); err != nil {
		return fmt.Errorf("failed to format Operator: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Value.Accept(f); err != nil {
		return fmt.Errorf("failed to format Value: %w", err)
	}
	return nil
}

func (f *formatter) VisitBinary(node *ast.Binary) error {
	if err := node.LeftOperand.Accept(f); err != nil {
		return fmt.Errorf("failed to format Operand: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Operator.Accept(f); err != nil {
		return fmt.Errorf("failed to format Operator: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.RightOperand.Accept(f); err != nil {
		return fmt.Errorf("failed to format Operand: %w", err)
	}
	return nil
}

func (f *formatter) VisitCall(node *ast.Call) error {
	if err := node.Function.Accept(f); err != nil {
		return fmt.Errorf("failed to format Function: %w", err)
	}
	if err := node.Open.Accept(f); err != nil {
		return fmt.Errorf("failed to format Open: %w", err)
	}
	for i, argument := range node.Arguments {
		if i > 0 {
			if err := f.str(token.Comma.String()); err != nil {
				return fmt.Errorf("failed to format comma: %w", err)
			}
			if err := f.space(); err != nil {
				return fmt.Errorf("failed to format space: %w", err)
			}
		}
		if err := argument.Accept(f); err != nil {
			return fmt.Errorf("failed to format Argument: %w", err)
		}
	}
	if err := node.Close.Accept(f); err != nil {
		return fmt.Errorf("failed to format Close: %w", err)
	}
	return nil
}

func (f *formatter) VisitCast(node *ast.Cast) error {
	if err := node.Value.Accept(f); err != nil {
		return fmt.Errorf("failed to format Value: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Operator.Accept(f); err != nil {
		return fmt.Errorf("failed to format Operator: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Type.Accept(f); err != nil {
		return fmt.Errorf("failed to format Type: %w", err)
	}
	return nil
}

func (f *formatter) VisitDocComment(node *ast.DocComment) error {
	if err := node.Open.Accept(f); err != nil {
		return fmt.Errorf("failed for format Open: %w", err)
	}
	text := trimBytes(node.Text.SourceLocation().Text())
	if bytes.ContainsRune(text, '\n') {
		f.level++
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		for i, line := range bytes.Split(text, []byte{'\n'}) {
			if err := f.bytes(trimBytes(line)); err != nil {
				return fmt.Errorf("failed to format comment text: %w", err)
			}
			if i != len(line)-1 {
				f.newline()
			}
		}
		f.level--
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	} else {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.bytes(trimBytes(text)); err != nil {
			return fmt.Errorf("failed to format comment text: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
	}
	if err := node.Close.Accept(f); err != nil {
		return fmt.Errorf("failed for format Close: %w", err)
	}
	return nil
}

func (f *formatter) VisitBlockComment(node *ast.BlockComment) error {
	if err := node.Open.Accept(f); err != nil {
		return fmt.Errorf("failed for format Open: %w", err)
	}
	text := trimBytes(node.Text.SourceLocation().Text())
	if bytes.ContainsRune(text, '\n') {
		f.level++
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		for i, line := range bytes.Split(text, []byte{'\n'}) {
			if err := f.bytes(trimBytes(line)); err != nil {
				return fmt.Errorf("failed to format comment text: %w", err)
			}
			if i != len(line)-1 {
				f.newline()
			}
		}
		f.level--
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	} else {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.bytes(trimBytes(text)); err != nil {
			return fmt.Errorf("failed to format comment text: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
	}
	if err := node.Close.Accept(f); err != nil {
		return fmt.Errorf("failed for format Close: %w", err)
	}
	return nil
}

func (f *formatter) VisitLineComment(node *ast.LineComment) error {
	if err := node.Open.Accept(f); err != nil {
		return fmt.Errorf("failed for format Open: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := f.bytes(trimBytes(node.Text.SourceLocation().Text())); err != nil {
		return fmt.Errorf("failed to format comment text: %w", err)
	}
	return nil
}

func (f *formatter) VisitEvent(node *ast.Event) error {
	if err := node.Keyword.Accept(f); err != nil {
		return fmt.Errorf("failed for format Keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed for format Name: %w", err)
	}
	if err := node.Open.Accept(f); err != nil {
		return fmt.Errorf("failed for format Open: %w", err)
	}
	for i, parameter := range node.Parameters {
		if i > 0 {
			if err := f.str(token.Comma.String()); err != nil {
				return fmt.Errorf("failed to format comma: %w", err)
			}
			if err := f.space(); err != nil {
				return fmt.Errorf("failed to format space: %w", err)
			}
		}
		if err := parameter.Accept(f); err != nil {
			return fmt.Errorf("failed to format Parameter: %w", err)
		}
	}
	if err := node.Close.Accept(f); err != nil {
		return fmt.Errorf("failed for format Close: %w", err)
	}
	if len(node.Native) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(token.Native.String()); err != nil {
			return fmt.Errorf("failed to format Native: %w", err)
		}
	}
	if node.Comment != nil {
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := node.Comment.Accept(f); err != nil {
			return fmt.Errorf("failed to format Comment: %w", err)
		}
	}
	if len(node.Native) == 0 {
		if node.Comment != nil {
			if err := f.newline(); err != nil {
				return fmt.Errorf("failed to format newline: %w", err)
			}
		}
		f.level++
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		for i, statement := range node.Statements {
			if i > 0 {
				if err := f.newline(); err != nil {
					return fmt.Errorf("failed to format newline: %w", err)
				}
			}
			if err := statement.Accept(f); err != nil {
				return fmt.Errorf("failed to format Statement: %w", err)
			}
		}
		f.level--
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := node.EndKeyword.Accept(f); err != nil {
			return fmt.Errorf("failed for format EndKeyword: %w", err)
		}
	}
	return nil
}

func (f *formatter) VisitFunction(node *ast.Function) error {
	if node.ReturnType != nil {
		if err := node.ReturnType.Accept(f); err != nil {
			return fmt.Errorf("failed for format ReturnType: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
	}
	if err := node.Keyword.Accept(f); err != nil {
		return fmt.Errorf("failed for format Keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed for format Name: %w", err)
	}
	if err := node.Open.Accept(f); err != nil {
		return fmt.Errorf("failed for format Open: %w", err)
	}
	for i, parameter := range node.Parameters {
		if i > 0 {
			if err := f.str(token.Comma.String()); err != nil {
				return fmt.Errorf("failed to format comma: %w", err)
			}
			if err := f.space(); err != nil {
				return fmt.Errorf("failed to format space: %w", err)
			}
		}
		if err := parameter.Accept(f); err != nil {
			return fmt.Errorf("failed to format Parameter: %w", err)
		}
	}
	if err := node.Close.Accept(f); err != nil {
		return fmt.Errorf("failed for format Close: %w", err)
	}
	if len(node.Global) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(token.Global.String()); err != nil {
			return fmt.Errorf("failed to format Global: %w", err)
		}
	}
	if len(node.Native) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(token.Native.String()); err != nil {
			return fmt.Errorf("failed to format Native: %w", err)
		}
	}
	if node.Comment != nil {
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := node.Comment.Accept(f); err != nil {
			return fmt.Errorf("failed to format Comment: %w", err)
		}
	}
	if len(node.Native) == 0 {
		if node.Comment != nil {
			if err := f.newline(); err != nil {
				return fmt.Errorf("failed to format newline: %w", err)
			}
		}
		f.level++
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		for i, statement := range node.Statements {
			if i > 0 {
				if err := f.newline(); err != nil {
					return fmt.Errorf("failed to format newline: %w", err)
				}
			}
			if err := statement.Accept(f); err != nil {
				return fmt.Errorf("failed to format Statement: %w", err)
			}
		}
		f.level--
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := node.EndKeyword.Accept(f); err != nil {
			return fmt.Errorf("failed for format EndKeyword: %w", err)
		}
	}
	return nil
}

func (f *formatter) VisitIdentifier(node *ast.Identifier) error {
	return f.bytes(node.Location.Text())
}

func (f *formatter) VisitIf(node *ast.If) error {
	if err := node.Keyword.Accept(f); err != nil {
		return fmt.Errorf("failed for format Keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Condition.Accept(f); err != nil {
		return fmt.Errorf("failed for format Condition: %w", err)
	}
	f.level++
	if err := f.newline(); err != nil {
		return fmt.Errorf("failed to format newline: %w", err)
	}
	for i, statement := range node.Statements {
		if i > 0 {
			if err := f.newline(); err != nil {
				return fmt.Errorf("failed to format newline: %w", err)
			}
		}
		if err := statement.Accept(f); err != nil {
			return fmt.Errorf("failed to format Statement: %w", err)
		}
	}
	f.level--
	if err := f.newline(); err != nil {
		return fmt.Errorf("failed to format newline: %w", err)
	}
	for _, elseif := range node.ElseIfs {
		if err := elseif.Accept(f); err != nil {
			return fmt.Errorf("failed to format ElseIf: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	if node.Else != nil {
		if err := node.Else.Accept(f); err != nil {
			return fmt.Errorf("failed to format Else: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	if err := node.EndKeyword.Accept(f); err != nil {
		return fmt.Errorf("failed for format EndKeyword: %w", err)
	}
	return nil
}

func (f *formatter) VisitElseIf(node *ast.ElseIf) error {
	if err := node.Keyword.Accept(f); err != nil {
		return fmt.Errorf("failed for format Keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Condition.Accept(f); err != nil {
		return fmt.Errorf("failed for format Condition: %w", err)
	}
	f.level++
	if err := f.newline(); err != nil {
		return fmt.Errorf("failed to format newline: %w", err)
	}
	for i, statement := range node.Statements {
		if i > 0 {
			if err := f.newline(); err != nil {
				return fmt.Errorf("failed to format newline: %w", err)
			}
		}
		if err := statement.Accept(f); err != nil {
			return fmt.Errorf("failed to format Statement: %w", err)
		}
	}
	f.level--
	return nil
}

func (f *formatter) VisitElse(node *ast.Else) error {
	if err := node.Keyword.Accept(f); err != nil {
		return fmt.Errorf("failed for format Keyword: %w", err)
	}
	f.level++
	if err := f.newline(); err != nil {
		return fmt.Errorf("failed to format newline: %w", err)
	}
	for i, statement := range node.Statements {
		if i > 0 {
			if err := f.newline(); err != nil {
				return fmt.Errorf("failed to format newline: %w", err)
			}
		}
		if err := statement.Accept(f); err != nil {
			return fmt.Errorf("failed to format Statement: %w", err)
		}
	}
	f.level--
	return nil
}

func (f *formatter) VisitImport(node *ast.Import) error {
	if err := node.Keyword.Accept(f); err != nil {
		return fmt.Errorf("failed for format Keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed for format Name: %w", err)
	}
	return nil
}

func (f *formatter) VisitIndex(node *ast.Index) error {
	if err := node.Value.Accept(f); err != nil {
		return fmt.Errorf("failed for format Value: %w", err)
	}
	if err := node.Open.Accept(f); err != nil {
		return fmt.Errorf("failed for format Open: %w", err)
	}
	if err := node.Index.Accept(f); err != nil {
		return fmt.Errorf("failed for format Index: %w", err)
	}
	if err := node.Close.Accept(f); err != nil {
		return fmt.Errorf("failed for format Close: %w", err)
	}
	return nil
}

func (f *formatter) VisitBoolLiteral(node *ast.BoolLiteral) error {
	text := f.keywords.False
	if node.Value {
		text = f.keywords.True
	}
	if err := f.str(text); err != nil {
		return fmt.Errorf("failed to format BoolLiteral: %w", err)
	}
	return nil
}

func (f *formatter) VisitIntLiteral(node *ast.IntLiteral) error {
	if err := f.bytes(node.Location.Text()); err != nil {
		return fmt.Errorf("failed to format IntLiteral: %w", err)
	}
	return nil
}

func (f *formatter) VisitFloatLiteral(node *ast.FloatLiteral) error {
	if err := f.bytes(node.Location.Text()); err != nil {
		return fmt.Errorf("failed to format FloatLiteral: %w", err)
	}
	return nil
}

func (f *formatter) VisitStringLiteral(node *ast.StringLiteral) error {
	if err := f.bytes(node.Location.Text()); err != nil {
		return fmt.Errorf("failed to format StringLiteral: %w", err)
	}
	return nil
}

func (f *formatter) VisitNoneLiteral(node *ast.NoneLiteral) error {
	if err := f.str(f.keywords.None); err != nil {
		return fmt.Errorf("failed to format None: %w", err)
	}
	return nil
}

func (f *formatter) VisitParameter(node *ast.Parameter) error {
	if err := node.Type.Accept(f); err != nil {
		return fmt.Errorf("failed to format Type: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format Name: %w", err)
	}
	if node.Operator != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Operator.Accept(f); err != nil {
			return fmt.Errorf("failed to format Operator: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Value.Accept(f); err != nil {
			return fmt.Errorf("failed to format Value: %w", err)
		}
	}
	return nil
}

func (f *formatter) VisitParenthetical(node *ast.Parenthetical) error {
	if err := node.Open.Accept(f); err != nil {
		return fmt.Errorf("failed to format Open: %w", err)
	}
	if err := node.Value.Accept(f); err != nil {
		return fmt.Errorf("failed to format Value: %w", err)
	}
	if err := node.Close.Accept(f); err != nil {
		return fmt.Errorf("failed to format Close: %w", err)
	}
	return nil
}

func (f *formatter) VisitProperty(node *ast.Property) error {
	if err := node.Type.Accept(f); err != nil {
		return fmt.Errorf("failed to format Type: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Keyword.Accept(f); err != nil {
		return fmt.Errorf("failed to format Keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format Name: %w", err)
	}
	if node.Operator != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Operator.Accept(f); err != nil {
			return fmt.Errorf("failed to format Operator: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Value.Accept(f); err != nil {
			return fmt.Errorf("failed to format Value: %w", err)
		}
	}
	if node.Auto != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Auto.Accept(f); err != nil {
			return fmt.Errorf("failed to format Auto: %w", err)
		}
	}
	if node.AutoReadOnly != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.AutoReadOnly.Accept(f); err != nil {
			return fmt.Errorf("failed to format AutoReadOnly: %w", err)
		}
	}
	if len(node.Hidden) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.Hidden); err != nil {
			return fmt.Errorf("failed to format Hidden: %w", err)
		}
	}
	if len(node.Conditional) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.Conditional); err != nil {
			return fmt.Errorf("failed to format Conditional: %w", err)
		}
	}
	if node.Comment != nil {
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := node.Comment.Accept(f); err != nil {
			return fmt.Errorf("failed to format Comment: %w", err)
		}
	}
	if node.Get != nil || node.Set != nil {
		if node.Comment != nil {
			if err := f.newline(); err != nil {
				return fmt.Errorf("failed to format newline: %w", err)
			}
		}
		f.level++
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if node.Get != nil {
			if err := node.Get.Accept(f); err != nil {
				return fmt.Errorf("failed to format Get: %w", err)
			}
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if node.Set != nil {
			if err := node.Get.Accept(f); err != nil {
				return fmt.Errorf("failed to format Set: %w", err)
			}
		}
		f.level--
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := node.EndKeyword.Accept(f); err != nil {
			return fmt.Errorf("failed to format EndKeyword: %w", err)
		}
	}
	return nil
}

func (f *formatter) VisitReturn(node *ast.Return) error {
	if err := node.Keyword.Accept(f); err != nil {
		return fmt.Errorf("failed to format Keyword: %w", err)
	}
	if node.Value != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Value.Accept(f); err != nil {
			return fmt.Errorf("failed to format Value: %w", err)
		}
	}
	return nil
}

func (f *formatter) VisitScript(node *ast.Script) error {
	if err := node.Keyword.Accept(f); err != nil {
		return fmt.Errorf("failed to format Keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format Name: %w", err)
	}
	if node.Extends != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Extends.Accept(f); err != nil {
			return fmt.Errorf("failed to format Extends: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Parent.Accept(f); err != nil {
			return fmt.Errorf("failed to format Parent: %w", err)
		}
	}
	if len(node.Hidden) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.Hidden); err != nil {
			return fmt.Errorf("failed to format Hidden: %w", err)
		}
	}
	if len(node.Conditional) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.Conditional); err != nil {
			return fmt.Errorf("failed to format Conditional: %w", err)
		}
	}
	if err := f.newline(); err != nil {
		return fmt.Errorf("failed to format newline: %w", err)
	}
	if node.Comment != nil {
		if err := node.Comment.Accept(f); err != nil {
			return fmt.Errorf("failed to format Comment: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	if err := f.newline(); err != nil {
		return fmt.Errorf("failed to format newline: %w", err)
	}
	// Extract statements and prepare for formatting each one.
	var imports []*ast.Import
	var properties []*ast.Property
	var variables []*ast.ScriptVariable
	var states []*ast.State
	var invokables []ast.Invokable
	for _, stmt := range node.Statements {
		switch stmt := stmt.(type) {
		case *ast.Import:
			imports = append(imports, stmt)
		case *ast.Property:
			properties = append(properties, stmt)
		case *ast.ScriptVariable:
			variables = append(variables, stmt)
		case *ast.State:
			states = append(states, stmt)
		case ast.Invokable:
			invokables = append(invokables, stmt)
		default:
			panic(fmt.Errorf("unknown script statement type: %T", stmt))
		}
	}
	slices.SortFunc(imports, func(a, b *ast.Import) int { return strings.Compare(a.Name.Normalized, b.Name.Normalized) })
	for i, state := range states {
		// Always pull the Auto state to the top, but otherwise maintain order.
		if state.Auto != nil {
			if i == 0 {
				break
			}
			s := states[0]
			states[0] = states[i]
			states[i] = s
			break
		}
	}
	for _, i := range imports {
		if err := i.Accept(f); err != nil {
			return fmt.Errorf("failed to format Import: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	if len(imports) > 0 {
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	for _, p := range properties {
		if err := p.Accept(f); err != nil {
			return fmt.Errorf("failed to format Property: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	if len(properties) > 0 {
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	for _, v := range variables {
		if err := v.Accept(f); err != nil {
			return fmt.Errorf("failed to format ScriptVariable: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	if len(variables) > 0 {
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	for _, s := range states {
		if err := s.Accept(f); err != nil {
			return fmt.Errorf("failed to format State: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	for _, i := range invokables {
		if err := i.Accept(f); err != nil {
			return fmt.Errorf("failed to format Invokable: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	return nil
}

func (f *formatter) VisitState(node *ast.State) error {
	if node.Auto != nil {
		if err := node.Auto.Accept(f); err != nil {
			return fmt.Errorf("failed to format Auto: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
	}
	if err := node.Keyword.Accept(f); err != nil {
		return fmt.Errorf("failed to format Keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format Name: %w", err)
	}
	f.level++
	if err := f.newline(); err != nil {
		return fmt.Errorf("failed to format newline: %w", err)
	}
	for i, invokable := range node.Invokables {
		if i > 0 {
			if err := f.newline(); err != nil {
				return fmt.Errorf("failed to format newline: %w", err)
			}
			if err := f.newline(); err != nil {
				return fmt.Errorf("failed to format newline: %w", err)
			}
		}
		if err := invokable.Accept(f); err != nil {
			return fmt.Errorf("failed to format Name: %w", err)
		}
	}
	f.level--
	if err := f.newline(); err != nil {
		return fmt.Errorf("failed to format newline: %w", err)
	}
	if err := node.EndKeyword.Accept(f); err != nil {
		return fmt.Errorf("failed to format EndKeyword: %w", err)
	}
	return nil
}

func (f *formatter) VisitToken(node *ast.Token) error {
	if node.Kind.IsSymbol() {
		if err := f.str(node.Kind.Symbol()); err != nil {
			return fmt.Errorf("failed to format symbol token: %w", err)
		}
		return nil
	}
	if node.Kind.IsKeyword() {
		if err := f.str(f.keywords.Text(node.Kind)); err != nil {
			return fmt.Errorf("failed to format keyword token: %w", err)
		}
		return nil
	}
	return fmt.Errorf("unexpected token: %s", node.Kind)
}

func (f *formatter) VisitTypeLiteral(node *ast.TypeLiteral) error {
	baseType := node.Type
	if arrayType, ok := baseType.(types.Array); ok {
		baseType = arrayType.ElementType
	}
	text := ""
	switch baseType.(type) {
	case types.Bool:
		text = f.keywords.Bool
	case types.Int:
		text = f.keywords.Int
	case types.Float:
		text = f.keywords.Float
	case types.String:
		text = f.keywords.String
	case types.Object:
		text = string(node.Text.Location.Text())
	default:
		return fmt.Errorf("unexpected types: %T", baseType)
	}
	if err := f.str(text); err != nil {
		return fmt.Errorf("failed to format type: %w", err)
	}
	if node.Open != nil {
		if err := node.Open.Accept(f); err != nil {
			return fmt.Errorf("failed to format Open: %w", err)
		}
		if err := node.Close.Accept(f); err != nil {
			return fmt.Errorf("failed to format Close: %w", err)
		}
	}
	return nil
}

func (f *formatter) VisitUnary(node *ast.Unary) error {
	if err := node.Operator.Accept(f); err != nil {
		return fmt.Errorf("failed to format Operator: %w", err)
	}
	if err := node.Operand.Accept(f); err != nil {
		return fmt.Errorf("failed to format Operand: %w", err)
	}
	return nil
}

func (f *formatter) VisitScriptVariable(node *ast.ScriptVariable) error {
	if err := node.Type.Accept(f); err != nil {
		return fmt.Errorf("failed to format Type: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format Name: %w", err)
	}
	if node.Value != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Operator.Accept(f); err != nil {
			return fmt.Errorf("failed to format Operator: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Value.Accept(f); err != nil {
			return fmt.Errorf("failed to format Value: %w", err)
		}
	}
	if len(node.Conditional) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.Conditional); err != nil {
			return fmt.Errorf("failed to format Conditional: %w", err)
		}
	}
	return nil
}

func (f *formatter) VisitFunctionVariable(node *ast.FunctionVariable) error {
	if err := node.Type.Accept(f); err != nil {
		return fmt.Errorf("failed to format Type: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format Name: %w", err)
	}
	if node.Value != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Operator.Accept(f); err != nil {
			return fmt.Errorf("failed to format Operator: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Value.Accept(f); err != nil {
			return fmt.Errorf("failed to format Value: %w", err)
		}
	}
	return nil
}

func (f *formatter) VisitWhile(node *ast.While) error {
	return nil
}

func (f *formatter) VisitErrorScriptStatement(node *ast.ErrorScriptStatement) error {
	return fmt.Errorf("attempted to format error script statement: %s", node.Message)
}

func (f *formatter) VisitErrorFunctionStatement(node *ast.ErrorFunctionStatement) error {
	return fmt.Errorf("attempted to format error function statement: %s", node.Message)
}

var _ ast.Visitor = (*formatter)(nil)

func (f *formatter) space() error {
	return f.str(" ")
}

func (f *formatter) newline() error {
	if f.unixLineEndings {
		return f.str("\n")
	}
	if err := f.str("\r\n"); err != nil {
		return err
	}
	return f.indent()
}

func (f *formatter) indent() error {
	if f.level == 0 {
		return nil
	}
	if f.useTabs {
		return f.str(strings.Repeat("\t", f.level))
	}
	return f.str(strings.Repeat(" ", f.indentWidth*f.level))
}

func (f *formatter) str(text string) error {
	_, err := io.WriteString(f.out, text)
	return err
}

func (f *formatter) bytes(text []byte) error {
	_, err := f.out.Write(text)
	return err
}

func trimBytes(text []byte) []byte {
	return bytes.TrimSpace(bytes.TrimRight(text, "\r\n"))
}
