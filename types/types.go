// Package types defines the Papyrus type system.
package types

import (
	"fmt"
	"slices"
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

// Identical returns true if this type is
// identical to another type and false otherwise.
func Identical(a, b Type) bool {
	switch a := a.(type) {
	case nil:
		switch b.(type) {
		case nil:
			// Consider two nil types to be identical,
			// e.g. checking the parents of two objects.
			return true
		default:
			return false
		}
	case *Array:
		b, ok := b.(*Array)
		return ok && Identical(a.element, b.element) && a.length == b.length
	case *Basic:
		b, ok := b.(*Basic)
		return ok && a.kind == b.kind
	case *Object:
		b, ok := b.(*Object)
		return ok && a.normalized == b.normalized && Identical(a.parent, b.parent)
	case *Signature:
		b, ok := b.(*Signature)
		return ok && Identical(a.returnType, b.returnType) && len(a.params) == len(b.params) && slices.EqualFunc(a.params, b.params, Identical)
	}
	return false
}

// Assignable returns true if a value of type rhs can be assigned to a variable
// of type lhs without an explicit type conversion and false otherwise.
func Assignable(lhs, rhs Type) bool {
	switch lhs := lhs.(type) {
	case *Array:
		rhs, ok := rhs.(*Array)
		return ok && Identical(lhs.element, rhs.element)
	case *Basic:
		switch lhs.kind {
		case BoolKind, StringKind:
			return true
		case IntKind:
			return Identical(lhs, rhs)
		case FloatKind:
			rhs, ok := rhs.(*Basic)
			return ok && (rhs.kind == IntKind || rhs.kind == FloatKind)
		}
	case *Object:
		rhs, ok := rhs.(*Object)
		return ok && (Identical(lhs, rhs) || Assignable(lhs, rhs.parent))
	case *Signature:
		return false
	}
	return false
}

// Equatable returns true if a value of this type can be checked for equality
// (i.e. with '==' or '!=') with a value of another type and false otherwise.
func Equatable(lhs, rhs Type) bool {
	return autoCastsTo(lhs, rhs) || autoCastsTo(rhs, lhs)
}

// Comparable returns true if a value of this type can be compared (i.e. with
// '>', '<=', etc.) with a value of another type and false otherwise.
func Comparable(lhs, rhs Type) bool {
	// Only Int, Float, and String are comparable.
	l, lok := lhs.(*Basic)
	r, rok := rhs.(*Basic)
	if !lok && !rok {
		return false
	}
	lc := lok && l.kind != BoolKind && autoCastsTo(rhs, lhs)
	rc := rok && r.kind != BoolKind && autoCastsTo(lhs, rhs)
	return lc || rc
}

// ConvertibleTo returns true if a value of this type can be converted to
// a value of another type through an explicit cast or false otherwise.
func ConvertibleTo(src, dst Type) bool {
	switch dst := dst.(type) {
	case *Array, *Signature:
		return false
	case *Basic:
		switch dst.kind {
		case BoolKind, StringKind:
			return true
		case IntKind:
			_, ok := src.(*Basic)
			return ok
		case FloatKind:
			_, ok := src.(*Basic)
			return ok
		}
	case *Object:
		src, ok := src.(*Object)
		return ok && (Identical(src, dst) || ConvertibleTo(src.parent, dst) || ConvertibleTo(src, dst.parent))
	}
	return false
}

func autoCastsTo(src, dst Type) bool {
	switch dst := dst.(type) {
	case *Array, *Signature:
		return false
	case *Basic:
		switch dst.kind {
		case BoolKind, StringKind:
			return true
		case IntKind:
			src, ok := src.(*Basic)
			return ok && src.kind == IntKind
		case FloatKind:
			src, ok := src.(*Basic)
			return ok && (src.kind == IntKind || src.kind == FloatKind)
		}
	case *Object:
		src, ok := src.(*Object)
		return ok && (Identical(src, dst) || autoCastsTo(src.parent, dst))
	}
	return false
}

func normalize(name string) string {
	return strings.ToLower(name)
}
