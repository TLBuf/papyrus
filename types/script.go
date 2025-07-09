package types

import "github.com/TLBuf/papyrus/ast"

// Script describes a Papyrus script.
type Script struct {
	universe   *Scope
	scope      *Scope
	typ        Type
	name       string
	normalized string
	node       *ast.Script
	imports    []*Script
}

// Imports returns the scripts this script imports.
func (s *Script) Imports() []*Script {
	return s.imports
}

// Enclosing returns the root scope that contains this script.
func (s *Script) Enclosing() *Scope {
	return s.universe
}

// Scope returns the scope this script defines.
func (s *Script) Scope() *Scope {
	return s.scope
}

// Type returns the type of this object.
func (s *Script) Type() Type {
	return s.typ
}

// Name returns the declared name for this object.
func (s *Script) Name() string {
	return s.name
}

// Normalized returns the normalized name for this object.
func (s *Script) Normalized() string {
	return s.normalized
}

// Node returns the [ast.Node] that defines this object.
func (s *Script) Node() ast.Node {
	return s.node
}
