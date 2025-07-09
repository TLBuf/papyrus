package types

// BasicKind defines the various basic types.
type BasicKind uint8

const (
	// BoolKind represents the basic boolean type.
	BoolKind BasicKind = iota
	// IntKind represents the basic integer type.
	IntKind
	// FloatKind represents the basic floating-point type.
	FloatKind
	// StringKind represents the basic string type.
	StringKind
)

// Basic represents a single basic, built-in type.
type Basic struct {
	kind BasicKind
	name string
}

// Kind returns the kind of basic type.
func (b *Basic) Kind() BasicKind {
	return b.kind
}

// Name returns the standard name for the type.
func (b *Basic) Name() string {
	return b.name
}

func (b *Basic) String() string {
	return b.name
}

func (*Basic) types() {}

func (*Basic) scalar() {}

var _ Scalar = (*Basic)(nil)

var (
	// Bool is the boolean type.
	Bool = &Basic{BoolKind, "Bool"}
	// Int is the integer type.
	Int = &Basic{IntKind, "Int"}
	// Float is the floating-point type.
	Float = &Basic{FloatKind, "Float"}
	// String is the string type.
	String = &Basic{StringKind, "String"}
)
