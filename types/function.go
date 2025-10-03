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
func NewFunction(name string, returnType Value, native, global bool, params ...Parameter) *Invokable {
	return &Invokable{
		kind:       FunctionKind,
		name:       name,
		normalized: normalize(name),
		params:     params,
		returnType: returnType,
		native:     native,
		global:     global,
	}
}

// NewEvent returns a new event Invokable type with zero or more parameters.
func NewEvent(name string, native bool, params ...Parameter) *Invokable {
	return &Invokable{
		kind:       EventKind,
		name:       name,
		normalized: normalize(name),
		returnType: VoidType,
		params:     params,
		native:     native,
	}
}

// Invokable is a function or event type.
type Invokable struct {
	name, normalized string
	params           []Parameter
	returnType       Value
	native, global   bool
	kind             InvokableKind
}

// Kind returns the kind of this invokable.
func (i *Invokable) Kind() InvokableKind {
	return i.kind
}

// ReturnType returns the return type of
// the function or [VoidType] if there isn't one.
func (i *Invokable) ReturnType() Value {
	return i.returnType
}

// Parameters returns the parameters in declaration order.
func (i *Invokable) Parameters() []Parameter {
	return i.params
}

// Native returns true if this invokable is native and false otherwise.
func (i *Invokable) Native() bool {
	return i.native
}

// Global returns true if this invokable is global and false otherwise.
func (i *Invokable) Global() bool {
	return i.global
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
		i.native == o.native && i.global == o.global &&
		len(i.params) == len(o.params) &&
		slices.EqualFunc(i.params, o.params, func(a, b Parameter) bool {
			return a.Type().IsIdentical(b.Type()) && a.normalized == b.normalized
		})
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
	if i.native {
		_, _ = sb.WriteString(" Native")
	}
	if i.global {
		_, _ = sb.WriteString(" Global")
	}
	return sb.String()
}

func (*Invokable) types() {}

var _ Type = (*Invokable)(nil)

// Parameter is a function or event parameter (though not a type itself).
type Parameter struct {
	name, normalized string
	typ              Value
	def              bool
}

// NewParameter returns a new parameter with the given name, type and
// default flag (true if the parameter has a default value and false otherwise).
func NewParameter(name string, typ Value, def bool) Parameter {
	return Parameter{
		name:       name,
		normalized: normalize(name),
		typ:        typ,
		def:        def,
	}
}

// Name returns the declared name for the parameter.
func (p Parameter) Name() string {
	return p.name
}

// Normalized returns the normalized name for the parameter.
func (p Parameter) Normalized() string {
	return p.normalized
}

// Type returns the type of the parameter.
func (p Parameter) Type() Type {
	return p.typ
}

// Default returns true if this parameter has a default value.
func (p Parameter) Default() bool {
	return p.def
}

func (p Parameter) String() string {
	if p.def {
		return p.typ.String() + " = <default>"
	}
	return p.typ.String()
}
