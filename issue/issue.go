// Package issue provides a common set of tools for describing problems encountered during processing.
package issue

import (
	"fmt"
	"iter"
	"regexp"
	"strconv"
	"strings"

	"github.com/TLBuf/papyrus/source"
)

// Severity describes how serious an issue detected by a processing step is.
type Severity uint8

const (
	// Internal indicates an issue due to a fault in the system rather than user
	// input. The user is not expected to fix these issue, rather report them to
	// the system owner.
	//
	// Internal issues almost always interrupt processing.
	Internal Severity = iota
	// Error indicates an issue that the user must address. This likely indicates
	// that the input is invalid in some fundamental way (e.g. bad syntax).
	//
	// Error issues from one processing phase usually prevent progression onto the
	// next phase.
	Error
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
	case Internal:
		return "Internal"
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
	definition *Definition
	detail     string
	file       *source.File
	location   source.Location
	related    []Related
}

// New returns a new issue at a specific location with a formatted message.
func New(def *Definition, file *source.File, loc source.Location) *Issue {
	return &Issue{
		definition: def,
		file:       file,
		location:   loc,
	}
}

// Definition returns the metadata that defines what kind of issue this is.
func (i *Issue) Definition() *Definition {
	return i.definition
}

// Detail returns supplemental information that
// is required to clarify the issue if any.
func (i *Issue) Detail() string {
	return i.detail
}

// File returns the source file where the issue was found.
func (i *Issue) File() *source.File {
	return i.file
}

// Location returns the source location where the issue was found.
func (i *Issue) Location() source.Location {
	return i.location
}

// Related returns an iterator over all related locations.
func (i *Issue) Related() iter.Seq2[int, Related] {
	return func(yield func(int, Related) bool) {
		for x, r := range i.related {
			if !yield(x, r) {
				return
			}
		}
	}
}

// WithDetail adds formatted detail information to this issue and returns it.
func (i *Issue) WithDetail(msg string, args ...any) *Issue {
	i.detail = fmt.Sprintf(msg, args...)
	return i
}

// AppendRelated adds a related source location to this issue.
func (i *Issue) AppendRelated(file *source.File, loc source.Location, msg string, args ...any) *Issue {
	i.related = append(i.related, Related{
		file:     file,
		location: loc,
		detail:   fmt.Sprintf(msg, args...),
	})
	return i
}

func (i *Issue) String() string {
	var sb strings.Builder
	_, _ = sb.WriteString(i.definition.String())
	_, _ = sb.WriteString(" - ")
	_, _ = sb.WriteString(i.file.Path())
	_, _ = sb.WriteRune(':')
	_, _ = sb.Write(strconv.AppendUint(nil, uint64(i.file.StartLine(i.location)), 10))
	_, _ = sb.WriteRune(':')
	_, _ = sb.Write(strconv.AppendUint(nil, uint64(i.file.StartColumn(i.location)), 10))
	if i.detail != "" {
		_, _ = sb.WriteString(" - ")
		_, _ = sb.WriteString(i.detail)
	}
	_, _ = sb.WriteRune('\n')
	for _, r := range i.related {
		_, _ = sb.WriteString("  ")
		_, _ = sb.WriteString(r.String())
		_, _ = sb.WriteRune('\n')
	}
	return sb.String()
}

var definitions map[string]*Definition

// Definition is a definition of a particular type of issue.
type Definition struct {
	// ID is the unqiue identifier of this type of issue.
	id string
	// Serverity describes how serious the issue is.
	severity Severity
	// Description is a human-readable description of the issue.
	description string
}

// NewInternal returns a new [Definition] for an [Internal] issue.
func NewInternal(id, description string) *Definition {
	return def(Internal, id, description)
}

// NewError returns a new [Definition] for an [Error] issue.
func NewError(id, description string) *Definition {
	return def(Error, id, description)
}

// NewWarning returns a new [Definition] for a [Warning] issue.
func NewWarning(id, description string) *Definition {
	return def(Warning, id, description)
}

// NewInfo returns a new [Definition] for an [Info] issue.
func NewInfo(id, description string) *Definition {
	return def(Info, id, description)
}

func def(severity Severity, id, description string) *Definition {
	if !idRegexp.MatchString(id) {
		panic(fmt.Sprintf("%q is not a valid ID, must match %q", id, idRegexp))
	}
	if description == "" {
		panic("description is empty")
	}
	if definitions == nil {
		definitions = make(map[string]*Definition)
	}
	existing, ok := definitions[id]
	if ok {
		panic(fmt.Sprintf("definition already exists with ID %q: %v", id, existing))
	}
	def := &Definition{
		id:          id,
		severity:    severity,
		description: description,
	}
	definitions[id] = def
	return def
}

var idRegexp = regexp.MustCompile(`^[A-Z]{4}[0-9]{4}$`)

// Description returns the standard human-readable description of the issue.
func (d Definition) Description() string {
	return d.description
}

// Severity returns the severity of the issue.
func (d Definition) Severity() Severity {
	return d.severity
}

// ID returns the unqiue identifier of this type of issue.
func (d Definition) ID() string {
	return d.id
}

func (d Definition) String() string {
	return fmt.Sprintf("[%s] %s: %s", d.id, d.severity, d.description)
}

// Related is a separate piece of source code that relates to some issue.
type Related struct {
	file     *source.File
	location source.Location
	detail   string
}

// Detail returns supplemental information that is required
// to clarify how this location relates to the overall issue.
func (r Related) Detail() string {
	return r.detail
}

// File returns the source file where the related piece of information is found.
func (r Related) File() *source.File {
	return r.file
}

// Location returns the source location of a related piece of information.
func (r Related) Location() source.Location {
	return r.location
}

func (r Related) String() string {
	if r.detail == "" {
		return fmt.Sprintf("%s:%d:%d", r.file.Path(), r.file.StartLine(r.location), r.file.StartColumn(r.location))
	}
	return fmt.Sprintf(
		"%s:%d:%d - %s",
		r.file.Path(),
		r.file.StartLine(r.location),
		r.file.StartColumn(r.location),
		r.detail,
	)
}
