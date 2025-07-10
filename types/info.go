package types

import (
	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/literal"
)

// Info holds type information for a type-checked set of scripts.
type Info struct {
	Types  map[ast.Expression]Type
	Values map[ast.Literal]literal.Value
}
