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
	// ErrUnknownType indicates that a symbol could not be inserted into a scope
	// because the symbol's type could not be determined.
	ErrUnknownType = errors.New("unknown type")
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

// Namespace defines a boundary around different kinds
// of symbols within which names must be unqiue.
type Namespace uint16

const (
	// Values identifies the symbol namespace that is
	// shared by scripts, properties, variables, and parameters.
	Values = Namespace(Script | Property | Variable | Parameter)
	// Invokables identifies the symbol namespace
	// that is shared by functions and events.
	Invokables = Namespace(Function | Event)
	// States identifies the symbol namespace that is exclusive to states.
	States = Namespace(State)
)

func (n Namespace) String() string {
	switch n {
	case Values:
		return "Value"
	case Invokables:
		return "Invokable"
	case States:
		return "State"
	default:
		return fmt.Sprintf("Unknown Namespace (%d)", n)
	}
}

type key struct {
	name      string
	namespace Namespace
}

// NewGlobalScope returns an empty global (root) scope.
func NewGlobalScope(resolver *types.Resolver) *Scope {
	return &Scope{
		resolver: resolver,
		kind:     globalScope,
	}
}

// ResolveKind returns the symbol associated with the name of a specific
// kind in this scope or a parent scope or nil if the symbol does not exist.
func (s *Scope) ResolveKind(name string, kind Kind) *Symbol {
	var symbol *Symbol
	switch {
	case Namespace(kind)&Values > 0:
		symbol = s.Resolve(name, Values)
	case Namespace(kind)&Invokables > 0:
		symbol = s.Resolve(name, Invokables)
	case Namespace(kind)&States > 0:
		symbol = s.Resolve(name, States)
	default:
		return nil
	}
	if symbol == nil || symbol.kind&kind == 0 {
		return nil
	}
	return symbol
}

// LookupKind returns the symbol associated with the name of a specific
// kind in this scope or nil if the symbol does not exist.
//
// Note: LookupKind does not search parent scopes.
func (s *Scope) LookupKind(name string, kind Kind) *Symbol {
	var symbol *Symbol
	switch {
	case Namespace(kind)&Values > 0:
		symbol = s.Lookup(name, Values)
	case Namespace(kind)&Invokables > 0:
		symbol = s.Lookup(name, Invokables)
	case Namespace(kind)&States > 0:
		symbol = s.Lookup(name, States)
	default:
		return nil
	}
	if symbol == nil || symbol.kind&kind == 0 {
		return nil
	}
	return symbol
}

// Resolve returns the symbol associated with the name in this
// scope or a parent scope or nil if the symbol does not exist.
func (s *Scope) Resolve(name string, namespace Namespace) *Symbol {
	key := key{name: normalize(name), namespace: namespace}
	return s.resolveRecursive(key)
}

// Lookup returns the symbol associated with the name in this scope
// or nil and false if the symbol does not exist in this scope.
//
// Note: Lookup does not search parent scopes.
func (s *Scope) Lookup(name string, namespace Namespace) *Symbol {
	key := key{name: normalize(name), namespace: namespace}
	return s.symbols[key]
}

func (s *Scope) resolveRecursive(key key) *Symbol {
	if symbol, ok := s.symbols[key]; ok {
		return symbol
	}
	if s.parent != nil {
		if s.kind != scriptScope {
			return s.parent.resolveRecursive(key)
		}
		symbol := s.parent.resolveRecursive(key)
		if symbol.kind&Variable == 0 {
			// Variables cannot be seen by other scripts.
			return symbol
		}
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
//   - [ErrUnknownType] if a script node extends another
//     script that has not had a symbol created for it
//   - [ErrNoSymbol] if a node does not define a symbol
//   - [ErrAlreadyExists] if the symbol with the same
//     name and class already exists in the scope
//   - [ErrSymbolNotSupported] if this scope does
//     not support the node's class of symbol
func (s *Scope) Symbol(node ast.Node) (symbol *Symbol, err error) {
	switch node := node.(type) {
	case *ast.Script:
		return s.insertScript(node)
	case *ast.State:
		return s.insertState(node)
	case *ast.Function:
		return s.insertFunction(node)
	case *ast.Event:
		return s.insertEvent(node)
	case *ast.Property:
		return s.insertProperty(node)
	case *ast.Variable:
		return s.insertVariable(node)
	case *ast.Parameter:
		return s.insertParameter(node)
	default:
		return nil, fmt.Errorf("%w", ErrNoSymbol)
	}
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
		kind:       Script,
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
	key := key{name: symbol.normalized, namespace: Values}
	// Only check global scope.
	if existing, ok := s.symbols[key]; ok {
		return nil, fmt.Errorf("%v collides with %v: %w", symbol, existing, ErrAlreadyExists)
	}
	s.symbols[key] = symbol
	s.children = append(s.children, symbol.scope)
	return symbol, nil
}

func (s *Scope) insertState(node *ast.State) (*Symbol, error) {
	if s.kind != scriptScope {
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       State,
		name:       node.Name.Text,
		normalized: normalize(node.Name.Text),
		node:       node,
		scope: &Scope{
			resolver: s.resolver,
			parent:   s,
			kind:     stateScope,
			node:     node,
		},
	}
	key := key{name: symbol.normalized, namespace: States}
	// Only check current script scope, parent can define same states.
	if existing, ok := s.symbols[key]; ok {
		return nil, fmt.Errorf("%v collides with %v: %w", symbol, existing, ErrAlreadyExists)
	}
	s.symbols[key] = symbol
	s.children = append(s.children, symbol.scope)
	return symbol, nil
}

func (s *Scope) insertFunction(node *ast.Function) (*Symbol, error) {
	if s.kind != scriptScope && s.kind != stateScope {
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
	}
	typ, err := s.resolveType(node)
	if err != nil {
		return nil, err
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       Function,
		typ:        typ,
		name:       node.Name.Text,
		normalized: typ.Normalized(),
		node:       node,
		scope: &Scope{
			resolver: s.resolver,
			parent:   s,
			kind:     functionScope,
			node:     node,
		},
	}
	key := key{name: symbol.normalized, namespace: Invokables}
	if existing := s.resolveRecursive(key); existing != nil {
		return nil, fmt.Errorf("%v collides with %v: %w", symbol, existing, ErrAlreadyExists)
	}
	s.symbols[key] = symbol
	s.children = append(s.children, symbol.scope)
	return symbol, nil
}

func (s *Scope) insertEvent(node *ast.Event) (*Symbol, error) {
	if s.kind != scriptScope && s.kind != stateScope {
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
	}
	typ, err := s.resolveType(node)
	if err != nil {
		return nil, err
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       Event,
		typ:        typ,
		name:       node.Name.Text,
		normalized: typ.Normalized(),
		node:       node,
		scope: &Scope{
			resolver: s.resolver,
			parent:   s,
			kind:     eventScope,
			node:     node,
		},
	}
	key := key{name: symbol.normalized, namespace: Invokables}
	if existing := s.resolveRecursive(key); existing != nil {
		return nil, fmt.Errorf("%v collides with %v: %w", symbol, existing, ErrAlreadyExists)
	}
	s.symbols[key] = symbol
	s.children = append(s.children, symbol.scope)
	return symbol, nil
}

func (s *Scope) insertProperty(node *ast.Property) (*Symbol, error) {
	if s.kind != scriptScope {
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
	}
	typ, err := s.resolveType(node)
	if err != nil {
		return nil, err
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       Property,
		typ:        typ,
		name:       node.Name.Text,
		normalized: typ.Normalized(),
		node:       node,
		scope: &Scope{
			resolver: s.resolver,
			parent:   s,
			kind:     propertyScope,
			node:     node,
		},
	}
	key := key{name: symbol.normalized, namespace: Values}
	if existing := s.resolveRecursive(key); existing != nil {
		return nil, fmt.Errorf("%v collides with %v: %w", symbol, existing, ErrAlreadyExists)
	}
	s.symbols[key] = symbol
	s.children = append(s.children, symbol.scope)
	return symbol, nil
}

func (s *Scope) insertVariable(node *ast.Variable) (*Symbol, error) {
	switch s.kind {
	case scriptScope, functionScope, eventScope, ifScope, elseIfScope, elseScope, whileScope:
		// OK.
	default:
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
	}
	typ, err := s.resolveType(node)
	if err != nil {
		return nil, err
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       Variable,
		typ:        typ,
		name:       node.Name.Text,
		normalized: typ.Normalized(),
		node:       node,
	}
	key := key{name: symbol.normalized, namespace: Values}
	if existing := s.resolveRecursive(key); existing != nil {
		return nil, fmt.Errorf("%v collides with %v: %w", symbol, existing, ErrAlreadyExists)
	}
	s.symbols[key] = symbol
	return symbol, nil
}

func (s *Scope) insertParameter(node *ast.Parameter) (*Symbol, error) {
	if s.kind != functionScope && s.kind != eventScope {
		return nil, fmt.Errorf("%w", ErrSymbolNotSupported)
	}
	typ, err := s.resolveType(node)
	if err != nil {
		return nil, err
	}
	symbol := &Symbol{
		enclosing:  s,
		kind:       Parameter,
		typ:        typ,
		name:       node.Name.Text,
		normalized: typ.Normalized(),
		node:       node,
	}
	key := key{name: symbol.normalized, namespace: Values}
	if existing := s.resolveRecursive(key); existing != nil {
		return nil, fmt.Errorf("%v collides with %v: %w", symbol, existing, ErrAlreadyExists)
	}
	s.symbols[key] = symbol
	return symbol, nil
}

func (s *Scope) resolveType(node ast.Node) (types.Type, error) {
	typ, err := s.resolver.Resolve(node)
	if err != nil {
		var nferr types.NotFoundError
		_ = errors.As(err, &nferr) // Will never return types.ErrNotTyped.
		return nil, fmt.Errorf("parent %q: %w", nferr.Name, ErrUnknownType)
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
