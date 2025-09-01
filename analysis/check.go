// Package analysis defines the Papyrus static analysis API.
package analysis

import (
	"slices"
	"strings"

	"github.com/TLBuf/papyrus/analysis/symbol"
	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/issue"
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/types"
	"github.com/TLBuf/papyrus/value"
)

// Check performs type-checking over some number of
// scripts and collated summarized type information.
func Check(log *issue.Log, scripts ...*ast.Script) (*Info, bool) {
	var resolver types.Resolver
	global := symbol.NewGlobalScope(&resolver)
	checker := &checker{
		log:    log,
		types:  &resolver,
		global: global,
		info: &Info{
			Expressions: make(map[ast.Expression]types.Type),
			Values:      make(map[ast.Literal]value.Value),
			Scopes:      make(map[ast.Node]*symbol.Scope),
			Global:      global,
		},
		typeNames: make(map[string]types.Type),
		scope:     global,
	}
	for _, t := range []types.Type{types.Bool, types.BoolArray, types.Int, types.IntArray, types.Float, types.FloatArray, types.String, types.StringArray} {
		checker.typeNames[normalize(t.Name())] = t
	}
	checker.check(scripts)
	return checker.info, checker.failed
}

type checker struct {
	log       *issue.Log
	types     *types.Resolver
	info      *Info
	global    *symbol.Scope
	typeNames map[string]types.Type
	script    *symbol.Symbol
	state     *symbol.Symbol
	scope     *symbol.Scope
	failed    bool
}

func (c *checker) check(scripts []*ast.Script) {
	// Build script types.
	if ok := c.sortScripts(scripts); !ok {
		return
	}
	for _, script := range scripts {
		sym, err := c.global.Symbol(script)
		if err != nil {
			c.failInFile(internalInvalidState, script.File, script.Location(), "Symbol creation: %v", err)
			continue
		}
		c.info.Symbols[script] = sym
		c.info.Scopes[script] = sym.Scope()
		c.typeNames[sym.Normalized()] = sym.Type()
	}
	if c.failed {
		return // Don't try to keep going.
	}
	// Build all other types, except imports.
	for _, script := range scripts {
		scriptSymbol, err := c.global.Lookup(script.Name.Text, symbol.ScriptClass)
		if err != nil {
			c.failInFile(internalInvalidState, script.File, script.Location(), "Script symbol lookup: %v", err)
			continue
		}
		if c.state, err = scriptSymbol.Scope().Lookup("", symbol.StateClass); err != nil {
			c.failInFile(internalInvalidState, script.File, script.Location(), "Script empty state lookup: %v", err)
			continue
		}
		emptyStateScope := c.state.Scope()
		for _, statement := range script.Statements {
			switch statement := statement.(type) {
			case *ast.Import, *ast.CommentStatement, *ast.ErrorStatement:
				// Nothing
			case *ast.Variable:
				c.scope = scriptSymbol.Scope()
				c.ScriptVariable(statement)
			case *ast.Property:
				c.scope = scriptSymbol.Scope()
				c.Property(statement)
			case *ast.State:
				c.scope = scriptSymbol.Scope()
				c.State(statement)
			case *ast.Function:
				c.scope = emptyStateScope
				c.Function(statement)
			case *ast.Event:
				c.scope = emptyStateScope
				c.Event(statement)
			default:
				c.failWithDetail(internalInvalidState, statement.Location(), "Unknown script statement: %v", statement)
			}
		}
	}

	// Handle imported types.
}

func (c *checker) State(node *ast.State) {
	sym, err := c.script.Scope().Symbol(node)
	if err != nil {
		c.failWithDetail(errorStateNameCollision, node.Name.Location(), "%v", err)
		return
	}
	c.info.Symbols[node] = sym
	c.info.Scopes[node] = sym.Scope()
	prev := c.scope
	defer func() {
		c.scope = prev
	}()
	c.scope = sym.Scope()
	for _, invokable := range node.Invokables {
		switch invokable := invokable.(type) {
		case *ast.CommentStatement, *ast.ErrorStatement:
			// Nothing
		case *ast.Function:
			c.Function(invokable)
		case *ast.Event:
			c.Event(invokable)
		default:
			c.failWithDetail(internalInvalidState, invokable.Location(), "Unknown state invokable: %v", invokable)
		}
	}
}

func (c *checker) Property(node *ast.Property) {
	sym, err := c.script.Scope().Symbol(node)
	if err != nil {
		c.failWithDetail(errorValueNameCollision, node.Name.Location(), "%v", err)
		return
	}
	c.info.Symbols[node] = sym
	if sym.Scope() != nil {
		// Full Property
		c.info.Scopes[node] = sym.Scope()
		prev := c.scope
		defer func() {
			c.scope = prev
		}()
		c.scope = sym.Scope()
		if node.Get != nil {
			c.Function(node.Get)
		}
		if node.Set != nil {
			c.Function(node.Set)
		}
	}
}

func (c *checker) ScriptVariable(node *ast.Variable) {
	sym, err := c.script.Scope().Symbol(node)
	if err != nil {
		c.failWithDetail(errorValueNameCollision, node.Name.Location(), "%v", err)
		return
	}
	c.info.Symbols[node] = sym
	if node.Value == nil {
		return
	}
	if typ := c.Literal(node.Value.(ast.Literal)); typ != nil {
		c.info.Expressions[node.Value] = typ
	}
}

func (c *checker) Function(node *ast.Function) {
	sym, err := c.scope.Symbol(node)
	if err != nil {
		c.failWithDetail(errorFunctionNameCollision, node.Name.Location(), "%v", err)
		return
	}
	c.info.Symbols[node] = sym
	c.info.Scopes[node] = sym.Scope()
	prev := c.scope
	defer func() {
		c.scope = prev
	}()
	c.scope = sym.Scope()
	for _, p := range node.Parameters() {
		c.Parameter(p)
	}
	for _, s := range node.Statements {
		c.FunctionStatement(s)
	}
}

func (c *checker) Event(node *ast.Event) {
	sym, err := c.scope.Symbol(node)
	if err != nil {
		c.failWithDetail(errorFunctionNameCollision, node.Name.Location(), "%v", err)
	}
	c.info.Symbols[node] = sym
	c.info.Scopes[node] = sym.Scope()
	prev := c.scope
	defer func() {
		c.scope = prev
	}()
	c.scope = sym.Scope()
	for _, p := range node.Parameters() {
		c.Parameter(p)
	}
	for _, s := range node.Statements {
		c.FunctionStatement(s)
	}
}

func (c *checker) Parameter(node *ast.Parameter) {
	sym, err := c.scope.Symbol(node)
	if err != nil {
		c.failWithDetail(errorParameterNameCollision, node.Name.Location(), "%v", err)
	}
	c.info.Symbols[node] = sym
	if node.DefaultValue == nil {
		return
	}
	if typ := c.Literal(node.DefaultValue); typ != nil {
		c.info.Expressions[node.DefaultValue] = typ
	}
}

func (c *checker) FunctionStatement(node ast.FunctionStatement) {
	switch node := node.(type) {
	case *ast.CommentStatement, *ast.ErrorStatement:
		// Ignored.
	case *ast.ExpressionStatement:
		c.ExpressionStatement(node)
	case *ast.Return:
		c.Return(node)
	case *ast.Assignment:
		c.Assignment(node)
	case *ast.If:
		c.If(node)
	case *ast.While:
		c.While(node)
	case *ast.Variable:
		c.FunctionVariable(node)
	default:
		c.failWithDetail(internalInvalidState, node.Location(), "Unknown function statement: %v", node)
	}
}

func (c *checker) ExpressionStatement(node *ast.ExpressionStatement) {
	if typ := c.Expression(node.Expression); typ != nil {
		c.info.Expressions[node.Expression] = typ
	}
}

func (c *checker) Return(node *ast.Return) {
	if node.Value != nil {
		if typ := c.Expression(node.Value); typ != nil {
			c.info.Expressions[node.Value] = typ
		}
	}
}

func (c *checker) Assignment(node *ast.Assignment) {
	if typ := c.Expression(node.Assignee); typ != nil {
		c.info.Expressions[node.Assignee] = typ
	}
	if typ := c.Expression(node.Value); typ != nil {
		c.info.Expressions[node.Value] = typ
	}
}

func (c *checker) If(node *ast.If) {
	if typ := c.Expression(node.Condition); typ != nil {
		c.info.Expressions[node.Condition] = typ
	}
	scope, err := c.scope.AnonymousScope(node)
	if err != nil {
		c.failWithDetail(internalInvalidState, node.Location(), "Anonymous scope creation: %v", err)
		return
	}
	parent := c.scope
	c.scope = scope
	for _, statement := range node.Statements {
		c.FunctionStatement(statement)
	}
	c.scope = parent
	for _, elseIf := range node.ElseIfs {
		if typ := c.Expression(node.Condition); typ != nil {
			c.info.Expressions[node.Condition] = typ
		}
		scope, err := c.scope.AnonymousScope(elseIf)
		if err != nil {
			c.failWithDetail(internalInvalidState, node.Location(), "Anonymous scope creation: %v", err)
			return
		}
		c.scope = scope
		for _, statement := range node.Statements {
			c.FunctionStatement(statement)
		}
		c.scope = parent
	}
	if node.Else != nil {
		scope, err := c.scope.AnonymousScope(node.Else)
		if err != nil {
			c.failWithDetail(internalInvalidState, node.Location(), "Anonymous scope creation: %v", err)
			return
		}
		c.scope = scope
		for _, statement := range node.Statements {
			c.FunctionStatement(statement)
		}
		c.scope = parent
	}
}

func (c *checker) While(node *ast.While) {
	if typ := c.Expression(node.Condition); typ != nil {
		c.info.Expressions[node.Condition] = typ
	}
	scope, err := c.scope.AnonymousScope(node)
	if err != nil {
		c.failWithDetail(internalInvalidState, node.Location(), "Anonymous scope creation: %v", err)
		return
	}
	parent := c.scope
	c.scope = scope
	for _, statement := range node.Statements {
		c.FunctionStatement(statement)
	}
	c.scope = parent
}

func (c *checker) FunctionVariable(node *ast.Variable) {
	sym, err := c.scope.Symbol(node)
	if err != nil {
		c.failWithDetail(errorValueNameCollision, node.Name.Location(), "%v", err)
		return
	}
	c.info.Symbols[node] = sym
	if node.Value == nil {
		return
	}
	if typ := c.Expression(node.Value); typ != nil {
		c.info.Expressions[node.Value] = typ
	}
}

func (c *checker) Expression(node ast.Expression) types.Type {
	switch node := node.(type) {
	case *ast.Access:
		return c.Access(false, node)
	case *ast.ArrayCreation:
		return c.ArrayCreation(node)
	case *ast.Binary:
		return c.Binary(node)
	case *ast.Call:
		return c.Call(node)
	case *ast.Cast:
		return c.Cast(node)
	case *ast.Identifier:
		return c.Identifier(false, node)
	case *ast.Index:
		return c.Index(node)
	case ast.Literal:
		return c.Literal(node)
	case *ast.Parenthetical:
		return c.Parenthetical(node)
	case *ast.Unary:
		return c.Unary(node)
	}
	c.failWithDetail(internalInvalidState, node.Location(), "Unknown function statement: %v", node)
	return nil
}

func (c *checker) Access(call bool, node *ast.Access) types.Type {
	typ := c.Expression(node.Value)
	if typ == nil {
		return nil
	}
	c.info.Expressions[node.Value] = typ
	lookup := normalize(node.Name.Text)
	switch typ := typ.(type) {
	case *types.Array:
		if lookup != "length" {
			c.failWithDetail(errorInvalidArrayAccess, node.Name.Location(), "Expected 'Length', but encountered %q", node.Name.Text)
		}
		return types.Int
	case *types.Object:
		sym, ok := c.info.Symbols[typ.Node()]
		if !ok {
			c.failWithDetail(internalInvalidState, node.Value.Location(), "Script symbol lookup: %s", typ.Name())
			return nil
		}
		if call {
			// Function access.
			lSym, err := sym.Scope().Lookup(lookup, symbol.FunctionClass)
			if err == nil {
				c.failWithDetail(errorUnknownFunction, node.Name.Location(), "%s does not define a function named %q", sym.Name(), node.Name.Text)
			}
			if lSym.Kind() == symbol.EventKind {
				c.failWithDetail(errorCannotCallEvent, node.Value.Location(), "%q is an event, not a function", lSym.Name())
			}
			return lSym.Type()
		}
		// Property access.
		lSym, err := sym.Scope().Lookup(lookup, symbol.ValueClass)
		if err == nil {
			c.failWithDetail(errorUnknownProperty, node.Name.Location(), "%s does not define a property named %q", sym.Name(), node.Name.Text)
		}
		if lSym.Kind() == symbol.VariableKind && lSym != c.script {
			c.failWithDetail(errorCannotAccessVariable, node.Value.Location(), "%q is a variable, not a property", lSym.Name())
		}
		return lSym.Type()
	case *types.Primitive:
		switch typ.Kind() {
		case types.BoolKind:
			c.failWithDetail(errorCannotAccessBool, node.Location(), "Attempting to access %q", node.Name.Text)
		case types.IntKind:
			c.failWithDetail(errorCannotAccessInt, node.Location(), "Attempting to access %q", node.Name.Text)
		case types.FloatKind:
			c.failWithDetail(errorCannotAccessFloat, node.Location(), "Attempting to access %q", node.Name.Text)
		case types.StringKind:
			c.failWithDetail(errorCannotAccessString, node.Location(), "Attempting to access %q", node.Name.Text)
		default:
			c.failWithDetail(internalInvalidState, node.Value.Location(), "Unknown primitive type: %v", typ)
		}
	case *types.Invokable:
		switch typ.Kind() {
		case types.FunctionKind:
			c.failWithDetail(errorCannotAccessFunction, node.Location(), "Attempting to access %q", node.Name.Text)
		case types.EventKind:
			c.failWithDetail(errorCannotAccessEvent, node.Location(), "Attempting to access %q", node.Name.Text)
		default:
			c.failWithDetail(internalInvalidState, node.Value.Location(), "Unknown invokable type: %v", typ)
		}
	case types.None:
		c.failWithDetail(errorCannotAccessNone, node.Location(), "Attempting to access %q", node.Name.Text)
	default:
		c.failWithDetail(internalInvalidState, node.Value.Location(), "Unknown type: %v", typ)
	}
	return nil
}

func (c *checker) ArrayCreation(node *ast.ArrayCreation) types.Type {
	typ, err := c.types.Resolve(node.Type)
	if err != nil {
		c.failWithDetail(
			errorTypeReferencesUnknownScript,
			node.Type.Location(),
			"%q is not a known script",
			node.Type.Name.Text,
		)
		return nil
	}
	c.info.Expressions[node] = typ
	return typ
}

func (c *checker) Binary(node *ast.Binary) types.Type {
	left := c.Expression(node.LeftOperand)
	if left != nil {
		c.info.Expressions[node.LeftOperand] = left
	}
	right := c.Expression(node.RightOperand)
	if right != nil {
		c.info.Expressions[node.RightOperand] = right
	}
	return nil
}

func (c *checker) Call(node *ast.Call) types.Type {
	return nil
}

func (c *checker) Cast(node *ast.Cast) types.Type {
	expr := c.Expression(node.Value)
	if expr != nil {
		c.info.Expressions[node.Value] = expr
	}
	typ, err := c.types.Resolve(node.Type)
	if err != nil {
		c.failWithDetail(
			errorTypeReferencesUnknownScript,
			node.Type.Location(),
			"%q is not a known script",
			node.Type.Name.Text,
		)
		return nil
	}
	if !expr.IsConvertible(typ) {
		c.failWithDetail(errorCastNotConvertible, node.Type.Location(), "Values of type %v cannot be cast to %v", expr, typ)
		return nil
	}
	return typ
}

func (c *checker) Identifier(call bool, node *ast.Identifier) types.Type {
	return nil
}

func (c *checker) Index(node *ast.Index) types.Type {
	index := c.Expression(node.Index)
	if index != nil {
		c.info.Expressions[node.Index] = index
		if !index.IsIdentical(types.Int) {
			c.failWithDetail(errorIndexNotInt, node.Index.Location(), "Index expression is typed: %v", index)
			index = nil
		}
	}
	val := c.Expression(node.Value)
	if val != nil {
		c.info.Expressions[node.Value] = val
		if _, ok := val.(*types.Array); !ok {
			c.failWithDetail(errorIndexTargetNotArray, node.Value.Location(), "Indexed value is typed: %v", index)
			val = nil
		}
	}
	if index == nil || val == nil {
		return nil
	}
	return val.(*types.Array).Element()
}

func (c *checker) Parenthetical(node *ast.Parenthetical) types.Type {
	if typ := c.Expression(node.Value); typ != nil {
		c.info.Expressions[node.Value] = typ
		return typ
	}
	return nil
}

func (c *checker) Unary(node *ast.Unary) types.Type {
	typ := c.Expression(node.Operand)
	if typ == nil {
		return nil
	}
	c.info.Expressions[node.Operand] = typ
	if node.Kind == ast.LogicalNot {
		return typ
	}
	// Type of expression must be numeric.
	if typ.IsIdentical(types.Int) || typ.IsIdentical(types.Float) {
		return typ
	}
	c.failWithDetail(errorNegationNotNumeric, node.Operand.Location(), "Operand is typed: %v", typ)
	return nil
}

func (c *checker) Literal(node ast.Literal) types.Type {
	val, err := value.New(node)
	if err != nil {
		def := internalInvalidState
		switch node.(type) {
		case *ast.BoolLiteral:
			def = errorBoolParseLiteral
		case *ast.IntLiteral:
			def = errorIntParseLiteral
		case *ast.FloatLiteral:
			def = errorFloatParseLiteral
		case *ast.StringLiteral:
			def = errorStringParseLiteral
		}
		c.failWithDetail(def, node.Location(), "%v", err)
	}
	c.info.Values[node] = val
	return val.Type()
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

func (c *checker) failWithDetail(def *issue.Definition, loc source.Location, msg string, args ...any) {
	c.log.Append(issue.New(def, c.file(), loc).WithDetail(msg, args...))
	c.failed = true
}

func (c *checker) failInFile(def *issue.Definition, file *source.File, loc source.Location, msg string, args ...any) {
	c.log.Append(issue.New(def, file, loc).WithDetail(msg, args...))
	c.failed = true
}

func normalize(name string) string {
	return strings.ToLower(name)
}
