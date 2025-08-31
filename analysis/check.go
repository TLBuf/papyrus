// Package analysis defines the Papyrus static analysis API.
package analysis

import (
	"slices"
	"strings"

	"github.com/TLBuf/papyrus/analysis/symbol"
	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/issue"
	"github.com/TLBuf/papyrus/literal"
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/types"
)

// Check performs type-checking over some number of
// scripts and collated summarized type information.
func Check(log *issue.Log, scripts ...*ast.Script) (*Info, bool) {
	global := symbol.NewGlobalScope()
	checker := &checker{
		log:    log,
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
	if ok := checker.check(scripts); !ok {
		return nil, false
	}
	return checker.info, true
}

type checker struct {
	log       *issue.Log
	info      *Info
	global    *symbol.Scope
	typeNames map[string]types.Type
	script    *symbol.Symbol
	state     *symbol.Symbol
	scope     *symbol.Scope
}

func (c *checker) check(scripts []*ast.Script) bool {
	success := true
	// Build script types.
	if ok := c.sortScripts(scripts); !ok {
		return false
	}
	for _, script := range scripts {
		sym, err := c.global.Symbol(script)
		if err != nil {
			c.log.Append(
				issue.New(intenalInvalidState, script.File, script.Location()).WithDetail("Symbol creation: %v", err),
			)
			success = false
			continue
		}
		c.info.Scopes[script] = sym.Scope()
		c.typeNames[sym.Normalized()] = sym.Type()
	}
	if !success {
		return false // Don't try to keep going.
	}
	// Build all other types, except imports.
	for _, script := range scripts {
		scriptSymbol, err := c.global.Lookup(script.Name.Text, symbol.ScriptClass)
		if err != nil {
			c.log.Append(
				issue.New(intenalInvalidState, script.File, script.Location()).WithDetail("Script symbol lookup: %v", err),
			)
			success = false
			continue
		}
		_ = scriptSymbol
	}

	// Handle imported types.
	return success
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

func (c *checker) sortScripts(scripts []*ast.Script) bool {
	success := true
	slices.SortFunc(scripts, func(a, b *ast.Script) int {
		return strings.Compare(normalize(a.Name.Text), normalize(b.Name.Text))
	})
	byName := make(map[string]*ast.Script, len(scripts))
	for _, s := range scripts {
		name := normalize(s.Name.Text)
		if existing, ok := byName[name]; ok {
			c.log.Append(
				issue.New(
					errorScriptNameCollision,
					s.File,
					s.Name.Location(),
				).WithDetail(
					"%s name collides with %s.",
					existing.File.Path(), s.File.Path(),
				),
			)
			success = false
			continue
		}
		byName[name] = s
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
			c.log.Append(
				issue.New(
					errorScriptUnknownParent,
					s.File,
					s.Parent.Location(),
				).WithDetail(
					"%q is not a known script.",
					s.Parent.Text,
				),
			)
			success = false
		}
		children[parent] = append(children[parent], s)
	}
	for len(queue) > 0 {
		s := queue[0]
		queue = queue[1:]
		sorted, seen[s] = append(sorted, s), struct{}{}
		for _, child := range children[s] {
			if _, ok := seen[child]; ok {
				c.log.Append(issue.New(errorScriptCycle, s.File, s.Parent.Location()))
				success = false
			}
			sorted, seen[child] = append(sorted, child), struct{}{}
		}
	}
	copy(scripts, sorted)
	return success
}

func normalize(name string) string {
	return strings.ToLower(name)
}
