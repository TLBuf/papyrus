package types

// Object represents a named typed (i.e. a script).
type Object struct {
	name, normalized string
	parent           *Object
}

// Parent is the object this object extends or nil if it extends nothing.
func (o *Object) Parent() *Object {
	return o.parent
}

// Name returns the declared name for the type.
func (o *Object) Name() string {
	return o.name
}

// Normalized returns the normalized name for the type.
func (o *Object) Normalized() string {
	return o.normalized
}

// IsIdentical returns true if this type is
// identical to another type and false otherwise.
//
// This method is commutative; both of the following expressions will always
// evaluate to the same value:
//
//	a.IsIdentical(b)
//	b.IsIdentical(a)
func (o *Object) IsIdentical(other Type) bool {
	t, ok := other.(*Object)
	return ok && o.normalized == t.normalized && o.parent.IsIdentical(t.parent)
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
func (o *Object) IsAssignable(other Type) bool {
	t, ok := other.(*Object)
	return ok && o.IsIdentical(t) || o.IsAssignable(t.parent)
}

// IsComparable returns true if a value of this type can be compared (i.e.
// with '>', '<=', etc.) with a value of another type and false otherwise.
//
// This method is commutative; both of the following expressions will always
// evaluate to the same value:
//
//	a.IsComparable(b)
//	b.IsComparable(a)
func (o *Object) IsComparable(other Type) bool {
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
func (o *Object) IsEquatable(other Type) bool {
	t, ok := other.(*Object)
	return ok && (o.IsIdentical(t) || o.parent.IsEquatable(t) || t.parent.IsEquatable(o))
}

// IsConvertible returns true if a value of this type can be converted to
// a value of another type through an explicit cast and false otherwise.
//
// This method is NOT commutative; the following expressions may not evaluate
// to the same value:
//
//	a.IsConvertible(b)
//	b.IsConvertible(a)
func (o *Object) IsConvertible(other Type) bool {
	t, ok := other.(*Object)
	return ok && (o.IsIdentical(t) || o.IsConvertible(t.parent) || o.parent.IsConvertible(t))
}

func (o *Object) String() string {
	return o.name
}

func (*Object) types() {}

func (*Object) scalar() {}

func (*Object) value() {}

var _ Scalar = (*Object)(nil)
