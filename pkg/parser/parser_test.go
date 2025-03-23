package parser_test

import (
	"testing"

	"github.com/TLBuf/papyrus/pkg/ast"
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
					Text: "foo",
					Location: source.Range{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Location: source.Range{
					ByteOffset: 0,
					Length:     14,
				},
			},
		},
		{
			name:  "extends",
			input: "ScriptName Foo Extends Bar",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					Location: source.Range{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Extends: &ast.Identifier{
					Text: "bar",
					Location: source.Range{
						ByteOffset: 23,
						Length:     3,
					},
				},
				Location: source.Range{
					ByteOffset: 0,
					Length:     26,
				},
			},
		},
		{
			name:  "hidden",
			input: "ScriptName Foo Hidden",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					Location: source.Range{
						ByteOffset: 11,
						Length:     3,
					},
				},
				IsHidden: true,
				Location: source.Range{
					ByteOffset: 0,
					Length:     21,
				},
			},
		},
		{
			name:  "conditional",
			input: "ScriptName Foo Conditional",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					Location: source.Range{
						ByteOffset: 11,
						Length:     3,
					},
				},
				IsConditional: true,
				Location: source.Range{
					ByteOffset: 0,
					Length:     26,
				},
			},
		},
		{
			name:  "hidden_conditional",
			input: "ScriptName Foo Hidden Conditional",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					Location: source.Range{
						ByteOffset: 11,
						Length:     3,
					},
				},
				IsHidden:      true,
				IsConditional: true,
				Location: source.Range{
					ByteOffset: 0,
					Length:     33,
				},
			},
		},
		{
			name:  "conditional_hidden",
			input: "ScriptName Foo Conditional Hidden",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					Location: source.Range{
						ByteOffset: 11,
						Length:     3,
					},
				},
				IsHidden:      true,
				IsConditional: true,
				Location: source.Range{
					ByteOffset: 0,
					Length:     33,
				},
			},
		},
		{
			name:  "many_flags",
			input: "ScriptName Foo Conditional Hidden Conditional Hidden",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					Location: source.Range{
						ByteOffset: 11,
						Length:     3,
					},
				},
				IsHidden:      true,
				IsConditional: true,
				Location: source.Range{
					ByteOffset: 0,
					Length:     52,
				},
			},
		},
		{
			name:  "extends_many_flags",
			input: "ScriptName Foo Extends Bar Hidden Conditional Hidden Conditional",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					Location: source.Range{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Extends: &ast.Identifier{
					Text: "bar",
					Location: source.Range{
						ByteOffset: 23,
						Length:     3,
					},
				},
				IsHidden:      true,
				IsConditional: true,
				Location: source.Range{
					ByteOffset: 0,
					Length:     64,
				},
			},
		},
		{
			name: "import",
			input: `ScriptName Foo
			Import Bar`,
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					Location: source.Range{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Statements: []ast.ScriptStatement{
					&ast.Import{
						Name: &ast.Identifier{
							Text: "bar",
							Location: source.Range{
								ByteOffset: 25,
								Length:     3,
							},
						},
						Location: source.Range{
							ByteOffset: 18,
							Length:     10,
						},
					},
				},
				Location: source.Range{
					ByteOffset: 0,
					Length:     28,
				},
			},
		},
		{
			name: "state",
			input: `ScriptName Foo
			State Bar
			EndState`,
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					Location: source.Range{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Statements: []ast.ScriptStatement{
					&ast.State{
						Name: &ast.Identifier{
							Text: "bar",
							Location: source.Range{
								ByteOffset: 24,
								Length:     3,
							},
						},
						IsAuto: false,
						Location: source.Range{
							ByteOffset: 18,
							Length:     21,
						},
					},
				},
				Location: source.Range{
					ByteOffset: 0,
					Length:     39,
				},
			},
		},
		{
			name: "state_auto",
			input: `ScriptName Foo
			Auto State Bar
			EndState`,
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					Location: source.Range{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Statements: []ast.ScriptStatement{
					&ast.State{
						Name: &ast.Identifier{
							Text: "bar",
							Location: source.Range{
								ByteOffset: 29,
								Length:     3,
							},
						},
						IsAuto: true,
						Location: source.Range{
							ByteOffset: 18,
							Length:     26,
						},
					},
				},
				Location: source.Range{
					ByteOffset: 0,
					Length:     44,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := &source.File{Text: []byte(test.input)}
			p := parser.New()

			got, err := p.Parse(f)
			if err != nil {
				t.Errorf("ParseScript() returned an unexpected error: %v", err)
			}
			if got == nil {
				t.Fatalf("ParseScript() returned nil")
			}
			if diff := cmp.Diff(test.want, got, ignoreFields...); diff != "" {
				t.Errorf("ParseScript() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

var ignoreFields = []cmp.Option{
	// lexer_test pretty thoroughly covers these fields, if the
	// byte offset and length match, that's sufficent for this test.
	cmpopts.IgnoreFields(source.Range{}, "File", "StartLine", "StartColumn", "EndLine", "EndColumn", "PreambleLength", "PostambleLength"),
}
