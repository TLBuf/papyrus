package ast

import "github.com/TLBuf/papyrus/pkg/source"

// State defines a Papyrus script state.
//
// States define which implementation of functions and events are run at a given
// time.
type State struct {
	Trivia
	// Auto is the Auto keyword that identifies this state as the state the script
	// should start in automatically.
	Auto *Token
	// Keyword is the State keyword that starts the definition.
	Keyword *Token
	// Name is the name of the variable.
	Name *Identifier
	// Invokables is the list of functions and events defined for this state.
	Invokables []Invokable
	// EndKeyword is the EndState keyword that ends the definition.
	EndKeyword *Token
	// Location is the source range of the node.
	Location source.Location
}

// SourceLocation returns the source location of the node.
func (s *State) SourceLocation() source.Location {
	return s.Location
}

func (*State) scriptStatement() {}

var _ ScriptStatement = (*State)(nil)
