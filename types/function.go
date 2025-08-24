package types

import (
	"slices"
	"strings"
)

// NewFunction returns a new Function type with an
// optional return type and zero or more parameters.
func NewFunction(name string, returnType Value, params ...Value) *Function {
	return &Function{
		name:       name,
		normalized: normalize(name),
		params:     params,
		returnType: returnType,
	}
}

// Function is a function type.
type Function struct {
	name, normalized string
	params           []Value
	returnType       Value
}

// ReturnType returns the return type of
// the function or nil if there isn't one.
func (f *Function) ReturnType() Value {
	return f.returnType
}

// Parameters returns the parameters in declaration order.
func (f *Function) Parameters() []Value {
	return f.params
}

// Name returns the declared name for the type.
func (f *Function) Name() string {
	return f.name
}

// Normalized returns the normalized name for the type.
func (f *Function) Normalized() string {
	return f.normalized
}

// IsIdentical returns true if this type is
// identical to another type and false otherwise.
//
// This method is commutative; both of the following expressions will always
// evaluate to the same value:
//
//	a.IsIdentical(b)
//	b.IsIdentical(a)
func (f *Function) IsIdentical(other Type) bool {
	o, ok := other.(*Function)
	return ok && f.normalized == o.normalized && f.returnType.IsIdentical(o.returnType) && len(f.params) == len(o.params) && slices.EqualFunc(f.params, o.params, func(a, b Value) bool { return a.IsIdentical(b) })
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
func (f *Function) IsAssignable(other Type) bool {
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
func (f *Function) IsComparable(other Type) bool {
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
func (f *Function) IsEquatable(other Type) bool {
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
func (f *Function) IsConvertible(other Type) bool {
	return false
}

func (f *Function) String() string {
	var sb strings.Builder
	if f.returnType != nil {
		_, _ = sb.WriteString(f.returnType.String())
		_, _ = sb.Write([]byte{' '})
	}
	_, _ = sb.WriteString(f.name)
	_, _ = sb.Write([]byte{'('})
	for i, p := range f.params {
		if i > 0 {
			_, _ = sb.Write([]byte{','})
		}
		_, _ = sb.WriteString(p.String())
	}
	_, _ = sb.Write([]byte{')'})
	return sb.String()
}

func (*Function) types() {}

var _ Type = (*Function)(nil)
