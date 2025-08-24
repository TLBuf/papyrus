// Package literal provides representations for Papyrus literals.
package literal

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/TLBuf/papyrus/ast"
)

// Kind defines the different kinds of Papyrus literals.
type Kind uint8

const (
	// InvalidKind represents an invalid literal.
	InvalidKind Kind = iota
	// BoolKind represents a boolean literal.
	BoolKind
	// IntKind represents an interger literal.
	IntKind
	// FloatKind represents a floating-point literal.
	FloatKind
	// StringKind represents a string literal.
	StringKind
)

// Value represents the value of a Papyrus literal.
type Value interface {
	// Kind returns the kind of literal.
	Kind() Kind
	// String returns a string representation of the value, this
	// may not necessarily match how the value appeared in source.
	String() string

	value()
}

// NewValue returns a new value built from an [ast.Literal] or returns
// an [InvalidKind] if the node text cannot be parsed as a valid value.
func NewValue(node ast.Literal) (val Value) {
	text := string(node.Text())
	var err error
	switch node := node.(type) {
	case *ast.BoolLiteral:
		v := false
		v, err = parseBool(text)
		val = boolValue(v)
	case *ast.IntLiteral:
		v := int32(0)
		v, err = parseInt(text)
		val = intValue(v)
	case *ast.FloatLiteral:
		v := float32(0)
		v, err = parseFloat(text)
		val = floatValue(v)
	case *ast.StringLiteral:
		v := ""
		v, err = parseString(text)
		val = stringValue(v)
	default:
		err = fmt.Errorf("cannot create a Value from literal of type %T", node)
	}
	if err != nil {
		val = invalidValue{
			text: text,
			err:  err,
		}
	}
	return val
}

func parseBool(text string) (bool, error) {
	v, err := strconv.ParseBool(strings.ToLower(text))
	if err != nil {
		return false, fmt.Errorf("parse %q as a boolean value: %w", text, err)
	}
	return v, nil
}

func parseInt(text string) (int32, error) {
	v, err := strconv.ParseInt(strings.ToLower(text), 0, 32)
	if err != nil {
		return 0, fmt.Errorf("parse %q as a 32-bit integer value: %w", text, err)
	}
	return int32(v), nil
}

func parseFloat(text string) (float32, error) {
	v, err := strconv.ParseFloat(text, 32)
	if err != nil {
		return 0.0, fmt.Errorf("parse %q as a 32-bit floating-point value: %w", text, err)
	}
	return float32(v), nil
}

func parseString(text string) (string, error) {
	v, err := strconv.Unquote(text)
	if err != nil {
		return "", fmt.Errorf("parse %q as a string value: %w", text, err)
	}
	return v, nil
}

// Bool returns the value as a bool.
// Panics if the value isn't of [BoolKind].
func Bool(v Value) bool {
	if v, ok := v.(boolValue); ok {
		return bool(v)
	}
	panic(fmt.Sprintf("%v is not a Bool", v))
}

// Int returns the value as an int32.
// Panics if the value isn't of [IntKind].
func Int(v Value) int32 {
	if v, ok := v.(intValue); ok {
		return int32(v)
	}
	panic(fmt.Sprintf("%v is not an Int", v))
}

// Float returns the value as a float32.
// Panics if the value isn't of [FloatKind].
func Float(v Value) float32 {
	if v, ok := v.(floatValue); ok {
		return float32(v)
	}
	panic(fmt.Sprintf("%v is not a Float", v))
}

// String returns the value as a string.
// Panics if the value isn't of [StringKind].
func String(v Value) string {
	if v, ok := v.(stringValue); ok {
		return string(v)
	}
	panic(fmt.Sprintf("%v is not a String", v))
}

// ToBool returns the value as a Bool.
func ToBool(v Value) Value {
	return boolValue(asBool(v))
}

// ToInt returns the value as an Int.
func ToInt(v Value) Value {
	return intValue(asInt(v))
}

// ToFloat returns the value as a Float.
func ToFloat(v Value) Value {
	return floatValue(asFloat(v))
}

// ToString returns the value as a String.
func ToString(v Value) Value {
	return stringValue(asString(v))
}

func asBool(v Value) bool {
	switch v := v.(type) {
	case boolValue:
		return bool(v)
	case intValue:
		return v != 0
	case floatValue:
		return v != 0.0
	case stringValue:
		return v != ""
	}
	panic(fmt.Sprintf("%v cannot be converted to a Bool", v))
}

func asInt(v Value) int32 {
	switch v := v.(type) {
	case boolValue:
		if v {
			return 1
		}
		return 0
	case intValue:
		return int32(v)
	case floatValue:
		return int32(v)
	case stringValue:
		n, err := parseInt(string(v))
		if err != nil {
			return 0
		}
		return n
	}
	panic(fmt.Sprintf("%v cannot be converted to an Int", v))
}

func asFloat(v Value) float32 {
	switch v := v.(type) {
	case boolValue:
		if v {
			return 1.0
		}
		return 0.0
	case intValue:
		return float32(v)
	case floatValue:
		return float32(v)
	case stringValue:
		n, err := parseFloat(string(v))
		if err != nil {
			return 0.0
		}
		return n
	}
	panic(fmt.Sprintf("%v cannot be converted to a Float", v))
}

func asString(v Value) string {
	switch v := v.(type) {
	case boolValue:
		if v {
			return "True"
		}
		return "False"
	case intValue, floatValue:
		return fmt.Sprintf("%v", v)
	case stringValue:
		return string(v)
	}
	panic(fmt.Sprintf("%v cannot be converted to a String", v))
}

type invalidValue struct {
	text string
	err  error
}

func (invalidValue) Kind() Kind {
	return InvalidKind
}

func (v invalidValue) Error() error {
	return v.err
}

func (v invalidValue) String() string {
	return v.text
}

func (invalidValue) value() {}

type boolValue bool

func (boolValue) Kind() Kind {
	return BoolKind
}

func (v boolValue) String() string {
	return strconv.FormatBool(bool(v))
}

func (boolValue) value() {}

type intValue int32

func (intValue) Kind() Kind {
	return IntKind
}

func (v intValue) String() string {
	return strconv.FormatInt(int64(v), 10)
}

func (intValue) value() {}

type floatValue float32

func (floatValue) Kind() Kind {
	return FloatKind
}

func (v floatValue) String() string {
	return strconv.FormatFloat(float64(v), 'g', 6, 32)
}

func (floatValue) value() {}

type stringValue string

func (stringValue) Kind() Kind {
	return StringKind
}

func (v stringValue) String() string {
	return strconv.Quote(string(v))
}

func (stringValue) value() {}
