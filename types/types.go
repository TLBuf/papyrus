// Package types defines the Papyrus type system.
//
// Types are broken down into two main categories: [Invokable] and [Value]
// types; as the names imply, the former describes function type information and
// the latter value type information (e.g. for variables and parameters).
//
// Value types again break down into two categories: [Scalar] and [Array] types;
// the former represent single values while the latter represent sequence of
// some scalar type.
//
// Scalar types again break down into two categories: [Object] and [Primitive]
// types; the former representing an script object and the latter the four
// primitive types: [BoolType], [IntType], [FloatType], [StringType].
//
// There are also two special types: [NoneType] and [VoidType]. [None] is a
// scalar value type that is compatible with any object type and is used
// exclusively with the 'None' literal. [Void] is a value type that is used for
// functions that do not return a value.
package types

import (
	"fmt"
	"strings"
)

var (
	// BoolType is the boolean type.
	BoolType = &Primitive{BoolKind, "Bool", "bool"}
	// IntType is the integer type.
	IntType = &Primitive{IntKind, "Int", "int"}
	// FloatType is the floating-point type.
	FloatType = &Primitive{FloatKind, "Float", "float"}
	// StringType is the string type.
	StringType = &Primitive{StringKind, "String", "string"}
	// BoolArrayType is the boolean array type.
	BoolArrayType = &Array{BoolType}
	// IntArrayType is the integer array type.
	IntArrayType = &Array{IntType}
	// FloatArrayType is the floating-point array type.
	FloatArrayType = &Array{FloatType}
	// StringArrayType is the string array type.
	StringArrayType = &Array{StringType}
	// NoneType is the type that is compatible with any object or array
	// type. This is used exclusively with the 'None' literal.
	NoneType = None{}
	// VoidType is the type used for functions that do not return a value.
	VoidType = Void{}
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
