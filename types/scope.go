package types

import (
	"iter"
	"slices"
)

// Scope defines a scope of named entities.
type Scope struct {
	parent   *Scope
	children []*Scope
	elements map[string]Entity
}

// NewScope returns an empty scope that is a child of the parent scope.
func NewScope(parent *Scope) *Scope {
	s := &Scope{
		parent: parent,
	}
	if parent != nil {
		parent.children = append(parent.children, s)
	}
	return s
}

// Parent returns the scope that encloses this
// one or nil if this is the root scope.
func (s *Scope) Parent() *Scope {
	return s.parent
}

// Children returns an iterator over the scopes
// that are defined directly within this scope.
func (s *Scope) Children() iter.Seq[*Scope] {
	return slices.Values(s.children)
}

// Lookup returns the entity associated with the name in
// this scope or a parent scope or nil if there is not one.
func (s *Scope) Lookup(name string) Entity {
	return s.resolve(normalize(name))
}

func (s *Scope) resolve(name string) Entity {
	if e, ok := s.elements[name]; ok {
		return e
	}
	if s.parent != nil {
		return s.parent.resolve(name)
	}
	return nil
}
