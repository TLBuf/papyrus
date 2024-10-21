package lexer_test

import (
	"testing"

	"github.com/TLBuf/papyrus/pkg/lexer"
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
		wantType token.Type
		wantText string
	}{
		{token.ScriptName, "ScriptName"},
		{token.Identifier, "Foo"},
		{token.Extends, "Extends"},
		{token.Identifier, "Bar"},
		{token.Newline, "\n"},
		{token.DocComment, "{A muliline\n\nDoc Comment}"},
		{token.Newline, "\n"},
		{token.Newline, "\n"},
		{token.Auto, "Auto"},
		{token.State, "State"},
		{token.Identifier, "Waiting"},
		{token.Newline, "\n"},
		{token.Event, "Event"},
		{token.Identifier, "OnThing"},
		{token.LParen, "("},
		{token.Identifier, "Baz"},
		{token.Identifier, "arg"},
		{token.RParen, ")"},
		{token.Newline, "\n"},
		{token.BlockComment, ";/\n\t\t\tA\n\t\t\tBlock\n\t\t\tComment\n\t\t/;"},
		{token.Newline, "\n"},
		{token.Int, "Int"},
		{token.Identifier, "a"},
		{token.Assign, "="},
		{token.IntLiteral, "1"},
		{token.Add, "+"},
		{token.IntLiteral, "2"},
		{token.Newline, "\n"},
		{token.EndEvent, "EndEvent"},
		{token.Newline, "\n"},
		{token.Newline, "\n"},
		{token.Int, "Int"},
		{token.LBracket, "["},
		{token.RBracket, "]"},
		{token.Function, "Function"},
		{token.Identifier, "Foo"},
		{token.LParen, "("},
		{token.RParen, ")"},
		{token.Newline, "\n"},
		{token.Return, "Return"},
		{token.New, "New"},
		{token.Int, "Int"},
		{token.LBracket, "["},
		{token.IntLiteral, "2"},
		{token.RBracket, "]"},
		{token.Newline, "\n"},
		{token.EndFunction, "EndFunction"},
		{token.Newline, "\n"},
		{token.EndState, "EndState"},
		{token.LineComment, "; Comment"},
		{token.Newline, "\n"},
		{token.EOF, ""},
	}
	l := lexer.New([]byte(text))
	for i, tt := range tests {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatalf("unexpected error at token %d: %v", i, err)
		}
		if tok.Type != tt.wantType {
			t.Fatalf("token type mismatch at token %d, want: %v, got: %v", i, tt.wantType, tok.Type)
		}
		gotText := string(tok.Text)
		if gotText != tt.wantText {
			t.Fatalf("token text mismatch at token %d, want: %q, got: %q", i, tt.wantText, gotText)
		}
	}
}
