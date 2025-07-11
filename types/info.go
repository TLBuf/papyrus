package types

import (
	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/literal"
)

// Info holds type information for a type-checked set of scripts.
type Info struct {
	// Types maps expressions to their types.
	Types map[ast.Expression]Type

	// Values maps literals to their values.
	Values map[ast.Literal]literal.Value

	// Entities maps identifiers to the entity they identify.
	Entities map[*ast.Identifier]Entity

	// Scopes maps AST nodes to the scopes they define
	//
	// Scopes nest, with the Global scope being the outermost scope, enclosing
	// the Script scope, which contains zero or more other scopes.
	//
	// The following node types may appear in Scopes:
	//
	//    - [*ast.Script]
	//    - [*ast.Property]
	//    - [*ast.State]
	//    - [*ast.Function]
	//    - [*ast.Event]
	//    - [*ast.If]
	//    - [*ast.Else]
	//    - [*ast.ElseIf]
	//    - [*ast.While]
	//
	Scopes map[ast.Node]*Scope

	// Global is the outermost scope for all checked scripts.
	Global *Scope
}
