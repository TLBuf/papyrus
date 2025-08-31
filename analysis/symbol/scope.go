// Package symbol provides the API for tracking and looking up symbols.
package symbol

import (
	"errors"
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

var (
	// ErrNotFound indicates that a lookup for a symbol failed.
	ErrNotFound = errors.New("not found")
	// ErrAlreadyExists indicates that a symbol could not be inserted into a scope
	// because the scope already defines that symbol.
	ErrAlreadyExists = errors.New("already exists")
	// ErrSymbolNotSupported indicates that a symbol could not be inserted into a
	// scope because the scope does not support symbols of that class.
	ErrSymbolNotSupported = errors.New("symbol not supported by scope")
	// ErrNoSymbol indicates that a symbol could not be inserted into a scope
	// because the node does not define a symbol, i.e. it's not one of the
	// following:
	//
	//   - [ast.Script]
	//   - [ast.State]
	//   - [ast.Function]
	//   - [ast.Event]
	//   - [ast.Property]
	//   - [ast.Variable]
	//   - [ast.Parameter]
	ErrNoSymbol = errors.New("does not define a symbol")
	// ErrNoAnonymousScope indicates that an anonymous scope could not be inserted
	// into a scope because the node does not define one, i.e. it's not one of the
	// following:
	//
	//   - [ast.If]
	//   - [ast.ElseIf]
	//   - [ast.Else]
	//   - [ast.While]
	ErrNoAnonymousScope = errors.New("does not define an anonymous scope")
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
		return "Script"
	case StateClass:
		return "State"
	case FunctionClass:
		return "Function"
	case ValueClass:
		return "Value"
	default:
		return fmt.Sprintf("Unknown (%d)", c)
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

// Lookup returns the symbol associated with the name in this scope or a parent
// scope or and error wrapping [ErrNotFound].
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
		return nil, fmt.Errorf("%w", ErrNotFound)
	}
	key := key{name: normalize(name), namespace: namespace}
	symbol := s.resolve(key)
	if symbol == nil {
		return nil, fmt.Errorf("%w", ErrNotFound)
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

// AnonymousScope creates an anonymous scope for a node, inserts it into this
// scope, and returns it or an error that wraps [ErrNoAnonymousScope] if the
// node does not define an anonymous scope.
func (s *Scope) AnonymousScope(node ast.Node) (scope *Scope, err error) {
	scope = &Scope{
		resolver: s.resolver,
		parent:   s,
	}
	switch node.(type) {
	case *ast.If:
		scope.kind = ifScope
	case *ast.ElseIf:
		scope.kind = elseIfScope
	case *ast.Else:
		scope.kind = elseScope
	case *ast.While:
		scope.kind = whileScope
	default:
		return nil, fmt.Errorf("%w", ErrNoAnonymousScope)
	}
	s.children = append(s.children, scope)
	return scope, nil
}

// Symbol creates a symbol for a node (and scope if appropriate) and inserts it
// into this scope and returns the symbol or an error that wraps one of the
// following errors if this fails:
//
//   - [ErrNotFound] if a script node extends another
//     script that has not had a symbol created for it
//   - [ErrNoSymbol] if a node does not define a symbol
//   - [ErrAlreadyExists] if the symbol with the same
//     name and class already exists in the scope
//   - [ErrSymbolNotSupported] if this scope does
//     not support the node's class of symbol
func (s *Scope) Symbol(node ast.Node) (symbol *Symbol, err error) {
	var symbolKey key
	var class Class
	switch node := node.(type) {
	case *ast.Script:
		symbol, err = s.insertScript(node)
		symbolKey = key{name: symbol.normalized, namespace: valueNamespace}
		class = ScriptClass
	case *ast.State:
		symbol, err = s.insertState(node)
		symbolKey = key{name: symbol.normalized, namespace: stateNamespace}
		class = StateClass
	case *ast.Function:
		symbol, err = s.insertFunction(node)
		symbolKey = key{name: symbol.normalized, namespace: functionNamespace}
		class = FunctionClass
	case *ast.Event:
		symbol, err = s.insertEvent(node)
		symbolKey = key{name: symbol.normalized, namespace: functionNamespace}
		class = FunctionClass
	case *ast.Property:
		symbol, err = s.insertProperty(node)
		symbolKey = key{name: symbol.normalized, namespace: valueNamespace}
		class = ValueClass
	case *ast.Variable:
		symbol, err = s.insertVariable(node)
		symbolKey = key{name: symbol.normalized, namespace: valueNamespace}
		class = ValueClass
	case *ast.Parameter:
		symbol, err = s.insertParameter(node)
		symbolKey = key{name: symbol.normalized, namespace: valueNamespace}
		class = ValueClass
	default:
		return nil, fmt.Errorf("%w", ErrNoSymbol)
	}
	if err != nil {
		return nil, err
	}
	if existing, ok := s.symbols[symbolKey]; ok {
		return nil, fmt.Errorf("[%s] %s collides with %v: %w", class, symbolKey.name, existing, ErrAlreadyExists)
	}
	s.symbols[symbolKey] = symbol
	if symbol.scope != nil {
		s.children = append(s.children, symbol.scope)
	}
	return symbol, nil
}

func (s *Scope) insertScript(node *ast.Script) (*Symbol, error) {
	if s.kind != globalScope {
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
	}
	typ, err := s.resolveType(node)
	if err != nil {
		return nil, err
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
	if s.kind == globalScope || s.kind != scriptScope {
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
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
	if s.kind == globalScope || s.kind != stateScope {
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
	}
	typ, err := s.resolveType(node)
	if err != nil {
		return nil, err
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
	if s.kind == globalScope || s.kind != stateScope {
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
	}
	typ, err := s.resolveType(node)
	if err != nil {
		return nil, err
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
	if s.kind == globalScope || s.kind != scriptScope {
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
	}
	typ, err := s.resolveType(node)
	if err != nil {
		return nil, err
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
	switch s.kind {
	case globalScope, scriptScope, functionScope, eventScope, ifScope, elseIfScope, elseScope, whileScope:
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
	}
	typ, err := s.resolveType(node)
	if err != nil {
		return nil, err
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
	switch s.kind {
	case globalScope, functionScope, eventScope:
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
	}
	typ, err := s.resolveType(node)
	if err != nil {
		return nil, err
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

func (s *Scope) resolveType(node ast.Node) (types.Type, error) {
	typ, err := s.resolver.Resolve(node)
	if err != nil {
		var nferr types.NotFoundError
		_ = errors.As(err, &nferr) // Will never return types.ErrNotTyped.
		return nil, fmt.Errorf("parent %q: %w", nferr.Name, ErrNotFound)
	}
	return typ, nil
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
