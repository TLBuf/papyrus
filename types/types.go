// Package types defines the Papyrus type system.
package types

import (
	"fmt"
	"strings"
)

// Type is the common interface for all types.
type Type interface {
	fmt.Stringer
	types()
}

// Scalar is the common interface for all scalar (i.e. non-array) types.
type Scalar interface {
	Type
	scalar()
}

// BasicKind defines the various basic types.
type BasicKind uint8

const (
	// BoolKind represents the basic boolean type.
	BoolKind BasicKind = iota
	// IntKind represents the basic integer type.
	IntKind
	// FloatKind represents the basic floating-point type.
	FloatKind
	// StringKind represents the basic string type.
	StringKind
)

// Basic represents a single basic, built-in type.
type Basic struct {
	kind BasicKind
	name string
}

// Kind returns the kind of basic type.
func (b *Basic) Kind() BasicKind {
	return b.kind
}

// Name returns the standard name for the type.
func (b *Basic) Name() string {
	return b.name
}

func (b *Basic) String() string {
	return b.name
}

func (*Basic) types() {}

func (*Basic) scalar() {}

var _ Scalar = (*Basic)(nil)

var (
	// Bool is the boolean type.
	Bool = &Basic{BoolKind, "Bool"}
	// Int is the integer type.
	Int = &Basic{IntKind, "Int"}
	// Float is the floating-point type.
	Float = &Basic{FloatKind, "Float"}
	// String is the string type.
	String = &Basic{StringKind, "String"}
)

// NewNamed returns a new Named type with a given declared name.
func NewNamed(name string) *Named {
	return &Named{
		name:       name,
		normalized: strings.ToLower(name),
	}
}

// Named represents a named typed (i.e. a script).
type Named struct {
	// Name is the normalized object type name.
	name       string
	normalized string
}

// Name returns the declared name for the type.
func (n *Named) Name() string {
	return n.name
}

// Normalized returns the normalized name for the type.
func (n *Named) Normalized() string {
	return n.normalized
}

func (n *Named) String() string {
	return n.name
}

func (*Named) types() {}

func (*Named) scalar() {}

var _ Scalar = (*Named)(nil)

// Array represents an array type with an optionally known length.
type Array struct {
	elem Scalar
	len  uint8
}

// Element returns the scalar type for elements of the array.
func (a *Array) Element() Scalar {
	return a.elem
}

// Length returns the length of the array in the
// range [1, 128] if known or zero if not known.
func (a *Array) Length() uint8 {
	return a.len
}

func (a *Array) String() string {
	return a.elem.String() + "[]"
}

func (*Array) types() {}

var _ Type = (*Array)(nil)
