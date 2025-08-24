package types

// PrimitiveKind defines the various primitive types.
type PrimitiveKind uint8

const (
	// BoolKind represents the primitive boolean type.
	BoolKind PrimitiveKind = iota
	// IntKind represents the primitive integer type.
	IntKind
	// FloatKind represents the primitive floating-point type.
	FloatKind
	// StringKind represents the primitive string type.
	StringKind
)

// Primitive represents a single Primitive, built-in type.
type Primitive struct {
	kind             PrimitiveKind
	name, normalized string
}

// Kind returns the kind of basic type.
func (p *Primitive) Kind() PrimitiveKind {
	return p.kind
}

// Name returns the standard name for the type.
func (p *Primitive) Name() string {
	return p.name
}

// Normalized returns the normalized name for the type.
func (p *Primitive) Normalized() string {
	return p.normalized
}

// IsIdentical returns true if this type is
// identical to another type and false otherwise.
//
// This method is commutative; both of the following expressions will always
// evaluate to the same value:
//
//	a.IsIdentical(b)
//	b.IsIdentical(a)
func (p *Primitive) IsIdentical(other Type) bool {
	o, ok := other.(*Primitive)
	return ok && p.kind == o.kind
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
func (p *Primitive) IsAssignable(other Type) bool {
	switch p.kind {
	case BoolKind, StringKind:
		return true
	case IntKind:
		return p.IsIdentical(other)
	case FloatKind:
		o, ok := other.(*Primitive)
		return ok && (o.kind == IntKind || o.kind == FloatKind)
	}
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
func (p *Primitive) IsComparable(other Type) bool {
	o, ok := other.(*Primitive)
	if ok {
		return (p.kind != BoolKind && p.IsAssignable(o)) || (o.kind != BoolKind && o.IsAssignable(p))
	}
	return p.kind != BoolKind && other.IsAssignable(p)
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
func (p *Primitive) IsEquatable(other Type) bool {
	return p.IsAssignable(other) || other.IsAssignable(p)
}

// IsConvertible returns true if a value of this type can be converted to
// a value of another type through an explicit cast and false otherwise.
//
// This method is NOT commutative; the following expressions may not evaluate
// to the same value:
//
//	a.IsConvertible(b)
//	b.IsConvertible(a)
func (p *Primitive) IsConvertible(other Type) bool {
	switch p.kind {
	case BoolKind, StringKind:
		return true
	case IntKind, FloatKind:
		_, ok := other.(*Primitive)
		return ok
	}
	return false
}

func (b *Primitive) String() string {
	return b.Name()
}

func (*Primitive) types() {}

func (*Primitive) scalar() {}

func (*Primitive) value() {}

var _ Scalar = (*Primitive)(nil)
