// Package issue provides a common set of tools for describing problems encountered during processing.
package issue

import (
	"fmt"
	"strings"

	"github.com/TLBuf/papyrus/source"
)

// Severity describes how serious an issue detected by a processing step is.
type Severity uint8

const (
	// Error indicates an issue that the user must address. This likely indicates
	// that the input is invalid in some fundamental way (e.g. bad syntax).
	//
	// Error issues from one processing phase usually prevent progression onto the
	// next phase.
	Error Severity = iota
	// Warning indicates an issue that the user should address. While technically
	// valid, the input may not function as the user intended (e.g. unused
	// variables).
	//
	// Warning issues should never prevent further processing.
	Warning
	// Info indicates an issue that the user may address. This may be suggestions
	// to improve style, efficiency, etc.
	//
	// Info issues should never prevent further processing.
	Info
)

func (s Severity) String() string {
	switch s {
	case Error:
		return "Error"
	case Warning:
		return "Warning"
	case Info:
		return "Info"
	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}

// Issue describes an issue found while processing input.
//
// Issues never represent internal errors encountered in processing, those are
// conveyed via normal [error] returns.
type Issue struct {
	// Key identifies the overall kind of issue.
	Key Key
	// File is the source file where the issue was found.
	File *source.File
	// Location is the source location of the issue.
	Location source.Location
	// Message is a human-readable message describing the issue.
	Message string
	// Related zero or more additional locations with associated messages.
	Related []Related
}

// NewIssue returns a new issue at a specific location with a formatted message.
func NewIssue(key Key, file *source.File, loc source.Location, msg string, args ...any) *Issue {
	return &Issue{
		Key:      key,
		File:     file,
		Location: loc,
		Message:  fmt.Sprintf(msg, args...),
	}
}

// AppendRelated adds a related source location to this issue.
func (i *Issue) AppendRelated(file *source.File, loc *source.Location, msg string, args ...any) {
	i.Related = append(i.Related, Related{
		File:     file,
		Location: *loc,
		Message:  fmt.Sprintf(msg, args...),
	})
}

func (i *Issue) String() string {
	base := fmt.Sprintf(
		"%s %s:%d:%d - %s\n",
		i.Key,
		i.File.Path(),
		i.File.StartLine(i.Location),
		i.File.StartColumn(i.Location),
		i.Message,
	)
	if len(i.Related) == 0 {
		return base
	}
	var sb strings.Builder
	_, _ = sb.WriteString(base)
	for _, r := range i.Related {
		_, _ = sb.WriteString(r.String())
		_, _ = sb.WriteRune('\n')
	}
	return sb.String()
}

// Key is a unique identifier of a particular type of issue.
type Key struct {
	// Serverity describes how serious the issue is.
	Severity Severity
	// Namespace defines the namespace of the issue, usually a system identifier.
	Namespace string
	// ID is the unqiue identifier of the issue within the namespace.
	ID string
}

func (k Key) String() string {
	return fmt.Sprintf("%s [%s:%s]", k.Severity, k.Namespace, k.ID)
}

// Related is a separate piece of source code that relates to some issue.
type Related struct {
	// File is the source file where the related piece of information is found.
	File *source.File
	// Location is the source location of a related piece of information.
	Location source.Location
	// Messate is a human-readable description of how this location relates.
	Message string
}

func (r Related) String() string {
	return fmt.Sprintf(
		"%s:%d:%d - %s",
		r.File.Path(),
		r.File.StartLine(r.Location),
		r.File.StartColumn(r.Location),
		r.Message,
	)
}
