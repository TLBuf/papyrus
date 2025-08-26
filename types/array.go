package types

// Array represents an array type with an optionally known length.
type Array struct {
	element Scalar
}

// Element returns the scalar type for elements of the array.
func (a *Array) Element() Scalar {
	return a.element
}

// Name returns the declared name for the type.
func (a *Array) Name() string {
	return a.element.Name() + "[]"
}

// Normalized returns the normalized name for the type.
func (a *Array) Normalized() string {
	return a.element.Normalized() + "[]"
}

// IsIdentical returns true if this type is
// identical to another type and false otherwise.
//
// This method is commutative; both of the following expressions will always
// evaluate to the same value:
//
//	a.IsIdentical(b)
//	b.IsIdentical(a)
func (a *Array) IsIdentical(other Type) bool {
	o, ok := other.(*Array)
	return ok && a.element.IsIdentical(o.element)
}

// IsAssignable returns true if a value of another type can be assigned to a
// variable of this type without an explicit type conversion and false
// otherwise.
//
// This method is NOT commutative; the following expressions may not evaluate
// to the same value:
//
//	a.IsAssignable(b)
//	b.IsAssignable(a)
func (a *Array) IsAssignable(other Type) bool {
	o, ok := other.(*Array)
	return ok && a.element.IsIdentical(o.element)
}

// IsComparable returns true if a value of this type can be compared (i.e.
// with '>', '<=', etc.) with a value of another type and false otherwise.
//
// This method is commutative; both of the following expressions will always
// evaluate to the same value:
//
//	a.IsComparable(b)
//	b.IsComparable(a)
func (*Array) IsComparable(Type) bool {
	return false
}

// IsEquatable returns true if a value of this type can be checked for
// equality (i.e. with '==' or '!=') with a value of another type and false
// otherwise.
//
// This method is commutative; both of the following expressions will always
// evaluate to the same value:
//
//	a.IsEquatable(b)
//	b.IsEquatable(a)
func (*Array) IsEquatable(Type) bool {
	return false
}

// IsConvertible returns true if a value of this type can be converted to
// a value of another type through an explicit cast and false otherwise.
//
// This method is NOT commutative; the following expressions may not evaluate
// to the same value:
//
//	a.IsConvertible(b)
//	b.IsConvertible(a)
func (*Array) IsConvertible(Type) bool {
	return false
}

func (a *Array) String() string {
	return a.Name()
}

func (*Array) types() {}

func (*Array) value() {}

var _ Type = (*Array)(nil)
