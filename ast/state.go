package ast

import "github.com/TLBuf/papyrus/source"

// State defines a Papyrus script state.
//
// States define which implementation of functions and events are run at a given
// time.
type State struct {
	LineTrivia
	// IsAuto is true if the script should start in this state automatically.
	IsAuto bool
	// AutoLocation is the location of the Auto keyword that identifies this state
	// as the state the script should start in automatically.
	//
	// This is only valid if IsAuto is true.
	AutoLocation source.Location
	// StartKeywordLocation is the location of the State keyword that starts
	// the statement.
	StartKeywordLocation source.Location
	// Name is the name of the variable.
	Name *Identifier
	// Invokables is the list of functions and events defined for this state.
	Invokables []Invokable
	// EndKeywordLocation is the location of the EndState keyword that ends the
	// statement.
	EndKeywordLocation source.Location
}

// Trivia returns the [LineTrivia] assocaited with this node.
func (s *State) Trivia() LineTrivia {
	return s.LineTrivia
}

// Accept calls the appropriate visitor method for the node.
func (s *State) Accept(v Visitor) error {
	return v.VisitState(s)
}

// Location returns the source location of the node.
func (s *State) Location() source.Location {
	if s.IsAuto {
		return source.Span(s.AutoLocation, s.EndKeywordLocation)
	}
	return source.Span(s.StartKeywordLocation, s.EndKeywordLocation)
}

func (*State) statement() {}

func (*State) scriptStatement() {}

var _ ScriptStatement = (*State)(nil)
