package parser

import (
	"fmt"

	"github.com/TLBuf/papyrus/pkg/source"
)

// Error defines an error raised by the parser.
type Error struct {
	// A human-readable message describing what went wrong.
	Message string
	// SourceRange is the source range of the segment of input text that caused an
	// error.
	Location source.Range
}

// Error implments the error interface.
func (e Error) Error() string {
	return e.Message
}

func newError(location source.Range, msg string, args ...any) Error {
	return Error{
		Message:  fmt.Sprintf(msg, args...),
		Location: location,
	}
}
