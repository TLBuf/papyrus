package types

// None is the special type used for the 'None' literal.
type None struct{}

// Name returns the declared name for the type.
func (None) Name() string {
	return "<Any>"
}

// Normalized returns the normalized name for the type.
func (None) Normalized() string {
	return "<Any>"
}

// IsIdentical returns true if this type is
// identical to another type and false otherwise.
//
// This method is commutative; both of the following expressions will always
// evaluate to the same value:
//
//	a.IsIdentical(b)
//	b.IsIdentical(a)
func (None) IsIdentical(Type) bool {
	return false
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
func (None) IsAssignable(Type) bool {
	return false
}

// IsComparable returns true if a value of this type can be compared (i.e.
// with '>', '<=', etc.) with a value of another type and false otherwise.
//
// This method is commutative; both of the following expressions will always
// evaluate to the same value:
//
//	a.IsComparable(b)
//	b.IsComparable(a)
func (None) IsComparable(Type) bool {
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
func (None) IsEquatable(other Type) bool {
	switch other.(type) {
	case *Object, None:
		return true
	}
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
func (None) IsConvertible(Type) bool {
	return false
}

func (None) String() string {
	return "<Any>"
}

func (None) types() {}

func (None) scalar() {}

func (None) value() {}

var _ Scalar = None{}
