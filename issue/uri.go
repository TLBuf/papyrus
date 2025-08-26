package issue

import (
	"fmt"
	"regexp"
	"strings"
)

const scheme = "issue"

// URI is an RFC 3986 Uniform Resource Identifier for an issue with additional
// restrictions applied to simplify it.
//
// Issue URIs are expected to follow the following format:
//   - issue://<path element>/<path element>
//   - issue://<path element>/<path element>#<specifier>
//
// Where a <path element> and <specifier> are strings starting with a lower
// (latin) alphabet character, followed by one or more strings of lower
// alphanumeric strings with an optional leading underscore, e.g. "a", "a_b",
// "a_b_c". Every URI must have at least two path elements, each seperated by a
// '/' character and optionally a trailing specifier seprated by a '#'
// character.
type URI string

// NewURI returns a new URI with two or more path elements. Panics if any path
// element is not valid.
func NewURI(primary, secondary string, path ...string) URI {
	if !elementRegexp.MatchString(primary) {
		panic(fmt.Sprintf("path element %q is not valid", primary))
	}
	if !elementRegexp.MatchString(secondary) {
		panic(fmt.Sprintf("path element %q is not valid", secondary))
	}
	for _, p := range path {
		if !elementRegexp.MatchString(p) {
			panic(fmt.Sprintf("path element %q is not valid", p))
		}
	}
	base := fmt.Sprintf("%s://%s/%s", scheme, primary, secondary)
	if len(path) == 0 {
		return URI(base)
	}
	return URI(fmt.Sprintf("%s/%s", base, strings.Join(path, "/")))
}

// AppendSpecifier appends a specifier to a URI and returns the resulting URI.
// Panics if the URI already has a specifier or the specifier is invalid.
func (u URI) AppendSpecifier(specifier string) URI {
	if !u.IsValid() {
		panic(fmt.Sprintf("URI %q is not valid", u))
	}
	if strings.ContainsRune(string(u), '#') {
		panic(fmt.Sprintf("URI %q already has a specifier: %q", u, specifier))
	}
	if !elementRegexp.MatchString(specifier) {
		panic(fmt.Sprintf("specifier %q is not valid", specifier))
	}
	return URI(fmt.Sprintf("%s#%s", u, specifier))
}

// AppendPath appends some number of path elements to a URI and returns the
// resulting URI. Panics if the URI has a specifier (meaning it cannot have more
// path elements appended) or any element is invalid.
func (u URI) AppendPath(path ...string) URI {
	if !u.IsValid() {
		panic(fmt.Sprintf("URI %q is not valid", u))
	}
	if len(path) == 0 {
		return u
	}
	if strings.ContainsRune(string(u), '#') {
		panic(fmt.Sprintf("URI %q has a specifier: %q", u, path))
	}
	for _, p := range path {
		if !elementRegexp.MatchString(p) {
			panic(fmt.Sprintf("path element %q is not valid", p))
		}
	}
	return URI(fmt.Sprintf("%s/%s", u, strings.Join(path, "/")))
}

// IsValid returns true if this URI is valid.
func (u URI) IsValid() bool {
	return uriRegexp.MatchString(string(u))
}

func (u URI) String() string {
	return string(u)
}

var (
	elementRegexp = regexp.MustCompile(`[a-z](?:_?[a-z0-9])+`)
	uriRegexp     = regexp.MustCompile(`issue://[a-z](?:_?[a-z0-9])+(?:/[a-z](?:_?[a-z0-9])+)+(?:#[a-z](?:_?[a-z0-9])+)`)
)
