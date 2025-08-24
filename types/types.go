// Package types defines the Papyrus type system.
//
// Types are broken down into two main categories: [Function] and [Value] types;
// as the names imply, the former describes function type information and the
// latter value type information (e.g. for variables and parameters).
//
// Value types again break down into two categories: [Scalar] and [Array] types;
// the former represent single values while the latter represent sequence of
// some scalar type.
//
// Scalar types again break down into two categories: [Object] and [Primitive]
// types; the former representing an script object and the latter the four
// primitive types: [Bool], [Int], [Float], [String].
package types

import (
	"fmt"
	"strings"
)

var (
	// Bool is the boolean type.
	Bool = &Primitive{BoolKind, "Bool", "bool"}
	// Int is the integer type.
	Int = &Primitive{IntKind, "Int", "int"}
	// Float is the floating-point type.
	Float = &Primitive{FloatKind, "Float", "float"}
	// String is the string type.
	String = &Primitive{StringKind, "String", "string"}
	// BoolArray is the boolean array type.
	BoolArray = &Array{Bool}
	// IntArray is the integer array type.
	IntArray = &Array{Int}
	// FloatArray is the floating-point array type.
	FloatArray = &Array{Float}
	// StringArray is the string array type.
	StringArray = &Array{String}
)

// Type is the common interface for all types.
type Type interface {
	fmt.Stringer

	// Name returns the standard name for the type.
	Name() string

	// Normalized returns the normalized name for the type.
	Normalized() string

	// IsIdentical returns true if this type is
	// identical to another type and false otherwise.
	//
	// This method is commutative; both of the following expressions will always
	// evaluate to the same value:
	//
	//  a.IsIdentical(b)
	//  b.IsIdentical(a)
	IsIdentical(Type) bool

	// IsAssignable returns true if a value of another type can be assigned to a
	// variable of this type without an explicit type conversion and false
	// otherwise.
	//
	// This method is NOT commutative; the following expressions may not evaluate
	// to the same value:
	//
	//  a.IsAssignable(b)
	//  b.IsAssignable(a)
	IsAssignable(Type) bool

	// IsComparable returns true if a value of this type can be compared (i.e.
	// with '>', '<=', etc.) with a value of another type and false otherwise.
	//
	// This method is commutative; both of the following expressions will always
	// evaluate to the same value:
	//
	//  a.IsComparable(b)
	//  b.IsComparable(a)
	IsComparable(Type) bool

	// IsEquatable returns true if a value of this type can be checked for
	// equality (i.e. with '==' or '!=') with a value of another type and false
	// otherwise.
	//
	// This method is commutative; both of the following expressions will always
	// evaluate to the same value:
	//
	//  a.IsEquatable(b)
	//  b.IsEquatable(a)
	IsEquatable(Type) bool

	// IsConvertible returns true if a value of this type can be converted to
	// a value of another type through an explicit cast and false otherwise.
	//
	// This method is NOT commutative; the following expressions may not evaluate
	// to the same value:
	//
	//  a.IsConvertible(b)
	//  b.IsConvertible(a)
	IsConvertible(Type) bool

	types()
}

// Value is a common interface for all value types (scalars and arrays).
type Value interface {
	Type

	value()
}

// Scalar is the common interface for all scalar (i.e. non-array) types.
type Scalar interface {
	Value

	scalar()
}

func normalize(name string) string {
	return strings.ToLower(name)
}
