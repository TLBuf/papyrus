package lexer_test

import (
	"testing"

	"github.com/TLBuf/papyrus/pkg/lexer"
	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/TLBuf/papyrus/pkg/token"
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
		wantType   token.Type
		wantText   string
		wantOffset int
		wantLine   int
		wantColumn int
	}{
		{token.ScriptName, "ScriptName", 0, 1, 1},
		{token.Identifier, "Foo", 11, 1, 12},
		{token.Extends, "Extends", 15, 1, 16},
		{token.Identifier, "Bar", 23, 1, 24},
		{token.Newline, "\n", 26, 1, 27},
		{token.DocComment, "{A muliline\n\nDoc Comment}", 27, 2, 1},
		{token.Newline, "\n", 52, 4, 13},
		{token.Newline, "\n", 53, 5, 1},
		{token.Auto, "Auto", 54, 6, 1},
		{token.State, "State", 59, 6, 6},
		{token.Identifier, "Waiting", 65, 6, 12},
		{token.Newline, "\n", 72, 6, 19},
		{token.Event, "Event", 74, 7, 2},
		{token.Identifier, "OnThing", 80, 7, 8},
		{token.LParen, "(", 87, 7, 15},
		{token.Identifier, "Baz", 88, 7, 16},
		{token.Identifier, "arg", 92, 7, 20},
		{token.RParen, ")", 95, 7, 23},
		{token.Newline, "\n", 96, 7, 24},
		{token.BlockComment, ";/\n\t\t\tA\n\t\t\tBlock\n\t\t\tComment\n\t\t/;", 99, 8, 3},
		{token.Newline, "\n", 131, 12, 5},
		{token.Int, "Int", 134, 13, 3},
		{token.Identifier, "a", 138, 13, 7},
		{token.Assign, "=", 140, 13, 9},
		{token.IntLiteral, "1", 142, 13, 11},
		{token.Add, "+", 144, 13, 13},
		{token.IntLiteral, "2", 146, 13, 15},
		{token.Newline, "\n", 147, 13, 16},
		{token.EndEvent, "EndEvent", 149, 14, 2},
		{token.Newline, "\n", 157, 14, 10},
		{token.Newline, "\n", 158, 15, 1},
		{token.Int, "Int", 160, 16, 2},
		{token.LBracket, "[", 163, 16, 5},
		{token.RBracket, "]", 164, 16, 6},
		{token.Function, "Function", 166, 16, 8},
		{token.Identifier, "Foo", 175, 16, 17},
		{token.LParen, "(", 178, 16, 20},
		{token.RParen, ")", 179, 16, 21},
		{token.Newline, "\n", 180, 16, 22},
		{token.Return, "Return", 183, 17, 3},
		{token.New, "New", 190, 17, 10},
		{token.Int, "Int", 194, 17, 14},
		{token.LBracket, "[", 197, 17, 17},
		{token.IntLiteral, "2", 198, 17, 18},
		{token.RBracket, "]", 199, 17, 19},
		{token.Newline, "\n", 200, 17, 20},
		{token.EndFunction, "EndFunction", 202, 18, 2},
		{token.Newline, "\n", 213, 18, 13},
		{token.EndState, "EndState", 214, 19, 1},
		{token.LineComment, "; Comment", 223, 19, 10},
		{token.Newline, "\n", 232, 19, 19},
		{token.EOF, "", 233, 20, 1},
	}
	file := &source.File{
		Text: []byte(text),
	}
	l := lexer.New(file)
	for i, tt := range tests {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("unexpected error at token %d: %v", i, err)
		}
		if tok.Type != tt.wantType {
			t.Errorf("token type mismatch at token %d, want: %v, got: %v", i, tt.wantType, tok.Type)
		}
		gotText := string(tok.SourceRange.Text())
		if gotText != tt.wantText {
			t.Errorf("token text mismatch at token %d, want: %q, got: %q", i, tt.wantText, gotText)
		}
		if tok.SourceRange.ByteOffset != tt.wantOffset {
			t.Errorf("token byte offset mismatch at token %d, want: %d, got: %d", i, tt.wantOffset, tok.SourceRange.ByteOffset)
		}
		if tok.SourceRange.Line != tt.wantLine {
			t.Errorf("token line mismatch at token %d, want: %d, got: %d", i, tt.wantLine, tok.SourceRange.Line)
		}
		if tok.SourceRange.Column != tt.wantColumn {
			t.Errorf("token column mismatch at token %d, want: %d, got: %d", i, tt.wantColumn, tok.SourceRange.Column)
		}
	}
}
