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
			Expressions: make(map[ast.Expression]types.Type),
			Values:      make(map[ast.Literal]literal.Value),
			Scopes:      make(map[ast.Node]*symbol.Scope),
			Global:      global,
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
		c.info.Symbols[script] = sym
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
		c.scope = scriptSymbol.Scope()
		if c.state, err = scriptSymbol.Scope().Lookup("", symbol.StateClass); err != nil {
			c.log.Append(
				issue.New(intenalInvalidState, script.File, script.Location()).WithDetail("Script empty state lookup: %v", err),
			)
		}
		c.scope = c.state.Scope()
	}

	// Handle imported types.
	return success
}

func (c *checker) State(node *ast.State) bool {
	sym, err := c.script.Scope().Symbol(node)
	if err != nil {
		c.log.Append(issue.New(errorStateNameCollision, c.file(), node.Name.Location()).WithDetail("%v", err))
		return false
	}
	c.info.Symbols[node] = sym
	c.info.Scopes[node] = sym.Scope()
	prev := c.scope
	defer func() {
		c.scope = prev
	}()
	c.scope = sym.Scope()
	success := true
	for _, invokable := range node.Invokables {
		switch invokable := invokable.(type) {
		case *ast.CommentStatement:
			// Nothing
		case *ast.Function:
			success = c.Function(invokable) && success
		case *ast.Event:
			success = c.Event(invokable) && success
		default:
			c.log.Append(issue.New(intenalInvalidState, c.file(), node.Name.Location()).WithDetail("Unknown state invokable: %v", invokable))
			success = false
		}
	}
	return success
}

func (c *checker) Property(node *ast.Property) bool {
	sym, err := c.script.Scope().Symbol(node)
	if err != nil {
		c.log.Append(issue.New(errorScriptValueNameCollision, c.file(), node.Name.Location()).WithDetail("%v", err))
		return false
	}
	c.info.Symbols[node] = sym
	success := true
	if sym.Scope() != nil {
		// Full Property
		c.info.Scopes[node] = sym.Scope()
		prev := c.scope
		defer func() {
			c.scope = prev
		}()
		c.scope = sym.Scope()
		if node.Get != nil {
			success = c.Function(node.Get) && success
		}
		if node.Set != nil {
			success = c.Function(node.Set) && success
		}
	}
	return success
}

func (c *checker) ScriptVariable(node *ast.Variable) bool {
	sym, err := c.script.Scope().Symbol(node)
	if err != nil {
		c.log.Append(issue.New(errorScriptValueNameCollision, c.file(), node.Name.Location()).WithDetail("%v", err))
		return false
	}
	c.info.Symbols[node] = sym
	if node.Value != nil {
		return c.Literal(node.Value.(ast.Literal))
	}
	return true
}

func (c *checker) Function(node *ast.Function) bool {
	sym, err := c.scope.Symbol(node)
	if err != nil {
		c.log.Append(issue.New(errorFunctionNameCollision, c.file(), node.Name.Location()).WithDetail("%v", err))
		return false
	}
	c.info.Symbols[node] = sym
	c.info.Scopes[node] = sym.Scope()
	prev := c.scope
	defer func() {
		c.scope = prev
	}()
	c.scope = sym.Scope()
	success := true
	for _, p := range node.Parameters() {
		success = c.Parameter(p) && success
	}
	for _, s := range node.Statements {
		success = c.FunctionStatement(s) && success
	}
	return success
}

func (c *checker) Event(node *ast.Event) bool {
	sym, err := c.scope.Symbol(node)
	if err != nil {
		c.log.Append(issue.New(errorFunctionNameCollision, c.file(), node.Name.Location()).WithDetail("%v", err))
		return false
	}
	c.info.Symbols[node] = sym
	c.info.Scopes[node] = sym.Scope()
	prev := c.scope
	defer func() {
		c.scope = prev
	}()
	c.scope = sym.Scope()
	success := true
	for _, p := range node.Parameters() {
		success = c.Parameter(p) && success
	}
	for _, s := range node.Statements {
		success = c.FunctionStatement(s) && success
	}
	return success
}

func (c *checker) Parameter(node *ast.Parameter) bool {
	sym, err := c.scope.Symbol(node)
	if err != nil {
		c.log.Append(issue.New(errorParameterNameCollision, c.file(), node.Name.Location()).WithDetail("%v", err))
		return false
	}
	c.info.Symbols[node] = sym
	if node.DefaultValue != nil && !c.Literal(node.DefaultValue) {
		return false
	}
	return true
}

func (*checker) FunctionStatement(ast.FunctionStatement) bool {
	return false
}

func (*checker) Literal(ast.Literal) bool {
	return false
}

// file returns the current source file being checked.
func (c *checker) file() *source.File {
	return c.script.Node().(*ast.Script).File
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
