// Package format provides utilities for writing formatted Papyrus code.
package format

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/source"
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

var (
	fragmentHeader = []byte("BEGIN FRAGMENT CODE")
	fragmentFooter = []byte("END FRAGMENT CODE")
)

// Option defines a format option.
type Option func(f *formatter) error

// WithIndentWidth returns an [Option] that sets the number of spaces used
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
		file:            script.File,
		out:             w,
		indentWidth:     DefaultIndentWidth,
		useTabs:         DefaultUseTabs,
		unixLineEndings: DefaultUnixLineEndings,
		keywords:        defaultKeywords,
		level:           0,
	}
	for _, opt := range opts {
		if err := opt(f); err != nil {
			return fmt.Errorf("invalid option: %w", err)
		}
	}
	return f.VisitScript(script)
}

type formatter struct {
	file            source.File
	out             io.Writer
	indentWidth     int
	useTabs         bool
	unixLineEndings bool
	keywords        Keywords
	level           int
}

func (f *formatter) visitPrefixComments(node interface{ Comments() *ast.Comments }) error {
	if node.Comments() == nil {
		return nil
	}
	for _, c := range node.Comments().PrefixComments {
		if err := c.Accept(f); err != nil {
			return fmt.Errorf("failed to format prefix comment: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
	}
	return nil
}

func (f *formatter) visitSuffixComments(node interface{ Comments() *ast.Comments }) error {
	if node.Comments() == nil {
		return nil
	}
	for _, c := range node.Comments().SuffixComments {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := c.Accept(f); err != nil {
			return fmt.Errorf("failed to format suffix comment: %w", err)
		}
	}
	return nil
}

func (f *formatter) VisitAccess(node *ast.Access) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := node.Value.Accept(f); err != nil {
		return fmt.Errorf("failed to format value: %w", err)
	}
	if err := f.str(token.Dot.Symbol()); err != nil {
		return fmt.Errorf("failed to format dot operator: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format name: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitArgument(node *ast.Argument) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if node.Name != nil {
		if err := node.Name.Accept(f); err != nil {
			return fmt.Errorf("failed to format name: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(token.Assign.Symbol()); err != nil {
			return fmt.Errorf("failed to format assign operator: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
	}
	if err := node.Value.Accept(f); err != nil {
		return fmt.Errorf("failed to format value: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitArrayCreation(node *ast.ArrayCreation) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.str(f.keywords.New); err != nil {
		return fmt.Errorf("failed to format new operator: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Type.Accept(f); err != nil {
		return fmt.Errorf("failed to format type: %w", err)
	}
	if err := f.str(token.BracketOpen.Symbol()); err != nil {
		return fmt.Errorf("failed to format open bracket: %w", err)
	}
	if err := node.Size.Accept(f); err != nil {
		return fmt.Errorf("failed to format size: %w", err)
	}
	if err := f.str(token.BracketClose.Symbol()); err != nil {
		return fmt.Errorf("failed to format close bracket: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitAssignment(node *ast.Assignment) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := node.Assignee.Accept(f); err != nil {
		return fmt.Errorf("failed to format assignee: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := f.str(node.Kind.Symbol()); err != nil {
		return fmt.Errorf("failed to format operator: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Value.Accept(f); err != nil {
		return fmt.Errorf("failed to format value: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitBinary(node *ast.Binary) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := node.LeftOperand.Accept(f); err != nil {
		return fmt.Errorf("failed to format left operand: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := f.str(node.Kind.Symbol()); err != nil {
		return fmt.Errorf("failed to format operator: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.RightOperand.Accept(f); err != nil {
		return fmt.Errorf("failed to format right operand: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitCall(node *ast.Call) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := node.Function.Accept(f); err != nil {
		return fmt.Errorf("failed to format function: %w", err)
	}
	if err := f.str(token.ParenthesisOpen.Symbol()); err != nil {
		return fmt.Errorf("failed to format open parenthesis: %w", err)
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
			return fmt.Errorf("failed to format argument: %w", err)
		}
	}
	if err := f.str(token.ParenthesisClose.Symbol()); err != nil {
		return fmt.Errorf("failed to format close parenthesis: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitCast(node *ast.Cast) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := node.Value.Accept(f); err != nil {
		return fmt.Errorf("failed to format value: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := f.str(f.keywords.As); err != nil {
		return fmt.Errorf("failed to format as operator: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Type.Accept(f); err != nil {
		return fmt.Errorf("failed to format type: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitDocumentation(node *ast.Documentation) error {
	if err := f.str(token.BraceOpen.Symbol()); err != nil {
		return fmt.Errorf("failed for format open brace: %w", err)
	}
	text := bytes.TrimSpace(node.TextLocation.Text(f.file))
	if bytes.ContainsRune(text, '\n') {
		f.level++
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		for i, line := range bytes.Split(text, []byte{'\n'}) {
			if i > 0 {
				if err := f.newline(); err != nil {
					return fmt.Errorf("failed to format newline: %w", err)
				}
			}
			if err := f.bytes(trimRight(line)); err != nil {
				return fmt.Errorf("failed to format comment text: %w", err)
			}
		}
		f.level--
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	} else {
		if err := f.bytes(bytes.TrimSpace(text)); err != nil {
			return fmt.Errorf("failed to format comment text: %w", err)
		}
	}
	if err := f.str(token.BraceClose.Symbol()); err != nil {
		return fmt.Errorf("failed for format close brace: %w", err)
	}
	return nil
}

func (f *formatter) VisitBlockComment(node *ast.BlockComment) error {
	if node.HasLeadingBlankLine {
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	if err := f.str(token.BlockCommentOpen.Symbol()); err != nil {
		return fmt.Errorf("failed for format block comment open: %w", err)
	}
	text := trimRight(node.TextLocation.Text(f.file))
	if bytes.ContainsRune(text, '\n') {
		f.level++
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		for i, line := range bytes.Split(text, []byte{'\n'}) {
			if i > 0 {
				if err := f.newline(); err != nil {
					return fmt.Errorf("failed to format newline: %w", err)
				}
			}
			if err := f.bytes(trimRight(line)); err != nil {
				return fmt.Errorf("failed to format comment text: %w", err)
			}
		}
		f.level--
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	} else {
		text = bytes.TrimSpace(text)
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.bytes(text); err != nil {
			return fmt.Errorf("failed to format comment text: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
	}
	if err := f.str(token.BlockCommentClose.Symbol()); err != nil {
		return fmt.Errorf("failed for format block comment close: %w", err)
	}
	return nil
}

func (f *formatter) VisitCommentStatement(node *ast.CommentStatement) error {
	for i, e := range node.Elements {
		if i > 0 {
			if err := f.newline(); err != nil {
				return fmt.Errorf("failed to format newline: %w", err)
			}
			if e.LeadingBlankLine() {
				if err := f.newline(); err != nil {
					return fmt.Errorf("failed to format newline: %w", err)
				}
			}
		}
		if err := e.Accept(f); err != nil {
			return fmt.Errorf("failed to format comment element: %w", err)
		}
	}
	return nil
}

func (f *formatter) VisitLineComment(node *ast.LineComment) error {
	if err := f.str(token.Semicolon.Symbol()); err != nil {
		return fmt.Errorf("failed for format semicolon: %w", err)
	}
	text := trimRight(node.TextLocation.Text(f.file))
	if bytes.HasPrefix(text, fragmentHeader) {
		f.level--
	}
	if bytes.HasPrefix(text, fragmentFooter) {
		f.level++
	}
	if err := f.bytes(text); err != nil {
		return fmt.Errorf("failed to format comment text: %w", err)
	}
	return nil
}

func (f *formatter) VisitEvent(node *ast.Event) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.str(f.keywords.Event); err != nil {
		return fmt.Errorf("failed for format start keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed for format name: %w", err)
	}
	if err := f.str(token.ParenthesisOpen.Symbol()); err != nil {
		return fmt.Errorf("failed for format open parenthesis: %w", err)
	}
	for i, parameter := range node.ParameterList {
		if i > 0 {
			if err := f.str(token.Comma.String()); err != nil {
				return fmt.Errorf("failed to format comma: %w", err)
			}
			if err := f.space(); err != nil {
				return fmt.Errorf("failed to format space: %w", err)
			}
		}
		if err := parameter.Accept(f); err != nil {
			return fmt.Errorf("failed to format parameter: %w", err)
		}
	}
	if err := f.str(token.ParenthesisClose.Symbol()); err != nil {
		return fmt.Errorf("failed for format close parenthesis: %w", err)
	}
	if len(node.NativeLocations) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(token.Native.String()); err != nil {
			return fmt.Errorf("failed to format native keyword: %w", err)
		}
	}
	if node.Documentation != nil {
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := node.Documentation.Accept(f); err != nil {
			return fmt.Errorf("failed to format documentation: %w", err)
		}
	}
	if len(node.NativeLocations) == 0 {
		f.level++
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		for i, statement := range node.Statements {
			if i > 0 {
				if err := f.newline(); err != nil {
					return fmt.Errorf("failed to format newline: %w", err)
				}
				if statement.LeadingBlankLine() {
					if err := f.newline(); err != nil {
						return fmt.Errorf("failed to format newline: %w", err)
					}
				}
			}
			if err := statement.Accept(f); err != nil {
				return fmt.Errorf("failed to format statement: %w", err)
			}
		}
		f.level--
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := f.str(f.keywords.EndEvent); err != nil {
			return fmt.Errorf("failed for format end keyword: %w", err)
		}
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitFunction(node *ast.Function) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if node.ReturnType != nil {
		if err := node.ReturnType.Accept(f); err != nil {
			return fmt.Errorf("failed for format ReturnType: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
	}
	if err := f.str(f.keywords.Function); err != nil {
		return fmt.Errorf("failed for format start keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed for format name: %w", err)
	}
	if err := f.str(token.ParenthesisOpen.Symbol()); err != nil {
		return fmt.Errorf("failed for format open parenthesis: %w", err)
	}
	for i, parameter := range node.ParameterList {
		if i > 0 {
			if err := f.str(token.Comma.String()); err != nil {
				return fmt.Errorf("failed to format comma: %w", err)
			}
			if err := f.space(); err != nil {
				return fmt.Errorf("failed to format space: %w", err)
			}
		}
		if err := parameter.Accept(f); err != nil {
			return fmt.Errorf("failed to format parameter: %w", err)
		}
	}
	if err := f.str(token.ParenthesisClose.Symbol()); err != nil {
		return fmt.Errorf("failed for format close parenthesis: %w", err)
	}
	if len(node.GlobalLocations) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(token.Global.String()); err != nil {
			return fmt.Errorf("failed to format global keyword: %w", err)
		}
	}
	if len(node.NativeLocations) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(token.Native.String()); err != nil {
			return fmt.Errorf("failed to format native keyword: %w", err)
		}
	}
	if node.Documentation != nil {
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := node.Documentation.Accept(f); err != nil {
			return fmt.Errorf("failed to format documentation: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	if len(node.NativeLocations) == 0 {
		f.level++
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		for i, statement := range node.Statements {
			if i > 0 {
				if err := f.newline(); err != nil {
					return fmt.Errorf("failed to format newline: %w", err)
				}
				if statement.LeadingBlankLine() {
					if err := f.newline(); err != nil {
						return fmt.Errorf("failed to format newline: %w", err)
					}
				}
			}
			if err := statement.Accept(f); err != nil {
				return fmt.Errorf("failed to format statement: %w", err)
			}
		}
		f.level--
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := f.str(f.keywords.EndFunction); err != nil {
			return fmt.Errorf("failed for format end keyword: %w", err)
		}
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitIdentifier(node *ast.Identifier) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	return f.bytes(node.Location().Text(f.file))
}

func (f *formatter) VisitIf(node *ast.If) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.str(f.keywords.If); err != nil {
		return fmt.Errorf("failed for format start keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Condition.Accept(f); err != nil {
		return fmt.Errorf("failed for format condition: %w", err)
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
			if statement.LeadingBlankLine() {
				if err := f.newline(); err != nil {
					return fmt.Errorf("failed to format newline: %w", err)
				}
			}
		}
		if err := statement.Accept(f); err != nil {
			return fmt.Errorf("failed to format statement: %w", err)
		}
	}
	f.level--
	if err := f.newline(); err != nil {
		return fmt.Errorf("failed to format newline: %w", err)
	}
	for _, elseif := range node.ElseIfs {
		if err := elseif.Accept(f); err != nil {
			return fmt.Errorf("failed to format elseif: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	if node.Else != nil {
		if err := node.Else.Accept(f); err != nil {
			return fmt.Errorf("failed to format else: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	if err := f.str(f.keywords.EndIf); err != nil {
		return fmt.Errorf("failed for format end keyword: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitElseIf(node *ast.ElseIf) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.str(f.keywords.ElseIf); err != nil {
		return fmt.Errorf("failed for format keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Condition.Accept(f); err != nil {
		return fmt.Errorf("failed for format condition: %w", err)
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
			if statement.LeadingBlankLine() {
				if err := f.newline(); err != nil {
					return fmt.Errorf("failed to format newline: %w", err)
				}
			}
		}
		if err := statement.Accept(f); err != nil {
			return fmt.Errorf("failed to format statement: %w", err)
		}
	}
	if err := f.visitSuffixComments(node); err != nil {
		return err
	}
	f.level--
	return nil
}

func (f *formatter) VisitElse(node *ast.Else) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.str(f.keywords.Else); err != nil {
		return fmt.Errorf("failed for format keyword: %w", err)
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
			if statement.LeadingBlankLine() {
				if err := f.newline(); err != nil {
					return fmt.Errorf("failed to format newline: %w", err)
				}
			}
		}
		if err := statement.Accept(f); err != nil {
			return fmt.Errorf("failed to format statement: %w", err)
		}
	}
	if err := f.visitSuffixComments(node); err != nil {
		return err
	}
	f.level--
	return nil
}

func (f *formatter) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := node.Expression.Accept(f); err != nil {
		return fmt.Errorf("failed for format expression: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitImport(node *ast.Import) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.str(f.keywords.Import); err != nil {
		return fmt.Errorf("failed for format keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed for format Name: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitIndex(node *ast.Index) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := node.Value.Accept(f); err != nil {
		return fmt.Errorf("failed for format value: %w", err)
	}
	if err := f.str(token.BracketOpen.Symbol()); err != nil {
		return fmt.Errorf("failed to format open bracket: %w", err)
	}
	if err := node.Index.Accept(f); err != nil {
		return fmt.Errorf("failed for format index: %w", err)
	}
	if err := f.str(token.BracketClose.Symbol()); err != nil {
		return fmt.Errorf("failed to format close bracket: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitBoolLiteral(node *ast.BoolLiteral) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	text := f.keywords.False
	if node.Value {
		text = f.keywords.True
	}
	if err := f.str(text); err != nil {
		return fmt.Errorf("failed to format text: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitIntLiteral(node *ast.IntLiteral) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.bytes(node.Location().Text(f.file)); err != nil {
		return fmt.Errorf("failed to format text: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitFloatLiteral(node *ast.FloatLiteral) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.bytes(node.Location().Text(f.file)); err != nil {
		return fmt.Errorf("failed to format text: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitStringLiteral(node *ast.StringLiteral) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.bytes(node.Location().Text(f.file)); err != nil {
		return fmt.Errorf("failed to format text: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitNoneLiteral(node *ast.NoneLiteral) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.str(f.keywords.None); err != nil {
		return fmt.Errorf("failed to format None keyword: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitParameter(node *ast.Parameter) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := node.Type.Accept(f); err != nil {
		return fmt.Errorf("failed to format type: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format name: %w", err)
	}
	if node.DefaultValue != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(token.Assign.Symbol()); err != nil {
			return fmt.Errorf("failed to format assign operator: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.DefaultValue.Accept(f); err != nil {
			return fmt.Errorf("failed to format value: %w", err)
		}
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitParenthetical(node *ast.Parenthetical) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.str(token.ParenthesisOpen.Symbol()); err != nil {
		return fmt.Errorf("failed for format open parenthesis: %w", err)
	}
	if err := node.Value.Accept(f); err != nil {
		return fmt.Errorf("failed to format Value: %w", err)
	}
	if err := f.str(token.ParenthesisClose.Symbol()); err != nil {
		return fmt.Errorf("failed for format close parenthesis: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitProperty(node *ast.Property) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := node.Type.Accept(f); err != nil {
		return fmt.Errorf("failed to format type: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := f.str(f.keywords.Property); err != nil {
		return fmt.Errorf("failed for format start keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format name: %w", err)
	}
	if node.Value != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(token.Assign.Symbol()); err != nil {
			return fmt.Errorf("failed to format assign operator: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Value.Accept(f); err != nil {
			return fmt.Errorf("failed to format value: %w", err)
		}
	}
	if node.Kind == ast.Auto {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.Auto); err != nil {
			return fmt.Errorf("failed to format Auto keyword: %w", err)
		}
	}
	if node.Kind == ast.AutoReadOnly {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.AutoReadOnly); err != nil {
			return fmt.Errorf("failed to format AutoReadOnly keyword: %w", err)
		}
	}
	if len(node.HiddenLocations) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.Hidden); err != nil {
			return fmt.Errorf("failed to format Hidden keyword: %w", err)
		}
	}
	if len(node.ConditionalLocations) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.Conditional); err != nil {
			return fmt.Errorf("failed to format Conditional keyword: %w", err)
		}
	}
	if node.Documentation != nil {
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
		if err := node.Documentation.Accept(f); err != nil {
			return fmt.Errorf("failed to format documentation: %w", err)
		}
	}
	if node.Get != nil || node.Set != nil {
		if node.Documentation != nil {
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
		if err := f.str(f.keywords.EndProperty); err != nil {
			return fmt.Errorf("failed for format end keyword: %w", err)
		}
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitReturn(node *ast.Return) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.str(f.keywords.Return); err != nil {
		return fmt.Errorf("failed for format keyword: %w", err)
	}
	if node.Value != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Value.Accept(f); err != nil {
			return fmt.Errorf("failed to format Value: %w", err)
		}
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitScript(node *ast.Script) error {
	for _, c := range node.HeaderComments {
		if err := c.Accept(f); err != nil {
			return fmt.Errorf("failed to format header comment: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	if err := f.str(f.keywords.ScriptName); err != nil {
		return fmt.Errorf("failed for format ScriptName keyword: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format name: %w", err)
	}
	if node.Parent != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.Extends); err != nil {
			return fmt.Errorf("failed for format Extends keyword: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Parent.Accept(f); err != nil {
			return fmt.Errorf("failed to format parent: %w", err)
		}
	}
	if len(node.HiddenLocations) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.Hidden); err != nil {
			return fmt.Errorf("failed to format Hidden keyword: %w", err)
		}
	}
	if len(node.ConditionalLocations) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.Conditional); err != nil {
			return fmt.Errorf("failed to format Conditional keyword: %w", err)
		}
	}
	if err := f.newline(); err != nil {
		return fmt.Errorf("failed to format newline: %w", err)
	}
	if node.Documentation != nil {
		if err := node.Documentation.Accept(f); err != nil {
			return fmt.Errorf("failed to format documentation: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	if err := f.newline(); err != nil {
		return fmt.Errorf("failed to format newline: %w", err)
	}
	for i, stmt := range node.Statements {
		if i > 0 && stmt.LeadingBlankLine() {
			if err := f.newline(); err != nil {
				return fmt.Errorf("failed to format newline: %w", err)
			}
		}
		if err := stmt.Accept(f); err != nil {
			return fmt.Errorf("failed to format script statement: %w", err)
		}
		if err := f.newline(); err != nil {
			return fmt.Errorf("failed to format newline: %w", err)
		}
	}
	return nil
}

func (f *formatter) VisitState(node *ast.State) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if node.IsAuto {
		if err := f.str(f.keywords.Auto); err != nil {
			return fmt.Errorf("failed to format Auto keyword: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
	}
	if err := f.str(f.keywords.State); err != nil {
		return fmt.Errorf("failed for format start keyword: %w", err)
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
			if invokable.LeadingBlankLine() {
				if err := f.newline(); err != nil {
					return fmt.Errorf("failed to format newline: %w", err)
				}
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
	if err := f.str(f.keywords.EndState); err != nil {
		return fmt.Errorf("failed for format end keyword: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitTypeLiteral(node *ast.TypeLiteral) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
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
		text = string(node.Location().Text(f.file))
	default:
		return fmt.Errorf("unexpected types: %T", baseType)
	}
	if err := f.str(text); err != nil {
		return fmt.Errorf("failed to format type: %w", err)
	}
	if _, ok := node.Type.(types.Array); ok {
		if err := f.str("[]"); err != nil {
			return fmt.Errorf("failed to format array type: %w", err)
		}
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitUnary(node *ast.Unary) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.str(node.Kind.Symbol()); err != nil {
		return fmt.Errorf("failed to format operator: %w", err)
	}
	if err := node.Operand.Accept(f); err != nil {
		return fmt.Errorf("failed to format operand: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitScriptVariable(node *ast.ScriptVariable) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := node.Type.Accept(f); err != nil {
		return fmt.Errorf("failed to format Type: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format name: %w", err)
	}
	if node.Value != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(token.Assign.Symbol()); err != nil {
			return fmt.Errorf("failed to format assign operator: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Value.Accept(f); err != nil {
			return fmt.Errorf("failed to format value: %w", err)
		}
	}
	if len(node.ConditionalLocations) > 0 {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(f.keywords.Conditional); err != nil {
			return fmt.Errorf("failed to format Conditional keyword: %w", err)
		}
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitFunctionVariable(node *ast.FunctionVariable) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := node.Type.Accept(f); err != nil {
		return fmt.Errorf("failed to format Type: %w", err)
	}
	if err := f.space(); err != nil {
		return fmt.Errorf("failed to format space: %w", err)
	}
	if err := node.Name.Accept(f); err != nil {
		return fmt.Errorf("failed to format name: %w", err)
	}
	if node.Value != nil {
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := f.str(token.Assign.Symbol()); err != nil {
			return fmt.Errorf("failed to format assign operator: %w", err)
		}
		if err := f.space(); err != nil {
			return fmt.Errorf("failed to format space: %w", err)
		}
		if err := node.Value.Accept(f); err != nil {
			return fmt.Errorf("failed to format value: %w", err)
		}
	}
	return f.visitSuffixComments(node)
}

func (f *formatter) VisitWhile(node *ast.While) error {
	if err := f.visitPrefixComments(node); err != nil {
		return err
	}
	if err := f.str(f.keywords.While); err != nil {
		return fmt.Errorf("failed for format start keyword: %w", err)
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
			if statement.LeadingBlankLine() {
				if err := f.newline(); err != nil {
					return fmt.Errorf("failed to format newline: %w", err)
				}
			}
		}
		if err := statement.Accept(f); err != nil {
			return fmt.Errorf("failed to format statement: %w", err)
		}
	}
	f.level--
	if err := f.newline(); err != nil {
		return fmt.Errorf("failed to format newline: %w", err)
	}
	if err := f.str(f.keywords.EndWhile); err != nil {
		return fmt.Errorf("failed for format start keyword: %w", err)
	}
	return f.visitSuffixComments(node)
}

func (*formatter) VisitErrorStatement(node *ast.ErrorStatement) error {
	return fmt.Errorf("attempted to format error statement: %s", node.ErrorMessage)
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
	if f.level <= 0 {
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

func trimRight(text []byte) []byte {
	return bytes.TrimRight(text, " \r\n\t")
}

func trim(text []byte) []byte {
	return bytes.Trim(text, " \r\n\t")
}
