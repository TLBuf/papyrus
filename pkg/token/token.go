// Package token defines the Papyrus tokens understood by the parser.
package token

import (
	"strings"

	"github.com/TLBuf/papyrus/pkg/source"
)

// Kind is the type of token.
type Kind byte

// The set of tokens for Papyrus.
const (
	Illegal Kind = iota
	EOF
	As
	Assign
	AssignAdd
	AssignDivide
	AssignModulo
	AssignMultiply
	AssignSubtract
	Auto
	AutoReadOnly
	BlockCommentClose
	BlockCommentOpen
	Bool
	BraceClose
	BraceOpen
	BracketClose
	BracketOpen
	Comma
	Comment
	Conditional
	Divide
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
	Length
	Less
	LessOrEqual
	LogicalAnd
	LogicalNot
	LogicalOr
	Minus
	Modulo
	Multiply
	Native
	New
	Newline
	None
	NotEqual
	Parent
	ParenthesisClose
	ParenthesisOpen
	Plus
	Property
	Return
	ScriptName
	Self
	Semicolon
	State
	String
	StringLiteral
	True
	While
)

// String implements the [fmt.Stringer] interface.
func (t Kind) String() string {
	name, ok := names[t]
	if ok {
		return name
	}
	return "<unknown>"
}

// Token encodes a single lexical token in the Papyrus language.
//
// Each token has a [Kind] and information about where it is located
// in a source file.
type Token struct {
	// Kind defines the specific kind of token the text represents.
	Kind Kind
	// Location describes exactly where in a file the token is.
	Location source.Location
}

// String implements the [fmt.Stringer] interface.
func (t Token) String() string {
	return string(t.Location.Text())
}

// LookupIdentifier returns the [Kind] of the given identifier or keyword.
func LookupIdentifier(ident string) Kind {
	if t, ok := keywords[strings.ToLower(ident)]; ok {
		return t
	}
	return Identifier
}

var keywords = map[string]Kind{
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

// Returns the symbol representation
func Symbol(k Kind) string {
	if s, ok := symbols[k]; ok {
		return s
	}
	return ""
}

var symbols = map[Kind]string{
	Assign:            "=",
	AssignAdd:         "+=",
	AssignDivide:      "/=",
	AssignModulo:      "%=",
	AssignMultiply:    "*=",
	AssignSubtract:    "-=",
	BlockCommentClose: "/;",
	BlockCommentOpen:  ";/",
	BraceClose:        "}",
	BraceOpen:         "{",
	BracketClose:      "]",
	BracketOpen:       "[",
	Comma:             ",",
	Divide:            "/",
	Dot:               ".",
	Equal:             "==",
	Greater:           ">",
	GreaterOrEqual:    ">=",
	LogicalAnd:        "&&",
	LogicalNot:        "!",
	LogicalOr:         "||",
	Minus:             "-",
	Modulo:            "%",
	Multiply:          "*",
	NotEqual:          "!=",
	ParenthesisClose:  ")",
	ParenthesisOpen:   "(",
	Plus:              "+",
	Semicolon:         ";",
}

var names = map[Kind]string{
	Illegal:           "Illegal",
	EOF:               "EOF",
	As:                "As",
	Assign:            "Assign",
	AssignAdd:         "AssignAdd",
	AssignDivide:      "AssignDivide",
	AssignModulo:      "AssignModulo",
	AssignMultiply:    "AssignMultiply",
	AssignSubtract:    "AssignSubtract",
	Auto:              "Auto",
	AutoReadOnly:      "AutoReadOnly",
	BlockCommentClose: "BlockCommentClose",
	BlockCommentOpen:  "BlockCommentOpen",
	Bool:              "Bool",
	BraceClose:        "BraceClose",
	BraceOpen:         "BraceOpen",
	BracketClose:      "BracketClose",
	BracketOpen:       "BracketOpen",
	Comma:             "Comma",
	Comment:           "Comment",
	Conditional:       "Conditional",
	Divide:            "Divide",
	Dot:               "Dot",
	Else:              "Else",
	ElseIf:            "ElseIf",
	EndEvent:          "EndEvent",
	EndFunction:       "EndFunction",
	EndIf:             "EndIf",
	EndProperty:       "EndProperty",
	EndState:          "EndState",
	EndWhile:          "EndWhile",
	Equal:             "Equal",
	Event:             "Event",
	Extends:           "Extends",
	False:             "False",
	Float:             "Float",
	FloatLiteral:      "FloatLiteral",
	Function:          "Function",
	Global:            "Global",
	Greater:           "Greater",
	GreaterOrEqual:    "GreaterOrEqual",
	Hidden:            "Hidden",
	Identifier:        "Identifier",
	If:                "If",
	Import:            "Import",
	Int:               "Int",
	IntLiteral:        "IntLiteral",
	Length:            "Length",
	Less:              "Less",
	LessOrEqual:       "LessOrEqual",
	LogicalAnd:        "LogicalAnd",
	LogicalNot:        "LogicalNot",
	LogicalOr:         "LogicalOr",
	Minus:             "Minus",
	Modulo:            "Modulo",
	Multiply:          "Multiply",
	Native:            "Native",
	New:               "New",
	Newline:           "Newline",
	None:              "None",
	NotEqual:          "NotEqual",
	Parent:            "Parent",
	ParenthesisClose:  "ParenthesisClose",
	ParenthesisOpen:   "ParenthesisOpen",
	Plus:              "Plus",
	Property:          "Property",
	Return:            "Return",
	ScriptName:        "ScriptName",
	Self:              "Self",
	Semicolon:         "Semicolon",
	State:             "State",
	String:            "String",
	StringLiteral:     "StringLiteral",
	True:              "True",
	While:             "While",
}
