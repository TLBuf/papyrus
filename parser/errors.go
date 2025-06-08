package parser

import (
	"fmt"

	"github.com/TLBuf/papyrus/source"
)

// Error defines an error raised by the parser.
type Error struct {
	// The underlying error.
	Err error
	// Location identifies the place in the source that caused the error.
	Location source.Location
}

// Error implments the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Location, e.Err)
}

// Unwrap returns the underlying error.
func (e Error) Unwrap() error {
	return e.Err
}

func newError(location source.Location, msg string, args ...any) Error {
	return Error{
		Err:      fmt.Errorf(msg, args...),
		Location: location,
	}
}
