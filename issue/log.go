package issue

import (
	"iter"
)

// Log is an ordered record of issues collected during processing.
type Log struct {
	issues   []Issue
	errors   int
	warnings int
	infos    int
}

// NewLog returns an empty log.
func NewLog() *Log {
	return &Log{}
}

// Append appends an issue to the log.
func (l *Log) Append(issue Issue) {
	switch issue.Severity {
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

// HasErrors returns true if the log has at least one [Error] issue.
func (l *Log) HasErrors() bool {
	return l.errors > 0
}

// All returns an iterator over all issues in the log.
func (l *Log) All() iter.Seq2[int, Issue] {
	return func(yield func(int, Issue) bool) {
		for i, u := range l.issues {
			if !yield(i, u) {
				return
			}
		}
	}
}

// Errors returns an iterator over all [Error] issues in the log.
func (l *Log) Errors() iter.Seq[Issue] {
	if l.errors == 0 {
		return func(func(Issue) bool) {}
	}
	return l.iter(Error)
}

// Warnings returns an iterator over all [Warning] issues in the log.
func (l *Log) Warnings() iter.Seq[Issue] {
	if l.warnings == 0 {
		return func(func(Issue) bool) {}
	}
	return l.iter(Warning)
}

// Infos returns an iterator over all [Info] issues in the log.
func (l *Log) Infos() iter.Seq[Issue] {
	if l.infos == 0 {
		return func(func(Issue) bool) {}
	}
	return l.iter(Info)
}

func (l *Log) iter(severity Severity) iter.Seq[Issue] {
	return func(yield func(Issue) bool) {
		for _, i := range l.issues {
			if i.Severity == severity && !yield(i) {
				return
			}
		}
	}
}
