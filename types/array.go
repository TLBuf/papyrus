package types

// Array represents an array type with an optionally known length.
type Array struct {
	element Scalar
	length  uint8
}

// NewArray returns a new array with a given element type and optional length.
//
// Length must be in range [1, 128] if known or zero otherwise.
func NewArray(element Scalar, length uint8) *Array {
	if length > 128 {
		panic("array length out of range")
	}
	return &Array{
		element: element,
		length:  length,
	}
}

// Element returns the scalar type for elements of the array.
func (a *Array) Element() Scalar {
	return a.element
}

// Length returns the length of the array in the
// range [1, 128] if known or zero if not known.
func (a *Array) Length() uint8 {
	return a.length
}

func (a *Array) String() string {
	return a.element.String() + "[]"
}

func (*Array) types() {}

var _ Type = (*Array)(nil)
