package analysis

import (
	"fmt"
	"strings"

	"github.com/TLBuf/papyrus/analysis/symbol"
	"github.com/TLBuf/papyrus/ast"
	"github.com/TLBuf/papyrus/literal"
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/types"
)

// Error defines an error raised by the type checker.
type Error struct {
	// The underlying error.
	Err error
	// Location identifies the place in the source that caused the error.
	Location source.Location
}

// Error implments the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Location, e.Err)
}

// Unwrap returns the underlying error.
func (e Error) Unwrap() error {
	return e.Err
}

func newError(location source.Location, msg string, args ...any) Error {
	return Error{
		Err:      fmt.Errorf(msg, args...),
		Location: location,
	}
}

// Check performs type-checking over some number of
// scripts and collated summarized type information.
func Check(scripts ...*ast.Script) (*Info, error) {
	global := symbol.NewGlobalScope()
	checker := &checker{
		global: global,
		info: &Info{
			Types:   make(map[ast.Expression]types.Type),
			Values:  make(map[ast.Literal]literal.Value),
			Symbols: make(map[*ast.Identifier]*symbol.Symbol),
			Scopes:  make(map[ast.Node]*symbol.Scope),
			Global:  global,
		},
		typeNames: make(map[string]types.Type),
		scope:     global,
	}
	for _, t := range []types.Type{types.Bool, types.BoolArray, types.Int, types.IntArray, types.Float, types.FloatArray, types.String, types.StringArray} {
		checker.typeNames[normalize(t.Name())] = t
	}
	if err := checker.check(scripts); err != nil {
		return nil, err
	}
	return checker.info, nil
}

type checker struct {
	info      *Info
	global    *symbol.Scope
	typeNames map[string]types.Type
	script    *symbol.Symbol
	state     *symbol.Symbol
	scope     *symbol.Scope
}

// file returns the current source file being checked.
func (c *checker) file() *source.File {
	return c.script.Node().(*ast.Script).File
}

// self returns the type of the current script being checked.
func (c *checker) self() *types.Object {
	return c.script.Type().(*types.Object)
}

func (c *checker) enterScope(s *symbol.Scope) {
	c.scope = s
}

func (c *checker) leaveScope() {
	c.scope = c.scope.Parent()
}

func (c *checker) check(scripts []*ast.Script) error {
	// TODO: Implement.
	return nil
}

func normalize(name string) string {
	return strings.ToLower(name)
}
