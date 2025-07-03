package ast

import "github.com/TLBuf/papyrus/source"

// State defines a Papyrus script state.
//
// States define which implementation of functions and events are run at a given
// time.
type State struct {
	LineTrivia
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

// Trivia returns the [LineTrivia] assocaited with this node.
func (s *State) Trivia() LineTrivia {
	return s.LineTrivia
}

// Accept calls the appropriate visitor method for the node.
func (s *State) Accept(v Visitor) error {
	return v.VisitState(s)
}

// SourceLocation returns the source location of the node.
func (s *State) SourceLocation() source.Location {
	return s.Location
}

func (*State) statement() {}

func (*State) scriptStatement() {}

var _ ScriptStatement = (*State)(nil)
