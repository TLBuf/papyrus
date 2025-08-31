package analysis

import (
	"github.com/TLBuf/papyrus/analysis/symbol"
	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/literal"
	"github.com/TLBuf/papyrus/types"
)

// Info holds type information for a type-checked set of scripts.
type Info struct {
	// Expressions maps expressions to their types.
	Expressions map[ast.Expression]types.Type

	// Values maps literals to their values.
	Values map[ast.Literal]literal.Value

	// Scopes maps AST nodes to the scopes they define.
	//
	// Scopes nest, with the Global scope being the outermost scope, enclosing
	// the Script scope, which contains zero or more other scopes.
	//
	// The following node types may appear in Scopes:
	//
	//   - [ast.Script]
	//   - [ast.Property]
	//   - [ast.State]
	//   - [ast.Function]
	//   - [ast.Event]
	//   - [ast.If]
	//   - [ast.Else]
	//   - [ast.ElseIf]
	//   - [ast.While]
	Scopes map[ast.Node]*symbol.Scope

	// Symbols maps AST nodes to the symbols they define.
	//
	// The following node types may appear in Scopes:
	//
	//   - [ast.Script]
	//   - [ast.State]
	//   - [ast.Function]
	//   - [ast.Event]
	//   - [ast.Property]
	//   - [ast.Variable]
	//   - [ast.Parameter]
	Symbols map[ast.Node]*symbol.Symbol

	// Global is the outermost scope for all checked scripts.
	Global *symbol.Scope
}
