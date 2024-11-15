package parser_test

import (
	"testing"

	"github.com/TLBuf/papyrus/pkg/ast"
	"github.com/TLBuf/papyrus/pkg/lexer"
	"github.com/TLBuf/papyrus/pkg/parser"
	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestHeader(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  *ast.Script
	}{
		{
			name:  "basic",
			input: "ScriptName Foo",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "Foo",
					SourceRange: source.Range{
						ByteOffset: 11,
						Length:     3,
						Line:       1,
						Column:     12,
					},
				},
				SourceRange: source.Range{
					ByteOffset: 0,
					Length:     14,
					Line:       1,
					Column:     1,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := &source.File{Text: []byte(test.input)}
			l := lexer.New(f)
			p := parser.New(l)

			got, err := p.ParseScript()
			if err != nil {
				t.Errorf("ParseScript() returned an unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("ParseScript() returned nil")
			}
			if diff := cmp.Diff(test.want, got, cmpopts.IgnoreFields(source.Range{}, "File")); diff != "" {
				t.Errorf("ParseScript() mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
