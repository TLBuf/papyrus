package symbol

import (
	"fmt"

	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/types"
)

// Kind defines the different kinds of symbols.
type Kind uint16

const (
	// ScriptKind is the kind of symbol defined by an [*ast.Script], always
	// defines a [Scope], and is always of type [*types.Object].
	ScriptKind Kind = 1 << iota
	// StateKind is the kind of symbol defined by an [*ast.State], always defines
	// a [Scope], and unlike all other symbols, has no type.
	StateKind
	// FunctionKind is the kind of symbol defined by an [*ast.Function], defines a
	// scope if the function is not native, and is always of type
	// [*types.FunctionKind].
	FunctionKind
	// EventKind is the kind of symbol defined by an [*ast.Event], defines a scope
	// if the event is not native, and is always of type [*types.Function].
	EventKind
	// PropertyKind is the kind of symbol defined by an [*ast.Property], defines a
	// scope only if the definition is a full property, and is always of type
	// [types.Value].
	PropertyKind
	// VariableKind is the kind of symbol defined by an [*ast.Variable], never
	// defines a scope, and is always of type [types.Value].
	VariableKind
	// ParameterKind is the kind of symbol defined by an [*ast.Parameter], never
	// defines a scope, and is always of type [types.Value].
	ParameterKind
)

func (k Kind) String() string {
	switch k {
	case ScriptKind:
		return "Script"
	case StateKind:
		return "State"
	case FunctionKind:
		return "Function"
	case EventKind:
		return "Event"
	case PropertyKind:
		return "Property"
	case VariableKind:
		return "Variable"
	case ParameterKind:
		return "Parameter"
	default:
		return fmt.Sprintf("Unknown Kind (%d)", k)
	}
}

// Symbol defines a named, usually typed, Papyrus entity.
type Symbol struct {
	enclosing  *Scope
	scope      *Scope
	typ        types.Type
	node       ast.Node
	name       string
	normalized string
	kind       Kind
}

// Kind returns the kind of this symbol.
func (s *Symbol) Kind() Kind {
	return s.kind
}

// Enclosing returns the scope in which this symbol is defined.
func (s *Symbol) Enclosing() *Scope {
	return s.enclosing
}

// Scope returns the scope defined by this
// symbol or nil if it doesn't define one.
func (s *Symbol) Scope() *Scope {
	return s.scope
}

// Type returns the type of this symbol.
func (s *Symbol) Type() types.Type {
	return s.typ
}

// Name returns the declared name of this symbol.
func (s *Symbol) Name() string {
	return s.name
}

// Normalized returns the normalized name of this symbol.
func (s *Symbol) Normalized() string {
	return s.normalized
}

// Node returns the [ast.Node] that defines this symbol or
// nil if this symbol is implicit (i.e. the empty state).
func (s *Symbol) Node() ast.Node {
	return s.node
}

func (s Symbol) String() string {
	return fmt.Sprintf("%v: %q", s.kind, s.name)
}
