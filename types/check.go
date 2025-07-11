package types

import (
	"fmt"

	"github.com/TLBuf/papyrus/ast"
)

// Check performs type-checking over some number of
// scripts and collated summarized type information.
func Check(scripts ...*ast.Script) (*Info, error) {
	return nil, fmt.Errorf("not implemented")
}
