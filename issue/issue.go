// Package issue provides a common set of tools for describing problems encountered during processing.
package issue

import (
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
	// Informational indicates an issue that the user may address. This may be
	// suggestions to improve style, efficiency, etc.
	//
	// Informational issues should never prevent further processing.
	Informational
)

// Issue describes an issue found while processing input.
//
// Issues never represent internal errors encountered in proccessing, those are
// conveyed via normal [error] returns.
type Issue struct {
	// URI is the URI for this issue.
	URI URI
	// Serverity describes how serious the issue is.
	Severity Severity
	// File is the source file where the issue was found.
	File *source.File
	// Location is the source location of the issue.
	Location source.Location
	// Message is a human-readable message describing the issue.
	Message string
	// Related zero or more additional locations with associated messages.
	Related []Related
}

// Related is a seperate piece of source code that relates to some issue.
type Related struct {
	// Location is the source location of a related piece of information.
	Location source.Location
	// Messate is a human-readable description of how this location relates.
	Message string
}
