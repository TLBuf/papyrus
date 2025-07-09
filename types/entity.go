package types

import (
	"strings"

	"github.com/TLBuf/papyrus/ast"
)

// Entity is a named papyrus entity, e.g. a [Function] or [Property].
type Entity interface {
	// Enclosing returns the scope in which this object is defined.
	Enclosing() *Scope
	// Scope returns the scope defined by this
	// entity or nil if it doesn't define one.
	Scope() *Scope
	// Type returns the type of this object.
	Type() Type
	// Name returns the declared name for this object.
	Name() string
	// Node returns the [ast.Node] that defines this object.
	Node() ast.Node
}

// Function is a entity associated with a [ast.Function].
type Function struct {
	entity
}

// NewFunction returns a new function entity.
func NewFunction(name string, node *ast.Function, signature *Signature, enclosing *Scope) *Function {
	return &Function{
		entity: entity{
			enclosing: enclosing,
			scope:     NewScope(enclosing),
			typ:       signature,
			name:      name,
			node:      node,
		},
	}
}

// Global returns whether or not this function is global.
func (f *Function) Global() bool {
	return len(f.node.(*ast.Function).GlobalLocations) > 0
}

// Signature returns the signature of this function.
func (f *Function) Signature() *Signature {
	if s, ok := f.typ.(*Signature); ok {
		return s
	}
	return nil
}

func (f *Function) String() string {
	var sb strings.Builder
	f.Signature().write(&sb, f.name)
	return sb.String()
}

var _ Entity = (*Function)(nil)

// Event is a entity associated with an [ast.Event].
type Event struct {
	entity
}

// NewEvent returns a new event entity.
func NewEvent(name string, node *ast.Event, signature *Signature, enclosing *Scope) *Event {
	return &Event{
		entity: entity{
			enclosing: enclosing,
			scope:     NewScope(enclosing),
			typ:       signature,
			name:      name,
			node:      node,
		},
	}
}

// Signature returns the signature of this event.
func (e *Event) Signature() *Signature {
	if s, ok := e.typ.(*Signature); ok {
		return s
	}
	return nil
}

func (e *Event) String() string {
	var sb strings.Builder
	e.Signature().write(&sb, e.name)
	return sb.String()
}

// Property is a entity associated with a [ast.Property].
type Property struct {
	entity
}

// NewProperty returns a new property entity.
func NewProperty(name string, node *ast.Property, typ Type, enclosing *Scope) *Property {
	return &Property{
		entity: entity{
			enclosing: enclosing,
			scope:     NewScope(enclosing),
			typ:       typ,
			name:      name,
			node:      node,
		},
	}
}

// Readable returns true if this property is readable.
func (p *Property) Readable() bool {
	if n, ok := p.node.(*ast.Property); ok {
		return n.Kind != ast.Full || n.Get != nil
	}
	return false
}

// Writable returns true if this property is writable.
func (p *Property) Writable() bool {
	if n, ok := p.node.(*ast.Property); ok {
		return n.Kind == ast.Auto || n.Set != nil
	}
	return false
}

// Variable is a entity associated with a [ast.Variable].
type Variable struct {
	entity
}

// NewVariable returns a new variable entity.
func NewVariable(name string, node *ast.Variable, typ Type, enclosing *Scope) *Variable {
	return &Variable{
		entity: entity{
			enclosing: enclosing,
			scope:     NewScope(enclosing),
			typ:       typ,
			name:      name,
			node:      node,
		},
	}
}

// Parameter is a entity associated with a [ast.Parameter].
type Parameter struct {
	entity
}

// NewParameter returns a new parameter entity.
func NewParameter(name string, node *ast.Parameter, typ Type, enclosing *Scope) *Parameter {
	return &Parameter{
		entity: entity{
			enclosing: enclosing,
			scope:     NewScope(enclosing),
			typ:       typ,
			name:      name,
			node:      node,
		},
	}
}

type entity struct {
	enclosing *Scope
	scope     *Scope
	typ       Type
	name      string
	node      ast.Node
}

// Enclosing returns the scope in which this object is defined.
func (o *entity) Enclosing() *Scope {
	return o.enclosing
}

// Scope returns the scope defined by this
// entity or nil if it doesn't define one.
func (o *entity) Scope() *Scope {
	return o.scope
}

// Type returns the type of this object.
func (o *entity) Type() Type {
	return o.typ
}

// Name returns the declared name for this object.
func (o *entity) Name() string {
	return o.name
}

// Node returns the [ast.Node] that defines this object.
func (o *entity) Node() ast.Node {
	return o.node
}
