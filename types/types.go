// Package types defines the Papyrus type system.
package types

// Type is the common interface for all types.
type Type interface {
	types()
}

// Scalar is the common interface for all scalar (i.e. non-array) types.
type Scalar interface {
	Type
	scalar()
}

// Bool represents the boolean (i.e. true or false) type.
type Bool struct{}

func (b Bool) types() {}

func (b Bool) scalar() {}

var _ Scalar = Bool{}

// Int represents the signed 32-bit integer type.
type Int struct{}

func (i Int) types() {}

func (i Int) scalar() {}

var _ Scalar = Int{}

// Float represents the signed 32-bit floating-point type.
type Float struct{}

func (f Float) types() {}

func (f Float) scalar() {}

var _ Scalar = Float{}

// String represents the string type.
type String struct{}

func (s String) types() {}

func (s String) scalar() {}

var _ Scalar = String{}

// Object represents the object type.
type Object struct {
	// Name is the normalized object type name.
	Name string
}

func (o Object) types() {}

func (o Object) scalar() {}

var _ Scalar = Object{}

// Array represents the array type.
type Array struct {
	ElementType Scalar
}

func (a Array) types() {}

var _ Type = Array{}
