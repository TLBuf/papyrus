package lexer_test

import (
	"strings"
	"testing"

	"github.com/TLBuf/papyrus/lexer"
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/token"
)

func TestNextToken(t *testing.T) {
	text := `ScriptName Foo Extends Bar
{A muliline

Doc Comment}

Auto State Waiting
	Event OnThing(Baz arg)
		;/
			A
			Block
			Comment
		/;
		Int a = 1 + 2
	EndEvent

	Int[] Function Foo()
		Return New Int[2]
	EndFunction
EndState ; Comment
`
	tests := []struct {
		wantType            token.Kind
		wantText            string
		wantOffset          uint32
		wantLength          uint32
		wantStartLine       uint32
		wantStartColumn     uint32
		wantEndLine         uint32
		wantEndColumn       uint32
		wantPreambleLength  uint32
		wantPostambleLength uint32
	}{
		{token.ScriptName, "ScriptName", 0, 10, 1, 1, 1, 10, 0, 16},
		{token.Identifier, "Foo", 11, 3, 1, 12, 1, 14, 11, 12},
		{token.Extends, "Extends", 15, 7, 1, 16, 1, 22, 15, 4},
		{token.Identifier, "Bar", 23, 3, 1, 24, 1, 26, 23, 0},
		{token.Newline, "\r\n", 26, 2, 1, 27, 1, 28, 26, 0},
		{token.BraceOpen, "{", 28, 1, 2, 1, 2, 1, 0, 10},
		{token.Comment, "A muliline\r\n\r\nDoc Comment", 29, 25, 2, 2, 4, 11, 1, 1},
		{token.BraceClose, "}", 54, 1, 4, 12, 4, 12, 11, 0},
		{token.Newline, "\r\n", 55, 2, 4, 13, 4, 14, 12, 0},
		{token.Newline, "\r\n", 57, 2, 5, 1, 5, 2, 0, 0},
		{token.Auto, "Auto", 59, 4, 6, 1, 6, 4, 0, 14},
		{token.State, "State", 64, 5, 6, 6, 6, 10, 5, 8},
		{token.Identifier, "Waiting", 70, 7, 6, 12, 6, 18, 11, 0},
		{token.Newline, "\r\n", 77, 2, 6, 19, 6, 20, 18, 0},
		{token.Event, "Event", 80, 5, 7, 2, 7, 6, 1, 17},
		{token.Identifier, "OnThing", 86, 7, 7, 8, 7, 14, 7, 9},
		{token.ParenthesisOpen, "(", 93, 1, 7, 15, 7, 15, 14, 8},
		{token.Identifier, "Baz", 94, 3, 7, 16, 7, 18, 15, 5},
		{token.Identifier, "arg", 98, 3, 7, 20, 7, 22, 19, 1},
		{token.ParenthesisClose, ")", 101, 1, 7, 23, 7, 23, 22, 0},
		{token.Newline, "\r\n", 102, 2, 7, 24, 7, 25, 23, 0},
		{token.BlockCommentOpen, ";/", 106, 2, 8, 3, 8, 4, 2, 0},
		{token.Comment, "\r\n\t\t\tA\r\n\t\t\tBlock\r\n\t\t\tComment\r\n\t\t", 108, 32, 8, 5, 12, 2, 4, 2},
		{token.BlockCommentClose, "/;", 140, 2, 12, 3, 12, 4, 2, 0},
		{token.Newline, "\r\n", 142, 2, 12, 5, 12, 6, 4, 0},
		{token.Int, "Int", 146, 3, 13, 3, 13, 5, 2, 10},
		{token.Identifier, "a", 150, 1, 13, 7, 13, 7, 6, 8},
		{token.Assign, "=", 152, 1, 13, 9, 13, 9, 8, 6},
		{token.IntLiteral, "1", 154, 1, 13, 11, 13, 11, 10, 4},
		{token.Plus, "+", 156, 1, 13, 13, 13, 13, 12, 2},
		{token.IntLiteral, "2", 158, 1, 13, 15, 13, 15, 14, 0},
		{token.Newline, "\r\n", 159, 2, 13, 16, 13, 17, 15, 0},
		{token.EndEvent, "EndEvent", 162, 8, 14, 2, 14, 9, 1, 0},
		{token.Newline, "\r\n", 170, 2, 14, 10, 14, 11, 9, 0},
		{token.Newline, "\r\n", 172, 2, 15, 1, 15, 2, 0, 0},
		{token.Int, "Int", 175, 3, 16, 2, 16, 4, 1, 17},
		{token.ArrayType, "[]", 178, 2, 16, 5, 16, 6, 4, 15},
		{token.Function, "Function", 181, 8, 16, 8, 16, 15, 7, 6},
		{token.Identifier, "Foo", 190, 3, 16, 17, 16, 19, 16, 2},
		{token.ParenthesisOpen, "(", 193, 1, 16, 20, 16, 20, 19, 1},
		{token.ParenthesisClose, ")", 194, 1, 16, 21, 16, 21, 20, 0},
		{token.Newline, "\r\n", 195, 2, 16, 22, 16, 23, 21, 0},
		{token.Return, "Return", 199, 6, 17, 3, 17, 8, 2, 11},
		{token.New, "New", 206, 3, 17, 10, 17, 12, 9, 7},
		{token.Int, "Int", 210, 3, 17, 14, 17, 16, 13, 3},
		{token.BracketOpen, "[", 213, 1, 17, 17, 17, 17, 16, 2},
		{token.IntLiteral, "2", 214, 1, 17, 18, 17, 18, 17, 1},
		{token.BracketClose, "]", 215, 1, 17, 19, 17, 19, 18, 0},
		{token.Newline, "\r\n", 216, 2, 17, 20, 17, 21, 19, 0},
		{token.EndFunction, "EndFunction", 219, 11, 18, 2, 18, 12, 1, 0},
		{token.Newline, "\r\n", 230, 2, 18, 13, 18, 14, 12, 0},
		{token.EndState, "EndState", 232, 8, 19, 1, 19, 8, 0, 10},
		{token.Semicolon, ";", 241, 1, 19, 10, 19, 10, 9, 8},
		{token.Comment, " Comment", 242, 8, 19, 11, 19, 18, 10, 0},
		{token.Newline, "\r\n", 250, 2, 19, 19, 19, 20, 18, 0},
		{token.EOF, "", 252, 0, 20, 1, 20, 1, 0, 0},
	}
	// Papyrus uses Windows line endings.
	text = strings.ReplaceAll(text, "\n", "\r\n")
	file, _ := source.NewFile("test.psc", []byte(text))
	lex, err := lexer.New(file)
	if err != nil {
		t.Fatalf("New() returned an unexpected error: %v", err)
	}
	for i, tt := range tests {
		tok, err := lex.Next()
		if err != nil {
			t.Fatalf("NextToken() returned an unexpected error at token %d: %v", i, err)
		}
		if tok.Kind != tt.wantType {
			t.Errorf("token type mismatch at token %d %q, want: %v, got: %v", i, tok, tt.wantType, tok.Kind)
		}
		gotText := string(file.Bytes(tok.Location))
		if gotText != tt.wantText {
			t.Errorf("token text mismatch at token %d %q, want: %q, got: %q", i, tok, tt.wantText, gotText)
		}
		if tok.Location.ByteOffset != tt.wantOffset {
			t.Errorf(
				"token byte offset mismatch at token %d %q, want: %d, got: %d",
				i,
				tok,
				tt.wantOffset,
				tok.Location.ByteOffset,
			)
		}
		if tok.Location.Length != tt.wantLength {
			t.Errorf("token length mismatch at token %d %q, want: %d, got: %d", i, tok, tt.wantLength, tok.Location.Length)
		}
		if tok.Location.StartLine != tt.wantStartLine {
			t.Errorf(
				"token start line mismatch at token %d %q, want: %d, got: %d",
				i,
				tok,
				tt.wantStartLine,
				tok.Location.StartLine,
			)
		}
		if tok.Location.StartColumn != tt.wantStartColumn {
			t.Errorf(
				"token start column mismatch at token %d %q, want: %d, got: %d",
				i,
				tok,
				tt.wantStartColumn,
				tok.Location.StartColumn,
			)
		}
		if tok.Location.EndLine != tt.wantEndLine {
			t.Errorf(
				"token end line mismatch at token %d %q, want: %d, got: %d",
				i,
				tok,
				tt.wantEndLine,
				tok.Location.EndLine,
			)
		}
		if tok.Location.EndColumn != tt.wantEndColumn {
			t.Errorf(
				"token end column mismatch at token %d %q, want: %d, got: %d",
				i,
				tok,
				tt.wantEndColumn,
				tok.Location.EndColumn,
			)
		}
		if tok.Location.PreambleLength != tt.wantPreambleLength {
			t.Errorf(
				"token preamble length mismatch at token %d %q, want: %d, got: %d",
				i,
				tok,
				tt.wantPreambleLength,
				tok.Location.PreambleLength,
			)
		}
		if tok.Location.PostambleLength != tt.wantPostambleLength {
			t.Errorf(
				"token postamble length mismatch at token %d %q, want: %d, got: %d",
				i,
				tok,
				tt.wantPostambleLength,
				tok.Location.PostambleLength,
			)
		}
	}
}
