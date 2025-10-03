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
	for _, t := range []types.Type{types.BoolType, types.BoolArrayType, types.IntType, types.IntArrayType, types.FloatType, types.FloatArrayType, types.StringType, types.StringArrayType} {
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
	for _, script := range scripts {
		scriptSymbol := c.global.Resolve(script.Name.Text, symbol.Values)
		if scriptSymbol == nil {
			c.failInFile(internalInvalidState, script.File, script.Location(), "Script symbol lookup: %q", script.Name.Text)
			continue
		}
		c.scope = scriptSymbol.Scope()
		c.script = scriptSymbol
		// Scan for imports and resolve global functions, but do not insert them
		// into the script scope yet. We need to verify that there are no functions
		// or events declared in this script that may override them first.
		imports := make(map[string]*symbol.Symbol)
		clashes := make(map[string]struct{})
		for _, stmt := range script.Statements {
			imp, ok := stmt.(*ast.Import)
			if !ok {
				continue
			}
			imported := c.global.ResolveKind(imp.Name.Text, symbol.Script)
			if imported == nil {
				c.failWithDetail(errorImportUnknown, imp.Name.Location(), "Unknown script %q", imp.Name.Text)
				continue
			}
			for sym := range imported.Scope().Symbols() {
				if sym.Kind() != symbol.Function {
					continue
				}
				function := sym
				// revive:disable-next-line:unchecked-type-assertion
				if !function.Type().(*types.Invokable).Global() {
					continue
				}
				// Detect name clashes: If two imports have global functions with the same
				// name, neither can be imported. This is not an error in the script.
				if _, ok := clashes[function.Normalized()]; ok {
					continue
				}
				if _, ok := imports[function.Normalized()]; ok {
					clashes[function.Normalized()] = struct{}{}
					delete(imports, function.Normalized())
					continue
				}
				imports[function.Normalized()] = function
			}
		}
		if len(imports) > 0 {
			// Remove any imports that are overridden
			// by functions or events in this script.
			for _, stmt := range script.Statements {
				switch stmt := stmt.(type) {
				case *ast.Function:
					delete(imports, normalize(stmt.Name.Text))
				case *ast.Event:
					delete(imports, normalize(stmt.Name.Text))
				case *ast.State:
					for _, invokable := range stmt.Invokables {
						switch invokable := invokable.(type) {
						case *ast.Function:
							delete(imports, normalize(invokable.Name.Text))
						case *ast.Event:
							delete(imports, normalize(invokable.Name.Text))
						}
					}
				}
			}
			// Remove any imports that are already defined in
			// the script scope (e.g. from a parent script).
			for name := range imports {
				if c.scope.Resolve(name, symbol.Invokables) != nil {
					delete(imports, name)
				}
			}
		}
		// Inject global functions into this script's scope.
		for _, function := range imports {
			_, err := c.scope.Symbol(function.Node())
			if err != nil {
				c.failWithDetail(
					internalInvalidState,
					script.Location(),
					"insert imported global function into script scope: %v",
					err,
				)
				return
			}
		}
		// Process the rest of the script.
		var states []*ast.State
		for _, statement := range script.Statements {
			switch statement := statement.(type) {
			case *ast.Import, *ast.CommentStatement, *ast.ErrorStatement:
				// Nothing
			case *ast.Variable:
				c.ScriptVariable(statement)
			case *ast.Property:
				c.Property(statement)
			case *ast.Function:
				c.Function(statement)
			case *ast.Event:
				c.Event(statement)
			case *ast.State:
				states = append(states, statement)
			default:
				c.failWithDetail(internalInvalidState, statement.Location(), "Unknown script statement: %v", statement)
			}
		}
		// Process states last so we can verify that all
		// functions and events exist in the empty state.
		for _, state := range states {
			c.State(state)
		}
	}
}

func (c *checker) State(node *ast.State) {
	sym, err := c.script.Scope().Symbol(node)
	if err != nil {
		c.failWithDetail(errorStateNameCollision, node.Name.Location(), "%v", err)
		return
	}
	c.info.Symbols[node] = sym
	c.info.Scopes[node] = sym.Scope()
	parent := c.scope
	defer func() {
		c.scope = parent
	}()
	c.scope = sym.Scope()
	for _, invokable := range node.Invokables {
		switch invokable := invokable.(type) {
		case *ast.CommentStatement, *ast.ErrorStatement:
			// Nothing
		case *ast.Function:
			c.Function(invokable)
			function := c.scope.LookupKind(normalize(invokable.Name.Text), symbol.Function)
			if function == nil {
				continue // Error already reported.
			}
			if function.Type().(*types.Invokable).Native() {
				c.fail(errorStateNativeFunction, invokable.SignatureLocation())
			}
			if function.Type().(*types.Invokable).Global() {
				c.fail(errorStateGlobalFunction, invokable.SignatureLocation())
			}
			sym := parent.Lookup(normalize(invokable.Name.Text), symbol.Invokables)
			if sym == nil {
				c.failWithDetail(errorStateFunctionNoEmptyStateDefintion, invokable.SignatureLocation(), "Function %s", invokable.Name.Text)
				continue
			}
			if sym.Kind() == symbol.Event {
				//revive:disable-next-line:unchecked-type-assertion
				event := sym.Node().(*ast.Event)
				c.log.Append(issue.New(errorStateFunctionEventNameCollision, c.file(), invokable.SignatureLocation()).
					WithDetail("Function %s", invokable.Name.Text).
					AppendRelated(c.file(), event.SignatureLocation(), "Clashes with Event %s", event.Name.Text))
				c.failed = true
				continue
			}
			if !sym.Type().IsIdentical(function.Type()) {
				//revive:disable-next-line:unchecked-type-assertion
				other := sym.Node().(*ast.Function)
				c.log.Append(issue.New(errorStateFunctionSignatureMismatch, c.file(), invokable.SignatureLocation()).
					WithDetail("Signature %s", function.Type()).
					AppendRelated(c.file(), other.SignatureLocation(), "Does not match %s", sym.Type()))
				c.failed = true
			}
		case *ast.Event:
			c.Event(invokable)
			event := c.scope.LookupKind(normalize(invokable.Name.Text), symbol.Function)
			if event == nil {
				continue // Error already reported.
			}
			if event.Type().(*types.Invokable).Native() {
				c.fail(errorStateNativeEvent, invokable.SignatureLocation())
			}
			sym := parent.Lookup(normalize(invokable.Name.Text), symbol.Invokables)
			if sym == nil {
				c.failWithDetail(errorStateEventNoEmptyStateDefintion, invokable.SignatureLocation(), "Event %s", invokable.Name.Text)
				continue
			}
			if sym.Kind() == symbol.Function {
				//revive:disable-next-line:unchecked-type-assertion
				function := sym.Node().(*ast.Function)
				c.log.Append(issue.New(errorStateEventFunctionNameCollision, c.file(), invokable.SignatureLocation()).
					WithDetail("Event %s", invokable.Name.Text).
					AppendRelated(c.file(), function.SignatureLocation(), "Clashes with Function %s", function.Name.Text))
				c.failed = true
				continue
			}
			if !sym.Type().IsIdentical(event.Type()) {
				//revive:disable-next-line:unchecked-type-assertion
				other := sym.Node().(*ast.Event)
				c.log.Append(issue.New(errorStateEventSignatureMismatch, c.file(), invokable.SignatureLocation()).
					WithDetail("Signature %s", event.Type()).
					AppendRelated(c.file(), other.SignatureLocation(), "Does not match %s", sym.Type()))
				c.failed = true
			}
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
	decl, err := c.types.Resolve(node.Type)
	if err != nil {
		c.failWithDetail(
			errorTypeReferencesUnknownScript,
			node.Type.Location(),
			"%q is not a known script",
			node.Type.Name.Text,
		)
	}
	typ := c.Literal(node.Value.(ast.Literal))
	if typ != nil {
		c.info.Expressions[node.Value] = typ
	}
	if typ != nil && decl != nil && !typ.IsAssignable(decl) {
		c.failWithDetail(errorVariableTypeMismatch, node.Value.Location(), "%v is not assignable to %v", typ, decl)
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
	var walker returnWalker
	walker.Walk(node.Statements)
	returnType := sym.Type().(*types.Invokable).ReturnType()
	for _, ret := range walker.returns {
		if returnType == nil && ret.Value == nil {
			continue
		}
		if returnType == nil && ret.Value != nil {
			c.fail(errorFunctionReturnValueUnexpected, ret.Value.Location())
			continue
		}
		if returnType != nil && ret.Value == nil {
			c.failWithDetail(errorFunctionReturnValueMissing, ret.Location(), "Expected to return a %v", returnType)
			continue
		}
		retType, ok := c.info.Expressions[ret.Value]
		if !ok {
			continue
		}
		if !retType.IsAssignable(returnType) {
			c.failWithDetail(
				errorFunctionReturnTypeMismatch,
				ret.Value.Location(),
				"%v is not assignable to %v",
				retType,
				returnType,
			)
		}
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
	var walker returnWalker
	walker.Walk(node.Statements)
	for _, ret := range walker.returns {
		if ret.Value != nil {
			c.fail(errorEventReturnValueUnexpected, ret.Value.Location())
		}
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
	decl, err := c.types.Resolve(node.Type)
	if err != nil {
		c.failWithDetail(
			errorTypeReferencesUnknownScript,
			node.Type.Location(),
			"%q is not a known script",
			node.Type.Name.Text,
		)
	}
	typ := c.Literal(node.DefaultValue)
	if typ != nil {
		c.info.Expressions[node.DefaultValue] = typ
	}
	if typ != nil && decl != nil && !typ.IsAssignable(decl) {
		c.failWithDetail(
			errorParameterDefaultValueTypeMismatch,
			node.DefaultValue.Location(),
			"%v is not assignable to %v",
			typ,
			decl,
		)
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
	assigneeType := c.Expression(node.Assignee)
	if assigneeType != nil {
		c.info.Expressions[node.Assignee] = assigneeType
	}
	valueType := c.Expression(node.Value)
	if valueType != nil {
		c.info.Expressions[node.Value] = valueType
	}
	if assigneeType == nil || valueType == nil {
		return
	}
	switch node.Kind {
	case ast.Assign:
		if !valueType.IsAssignable(assigneeType) {
			c.failWithDetail(
				errorAssignmentTypeMismatch,
				node.Value.Location(),
				"%v is not assignable to %v",
				valueType,
				assigneeType,
			)
		}
	case ast.AssignAdd:
		// Either string concatenation or numeric addition.
		if assigneeType.IsIdentical(types.StringType) {
			break // String concatenation.
		}
		fallthrough
	case ast.AssignSubtract, ast.AssignMultiply, ast.AssignDivide:
		if !assigneeType.IsIdentical(types.IntType) && !assigneeType.IsIdentical(types.FloatType) {
			c.failWithDetail(
				errorAssignmentArithmeticAssigneeNotNumeric,
				node.Assignee.Location(),
				"Assignee is typed: %v",
				assigneeType,
			)
		}
		if !valueType.IsIdentical(types.IntType) && !valueType.IsIdentical(types.FloatType) {
			c.failWithDetail(errorAssignmentArithmeticValueNotNumeric, node.Value.Location(), "Value is typed: %v", valueType)
		}
	case ast.AssignModulo:
		if !assigneeType.IsIdentical(types.IntType) {
			c.failWithDetail(
				errorAssignmentModuloAssigneeNotInt,
				node.Assignee.Location(),
				"Assignee is typed: %v",
				assigneeType,
			)
		}
		if !valueType.IsIdentical(types.IntType) {
			c.failWithDetail(errorAssignmentModuloValueNotInt, node.Value.Location(), "Value is typed: %v", valueType)
		}
	}
}

func (c *checker) If(node *ast.If) {
	if typ := c.Expression(node.Condition); typ != nil {
		if typ.IsAssignable(types.BoolType) {
			c.info.Expressions[node.Condition] = typ
		} else {
			c.failWithDetail(errorIfConditionNotBool, node.Condition.Location(), "Expression is typed: %v", typ)
		}
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
			if typ.IsAssignable(types.BoolType) {
				c.info.Expressions[node.Condition] = typ
			} else {
				c.failWithDetail(errorElseIfConditionNotBool, node.Condition.Location(), "Expression is typed: %v", typ)
			}
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
		if typ.IsAssignable(types.BoolType) {
			c.info.Expressions[node.Condition] = typ
		} else {
			c.failWithDetail(errorWhileConditionNotBool, node.Condition.Location(), "Expression is typed: %v", typ)
		}
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
	decl, err := c.types.Resolve(node.Type)
	if err != nil {
		c.failWithDetail(
			errorTypeReferencesUnknownScript,
			node.Type.Location(),
			"%q is not a known script",
			node.Type.Name.Text,
		)
	}
	typ := c.Expression(node.Value)
	if typ != nil {
		c.info.Expressions[node.Value] = typ
	}
	if typ != nil && decl != nil && !typ.IsAssignable(decl) {
		c.failWithDetail(errorVariableTypeMismatch, node.Value.Location(), "%v is not assignable to %v", typ, decl)
	}
}

func (c *checker) Expression(node ast.Expression) types.Type {
	switch node := node.(type) {
	case *ast.Access:
		return c.ValueAccess(node)
	case *ast.ArrayCreation:
		return c.ArrayCreation(node)
	case *ast.Binary:
		return c.Binary(node)
	case *ast.Call:
		return c.Call(node)
	case *ast.Cast:
		return c.Cast(node)
	case *ast.Identifier:
		return c.Identifier(node)
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

func (c *checker) ValueAccess(node *ast.Access) types.Type {
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
			return nil
		}
		return types.IntType
	case *types.Object:
		sym, ok := c.info.Symbols[typ.Node()]
		if !ok {
			c.failWithDetail(internalInvalidState, node.Value.Location(), "Script symbol lookup: %s", typ.Name())
			return nil
		}
		pSym := sym.Scope().Resolve(lookup, symbol.Values)
		if pSym == nil {
			c.failWithDetail(errorUnknownProperty, node.Name.Location(), "%s does not define a property named %q", sym.Name(), node.Name.Text)
			return nil
		}
		if pSym.Kind() == symbol.Variable && pSym != c.script {
			c.failWithDetail(errorCannotAccessVariable, node.Value.Location(), "%q is a variable, not a property", pSym.Name())
			return nil
		}
		return pSym.Type()
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

func (c *checker) FunctionAccess(node *ast.Access) *symbol.Symbol {
	typ := c.Expression(node.Value)
	if typ == nil {
		return nil
	}
	c.info.Expressions[node.Value] = typ
	lookup := normalize(node.Name.Text)
	switch typ := typ.(type) {
	case *types.Object:
		sym, ok := c.info.Symbols[typ.Node()]
		if !ok {
			c.failWithDetail(internalInvalidState, node.Value.Location(), "Script symbol lookup: %s", typ.Name())
			return nil
		}
		lSym := sym.Scope().Resolve(lookup, symbol.Invokables)
		if lSym == nil {
			c.failWithDetail(errorAccessFunctionUnknown, node.Name.Location(), "%s does not define or inherit a function named %s", sym.Name(), node.Name.Text)
			return nil
		}
		if lSym.Kind() == symbol.Event {
			//revive:disable-next-line:unchecked-type-assertion
			c.log.Append(issue.New(errorAccessEvent, c.file(), node.Value.Location()).
				WithDetail("%s is not a function", sym.Name()).
				AppendRelated(c.file(), sym.Node().(*ast.Event).SignatureLocation(), "%s is an event", sym.Name()))
			c.failed = true
			return nil
		}
		//revive:disable-next-line:unchecked-type-assertion
		return lSym
	default:
		c.failWithDetail(errorCallRecieverNotObject, node.Value.Location(), "Value of type %v", typ)
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
	return types.ArrayOf(typ.(types.Scalar))
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
	if left == nil || right == nil {
		return nil
	}
	var typ types.Type
	switch node.Kind {
	case ast.Equal, ast.NotEqual:
		if !left.IsEquatable(right) {
			c.failWithDetail(
				errorBinaryOperandsNotEquatable,
				node.Location(),
				"%v is not equatable to %v (nor vice versa)",
				left,
				right,
			)
		}
		typ = types.BoolType
	case ast.Less, ast.LessOrEqual, ast.Greater, ast.GreaterOrEqual:
		if !left.IsComparable(right) {
			c.failWithDetail(
				errorBinaryOperandsNotComparable,
				node.Location(),
				"%v is not comparable to %v (nor vice versa)",
				left,
				right,
			)
		}
		typ = types.BoolType
	case ast.LogicalAnd, ast.LogicalOr:
		if !left.IsAssignable(types.BoolType) {
			c.failWithDetail(errorBinaryLogicalOperandNotBool, node.LeftOperand.Location(), "Left operand is typed: %v", left)
		}
		if !right.IsAssignable(types.BoolType) {
			c.failWithDetail(
				errorBinaryLogicalOperandNotBool,
				node.RightOperand.Location(),
				"Right operand is typed: %v",
				right,
			)
		}
		typ = types.BoolType
	case ast.Add:
		// Addition can be numeric or string concatenation (when the left is a string).
		if left.IsIdentical(types.StringType) {
			// String concatenation, right can technically be anything.
			typ = types.StringType
			break
		}
		fallthrough // Numeric addition.
	case ast.Subtract, ast.Multiply, ast.Divide:
		if !left.IsIdentical(types.IntType) && !left.IsIdentical(types.FloatType) {
			c.failWithDetail(
				errorBinaryArithmeticOperandNotNumeric,
				node.LeftOperand.Location(),
				"Left operand is typed: %v",
				left,
			)
			break
		}
		if !right.IsIdentical(types.IntType) && !right.IsIdentical(types.FloatType) {
			c.failWithDetail(
				errorBinaryArithmeticOperandNotNumeric,
				node.RightOperand.Location(),
				"Right operand is typed: %v",
				right,
			)
			break
		}
		// Result type is float if either side is float, int otherwise.
		if left.IsIdentical(types.FloatType) || right.IsIdentical(types.FloatType) {
			typ = types.FloatType
		} else {
			typ = types.IntType
		}
	case ast.Modulo:
		if !left.IsIdentical(types.IntType) {
			c.failWithDetail(errorBinaryModuloOperandNotInt, node.LeftOperand.Location(), "Left operand is typed: %v", left)
			break
		}
		if !right.IsIdentical(types.IntType) {
			c.failWithDetail(
				errorBinaryModuloOperandNotInt,
				node.RightOperand.Location(),
				"Right operand is typed: %v",
				right,
			)
			break
		}
		typ = types.IntType
	}
	return typ
}

func (c *checker) Call(node *ast.Call) types.Type {
	// The expression denoting the function being called must be either:
	//   - An Access
	//   - An Identifier
	//
	// We don't use Expression here because we need to resolve certain identifiers
	// as function calls, not variable accesses.
	var function *symbol.Symbol
	switch expr := node.Function.(type) {
	case *ast.Access:
		if function = c.FunctionAccess(expr); function == nil {
			// Error already reported.
			return nil
		}
	case *ast.Identifier:
		sym := c.scope.Resolve(normalize(expr.Text), symbol.Invokables)
		if sym == nil {
			c.failWithDetail(errorCallFunctionUnknown, expr.Location(), "%s does not define or inherit a function named %s", sym.Name(), expr.Text)
			return nil
		}
		if sym.Kind() == symbol.Event {
			//revive:disable-next-line:unchecked-type-assertion
			c.log.Append(issue.New(errorCallEvent, c.file(), expr.Location()).
				WithDetail("%s is not a function", expr.Text).
				AppendRelated(c.file(), sym.Node().(*ast.Event).SignatureLocation(), "%s is an event", sym.Name()))
			c.failed = true
			return nil
		}
		function = sym
	default:
		c.failWithDetail(errorCallMalformed, expr.Location(), "%T is not an identifier or access ending in an identifier", expr)
		return nil
	}
	//revive:disable-next-line:unchecked-type-assertion
	funcType := function.Type().(*types.Invokable)
	//revive:disable-next-line:unchecked-type-assertion
	funcNode := function.Node().(*ast.Function)
	// Match arguments to parameters.
	named := false
	for _, arg := range node.Arguments {
		if arg.Name != nil {
			named = true
			break
		}
	}
	if !named {
		// Positional arguments only, match in order. If there are too few
		// arguments, all remaining parameters must have defaults.
		for i, param := range funcNode.Parameters() {
			if i >= len(node.Arguments) {
				if param.DefaultValue != nil {
					continue
				}
				// Missing argument.
				c.log.Append(issue.New(errorCallArgumentMissing, c.file(), node.Location()).
					WithDetail("Expected an argument at position %d for parameter %s", i+1, param.Name.Text).
					AppendRelated(c.file(), param.Location(), "%s does not have a default value", param.Name.Text))
				c.failed = true
				continue
			}
			arg := node.Arguments[i]
			argType := c.Expression(arg.Value)
			if argType == nil {
				// Error already reported.
				continue
			}
			c.info.Expressions[arg.Value] = argType
			paramType := funcType.Parameters()[i].Type()
			if !argType.IsAssignable(paramType) {
				c.log.Append(issue.New(errorCallArgumentTypeMismatch, c.file(), arg.Value.Location()).
					WithDetail("%v is not assignable to %v", argType, paramType).
					AppendRelated(c.file(), param.Location(), "Parameter %s is of type %v", param.Name.Text, paramType))
				c.failed = true
			}
		}
		// Log issues for extra arguments.
		for i := len(funcNode.Parameters()); i < len(node.Arguments); i++ {
			arg := node.Arguments[i]
			c.log.Append(issue.New(errorCallArgumentExtra, c.file(), arg.Value.Location()).
				WithDetail("Argument at position %d exceeds the number of declared parameters", i+1).
				AppendRelated(
					c.file(),
					funcNode.SignatureLocation(),
					"%s declared %d parameter(s)",
					funcNode.Name.Text,
					len(funcNode.Parameters()),
				))
			c.failed = true
		}
	} else {
		// There's at least one named argument. To match arguments to parameters:
		//   - Align positional arguments in order until we reach a named argument.
		//   - Match named arguments to parameters by name.
		//   - Any unmatched parameters must have defaults.
		paramNodes := funcNode.Parameters()
		paramTypes := funcType.Parameters()
		params := make(map[string]int)
		for i, p := range paramTypes {
			params[p.Normalized()] = i
		}
		matched := make(map[string]*ast.Argument)
		var firstNamedArgument *ast.Argument
		for i, arg := range node.Arguments {
			if arg.Name != nil {
				// Named argument
				if firstNamedArgument == nil {
					firstNamedArgument = arg
				}
				argName := normalize(arg.Name.Text)
				paramIndex, ok := params[argName]
				if !ok {
					c.log.Append(issue.New(errorCallArgumentUnknownNamed, c.file(), arg.Name.Location()).
						WithDetail("Function %s does not have a parameter named %s", funcNode.Name.Text, arg.Name.Text).
						AppendRelated(c.file(), funcNode.SignatureLocation(), "Function definition"))
					c.failed = true
					continue
				}
				if existing, matched := matched[argName]; matched {
					c.log.Append(issue.New(errorCallArgumentNamedDuplicate, c.file(), arg.Name.Location()).
						WithDetail("Argument for parameter %s provided more than once", arg.Name.Text).
						AppendRelated(c.file(), existing.Location(), "Argument already associated with parameter %s", arg.Name.Text))
					c.failed = true
					continue
				}
				paramType := paramTypes[paramIndex].Type()
				argType := c.Expression(arg.Value)
				if argType == nil {
					// Error already reported.
					continue
				}
				c.info.Expressions[arg.Value] = argType
				if !argType.IsAssignable(paramType) {
					param := paramNodes[paramIndex]
					c.log.Append(issue.New(errorCallArgumentTypeMismatch, c.file(), arg.Value.Location()).
						WithDetail("%v is not assignable to %v", argType, paramType).
						AppendRelated(c.file(), param.Location(), "Parameter %s is of type %v", param.Name.Text, paramType))
					c.failed = true
				}
				matched[argName] = arg
			} else {
				// Positional argument
				if firstNamedArgument != nil {
					c.log.Append(issue.New(errorCallPositionalAfterNamed, c.file(), arg.Location()).
						WithDetail("Argument at position %d", i+1).
						AppendRelated(c.file(), firstNamedArgument.Location(), "First named argument"))
					c.failed = true
					continue
				}
				if i >= len(paramNodes) {
					c.log.Append(issue.New(errorCallArgumentExtra, c.file(), arg.Value.Location()).
						WithDetail("Argument at position %d exceeds the number of declared parameters", i+1).
						AppendRelated(
							c.file(),
							funcNode.SignatureLocation(),
							"%s declared %d parameter(s)",
							funcNode.Name.Text,
							len(paramNodes),
						))
					c.failed = true
					continue
				}
				paramType := paramTypes[i].Type()
				argType := c.Expression(arg.Value)
				if argType == nil {
					continue
				}
				c.info.Expressions[arg.Value] = argType
				if !argType.IsAssignable(paramType) {
					param := paramNodes[i]
					c.log.Append(issue.New(errorCallArgumentTypeMismatch, c.file(), arg.Value.Location()).
						WithDetail("%v is not assignable to %v", argType, paramType).
						AppendRelated(c.file(), param.Location(), "Parameter %s is of type %v", param.Name.Text, paramType))
					c.failed = true
				}
				matched[paramType.Normalized()] = arg
			}
		}
		// Check for unmatched parameters that don't have default values.
		for _, param := range paramNodes {
			paramName := normalize(param.Name.Text)
			if _, matched := matched[paramName]; matched || param.DefaultValue != nil {
				continue
			}
			c.log.Append(issue.New(errorCallArgumentMissing, c.file(), node.Location()).
				WithDetail("Missing argument for parameter %s", param.Name.Text).
				AppendRelated(c.file(), param.Location(), "%s does not have a default value", param.Name.Text))
			c.failed = true
		}
	}
	return funcType.ReturnType()
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
		c.failWithDetail(errorCastNotConvertible, node.Type.Location(), "%v cannot be cast to %v", expr, typ)
		return nil
	}
	return typ
}

func (c *checker) Identifier(node *ast.Identifier) types.Type {
	// In this context, an identifier must resolve to a Value.
	sym := c.scope.Resolve(normalize(node.Text), symbol.Values)
	if sym == nil {
		c.failWithDetail(errorIdentifierUnknown, node.Location(), "%s is not defined in this scope", node.Text)
		return nil
	}
	return sym.Type()
}

func (c *checker) Index(node *ast.Index) types.Type {
	index := c.Expression(node.Index)
	if index != nil {
		c.info.Expressions[node.Index] = index
		if !index.IsIdentical(types.IntType) {
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
		if !typ.IsAssignable(types.BoolType) {
			c.failWithDetail(errorUnaryNegationNotBool, node.Operand.Location(), "Operand is typed: %v", typ)
			return nil
		}
		return types.BoolType
	}
	// Numeric negation. Type of expression must be numeric.
	if typ.IsIdentical(types.IntType) || typ.IsIdentical(types.FloatType) {
		return typ
	}
	c.failWithDetail(errorUnaryNegationNotNumeric, node.Operand.Location(), "Operand is typed: %v", typ)
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

func (c *checker) fail(def *issue.Definition, loc source.Location) {
	c.log.Append(issue.New(def, c.file(), loc))
	c.failed = true
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

type returnWalker struct {
	returns []*ast.Return
}

func (w *returnWalker) Walk(statements []ast.FunctionStatement) {
	for _, s := range statements {
		switch s := s.(type) {
		case *ast.Return:
			w.returns = append(w.returns, s)
		case *ast.If:
			w.Walk(s.Statements)
			for _, elseIf := range s.ElseIfs {
				w.Walk(elseIf.Statements)
			}
			if s.Else != nil {
				w.Walk(s.Else.Statements)
			}
		case *ast.While:
			w.Walk(s.Statements)
		}
	}
}
