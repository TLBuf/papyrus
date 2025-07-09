package types

import (
	"io"
	"strings"
)

// NewSignature returns a new Signature type with an
// optional return type and zero or more parameters.
func NewSignature(returnType Type, params ...Type) *Signature {
	return &Signature{
		params:     params,
		returnType: returnType,
	}
}

// Signature is a function type.
type Signature struct {
	params     []Type
	returnType Type
}

func (s *Signature) String() string {
	var sb strings.Builder
	s.write(&sb, "function")
	return sb.String()
}

func (s *Signature) write(w io.Writer, name string) {
	if s.returnType != nil {
		_, _ = io.WriteString(w, s.returnType.String())
		_, _ = w.Write([]byte{' '})
	}
	_, _ = io.WriteString(w, name)
	_, _ = w.Write([]byte{'('})
	for i, p := range s.params {
		if i > 0 {
			_, _ = w.Write([]byte{','})
		}
		_, _ = io.WriteString(w, p.String())
	}
	_, _ = w.Write([]byte{')'})
}

func (*Signature) types() {}

var _ Type = (*Signature)(nil)
