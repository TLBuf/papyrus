package symbol

import (
	"fmt"
	"iter"
	"slices"
	"strings"

	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/types"
)

// Scope defines a scope of named entities.
type Scope struct {
	resolver *types.Resolver
	parent   *Scope
	children []*Scope
	symbols  map[key]*Symbol
	node     ast.Node
	kind     scopeKind
}

type scopeKind uint16

const (
	globalScope scopeKind = iota
	scriptScope
	stateScope
	functionScope
	eventScope
	propertyScope
	ifScope
	elseScope
	elseIfScope
	whileScope
)

// Class identifies the class of symbol being looked up.
type Class uint16

const (
	// ScriptClass identifies the symbol being looked up as a [Script].
	ScriptClass = Class(ScriptKind)
	// StateClass identifies the symbol being looked up as a [State].
	StateClass = Class(StateKind)
	// FunctionClass identifies the symbol being looked up as function-typed,
	// specifically a symbol of kind [Function] or [Event].
	FunctionClass = Class(FunctionKind | EventKind)
	// ValueClass identifies the symbol being looked up as value-typed,
	// specifically a symbol of kind [Property], [Variable], or [Parameter].
	ValueClass = Class(PropertyKind | VariableKind | ParameterKind)
)

func (c Class) String() string {
	switch c {
	case ScriptClass:
		return "Script Class"
	case StateClass:
		return "State Class"
	case FunctionClass:
		return "Function Class"
	case ValueClass:
		return "Value Class"
	default:
		return fmt.Sprintf("Unknown Class (%d)", c)
	}
}

type namespace uint8

const (
	valueNamespace    = namespace(ScriptClass | ValueClass)
	stateNamespace    = namespace(StateClass)
	functionNamespace = namespace(FunctionClass)
)

type key struct {
	name      string
	namespace namespace
}

// NewGlobalScope returns an empty global (root) scope.
func NewGlobalScope() *Scope {
	var resolver types.Resolver
	return &Scope{
		resolver: &resolver,
		kind:     globalScope,
	}
}

// Lookup returns the symbol associated with the name in
// this scope or a parent scope or nil if there is not one.
func (s *Scope) Lookup(name string, class Class) (*Symbol, error) {
	var namespace namespace
	switch class {
	case ScriptClass, ValueClass:
		namespace = valueNamespace
	case StateClass:
		namespace = stateNamespace
	case FunctionClass:
		namespace = functionNamespace
	default:
		return nil, fmt.Errorf("%v is not a valid class", class)
	}
	key := key{name: normalize(name), namespace: namespace}
	symbol := s.resolve(key)
	if symbol == nil {
		return nil, fmt.Errorf("%v symbol named %q not found", class, name)
	}
	return symbol, nil
}

func (s *Scope) resolve(key key) *Symbol {
	if e, ok := s.symbols[key]; ok {
		return e
	}
	if s.parent != nil {
		return s.parent.resolve(key)
	}
	return nil
}

// Scope creates an anonymous scope for a node, inserts it into this
// scope, and returns it or an error if this could not be done.
//
// Only the following nodes create anonymous scopes:
//
//   - [ast.If]
//   - [ast.ElseIf]
//   - [ast.Else]
//   - [ast.While]
func (s *Scope) Scope(node ast.Node) (scope *Scope, err error) {
	scope = &Scope{
		resolver: s.resolver,
		parent:   s,
	}
	switch node := node.(type) {
	case *ast.If:
		scope.kind = ifScope
	case *ast.ElseIf:
		scope.kind = elseIfScope
	case *ast.Else:
		scope.kind = elseScope
	case *ast.While:
		scope.kind = whileScope
	default:
		return nil, fmt.Errorf("%v does not define an anonymous scope", node)
	}
	s.children = append(s.children, scope)
	return scope, nil
}

// Symbol creates a symbol for a node and inserts it into this scope
// and returns the symbol or an error if this could not be done.
//
// Only the following nodes create symbols:
//
//   - [ast.Script]
//   - [ast.State]
//   - [ast.Function]
//   - [ast.Event]
//   - [ast.Property]
//   - [ast.Variable]
//   - [ast.Parameter]
func (s *Scope) Symbol(node ast.Node) (symbol *Symbol, err error) {
	var symbolKey key
	switch node := node.(type) {
	case *ast.Script:
		symbol, err = s.insertScript(node)
		symbolKey = key{name: symbol.normalized, namespace: valueNamespace}
	case *ast.State:
		symbol, err = s.insertState(node)
		symbolKey = key{name: symbol.normalized, namespace: stateNamespace}
	case *ast.Function:
		symbol, err = s.insertFunction(node)
		symbolKey = key{name: symbol.normalized, namespace: functionNamespace}
	case *ast.Event:
		symbol, err = s.insertEvent(node)
		symbolKey = key{name: symbol.normalized, namespace: functionNamespace}
	case *ast.Property:
		symbol, err = s.insertProperty(node)
		symbolKey = key{name: symbol.normalized, namespace: valueNamespace}
	case *ast.Variable:
		symbol, err = s.insertVariable(node)
		symbolKey = key{name: symbol.normalized, namespace: valueNamespace}
	case *ast.Parameter:
		symbol, err = s.insertParameter(node)
		symbolKey = key{name: symbol.normalized, namespace: valueNamespace}
	default:
		return nil, fmt.Errorf("%v does not define a symbol", node)
	}
	if err != nil {
		return nil, err
	}
	if existing, ok := s.symbols[symbolKey]; ok {
		return nil, fmt.Errorf("scope contains conflicting symbol: %v", existing)
	}
	s.symbols[symbolKey] = symbol
	if symbol.scope != nil {
		s.children = append(s.children, symbol.scope)
	}
	return symbol, nil
}

func (s *Scope) insertScript(node *ast.Script) (*Symbol, error) {
	if s.kind != globalScope {
		return nil, fmt.Errorf("insert to scope defined by %v, scripts can only belong to the global scope", s.node)
	}
	typ, err := s.resolver.Resolve(node)
	if err != nil {
		return nil, fmt.Errorf("type resolution: %w", err)
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       ScriptKind,
		typ:        typ,
		name:       node.Name.Text,
		normalized: typ.Normalized(),
		node:       node,
	}
	symbol.scope = &Scope{
		resolver: s.resolver,
		parent:   s,
		symbols:  make(map[key]*Symbol),
		kind:     scriptScope,
		node:     node,
	}
	emptyState := &Symbol{
		enclosing: symbol.scope,
		kind:      StateKind,
	}
	emptyState.scope = &Scope{
		resolver: s.resolver,
		parent:   symbol.scope,
		kind:     stateScope,
		node:     node,
	}
	symbol.scope.children = append(symbol.scope.children, emptyState.scope)
	symbol.scope.symbols[key{name: "", namespace: stateNamespace}] = emptyState
	return symbol, nil
}

func (s *Scope) insertState(node *ast.State) (*Symbol, error) {
	if s.kind == globalScope {
		return nil, fmt.Errorf("insert to global scope, states can only belong to a script scope")
	}
	if s.kind != scriptScope {
		return nil, fmt.Errorf("insert to scope defined by %v, states can only belong to a script scope", s.node)
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       StateKind,
		name:       node.Name.Text,
		normalized: normalize(node.Name.Text),
		node:       node,
	}
	symbol.scope = &Scope{
		resolver: s.resolver,
		parent:   s,
		kind:     stateScope,
		node:     node,
	}
	return symbol, nil
}

func (s *Scope) insertFunction(node *ast.Function) (*Symbol, error) {
	if s.kind == globalScope {
		return nil, fmt.Errorf("insert to global scope, functions can only belong to a state scope")
	}
	if s.kind != stateScope {
		return nil, fmt.Errorf("insert to scope defined by %v, functions can only belong to a state scope", s.node)
	}
	typ, err := s.resolver.Resolve(node)
	if err != nil {
		return nil, fmt.Errorf("type resolution: %w", err)
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       FunctionKind,
		typ:        typ,
		name:       node.Name.Text,
		normalized: typ.Normalized(),
		node:       node,
	}
	if len(node.NativeLocations) == 0 {
		symbol.scope = &Scope{
			resolver: s.resolver,
			parent:   s,
			kind:     functionScope,
			node:     node,
		}
	}
	return symbol, nil
}

func (s *Scope) insertEvent(node *ast.Event) (*Symbol, error) {
	if s.kind == globalScope {
		return nil, fmt.Errorf("insert to global scope, events can only belong to a state scope")
	}
	if s.kind != stateScope {
		return nil, fmt.Errorf("insert to scope defined by %v, events can only belong to a state scope", s.node)
	}
	typ, err := s.resolver.Resolve(node)
	if err != nil {
		return nil, fmt.Errorf("type resolution: %w", err)
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       EventKind,
		typ:        typ,
		name:       node.Name.Text,
		normalized: typ.Normalized(),
		node:       node,
	}
	if len(node.NativeLocations) == 0 {
		symbol.scope = &Scope{
			resolver: s.resolver,
			parent:   s,
			kind:     eventScope,
			node:     node,
		}
	}
	return symbol, nil
}

func (s *Scope) insertProperty(node *ast.Property) (*Symbol, error) {
	if s.kind == globalScope {
		return nil, fmt.Errorf("insert to global scope, properties can only belong to a script scope")
	}
	if s.kind != scriptScope {
		return nil, fmt.Errorf("insert to scope defined by %v, properties can only belong to a script scope", s.node)
	}
	typ, err := s.resolver.Resolve(node)
	if err != nil {
		return nil, fmt.Errorf("type resolution: %w", err)
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       PropertyKind,
		typ:        typ,
		name:       node.Name.Text,
		normalized: typ.Normalized(),
		node:       node,
	}
	if node.Kind == ast.Full {
		symbol.scope = &Scope{
			resolver: s.resolver,
			parent:   s,
			kind:     propertyScope,
			node:     node,
		}
	}
	return symbol, nil
}

func (s *Scope) insertVariable(node *ast.Variable) (*Symbol, error) {
	if s.kind == globalScope {
		return nil, fmt.Errorf("insert to global scope, variables can only belong to a script, function, or event scope")
	}
	switch s.kind {
	case scriptScope, functionScope, eventScope, ifScope, elseIfScope, elseScope, whileScope:
		return nil, fmt.Errorf("insert to scope defined by %v, variables can only belong to a script, function, event, if, else if, else, or while scope", s.node)
	}
	typ, err := s.resolver.Resolve(node)
	if err != nil {
		return nil, fmt.Errorf("type resolution: %w", err)
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       VariableKind,
		typ:        typ,
		name:       node.Name.Text,
		normalized: typ.Normalized(),
		node:       node,
	}
	return symbol, nil
}

func (s *Scope) insertParameter(node *ast.Parameter) (*Symbol, error) {
	if s.kind == globalScope {
		return nil, fmt.Errorf("insert to global scope, parameters can only belong to a function or event scope")
	}
	switch s.kind {
	case functionScope, eventScope:
		return nil, fmt.Errorf("insert to scope defined by %v, parameters can only belong to a function or event scope", s.node)
	}
	typ, err := s.resolver.Resolve(node)
	if err != nil {
		return nil, fmt.Errorf("type resolution: %w", err)
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       ParameterKind,
		typ:        typ,
		name:       node.Name.Text,
		normalized: typ.Normalized(),
		node:       node,
	}
	return symbol, nil
}

// Parent returns the scope that encloses this
// one or nil if this is the Global scope.
func (s *Scope) Parent() *Scope {
	return s.parent
}

// Children returns an iterator over the scopes
// that are defined directly within this scope.
func (s *Scope) Children() iter.Seq[*Scope] {
	return slices.Values(s.children)
}

func normalize(name string) string {
	return strings.ToLower(name)
}
