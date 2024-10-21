package ast

import "github.com/TLBuf/papyrus/pkg/source"

// State defines a Papyrus script state.
//
// States define which implementation of functions and events are run at a given time.
type State struct {
	// Name is the name of the variable.
	Name *Identifier
	// IsAuto
	IsAuto bool
	// Invokables is the list of functions and events defined for this state.
	Invokables []Invokable
	// SourceRange is the source range of the node.
	SourceRange source.Range
}

// Range returns the source range of the node.
func (s *State) Range() source.Range {
	return s.SourceRange
}

func (*State) scriptStatement() {}

var _ ScriptStatement = (*State)(nil)
