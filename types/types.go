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

func (Bool) types() {}

func (Bool) scalar() {}

var _ Scalar = Bool{}

// Int represents the signed 32-bit integer type.
type Int struct{}

func (Int) types() {}

func (Int) scalar() {}

var _ Scalar = Int{}

// Float represents the signed 32-bit floating-point type.
type Float struct{}

func (Float) types() {}

func (Float) scalar() {}

var _ Scalar = Float{}

// String represents the string type.
type String struct{}

func (String) types() {}

func (String) scalar() {}

var _ Scalar = String{}

// Object represents the object type.
type Object struct {
	// Name is the normalized object type name.
	Name string
}

func (Object) types() {}

func (Object) scalar() {}

var _ Scalar = Object{}

// Array represents the array type.
type Array struct {
	ElementType Scalar
}

func (Array) types() {}

var _ Type = Array{}
