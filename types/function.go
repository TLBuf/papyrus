package types

import (
	"slices"
	"strings"
)

// InvokableKind defines the different kinds of invokable types.
type InvokableKind uint8

const (
	// FunctionKind represents the function invokable type.
	FunctionKind InvokableKind = iota
	// EventKind represents the event invokable type.
	EventKind
)

// NewFunction returns a new function Invokable type with
// an optional return type and zero or more parameters.
func NewFunction(name string, returnType Value, params ...Value) *Invokable {
	return &Invokable{
		kind:       FunctionKind,
		name:       name,
		normalized: normalize(name),
		params:     params,
		returnType: returnType,
	}
}

// NewEvent returns a new event Invokable type with zero or more parameters.
func NewEvent(name string, params ...Value) *Invokable {
	return &Invokable{
		kind:       EventKind,
		name:       name,
		normalized: normalize(name),
		params:     params,
	}
}

// Invokable is a function or event type.
type Invokable struct {
	name, normalized string
	params           []Value
	returnType       Value
	kind             InvokableKind
}

// Kind returns the kind of this invokable.
func (i *Invokable) Kind() InvokableKind {
	return i.kind
}

// ReturnType returns the return type of
// the function or nil if there isn't one.
func (i *Invokable) ReturnType() Value {
	return i.returnType
}

// Parameters returns the parameters in declaration order.
func (i *Invokable) Parameters() []Value {
	return i.params
}

// Name returns the declared name for the type.
func (i *Invokable) Name() string {
	return i.name
}

// Normalized returns the normalized name for the type.
func (i *Invokable) Normalized() string {
	return i.normalized
}

// IsIdentical returns true if this type is
// identical to another type and false otherwise.
//
// This method is commutative; both of the following expressions will always
// evaluate to the same value:
//
//	a.IsIdentical(b)
//	b.IsIdentical(a)
func (i *Invokable) IsIdentical(other Type) bool {
	o, ok := other.(*Invokable)
	return ok && i.normalized == o.normalized && i.returnType.IsIdentical(o.returnType) &&
		len(i.params) == len(o.params) &&
		slices.EqualFunc(i.params, o.params, func(a, b Value) bool { return a.IsIdentical(b) })
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
func (*Invokable) IsAssignable(Type) bool {
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
func (*Invokable) IsComparable(Type) bool {
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
func (*Invokable) IsEquatable(Type) bool {
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
func (*Invokable) IsConvertible(Type) bool {
	return false
}

func (i *Invokable) String() string {
	var sb strings.Builder
	if i.returnType != nil {
		_, _ = sb.WriteString(i.returnType.String())
		_, _ = sb.Write([]byte{' '})
	}
	_, _ = sb.WriteString(i.name)
	_, _ = sb.Write([]byte{'('})
	for i, p := range i.params {
		if i > 0 {
			_, _ = sb.Write([]byte{','})
		}
		_, _ = sb.WriteString(p.String())
	}
	_, _ = sb.Write([]byte{')'})
	return sb.String()
}

func (*Invokable) types() {}

var _ Type = (*Invokable)(nil)
