package parser_test

import (
	"testing"

	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/parser"
	"github.com/TLBuf/papyrus/source"
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
				KeywordLocation: source.NewLocation(0, 10),
				Name: &ast.Identifier{
					Text:         "Foo",
					NodeLocation: source.NewLocation(11, 3),
				},
				NodeLocation: source.NewLocation(0, 14),
			},
		},
		{
			name:  "extends",
			input: "ScriptName Foo Extends Bar",
			want: &ast.Script{
				KeywordLocation: source.NewLocation(0, 10),
				Name: &ast.Identifier{
					Text:         "Foo",
					NodeLocation: source.NewLocation(11, 3),
				},
				ExtendsLocation: source.NewLocation(15, 7),
				Parent: &ast.Identifier{
					Text:         "Bar",
					NodeLocation: source.NewLocation(23, 3),
				},
				NodeLocation: source.NewLocation(0, 26),
			},
		},
		{
			name:  "hidden",
			input: "ScriptName Foo Hidden",
			want: &ast.Script{
				KeywordLocation: source.NewLocation(0, 10),
				Name: &ast.Identifier{
					Text:         "Foo",
					NodeLocation: source.NewLocation(11, 3),
				},
				HiddenLocations: []source.Location{source.NewLocation(15, 6)},
				NodeLocation:    source.NewLocation(0, 21),
			},
		},
		{
			name:  "conditional",
			input: "ScriptName Foo Conditional",
			want: &ast.Script{
				KeywordLocation: source.NewLocation(0, 10),
				Name: &ast.Identifier{
					Text:         "Foo",
					NodeLocation: source.NewLocation(11, 3),
				},
				ConditionalLocations: []source.Location{source.NewLocation(15, 11)},
				NodeLocation:         source.NewLocation(0, 26),
			},
		},
		{
			name:  "hidden_conditional",
			input: "ScriptName Foo Hidden Conditional",
			want: &ast.Script{
				KeywordLocation: source.NewLocation(0, 10),
				Name: &ast.Identifier{
					Text:         "Foo",
					NodeLocation: source.NewLocation(11, 3),
				},
				HiddenLocations:      []source.Location{source.NewLocation(15, 6)},
				ConditionalLocations: []source.Location{source.NewLocation(22, 11)},
				NodeLocation:         source.NewLocation(0, 33),
			},
		},
		{
			name:  "conditional_hidden",
			input: "ScriptName Foo Conditional Hidden",
			want: &ast.Script{
				KeywordLocation: source.NewLocation(0, 10),
				Name: &ast.Identifier{
					Text:         "Foo",
					NodeLocation: source.NewLocation(11, 3),
				},
				HiddenLocations:      []source.Location{source.NewLocation(27, 6)},
				ConditionalLocations: []source.Location{source.NewLocation(15, 11)},
				NodeLocation:         source.NewLocation(0, 33),
			},
		},
		{
			name:  "many_flags",
			input: "ScriptName Foo Conditional Hidden Conditional Hidden",
			want: &ast.Script{
				KeywordLocation: source.NewLocation(0, 10),
				Name: &ast.Identifier{
					Text:         "Foo",
					NodeLocation: source.NewLocation(11, 3),
				},
				HiddenLocations: []source.Location{
					source.NewLocation(27, 6), source.NewLocation(46, 6),
				},
				ConditionalLocations: []source.Location{
					source.NewLocation(15, 11), source.NewLocation(34, 11),
				},
				NodeLocation: source.NewLocation(0, 52),
			},
		},
		{
			name:  "extends_many_flags",
			input: "ScriptName Foo Extends Bar Hidden Conditional Hidden Conditional",
			want: &ast.Script{
				KeywordLocation: source.NewLocation(0, 10),
				Name: &ast.Identifier{
					Text:         "Foo",
					NodeLocation: source.NewLocation(11, 3),
				},
				ExtendsLocation: source.NewLocation(15, 7),
				Parent: &ast.Identifier{
					Text:         "Bar",
					NodeLocation: source.NewLocation(23, 3),
				},
				HiddenLocations: []source.Location{
					source.NewLocation(27, 6), source.NewLocation(46, 6),
				},
				ConditionalLocations: []source.Location{
					source.NewLocation(34, 11), source.NewLocation(53, 11),
				},
				NodeLocation: source.NewLocation(0, 64),
			},
		},
		{
			name: "import",
			input: `ScriptName Foo
			Import Bar`,
			want: &ast.Script{
				KeywordLocation: source.NewLocation(0, 10),
				Name: &ast.Identifier{
					Text:         "Foo",
					NodeLocation: source.NewLocation(11, 3),
				},
				Statements: []ast.ScriptStatement{
					&ast.Import{
						KeywordLocation: source.NewLocation(18, 6),
						Name: &ast.Identifier{
							Text:         "Bar",
							NodeLocation: source.NewLocation(25, 3),
						},
					},
				},
				NodeLocation: source.NewLocation(0, 28),
			},
		},
		{
			name: "state",
			input: `ScriptName Foo
			State Bar
			EndState`,
			want: &ast.Script{
				KeywordLocation: source.NewLocation(0, 10),
				Name: &ast.Identifier{
					Text:         "Foo",
					NodeLocation: source.NewLocation(11, 3),
				},
				Statements: []ast.ScriptStatement{
					&ast.State{
						StartKeywordLocation: source.NewLocation(18, 5),
						Name: &ast.Identifier{
							Text:         "Bar",
							NodeLocation: source.NewLocation(24, 3),
						},
						EndKeywordLocation: source.NewLocation(31, 8),
					},
				},
				NodeLocation: source.NewLocation(0, 39),
			},
		},
		{
			name: "state_auto",
			input: `ScriptName Foo
			Auto State Bar
			EndState`,
			want: &ast.Script{
				KeywordLocation: source.NewLocation(0, 10),
				Name: &ast.Identifier{
					Text:         "Foo",
					NodeLocation: source.NewLocation(11, 3),
				},
				Statements: []ast.ScriptStatement{
					&ast.State{
						IsAuto:               true,
						StartKeywordLocation: source.NewLocation(23, 5),
						Name: &ast.Identifier{
							Text:         "Bar",
							NodeLocation: source.NewLocation(29, 3),
						},
						AutoLocation:       source.NewLocation(18, 4),
						EndKeywordLocation: source.NewLocation(36, 8),
					},
				},
				NodeLocation: source.NewLocation(0, 44),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := source.NewFile("test.psc", []byte(test.input))
			if err != nil {
				t.Fatalf("NewFile() returned an unexpected error: %v", err)
			}
			got, err := parser.Parse(f)
			if err != nil {
				t.Errorf("ParseScript() returned an unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("ParseScript() returned nil")
			}
			if diff := cmp.Diff(test.want, got, ignoreFields...); diff != "" {
				t.Errorf("ParseScript() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

var ignoreFields = []cmp.Option{
	// lexer_test pretty thoroughly covers these fields, if the
	// byte offset and length match, that's sufficient for this test.
	cmpopts.IgnoreFields(
		ast.Script{},
		"File",
	),
}
