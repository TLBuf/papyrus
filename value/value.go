// Package value provides representations for Papyrus values.
package value

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/types"
)

var (
	// ErrParse indicates the text of a value could
	// not be parsed as an instance of that type.
	ErrParse = errors.New("parse")
	// ErrWrongKind indicates that an attempt to use the underlying
	// value failed because the value is not of that type.
	ErrWrongKind = errors.New("wrong kind")
)

// Kind defines the different kinds of Papyrus values.
type Kind uint8

const (
	// None represents the none value.
	None Kind = iota
	// Bool represents a boolean value.
	Bool
	// Int represents an interger value.
	Int
	// Float represents a floating-point value.
	Float
	// String represents a string value.
	String
)

// Value represents the value of a Papyrus expression.
type Value struct {
	kind Kind
	data any
}

// Kind returns the kind of this value.
func (v Value) Kind() Kind {
	return v.kind
}

// Type returns the type of this value.
func (v Value) Type() types.Type {
	switch v.kind {
	case Bool:
		return types.Bool
	case Int:
		return types.Int
	case Float:
		return types.Float
	case String:
		return types.String
	default:
		return types.Any
	}
}

// BoolValue returns this value as a bool or panics with an
// error that wraps [ErrWrongKind] if this value is not a boolean.
func (v Value) BoolValue() bool {
	if v.kind != Bool {
		panic(fmt.Errorf("%s is not a Bool: %w", v, ErrWrongKind))
	}
	//revive:disable-next-line:unchecked-type-assertion
	return v.data.(bool)
}

// IntValue returns this value as a int32 or panics with an
// error that wraps [ErrWrongKind] if this value is not an int.
func (v Value) IntValue() int32 {
	if v.kind != Int {
		panic(fmt.Errorf("%s is not a Int: %w", v, ErrWrongKind))
	}
	//revive:disable-next-line:unchecked-type-assertion
	return v.data.(int32)
}

// FloatValue returns this value as a float32 or panics with an
// error that wraps [ErrWrongKind] if this value is not a float.
func (v Value) FloatValue() float32 {
	if v.kind != Float {
		panic(fmt.Errorf("%s is not a Float: %w", v, ErrWrongKind))
	}
	//revive:disable-next-line:unchecked-type-assertion
	return v.data.(float32)
}

// StringValue returns this value as a string or panics with an
// error that wraps [ErrWrongKind] if this value is not a string.
func (v Value) StringValue() string {
	if v.kind != String {
		panic(fmt.Errorf("%s is not a String: %w", v, ErrWrongKind))
	}
	//revive:disable-next-line:unchecked-type-assertion
	return v.data.(string)
}

// String returns a string representation of the value, this
// may not necessarily match how the value appeared in source.
func (v Value) String() string {
	switch v.kind {
	case Bool:
		return strconv.FormatBool(v.data.(bool))
	case Int:
		return strconv.FormatInt(int64(v.data.(int32)), 10)
	case Float:
		return strconv.FormatFloat(float64(v.data.(float32)), 'g', 6, 32)
	case String:
		return fmt.Sprintf("%q", v.data.(string))
	default:
		return "None"
	}
}

// New returns a new value built from an [ast.Literal] or returns an error
// that wraps [ErrParse] if the literal's text couldn't be parsed as a value of
// the corresponding type.
func New(node ast.Literal) (Value, error) {
	text := string(node.Text())
	switch node.(type) {
	case *ast.BoolLiteral:
		v, err := parseBool(text)
		if err != nil {
			return Value{}, err
		}
		return Value{data: v, kind: Bool}, nil
	case *ast.IntLiteral:
		v, err := parseInt(text)
		if err != nil {
			return Value{}, err
		}
		return Value{data: v, kind: Int}, nil
	case *ast.FloatLiteral:
		v, err := parseFloat(text)
		if err != nil {
			return Value{}, err
		}
		return Value{data: v, kind: Float}, nil
	case *ast.StringLiteral:
		v, err := parseString(text)
		if err != nil {
			return Value{}, err
		}
		return Value{data: v, kind: String}, nil
	case *ast.NoneLiteral:
		return Value{kind: None}, nil
	}
	panic(fmt.Errorf("unknown literal type %T", node))
}

func parseBool(text string) (bool, error) {
	v, err := strconv.ParseBool(strings.ToLower(text))
	if err != nil {
		return false, fmt.Errorf("%w: %q as bool: %w", ErrParse, text, err)
	}
	return v, nil
}

func parseInt(text string) (int32, error) {
	v, err := strconv.ParseInt(strings.ToLower(text), 0, 32)
	if err != nil {
		return 0, fmt.Errorf("%w: %q as int32: %w", ErrParse, text, err)
	}
	return int32(v), nil
}

func parseFloat(text string) (float32, error) {
	v, err := strconv.ParseFloat(text, 32)
	if err != nil {
		return 0.0, fmt.Errorf("%w: %q as float32: %w", ErrParse, text, err)
	}
	return float32(v), nil
}

func parseString(text string) (string, error) {
	v, err := strconv.Unquote(text)
	if err != nil {
		return "", fmt.Errorf("%w: %q as string: %w", ErrParse, text, err)
	}
	return v, nil
}
