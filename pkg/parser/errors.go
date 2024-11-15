package parser

import (
	"fmt"

	"github.com/TLBuf/papyrus/pkg/source"
)

// Error defines an error raised by the parser.
type Error struct {
	// Issues contains one or more issues that were encountered while parsing.
	//
	// Issues are in the order encountered and should be handled in that order. In
	// other words, since an earlier issue may cause more issues (because error
	// recovery isn't perfect), fixing the first issue may also resolve multiple
	// other issues.
	Issues []*Issue
}

// Error implments the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("encountered %d issue(s) while parsing", len(e.Issues))
}

// Issue describes a single issue that was discovered while parser.
type Issue struct {
	// A human-readable message describing what went wrong.
	Message string
	// SourceRange is the source range of the segment of input text that caused an
	// error.
	Location source.Range
}

type status int

const (
	// statusOK indicates the operation completed successfully.
	statusOK status = iota
	// statusError indicates the operation encountered an issue, but error
	// recovery should be attempted.
	statusError
	// statusFatal indicates the operation encountered an issue and no error
	// recovery should be attempted.
	statusFatal
)

func (p *parser) issue(location source.Range, msg string, args ...any) {
	p.issues = append(p.issues, &Issue{
		Message:  fmt.Sprintf(msg, args...),
		Location: location,
	})
}