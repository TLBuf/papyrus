package analysis

import (
	"github.com/TLBuf/papyrus/analysis/symbol"
	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/literal"
	"github.com/TLBuf/papyrus/types"
)

// Info holds type information for a type-checked set of scripts.
type Info struct {
	// Types maps expressions to their types.
	Types map[ast.Expression]types.Type

	// Values maps literals to their values.
	Values map[ast.Literal]literal.Value

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
	Scopes map[ast.Node]*symbol.Scope

	// Global is the outermost scope for all checked scripts.
	Global *symbol.Scope
}
