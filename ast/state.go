package ast

import "github.com/TLBuf/papyrus/source"

// State defines a Papyrus script state.
//
// States define which implementation of functions and events are run at a given
// time.
type State struct {
	// IsAuto is true if the script should start in this state automatically.
	IsAuto bool
	// HasPrecedingBlankLine is true if this node was preceded by a blank line.
	HasPrecedingBlankLine bool
	// Name is the name of the variable.
	Name *Identifier
	// Invokables is the list of functions and events defined for this state.
	Invokables []Invokable
	// AutoLocation is the location of the Auto keyword that identifies this state
	// as the state the script should start in automatically.
	//
	// This is only valid if IsAuto is true.
	AutoLocation source.Location
	// StartKeywordLocation is the location of the State keyword that starts
	// the statement.
	StartKeywordLocation source.Location
	// EndKeywordLocation is the location of the EndState keyword that ends the
	// statement.
	EndKeywordLocation source.Location
	// NodeComments are the comments on lines before and/or after a
	// node or nil if the node has no comments associated with it.
	NodeComments *CrosslineComments
}

// PrecedingBlankLine returns true if this node was preceded by a blank line.
func (s *State) PrecedingBlankLine() bool {
	return s.HasPrecedingBlankLine
}

// Accept calls the appropriate visitor method for the node.
func (s *State) Accept(v Visitor) error {
	return v.VisitState(s)
}

// Comments returns the [CrosslineComments] associated
// with this node or nil if there are none.
func (s *State) Comments() *CrosslineComments {
	return s.NodeComments
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
