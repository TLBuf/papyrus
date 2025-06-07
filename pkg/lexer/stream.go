package lexer

import (
	"iter"
	"slices"

	"github.com/TLBuf/papyrus/pkg/token"
)

// TokenStream is an immutable stream of [Token] values lexed from a source
// file.
//
// Tokens are guaranteed to have monotonically increasing byte offsets. The last
// Token always has kind [token.EOF].
type TokenStream struct {
	tokens []token.Token
}

// Len returns the number of tokens in the stream.
func (s TokenStream) Len() int {
	return len(s.tokens)
}

// Get returns the [Token] at a specific index or false
// if the index is out of range for the stream.
func (s TokenStream) Get(index int) (token.Token, bool) {
	if index < 0 || index >= len(s.tokens) {
		return token.Token{}, false
	}
	return s.tokens[index], true
}

// All returns an iterator over all index-token pairs in the stream.
func (s TokenStream) All() iter.Seq2[int, token.Token] {
	return slices.All(s.tokens)
}

// Backward returns an iterator over index-token pairs in the stream, traversing
// it backward with descending indices starting at the index before the one
// specified.
//
// An empty iterator is returned if the index is out of range for the stream.
func (s TokenStream) Backward(index int) iter.Seq2[int, token.Token] {
	if index <= 0 || index > len(s.tokens) {
		return func(yield func(int, token.Token) bool) {}
	}
	return slices.Backward(s.tokens[:index])
}

// Forward returns an iterator over index-token pairs in the stream, traversing
// it forward with increasing indicies starting at the index specified.
//
// An empty iterator is returned if the index is out of range for the stream.
func (s TokenStream) Forward(index int) iter.Seq2[int, token.Token] {
	if index < 0 || index >= len(s.tokens) {
		return func(yield func(int, token.Token) bool) {}
	}
	return slices.All(s.tokens[index:])
}

// IndexAtByte returns the index of the token with the given byte offset and
// true or if there is not an exact match, the index where such a token would
// exist in the stream and false is returned.
//
// If this returns false, the returned value is either the index of the token
// with the smallest byte offset greater than the requested one if it exists or
// the number of tokens in the stream.
func (s TokenStream) IndexAtByte(offset int) (int, bool) {
	return slices.BinarySearchFunc(s.tokens, offset, func(e token.Token, t int) int {
		return e.Location.ByteOffset - t
	})
}
