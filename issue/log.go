package issue

import (
	"iter"
	"strings"
)

// Log is an ordered record of issues collected during processing.
type Log struct {
	issues   []*Issue
	internal int
	errors   int
	warnings int
	infos    int
}

// NewLog returns an empty log.
func NewLog() *Log {
	return &Log{}
}

// Append appends an issue to the log.
func (l *Log) Append(issue *Issue) {
	switch issue.Definition().Severity() {
	case Internal:
		l.internal++
	case Error:
		l.errors++
	case Warning:
		l.warnings++
	case Info:
		l.infos++
	}
	l.issues = append(l.issues, issue)
}

// Len returns the total number of issues in the log.
func (l *Log) Len() int {
	return len(l.issues)
}

// First returns the first issue in
// the log or nil if the log is empty.
func (l *Log) First() *Issue {
	if len(l.issues) == 0 {
		return nil
	}
	return l.issues[0]
}

// Last returns the last issue in the log (i.e. the
// most recently appended) or nil if the log is empty.
func (l *Log) Last() *Issue {
	if len(l.issues) == 0 {
		return nil
	}
	return l.issues[len(l.issues)-1]
}

// HasInternal returns true if the log has at least one [Internal] issue.
func (l *Log) HasInternal() bool {
	return l.internal > 0
}

// HasError returns true if the log has at least one [Error] issue.
func (l *Log) HasError() bool {
	return l.errors > 0
}

// HasWarning returns true if the log has at least one [Warning] issue.
func (l *Log) HasWarning() bool {
	return l.warnings > 0
}

// HasInfo returns true if the log has at least one [Info] issue.
func (l *Log) HasInfo() bool {
	return l.infos > 0
}

// All returns an iterator over all issues in the log.
func (l *Log) All() iter.Seq2[int, *Issue] {
	return func(yield func(int, *Issue) bool) {
		for i, u := range l.issues {
			if !yield(i, u) {
				return
			}
		}
	}
}

// Errors returns an iterator over all [Error] issues in the log.
func (l *Log) Errors() iter.Seq[*Issue] {
	if l.errors == 0 {
		return func(func(*Issue) bool) {}
	}
	return l.iter(Error)
}

// Warnings returns an iterator over all [Warning] issues in the log.
func (l *Log) Warnings() iter.Seq[*Issue] {
	if l.warnings == 0 {
		return func(func(*Issue) bool) {}
	}
	return l.iter(Warning)
}

// Infos returns an iterator over all [Info] issues in the log.
func (l *Log) Infos() iter.Seq[*Issue] {
	if l.infos == 0 {
		return func(func(*Issue) bool) {}
	}
	return l.iter(Info)
}

func (l *Log) iter(severity Severity) iter.Seq[*Issue] {
	return func(yield func(*Issue) bool) {
		for _, i := range l.issues {
			if i.Definition().Severity() == severity && !yield(i) {
				return
			}
		}
	}
}

func (l *Log) String() string {
	var sb strings.Builder
	for _, i := range l.issues {
		sb.WriteString(i.String())
		sb.WriteRune('\n')
	}
	return sb.String()
}
