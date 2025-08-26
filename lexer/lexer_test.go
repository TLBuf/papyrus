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
		wantType   token.Kind
		wantText   string
		wantOffset uint32
		wantLength uint32
	}{
		{token.ScriptName, "ScriptName", 0, 10},
		{token.Identifier, "Foo", 11, 3},
		{token.Extends, "Extends", 15, 7},
		{token.Identifier, "Bar", 23, 3},
		{token.Newline, "\r\n", 26, 2},
		{token.BraceOpen, "{", 28, 1},
		{token.Comment, "A muliline\r\n\r\nDoc Comment", 29, 25},
		{token.BraceClose, "}", 54, 1},
		{token.Newline, "\r\n", 55, 2},
		{token.Newline, "\r\n", 57, 2},
		{token.Auto, "Auto", 59, 4},
		{token.State, "State", 64, 5},
		{token.Identifier, "Waiting", 70, 7},
		{token.Newline, "\r\n", 77, 2},
		{token.Event, "Event", 80, 5},
		{token.Identifier, "OnThing", 86, 7},
		{token.ParenthesisOpen, "(", 93, 1},
		{token.Identifier, "Baz", 94, 3},
		{token.Identifier, "arg", 98, 3},
		{token.ParenthesisClose, ")", 101, 1},
		{token.Newline, "\r\n", 102, 2},
		{token.BlockCommentOpen, ";/", 106, 2},
		{token.Comment, "\r\n\t\t\tA\r\n\t\t\tBlock\r\n\t\t\tComment\r\n\t\t", 108, 32},
		{token.BlockCommentClose, "/;", 140, 2},
		{token.Newline, "\r\n", 142, 2},
		{token.Int, "Int", 146, 3},
		{token.Identifier, "a", 150, 1},
		{token.Assign, "=", 152, 1},
		{token.IntLiteral, "1", 154, 1},
		{token.Plus, "+", 156, 1},
		{token.IntLiteral, "2", 158, 1},
		{token.Newline, "\r\n", 159, 2},
		{token.EndEvent, "EndEvent", 162, 8},
		{token.Newline, "\r\n", 170, 2},
		{token.Newline, "\r\n", 172, 2},
		{token.Int, "Int", 175, 3},
		{token.ArrayType, "[]", 178, 2},
		{token.Function, "Function", 181, 8},
		{token.Identifier, "Foo", 190, 3},
		{token.ParenthesisOpen, "(", 193, 1},
		{token.ParenthesisClose, ")", 194, 1},
		{token.Newline, "\r\n", 195, 2},
		{token.Return, "Return", 199, 6},
		{token.New, "New", 206, 3},
		{token.Int, "Int", 210, 3},
		{token.BracketOpen, "[", 213, 1},
		{token.IntLiteral, "2", 214, 1},
		{token.BracketClose, "]", 215, 1},
		{token.Newline, "\r\n", 216, 2},
		{token.EndFunction, "EndFunction", 219, 11},
		{token.Newline, "\r\n", 230, 2},
		{token.EndState, "EndState", 232, 8},
		{token.Semicolon, ";", 241, 1},
		{token.Comment, " Comment", 242, 8},
		{token.Newline, "\r\n", 250, 2},
		{token.EOF, "", 252, 0},
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
	}
}
