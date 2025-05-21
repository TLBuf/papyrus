package parser_test

import (
	"testing"

	"github.com/TLBuf/papyrus/pkg/ast"
	"github.com/TLBuf/papyrus/pkg/parser"
	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/TLBuf/papyrus/pkg/token"
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
				Keyword: &ast.Token{
					Kind: token.ScriptName,
					Location: source.Location{
						ByteOffset: 0,
						Length:     10,
					},
				},
				Name: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 11,
							Length:     3,
						},
					},
					Normalized: "foo",
					Location: source.Location{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Location: source.Location{
					ByteOffset: 0,
					Length:     14,
				},
			},
		},
		{
			name:  "extends",
			input: "ScriptName Foo Extends Bar",
			want: &ast.Script{
				Keyword: &ast.Token{
					Kind: token.ScriptName,
					Location: source.Location{
						ByteOffset: 0,
						Length:     10,
					},
				},
				Name: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 11,
							Length:     3,
						},
					},
					Normalized: "foo",
					Location: source.Location{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Extends: &ast.Token{
					Kind: token.Extends,
					Location: source.Location{
						ByteOffset: 15,
						Length:     7,
					},
				},
				Parent: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 23,
							Length:     3,
						},
					},
					Normalized: "bar",
					Location: source.Location{
						ByteOffset: 23,
						Length:     3,
					},
				},
				Location: source.Location{
					ByteOffset: 0,
					Length:     26,
				},
			},
		},
		{
			name:  "hidden",
			input: "ScriptName Foo Hidden",
			want: &ast.Script{
				Keyword: &ast.Token{
					Kind: token.ScriptName,
					Location: source.Location{
						ByteOffset: 0,
						Length:     10,
					},
				},
				Name: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 11,
							Length:     3,
						},
					},
					Normalized: "foo",
					Location: source.Location{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Hidden: []*ast.Token{{
					Kind: token.Hidden,
					Location: source.Location{
						ByteOffset: 15,
						Length:     6,
					},
				}},
				Location: source.Location{
					ByteOffset: 0,
					Length:     21,
				},
			},
		},
		{
			name:  "conditional",
			input: "ScriptName Foo Conditional",
			want: &ast.Script{
				Keyword: &ast.Token{
					Kind: token.ScriptName,
					Location: source.Location{
						ByteOffset: 0,
						Length:     10,
					},
				},
				Name: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 11,
							Length:     3,
						},
					},
					Normalized: "foo",
					Location: source.Location{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Conditional: []*ast.Token{{
					Kind: token.Conditional,
					Location: source.Location{
						ByteOffset: 15,
						Length:     11,
					},
				}},
				Location: source.Location{
					ByteOffset: 0,
					Length:     26,
				},
			},
		},
		{
			name:  "hidden_conditional",
			input: "ScriptName Foo Hidden Conditional",
			want: &ast.Script{
				Keyword: &ast.Token{
					Kind: token.ScriptName,
					Location: source.Location{
						ByteOffset: 0,
						Length:     10,
					},
				},
				Name: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 11,
							Length:     3,
						},
					},
					Normalized: "foo",
					Location: source.Location{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Hidden: []*ast.Token{{
					Kind: token.Hidden,
					Location: source.Location{
						ByteOffset: 15,
						Length:     6,
					},
				}},
				Conditional: []*ast.Token{{
					Kind: token.Conditional,
					Location: source.Location{
						ByteOffset: 22,
						Length:     11,
					},
				}},
				Location: source.Location{
					ByteOffset: 0,
					Length:     33,
				},
			},
		},
		{
			name:  "conditional_hidden",
			input: "ScriptName Foo Conditional Hidden",
			want: &ast.Script{
				Keyword: &ast.Token{
					Kind: token.ScriptName,
					Location: source.Location{
						ByteOffset: 0,
						Length:     10,
					},
				},
				Name: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 11,
							Length:     3,
						},
					},
					Normalized: "foo",
					Location: source.Location{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Hidden: []*ast.Token{{
					Kind: token.Hidden,
					Location: source.Location{
						ByteOffset: 27,
						Length:     6,
					},
				}},
				Conditional: []*ast.Token{{
					Kind: token.Conditional,
					Location: source.Location{
						ByteOffset: 15,
						Length:     11,
					},
				}},
				Location: source.Location{
					ByteOffset: 0,
					Length:     33,
				},
			},
		},
		{
			name:  "many_flags",
			input: "ScriptName Foo Conditional Hidden Conditional Hidden",
			want: &ast.Script{
				Keyword: &ast.Token{
					Kind: token.ScriptName,
					Location: source.Location{
						ByteOffset: 0,
						Length:     10,
					},
				},
				Name: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 11,
							Length:     3,
						},
					},
					Normalized: "foo",
					Location: source.Location{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Hidden: []*ast.Token{
					{
						Kind: token.Hidden,
						Location: source.Location{
							ByteOffset: 27,
							Length:     6,
						},
					}, {
						Kind: token.Hidden,
						Location: source.Location{
							ByteOffset: 46,
							Length:     6,
						},
					},
				},
				Conditional: []*ast.Token{
					{
						Kind: token.Conditional,
						Location: source.Location{
							ByteOffset: 15,
							Length:     11,
						},
					}, {
						Kind: token.Conditional,
						Location: source.Location{
							ByteOffset: 34,
							Length:     11,
						},
					},
				},
				Location: source.Location{
					ByteOffset: 0,
					Length:     52,
				},
			},
		},
		{
			name:  "extends_many_flags",
			input: "ScriptName Foo Extends Bar Hidden Conditional Hidden Conditional",
			want: &ast.Script{
				Keyword: &ast.Token{
					Kind: token.ScriptName,
					Location: source.Location{
						ByteOffset: 0,
						Length:     10,
					},
				},
				Name: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 11,
							Length:     3,
						},
					},
					Normalized: "foo",
					Location: source.Location{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Extends: &ast.Token{
					Kind: token.Extends,
					Location: source.Location{
						ByteOffset: 15,
						Length:     7,
					},
				},
				Parent: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 23,
							Length:     3,
						},
					},
					Normalized: "bar",
					Location: source.Location{
						ByteOffset: 23,
						Length:     3,
					},
				},
				Hidden: []*ast.Token{
					{
						Kind: token.Hidden,
						Location: source.Location{
							ByteOffset: 27,
							Length:     6,
						},
					}, {
						Kind: token.Hidden,
						Location: source.Location{
							ByteOffset: 46,
							Length:     6,
						},
					},
				},
				Conditional: []*ast.Token{
					{
						Kind: token.Conditional,
						Location: source.Location{
							ByteOffset: 34,
							Length:     11,
						},
					}, {
						Kind: token.Conditional,
						Location: source.Location{
							ByteOffset: 53,
							Length:     11,
						},
					},
				},
				Location: source.Location{
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
				Keyword: &ast.Token{
					Kind: token.ScriptName,
					Location: source.Location{
						ByteOffset: 0,
						Length:     10,
					},
				},
				Name: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 11,
							Length:     3,
						},
					},
					Normalized: "foo",
					Location: source.Location{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Statements: []ast.ScriptStatement{
					&ast.Import{
						Keyword: &ast.Token{
							Kind: token.Import,
							Location: source.Location{
								ByteOffset: 18,
								Length:     6,
							},
						},
						Name: &ast.Identifier{
							Text: &ast.Token{
								Kind: token.Identifier,
								Location: source.Location{
									ByteOffset: 25,
									Length:     3,
								},
							},
							Normalized: "bar",
							Location: source.Location{
								ByteOffset: 25,
								Length:     3,
							},
						},
						Location: source.Location{
							ByteOffset: 18,
							Length:     10,
						},
					},
				},
				Location: source.Location{
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
				Keyword: &ast.Token{
					Kind: token.ScriptName,
					Location: source.Location{
						ByteOffset: 0,
						Length:     10,
					},
				},
				Name: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 11,
							Length:     3,
						},
					},
					Normalized: "foo",
					Location: source.Location{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Statements: []ast.ScriptStatement{
					&ast.State{
						Keyword: &ast.Token{
							Kind: token.State,
							Location: source.Location{
								ByteOffset: 18,
								Length:     5,
							},
						},
						Name: &ast.Identifier{
							Text: &ast.Token{
								Kind: token.Identifier,
								Location: source.Location{
									ByteOffset: 24,
									Length:     3,
								},
							},
							Normalized: "bar",
							Location: source.Location{
								ByteOffset: 24,
								Length:     3,
							},
						},
						EndKeyword: &ast.Token{
							Kind: token.EndState,
							Location: source.Location{
								ByteOffset: 31,
								Length:     8,
							},
						},
						Location: source.Location{
							ByteOffset: 18,
							Length:     21,
						},
					},
				},
				Location: source.Location{
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
				Keyword: &ast.Token{
					Kind: token.ScriptName,
					Location: source.Location{
						ByteOffset: 0,
						Length:     10,
					},
				},
				Name: &ast.Identifier{
					Text: &ast.Token{
						Kind: token.Identifier,
						Location: source.Location{
							ByteOffset: 11,
							Length:     3,
						},
					},
					Normalized: "foo",
					Location: source.Location{
						ByteOffset: 11,
						Length:     3,
					},
				},
				Statements: []ast.ScriptStatement{
					&ast.State{
						Keyword: &ast.Token{
							Kind: token.State,
							Location: source.Location{
								ByteOffset: 23,
								Length:     5,
							},
						},
						Name: &ast.Identifier{
							Text: &ast.Token{
								Kind: token.Identifier,
								Location: source.Location{
									ByteOffset: 29,
									Length:     3,
								},
							},
							Normalized: "bar",
							Location: source.Location{
								ByteOffset: 29,
								Length:     3,
							},
						},
						Auto: &ast.Token{
							Kind: token.Auto,
							Location: source.Location{
								ByteOffset: 18,
								Length:     4,
							},
						},
						EndKeyword: &ast.Token{
							Kind: token.EndState,
							Location: source.Location{
								ByteOffset: 36,
								Length:     8,
							},
						},
						Location: source.Location{
							ByteOffset: 18,
							Length:     26,
						},
					},
				},
				Location: source.Location{
					ByteOffset: 0,
					Length:     44,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f := &source.File{Text: []byte(test.input)}
			got, err := parser.Parse(f)
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
	cmpopts.IgnoreFields(source.Location{}, "File", "StartLine", "StartColumn", "EndLine", "EndColumn", "PreambleLength", "PostambleLength"),
}
