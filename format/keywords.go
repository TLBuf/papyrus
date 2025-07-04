package format

import (
	"fmt"
	"strings"

	"github.com/TLBuf/papyrus/token"
)

// defaultKeywords is populated with all default keyword values.
var defaultKeywords, _ = keywords(Keywords{})

// Keywords defines the text that is used by
// the formatter when printing a keyword.
type Keywords struct {
	// As is the text used when formatting a [token.As] keyword.
	As string
	// Auto is the text used when formatting a [token.Auto] keyword.
	Auto string
	// AutoReadOnly is the text used when formatting a [token.AutoReadOnly] keyword.
	AutoReadOnly string
	// Bool is the text used when formatting a [token.Bool] keyword.
	Bool string
	// Conditional is the text used when formatting a [token.Conditional] keyword.
	Conditional string
	// Else is the text used when formatting a [token.Else] keyword.
	Else string
	// ElseIf is the text used when formatting a [token.ElseIf] keyword.
	ElseIf string
	// EndEvent is the text used when formatting a [token.EndEvent] keyword.
	EndEvent string
	// EndFunction is the text used when formatting a [token.EndFunction] keyword.
	EndFunction string
	// EndIf is the text used when formatting a [token.EndIf] keyword.
	EndIf string
	// EndProperty is the text used when formatting a [token.EndProperty] keyword.
	EndProperty string
	// EndState is the text used when formatting a [token.EndState] keyword.
	EndState string
	// EndWhile is the text used when formatting a [token.EndWhile] keyword.
	EndWhile string
	// Event is the text used when formatting a [token.Event] keyword.
	Event string
	// Extends is the text used when formatting a [token.Extends] keyword.
	Extends string
	// False is the text used when formatting a [token.False] keyword.
	False string
	// Float is the text used when formatting a [token.Float] keyword.
	Float string
	// Function is the text used when formatting a [token.Function] keyword.
	Function string
	// Global is the text used when formatting a [token.Global] keyword.
	Global string
	// Hidden is the text used when formatting a [token.Hidden] keyword.
	Hidden string
	// If is the text used when formatting a [token.If] keyword.
	If string
	// Import is the text used when formatting a [token.Import] keyword.
	Import string
	// Int is the text used when formatting a [token.Int] keyword.
	Int string
	// Length is the text used when formatting a [token.Length] keyword.
	Length string
	// Native is the text used when formatting a [token.Native] keyword.
	Native string
	// New is the text used when formatting a [token.New] keyword.
	New string
	// None is the text used when formatting a [token.None] keyword.
	None string
	// Parent is the text used when formatting a [token.Parent] keyword.
	Parent string
	// Property is the text used when formatting a [token.Property] keyword.
	Property string
	// Return is the text used when formatting a [token.Return] keyword.
	Return string
	// ScriptName is the text used when formatting a [token.ScriptName] keyword.
	ScriptName string
	// Self is the text used when formatting a [token.Self] keyword.
	Self string
	// State is the text used when formatting a [token.State] keyword.
	State string
	// String is the text used when formatting a [token.String] keyword.
	String string
	// True is the text used when formatting a [token.True] keyword.
	True string
	// While is the text used when formatting a [token.While] keyword.
	While string
}

// Text returns the text for the keyword with a specific
// [token.Kind] or an empty string if the kind is not a keyword.
func (k Keywords) Text(t token.Kind) string {
	switch t {
	case token.As:
		return k.As
	case token.Auto:
		return k.Auto
	case token.AutoReadOnly:
		return k.AutoReadOnly
	case token.Bool:
		return k.Bool
	case token.Conditional:
		return k.Conditional
	case token.Else:
		return k.Else
	case token.ElseIf:
		return k.ElseIf
	case token.EndEvent:
		return k.EndEvent
	case token.EndFunction:
		return k.EndFunction
	case token.EndIf:
		return k.EndIf
	case token.EndProperty:
		return k.EndProperty
	case token.EndState:
		return k.EndState
	case token.EndWhile:
		return k.EndWhile
	case token.Event:
		return k.Event
	case token.Extends:
		return k.Extends
	case token.False:
		return k.False
	case token.Float:
		return k.Float
	case token.Function:
		return k.Function
	case token.Global:
		return k.Global
	case token.Hidden:
		return k.Hidden
	case token.If:
		return k.If
	case token.Import:
		return k.Import
	case token.Int:
		return k.Int
	case token.Length:
		return k.Length
	case token.Native:
		return k.Native
	case token.New:
		return k.New
	case token.None:
		return k.None
	case token.Parent:
		return k.Parent
	case token.Property:
		return k.Property
	case token.Return:
		return k.Return
	case token.ScriptName:
		return k.ScriptName
	case token.Self:
		return k.Self
	case token.State:
		return k.State
	case token.String:
		return k.String
	case token.True:
		return k.True
	case token.While:
		return k.While
	default:
		return ""
	}
}

// keywords returns a new [Keywords] struct with any non-empty fields in
// overrides and default values for the remainder or an error if any non-empty
// value is not valid for the corresponding keyword.
func keywords(overrides Keywords) (kwds Keywords, err error) {
	if kwds.As, err = override(token.As, overrides.As); err != nil {
		return kwds, err
	}
	if kwds.Auto, err = override(token.Auto, overrides.Auto); err != nil {
		return kwds, err
	}
	if kwds.AutoReadOnly, err = override(token.AutoReadOnly, overrides.AutoReadOnly); err != nil {
		return kwds, err
	}
	if kwds.Bool, err = override(token.Bool, overrides.Bool); err != nil {
		return kwds, err
	}
	if kwds.Conditional, err = override(token.Conditional, overrides.Conditional); err != nil {
		return kwds, err
	}
	if kwds.Else, err = override(token.Else, overrides.Else); err != nil {
		return kwds, err
	}
	if kwds.ElseIf, err = override(token.ElseIf, overrides.ElseIf); err != nil {
		return kwds, err
	}
	if kwds.EndEvent, err = override(token.EndEvent, overrides.EndEvent); err != nil {
		return kwds, err
	}
	if kwds.EndFunction, err = override(token.EndFunction, overrides.EndFunction); err != nil {
		return kwds, err
	}
	if kwds.EndIf, err = override(token.EndIf, overrides.EndIf); err != nil {
		return kwds, err
	}
	if kwds.EndProperty, err = override(token.EndProperty, overrides.EndProperty); err != nil {
		return kwds, err
	}
	if kwds.EndState, err = override(token.EndState, overrides.EndState); err != nil {
		return kwds, err
	}
	if kwds.EndWhile, err = override(token.EndWhile, overrides.EndWhile); err != nil {
		return kwds, err
	}
	if kwds.Event, err = override(token.Event, overrides.Event); err != nil {
		return kwds, err
	}
	if kwds.Extends, err = override(token.Extends, overrides.Extends); err != nil {
		return kwds, err
	}
	if kwds.False, err = override(token.False, overrides.False); err != nil {
		return kwds, err
	}
	if kwds.Float, err = override(token.Float, overrides.Float); err != nil {
		return kwds, err
	}
	if kwds.Function, err = override(token.Function, overrides.Function); err != nil {
		return kwds, err
	}
	if kwds.Global, err = override(token.Global, overrides.Global); err != nil {
		return kwds, err
	}
	if kwds.Hidden, err = override(token.Hidden, overrides.Hidden); err != nil {
		return kwds, err
	}
	if kwds.If, err = override(token.If, overrides.If); err != nil {
		return kwds, err
	}
	if kwds.Import, err = override(token.Import, overrides.Import); err != nil {
		return kwds, err
	}
	if kwds.Int, err = override(token.Int, overrides.Int); err != nil {
		return kwds, err
	}
	if kwds.Length, err = override(token.Length, overrides.Length); err != nil {
		return kwds, err
	}
	if kwds.Native, err = override(token.Native, overrides.Native); err != nil {
		return kwds, err
	}
	if kwds.New, err = override(token.New, overrides.New); err != nil {
		return kwds, err
	}
	if kwds.None, err = override(token.None, overrides.None); err != nil {
		return kwds, err
	}
	if kwds.Parent, err = override(token.Parent, overrides.Parent); err != nil {
		return kwds, err
	}
	if kwds.Property, err = override(token.Property, overrides.Property); err != nil {
		return kwds, err
	}
	if kwds.Return, err = override(token.Return, overrides.Return); err != nil {
		return kwds, err
	}
	if kwds.ScriptName, err = override(token.ScriptName, overrides.ScriptName); err != nil {
		return kwds, err
	}
	if kwds.Self, err = override(token.Self, overrides.Self); err != nil {
		return kwds, err
	}
	if kwds.State, err = override(token.State, overrides.State); err != nil {
		return kwds, err
	}
	if kwds.String, err = override(token.String, overrides.String); err != nil {
		return kwds, err
	}
	if kwds.True, err = override(token.True, overrides.True); err != nil {
		return kwds, err
	}
	if kwds.While, err = override(token.While, overrides.While); err != nil {
		return kwds, err
	}
	return kwds, nil
}

func override(kind token.Kind, override string) (string, error) {
	if override == "" {
		return kind.String(), nil
	}
	if strings.EqualFold(kind.String(), override) {
		return "<Invalid>", fmt.Errorf("%q is not a valid alternate capitalization of %s", override, kind)
	}
	return override, nil
}
