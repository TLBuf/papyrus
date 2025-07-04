// Package token defines the Papyrus tokens understood by the parser.
package token

import (
	"strings"

	"github.com/TLBuf/papyrus/source"
)

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

// String returns the text of the token as it appears in the lexed source.
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

// Kind is the type of token.
type Kind byte

const (
	// Illegal is the zero value token denoting a missing or malformed token.
	Illegal Kind = iota
	// EOF denotes the end of a token stream.
	EOF
	// As is the casting keyword.
	As
	// Assign, '=', is the assignment operator symbol.
	Assign
	// AssignAdd, '+=', is the addition assignment operator symbol.
	AssignAdd
	// AssignDivide, '/=', is the division assignment operator symbol.
	AssignDivide
	// AssignModulo, '%=', is the modular assignment operator symbol.
	AssignModulo
	// AssignMultiply, '*=', is the multiplication assignment operator symbol.
	AssignMultiply
	// AssignSubtract, '-=', is the subtraction assignment operator symbol.
	AssignSubtract
	// Auto is the auto property keyword.
	Auto
	// AutoReadOnly is the auto read-only property keyword.
	AutoReadOnly
	// BlockCommentClose, '/;', is the symbol that ends a block comment.
	BlockCommentClose
	// BlockCommentClose, '/;', is the symbol that starts a block comment.
	BlockCommentOpen
	// Bool is the bool type keyword.
	Bool
	// BoolArray is the bool array type keyword.
	BoolArray
	// BraceClose, '}', is the symbol that ends a documentation comment.
	BraceClose
	// BraceOpen, '{', is the symbol that starts a documentation comment.
	BraceOpen
	// BracketClose, ']', is the closing bracket symbol used by arrays.
	BracketClose
	// BracketOpen, '[', is the openning bracket symbol used by arrays.
	BracketOpen
	// Comma, ',', is the symbol used to separate list elements.
	Comma
	// Comment denotes a token that contains all content of a comment.
	Comment
	// Conditional is the conditional flag keyword.
	Conditional
	// Divide, '/', is the division symbol.
	Divide
	// Dot, '.', is the access operator symbol.
	Dot
	// Else is the unconditional alternative keyword.
	Else
	// ElseIf is the conditional alternative keyword.
	ElseIf
	// EndEvent is the keyword that ends an event block.
	EndEvent
	// EndFunction is the keyword that ends a function block.
	EndFunction
	// EndIf is the keyword that ends a conditional block.
	EndIf
	// EndProperty is the keyword that ends a property block.
	EndProperty
	// EndState is the keyword that ends a state block.
	EndState
	// EndWhile is the keyword that ends a loop block.
	EndWhile
	// Equal, '==', is the equality comparison operator symbol.
	Equal
	// Event is the keyword that starts an event block.
	Event
	// Extends is the keyword used to extend another script.
	Extends
	// False is the boolean literal false value keyword.
	False
	// Float is the floating-point type keyword.
	Float
	// FloatArray is the floating-point array type keyword.
	FloatArray
	// FloatLiteral denotes a floating-point literal value.
	FloatLiteral
	// Function is the keyword that starts a function block.
	Function
	// Global is the global flag keyword.
	Global
	// Greater, '>', is the greater than comparison operator symbol.
	Greater
	// GreaterOrEqual, '>=', is the greater than or equal to comparison operator symbol.
	GreaterOrEqual
	// Hidden is the hidden flag keyword.
	Hidden
	// Identifier denotes a non-keyword identifier.
	Identifier
	// If is the keyword that starts a conditional block.
	If
	// Import is the keyword used to import another script into the current one.
	Import
	// Int is the integer type keyword.
	Int
	// IntArray is the integer array type keyword.
	IntArray
	// IntLiteral denotes an integer literal value.
	IntLiteral
	// Length is the array length keyword.
	Length
	// Less, '<', is the less than comparison operator symbol.
	Less
	// LessOrEqual, '<=', is the less than or equal to comparison operator symbol.
	LessOrEqual
	// LogicalAnd, '&&', is the logical AND operator symbol.
	LogicalAnd
	// LogicalNot, '!', is the logical NOT operator symbol.
	LogicalNot
	// LogicalOr, '||', is the logical OR operator symbol.
	LogicalOr
	// Minus, '-', is the subtraction or unary numeric negation operator symbol.
	Minus
	// Modulo, '%', is the modular operator symbol.
	Modulo
	// Multiply, '*', is the multiplication operator symbol.
	Multiply
	// Native is the native flag keyword.
	Native
	// New is the keyword used to create new arrays.
	New
	// Newline denotes a line break (end of statement).
	Newline
	// None is the object literal empty value keyword.
	None
	// NotEqual, '!=', is the negaive equality comparison operator symbol.
	NotEqual
	// ObjectArray is the object array type keyword.
	ObjectArray
	// Parent is the keyword used to refer to symbols in an extended script.
	Parent
	// ParenthesisClose, ')', is the is the closing symbol used by parameter and argument lists and parentheticals.
	ParenthesisClose
	// ParenthesisOpen, '(', is the is the openning symbol used by parameter and argument lists and parentheticals.
	ParenthesisOpen
	// Plus, '+', is the addition operator symbol.
	Plus
	// Property is the keyword that starts a property definition or block.
	Property
	// Return is the keyword used to end execution of a function or event.
	Return
	// ScriptName is the keyword used define a script's name.
	ScriptName
	// Self is the keyword used to refer to symbols in the script itself.
	Self
	// Semicolon is the keyword that denotes the start of a line comment.
	Semicolon
	// State is the keyword that starts a state block.
	State
	// String is the string type keyword.
	String
	// StringArray is the string array type keyword.
	StringArray
	// StringLiteral denotes a string literal value.
	StringLiteral
	// True is the boolean literal true value keyword.
	True
	// While is the keyword that starts a loop block.
	While
)

// Keyword returns the string representation of this
// kind or an empty string if it is not a keyword.
//
// This method will always return a standardized
// capitalization regardless of any lexed text.
func (k Kind) Keyword() string {
	if k.IsKeyword() {
		return names[k]
	}
	return ""
}

// IsKeyword returns true if this kind is an
// alphabetic keyword and false otherwise.
func (k Kind) IsKeyword() bool {
	switch k {
	case As,
		Auto,
		AutoReadOnly,
		Bool,
		BoolArray,
		Conditional,
		Else,
		ElseIf,
		EndEvent,
		EndFunction,
		EndIf,
		EndProperty,
		EndState,
		EndWhile,
		Event,
		Extends,
		False,
		Float,
		FloatArray,
		Function,
		Global,
		Hidden,
		If,
		Import,
		Int,
		IntArray,
		Length,
		Native,
		New,
		None,
		Parent,
		Property,
		Return,
		ScriptName,
		Self,
		State,
		String,
		StringArray,
		True,
		While:
		return true
	default:
		return false
	}
}

// Symbol returns the string representation of this
// kind or an empty string if it is not a symbol.
func (k Kind) Symbol() string {
	if k.IsSymbol() {
		return names[k]
	}
	return ""
}

// IsSymbol returns true if this kind is a
// non-alphabetic symbol and false otherwise.
func (k Kind) IsSymbol() bool {
	switch k {
	case Assign,
		AssignAdd,
		AssignDivide,
		AssignModulo,
		AssignMultiply,
		AssignSubtract,
		BlockCommentClose,
		BlockCommentOpen,
		BraceClose,
		BraceOpen,
		BracketClose,
		BracketOpen,
		Comma,
		Divide,
		Dot,
		Equal,
		Greater,
		GreaterOrEqual,
		Less,
		LessOrEqual,
		LogicalAnd,
		LogicalNot,
		LogicalOr,
		Minus,
		Modulo,
		Multiply,
		NotEqual,
		ParenthesisClose,
		ParenthesisOpen,
		Plus,
		Semicolon:
		return true
	default:
		return false
	}
}

// String returns the string representation of this Kind.
//
// If the Kind is a symbol or keyword, this returns the same string as
// [Kind.Symbol] and [Kind.Keyword] respectively, otherwise it returns the name
// of the token surrounded by angle brackets.
//
// This method will always return a standardized capitalization regardless of
// any lexed text.
func (k Kind) String() string {
	if int(k) < len(names) {
		return names[k]
	}
	return "<Unknown>"
}

var names = []string{
	"<Illegal>",
	"<EOF>",
	"As",
	"=",
	"+=",
	"/=",
	"%=",
	"*=",
	"-=",
	"Auto",
	"AutoReadOnly",
	"/;",
	";/",
	"Bool",
	"Bool[]",
	"}",
	"{",
	"]",
	"[",
	",",
	"<Comment>",
	"Conditional",
	"/",
	".",
	"Else",
	"ElseIf",
	"EndEvent",
	"EndFunction",
	"EndIf",
	"EndProperty",
	"EndState",
	"EndWhile",
	"==",
	"Event",
	"Extends",
	"False",
	"Float",
	"Float[]",
	"<FloatLiteral>",
	"Function",
	"Global",
	">",
	">=",
	"Hidden",
	"<Identifier>",
	"If",
	"Import",
	"Int",
	"Int[]",
	"<IntLiteral>",
	"Length",
	"<",
	"<=",
	"&&",
	"!",
	"||",
	"-",
	"%",
	"*",
	"Native",
	"New",
	"<Newline>",
	"None",
	"!=",
	"<Object[]>",
	"Parent",
	")",
	"(",
	"+",
	"Property",
	"Return",
	"ScriptName",
	"Self",
	";",
	"State",
	"String",
	"String[]",
	"<StringLiteral>",
	"True",
	"While",
}
