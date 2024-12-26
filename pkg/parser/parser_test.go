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
		{
			name:  "extends",
			input: "ScriptName Foo Extends Bar",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					SourceRange: source.Range{
						ByteOffset: 11,
						Length:     3,
						Line:       1,
						Column:     12,
					},
				},
				Extends: &ast.Identifier{
					Text: "bar",
					SourceRange: source.Range{
						ByteOffset: 23,
						Length:     3,
						Line:       1,
						Column:     24,
					},
				},
				SourceRange: source.Range{
					ByteOffset: 0,
					Length:     26,
					Line:       1,
					Column:     1,
				},
			},
		},
		{
			name:  "hidden",
			input: "ScriptName Foo Hidden",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					SourceRange: source.Range{
						ByteOffset: 11,
						Length:     3,
						Line:       1,
						Column:     12,
					},
				},
				IsHidden: true,
				SourceRange: source.Range{
					ByteOffset: 0,
					Length:     21,
					Line:       1,
					Column:     1,
				},
			},
		},
		{
			name:  "conditional",
			input: "ScriptName Foo Conditional",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					SourceRange: source.Range{
						ByteOffset: 11,
						Length:     3,
						Line:       1,
						Column:     12,
					},
				},
				IsConditional: true,
				SourceRange: source.Range{
					ByteOffset: 0,
					Length:     26,
					Line:       1,
					Column:     1,
				},
			},
		},
		{
			name:  "hidden_conditional",
			input: "ScriptName Foo Hidden Conditional",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					SourceRange: source.Range{
						ByteOffset: 11,
						Length:     3,
						Line:       1,
						Column:     12,
					},
				},
				IsHidden:      true,
				IsConditional: true,
				SourceRange: source.Range{
					ByteOffset: 0,
					Length:     33,
					Line:       1,
					Column:     1,
				},
			},
		},
		{
			name:  "conditional_hidden",
			input: "ScriptName Foo Conditional Hidden",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					SourceRange: source.Range{
						ByteOffset: 11,
						Length:     3,
						Line:       1,
						Column:     12,
					},
				},
				IsHidden:      true,
				IsConditional: true,
				SourceRange: source.Range{
					ByteOffset: 0,
					Length:     33,
					Line:       1,
					Column:     1,
				},
			},
		},
		{
			name:  "many_flags",
			input: "ScriptName Foo Conditional Hidden Conditional Hidden",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					SourceRange: source.Range{
						ByteOffset: 11,
						Length:     3,
						Line:       1,
						Column:     12,
					},
				},
				IsHidden:      true,
				IsConditional: true,
				SourceRange: source.Range{
					ByteOffset: 0,
					Length:     52,
					Line:       1,
					Column:     1,
				},
			},
		},
		{
			name:  "extends_many_flags",
			input: "ScriptName Foo Extends Bar Hidden Conditional Hidden Conditional",
			want: &ast.Script{
				Name: &ast.Identifier{
					Text: "foo",
					SourceRange: source.Range{
						ByteOffset: 11,
						Length:     3,
						Line:       1,
						Column:     12,
					},
				},
				Extends: &ast.Identifier{
					Text: "bar",
					SourceRange: source.Range{
						ByteOffset: 23,
						Length:     3,
						Line:       1,
						Column:     24,
					},
				},
				IsHidden:      true,
				IsConditional: true,
				SourceRange: source.Range{
					ByteOffset: 0,
					Length:     64,
					Line:       1,
					Column:     1,
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
					SourceRange: source.Range{
						ByteOffset: 11,
						Length:     3,
						Line:       1,
						Column:     12,
					},
				},
				Statements: []ast.ScriptStatement{
					&ast.Import{
						Name: &ast.Identifier{
							Text: "bar",
							SourceRange: source.Range{
								ByteOffset: 25,
								Length:     3,
								Line:       2,
								Column:     11,
							},
						},
						SourceRange: source.Range{
							ByteOffset: 18,
							Length:     10,
							Line:       2,
							Column:     2,
						},
					},
				},
				SourceRange: source.Range{
					ByteOffset: 0,
					Length:     28,
					Line:       1,
					Column:     1,
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
					SourceRange: source.Range{
						ByteOffset: 11,
						Length:     3,
						Line:       1,
						Column:     12,
					},
				},
				Statements: []ast.ScriptStatement{
					&ast.State{
						Name: &ast.Identifier{
							Text: "bar",
							SourceRange: source.Range{
								ByteOffset: 24,
								Length:     3,
								Line:       2,
								Column:     10,
							},
						},
						IsAuto: false,
						SourceRange: source.Range{
							ByteOffset: 18,
							Length:     21,
							Line:       2,
							Column:     2,
						},
					},
				},
				SourceRange: source.Range{
					ByteOffset: 0,
					Length:     39,
					Line:       1,
					Column:     1,
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
					SourceRange: source.Range{
						ByteOffset: 11,
						Length:     3,
						Line:       1,
						Column:     12,
					},
				},
				Statements: []ast.ScriptStatement{
					&ast.State{
						Name: &ast.Identifier{
							Text: "bar",
							SourceRange: source.Range{
								ByteOffset: 29,
								Length:     3,
								Line:       2,
								Column:     15,
							},
						},
						IsAuto: true,
						SourceRange: source.Range{
							ByteOffset: 18,
							Length:     26,
							Line:       2,
							Column:     2,
						},
					},
				},
				SourceRange: source.Range{
					ByteOffset: 0,
					Length:     44,
					Line:       1,
					Column:     1,
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
			if diff := cmp.Diff(test.want, got, cmpopts.IgnoreFields(source.Range{}, "File")); diff != "" {
				t.Errorf("ParseScript() mismatch (-want +got):\n%s", diff)
			}
		})
	}

}
