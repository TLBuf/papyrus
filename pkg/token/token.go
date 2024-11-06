// Package token defines the Papyrus tokens understood by the parser.
package token

import (
	"strings"

	"github.com/TLBuf/papyrus/pkg/source"
)

// Type is the type of token.
type Type byte

// The set of tokens for Papyrus.
const (
	Illegal Type = iota
	EOF
	Add
	As
	Assign
	AssignAdd
	AssignDivide
	AssignModulo
	AssignMultiply
	AssignSubtract
	Auto
	AutoReadOnly
	BlockComment
	Bool
	Comma
	Conditional
	Divide
	DocComment
	Dot
	Else
	ElseIf
	EndEvent
	EndFunction
	EndIf
	EndProperty
	EndState
	EndWhile
	Equal
	Event
	Extends
	False
	Float
	FloatLiteral
	Function
	Global
	Greater
	GreaterOrEqual
	Hidden
	Identifier
	If
	Import
	Int
	IntLiteral
	LBracket
	Length
	Less
	LessOrEqual
	LineComment
	LogicalAnd
	LogicalNot
	LogicalOr
	LParen
	Modulo
	Multiply
	Native
	New
	Newline
	None
	NotEqual
	Parent
	Property
	RBracket
	Return
	RParen
	ScriptName
	Self
	State
	String
	StringLiteral
	Subtract
	True
	While
)

func (t Type) String() string {
	name, ok := names[t]
	if ok {
		return name
	}
	return "<unknown>"
}

// Token encodes a single lexical token in the Papyrus language.
//
// Each token has a [Type] and information about where it is located
// in a source file.
type Token struct {
	Type        Type
	SourceRange source.Range
}

// LookupIdentifier returns the [Type] of the given identifier or keyword.
func LookupIdentifier(ident string) Type {
	if t, ok := keywords[strings.ToLower(ident)]; ok {
		return t
	}
	return Identifier
}

var keywords = map[string]Type{
	"as":           As,
	"auto":         Auto,
	"autoreadonly": AutoReadOnly,
	"bool":         Bool,
	"conditional":  Conditional,
	"else":         Else,
	"elseif":       ElseIf,
	"endevent":     EndEvent,
	"endfunction":  EndFunction,
	"endif":        EndIf,
	"endproperty":  EndProperty,
	"endstate":     EndState,
	"endwhile":     EndWhile,
	"event":        Event,
	"extends":      Extends,
	"false":        False,
	"float":        Float,
	"function":     Function,
	"global":       Global,
	"hidden":       Hidden,
	"if":           If,
	"import":       Import,
	"int":          Int,
	"length":       Length,
	"native":       Native,
	"new":          New,
	"none":         None,
	"parent":       Parent,
	"property":     Property,
	"return":       Return,
	"scriptname":   ScriptName,
	"self":         Self,
	"state":        State,
	"string":       String,
	"true":         True,
	"while":        While,
}

var names = map[Type]string{
	Illegal:        "Illegal",
	EOF:            "EOF",
	Add:            "Add",
	As:             "As",
	Assign:         "Assign",
	AssignAdd:      "AssignAdd",
	AssignDivide:   "AssignDivide",
	AssignModulo:   "AssignModulo",
	AssignMultiply: "AssignMultiply",
	AssignSubtract: "AssignSubtract",
	Auto:           "Auto",
	AutoReadOnly:   "AutoReadOnly",
	BlockComment:   "BlockComment",
	Bool:           "Bool",
	Comma:          "Comma",
	Conditional:    "Conditional",
	Divide:         "Divide",
	DocComment:     "DocComment",
	Dot:            "Dot",
	Else:           "Else",
	ElseIf:         "ElseIf",
	EndEvent:       "EndEvent",
	EndFunction:    "EndFunction",
	EndIf:          "EndIf",
	EndProperty:    "EndProperty",
	EndState:       "EndState",
	EndWhile:       "EndWhile",
	Equal:          "Equal",
	Event:          "Event",
	Extends:        "Extends",
	False:          "False",
	Float:          "Float",
	FloatLiteral:   "FloatLiteral",
	Function:       "Function",
	Global:         "Global",
	Greater:        "Greater",
	GreaterOrEqual: "GreaterOrEqual",
	Hidden:         "Hidden",
	Identifier:     "Identifier",
	If:             "If",
	Import:         "Import",
	Int:            "Int",
	IntLiteral:     "IntLiteral",
	LBracket:       "LBracket",
	Length:         "Length",
	Less:           "Less",
	LessOrEqual:    "LessOrEqual",
	LineComment:    "LineComment",
	LogicalAnd:     "LogicalAnd",
	LogicalNot:     "LogicalNot",
	LogicalOr:      "LogicalOr",
	LParen:         "LParen",
	Modulo:         "Modulo",
	Multiply:       "Multiply",
	Native:         "Native",
	New:            "New",
	Newline:        "Newline",
	None:           "None",
	NotEqual:       "NotEqual",
	Parent:         "Parent",
	Property:       "Property",
	RBracket:       "RBracket",
	Return:         "Return",
	RParen:         "RParen",
	ScriptName:     "ScriptName",
	Self:           "Self",
	State:          "State",
	String:         "String",
	StringLiteral:  "StringLiteral",
	Subtract:       "Subtract",
	True:           "True",
	While:          "While",
}
