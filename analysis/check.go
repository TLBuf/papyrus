// Package analysis defines the Papyrus static analysis API.
package analysis

import (
	"fmt"
	"slices"
	"strings"

	"github.com/TLBuf/papyrus/analysis/symbol"
	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/literal"
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/types"
)

// Error defines an error raised by the type checker.
type Error struct {
	// The underlying error.
	Err error
	// Location identifies the place in the source that caused the error.
	Location source.Location
}

// Error implments the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Location, e.Err)
}

// Unwrap returns the underlying error.
func (e Error) Unwrap() error {
	return e.Err
}

func newError(location source.Location, msg string, args ...any) Error {
	return Error{
		Err:      fmt.Errorf(msg, args...),
		Location: location,
	}
}

// Check performs type-checking over some number of
// scripts and collated summarized type information.
func Check(scripts ...*ast.Script) (*Info, error) {
	global := symbol.NewGlobalScope()
	checker := &checker{
		global: global,
		info: &Info{
			Types:  make(map[ast.Expression]types.Type),
			Values: make(map[ast.Literal]literal.Value),
			Scopes: make(map[ast.Node]*symbol.Scope),
			Global: global,
		},
		typeNames: make(map[string]types.Type),
		scope:     global,
	}
	for _, t := range []types.Type{types.Bool, types.BoolArray, types.Int, types.IntArray, types.Float, types.FloatArray, types.String, types.StringArray} {
		checker.typeNames[normalize(t.Name())] = t
	}
	if err := sortScripts(scripts); err != nil {
		return nil, fmt.Errorf("check script inheritance: %w", err)
	}
	for _, script := range scripts {
		sym, err := global.Symbol(script)
		if err != nil {
			return nil, fmt.Errorf("%v symbol: %w", script, err)
		}
		checker.info.Scopes[script] = sym.Scope()
		checker.typeNames[sym.Normalized()] = sym.Type()
	}
	if err := checker.check(scripts); err != nil {
		return nil, err
	}
	return checker.info, nil
}

type checker struct {
	info      *Info
	global    *symbol.Scope
	typeNames map[string]types.Type
	script    *symbol.Symbol
	state     *symbol.Symbol
	scope     *symbol.Scope
}

// file returns the current source file being checked.
func (c *checker) file() *source.File {
	return c.script.Node().(*ast.Script).File
}

// self returns the type of the current script being checked.
func (c *checker) self() *types.Object {
	return c.script.Type().(*types.Object)
}

func (c *checker) enterScope(s *symbol.Scope) {
	c.scope = s
}

func (c *checker) leaveScope() {
	c.scope = c.scope.Parent()
}

func (c *checker) check(scripts []*ast.Script) error {
	// TODO: Implement.
	return nil
}

func sortScripts(scripts []*ast.Script) error {
	slices.SortFunc(scripts, func(a, b *ast.Script) int {
		return strings.Compare(normalize(a.Name.Text), normalize(b.Name.Text))
	})
	byName := make(map[string]*ast.Script, len(scripts))
	for _, s := range scripts {
		byName[normalize(s.Name.Text)] = s
	}
	seen := make(map[*ast.Script]struct{}, len(scripts))
	children := make(map[*ast.Script][]*ast.Script, len(scripts))
	queue, sorted := make([]*ast.Script, 0, len(scripts)), make([]*ast.Script, 0, len(scripts))
	for _, s := range scripts {
		if s.Parent == nil {
			queue = append(queue, s)
			continue
		}
		parent, ok := byName[normalize(s.Parent.Text)]
		if !ok {
			return fmt.Errorf("%v extends unknown script %q", s, s.Parent.Text)
		}
		children[parent] = append(children[parent], s)
	}
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		sorted, seen[s] = append(sorted, s), struct{}{}
		for _, child := range children[s] {
			if _, ok := seen[child]; ok {
				return fmt.Errorf("%v extends a script that forms a cycle", child)
			}
			sorted, seen[child] = append(sorted, child), struct{}{}
		}
	}
	copy(scripts, sorted)
	return nil
}

func normalize(name string) string {
	return strings.ToLower(name)
}
