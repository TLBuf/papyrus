// Package lexer implements the Papyrus lexer.
package lexer

import (
	"fmt"
	"unicode/utf8"

	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/TLBuf/papyrus/pkg/token"
)

// Error defines an error raised by the lexer.
type Error struct {
	// A human-readable message describing what went wrong.
	Message string
	// Location is the source range of the segment of input text that caused an
	// error.
	Location source.Range
}

// Error implments the error interface.
func (e Error) Error() string {
	return e.Message
}

// Lexer provides the ability to lex a Papyrus script.
type Lexer struct {
	file      *source.File
	position  int
	next      int
	character rune
	column    int
	line      int
}

// New returns a [*Lexer] initialized for the given text.
func New(file *source.File) *Lexer {
	l := &Lexer{
		file:   file,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// NextToken scans the input for the next [token.Token].
//
// Returns an [Error] if the input could not be lexed as a token.
func (l *Lexer) NextToken() (token.Token, error) {
	var tok token.Token
	l.skipWhitespace()
	switch l.character {
	case 0:
		tok = token.Token{
			Type: token.EOF,
			SourceRange: source.Range{
				File:       l.file,
				ByteOffset: l.position,
				Length:     0,
				Line:       l.line,
				Column:     l.column,
			},
		}
	case '(':
		tok = l.newToken(token.LParen)
	case ')':
		tok = l.newToken(token.RParen)
	case '[':
		tok = l.newToken(token.LBracket)
	case ']':
		tok = l.newToken(token.RBracket)
	case ',':
		tok = l.newToken(token.Comma)
	case '.':
		tok = l.newToken(token.Dot)
	case '\n':
		tok = l.newToken(token.Newline)
	case '\r':
		column := l.column
		l.readChar()
		if l.character == '\n' {
			tok = l.newTokenWithRange(token.Newline, l.position-1, l.next-l.position+1, l.line, column)
		}
		errTok := l.newTokenWithRange(token.Illegal, l.position-1, 1, l.line, column)
		return errTok, Error{Message: "expected a newline after carriage return", Location: errTok.SourceRange}
	case '\\':
		column := l.column
		l.readChar()
		tok, err := l.NextToken()
		if err != nil {
			return tok, err
		}
		if tok.Type != token.Newline {
			errTok := l.newTokenWithRange(token.Illegal, tok.SourceRange.ByteOffset, 1, l.line, column)
			return errTok, Error{Message: "expected a newline immediately after '/'", Location: errTok.SourceRange}
		}
		return l.NextToken()
	case '=':
		column := l.column
		l.readChar()
		if l.character == '=' {
			tok = l.newTokenWithRange(token.Equal, l.position-1, 2, l.line, column)
		}
		return l.newTokenWithRange(token.Assign, l.position-1, 1, l.line, column), nil
	case '+':
		column := l.column
		l.readChar()
		if l.character == '=' {
			tok = l.newTokenWithRange(token.AssignAdd, l.position-1, 2, l.line, column)
		}
		return l.newTokenWithRange(token.Add, l.position-1, 1, l.line, column), nil
	case '-':
		column := l.column
		l.readChar()
		if l.character == '=' {
			tok = l.newTokenWithRange(token.AssignSubtract, l.position-1, 2, l.line, column)
		}
		return l.newTokenWithRange(token.Subtract, l.position-1, 1, l.line, column), nil
	case '*':
		column := l.column
		l.readChar()
		if l.character == '=' {
			tok = l.newTokenWithRange(token.AssignMultiply, l.position-1, 2, l.line, column)
		}
		return l.newTokenWithRange(token.Multiply, l.position-1, 1, l.line, column), nil
	case '/':
		column := l.column
		l.readChar()
		if l.character == '=' {
			tok = l.newTokenWithRange(token.AssignDivide, l.position-1, 2, l.line, column)
		}
		return l.newTokenWithRange(token.Divide, l.position-1, 1, l.line, column), nil
	case '%':
		column := l.column
		l.readChar()
		if l.character == '=' {
			tok = l.newTokenWithRange(token.AssignModulo, l.position-1, 2, l.line, column)
		}
		return l.newTokenWithRange(token.Modulo, l.position-1, 1, l.line, column), nil
	case '!':
		column := l.column
		l.readChar()
		if l.character == '=' {
			tok = l.newTokenWithRange(token.NotEqual, l.position-1, 2, l.line, column)
		}
		return l.newTokenWithRange(token.LogicalNot, l.position-1, 1, l.line, column), nil
	case '>':
		column := l.column
		l.readChar()
		if l.character == '=' {
			tok = l.newTokenWithRange(token.GreaterOrEqual, l.position-1, 2, l.line, column)
		}
		return l.newTokenWithRange(token.Greater, l.position-1, 1, l.line, column), nil
	case '<':
		column := l.column
		l.readChar()
		if l.character == '=' {
			tok = l.newTokenWithRange(token.LessOrEqual, l.position-1, 2, l.line, column)
		}
		return l.newTokenWithRange(token.Less, l.position-1, 1, l.line, column), nil
	case '|':
		column := l.column
		l.readChar()
		if l.character == '|' {
			tok = l.newTokenWithRange(token.LogicalOr, l.position-1, 2, l.line, column)
		}
		errTok := l.newTokenWithRange(token.Illegal, l.position-1, 1, l.line, column)
		return errTok, Error{Message: "'|' is not a valid operator", Location: errTok.SourceRange}
	case '&':
		column := l.column
		l.readChar()
		if l.character == '&' {
			tok = l.newTokenWithRange(token.LogicalAnd, l.position-1, 2, l.line, column)
		}
		errTok := l.newTokenWithRange(token.Illegal, l.position-1, 1, l.line, column)
		return errTok, Error{Message: "'&' is not a valid operator", Location: errTok.SourceRange}
	case '{', ';':
		return l.readComment()
	case '"':
		return l.readString()
	default:
		if isLetter(l.character) {
			return l.readIdentifier(), nil
		} else if isDigit(l.character) {
			return l.readNumber()
		} else {
			tok = l.newToken(token.Illegal)
			l.readChar()
			return tok, Error{Message: "failed to lex any token", Location: tok.SourceRange}
		}
	}
	l.readChar()
	return tok, nil
}

func (l *Lexer) newToken(t token.Type) token.Token {
	return token.Token{
		Type: t,
		SourceRange: source.Range{
			File:       l.file,
			ByteOffset: l.position,
			Length:     l.next - l.position,
			Line:       l.line,
			Column:     l.column,
		},
	}
}

func (l *Lexer) newTokenWithRange(t token.Type, byteOffset, length, line, column int) token.Token {
	return token.Token{
		Type: t,
		SourceRange: source.Range{
			File:       l.file,
			ByteOffset: byteOffset,
			Length:     length,
			Line:       line,
			Column:     column,
		},
	}
}

func (l *Lexer) readIdentifier() token.Token {
	start := l.position
	column := l.column
	l.readChar()
	for isLetter(l.character) || isDigit(l.character) {
		l.readChar()
	}
	text := l.file.Text[start:l.position]
	return l.newTokenWithRange(token.LookupIdentifier(string(text)), start, l.position-start, l.line, column)
}

func (l *Lexer) readNumber() (token.Token, error) {
	start := l.position
	first := l.character
	column := l.column
	l.readChar()
	if first == '0' && (l.character == 'x' || l.character == 'X') {
		// Hex Int
		l.readChar()
		for isHexDigit(l.character) {
			l.readChar()
		}
		tok := l.newTokenWithRange(token.IntLiteral, start, l.position-start, l.line, column)
		if l.file.Text[l.position-1] == 'x' || l.file.Text[l.position-1] == 'X' {
			tok.Type = token.Illegal
			return tok, Error{Message: fmt.Sprintf("expected a digit to follow the %s in a hex int literal", string(l.file.Text[l.position-1])), Location: tok.SourceRange}
		}
		return tok, nil
	}
	isFloat := false
	for isDigit(l.character) || l.character == '.' {
		if l.character == '.' {
			isFloat = true
		}
		l.readChar()
	}
	tok := l.newTokenWithRange(token.IntLiteral, start, l.position-start, l.line, column)
	if l.file.Text[l.position-1] == '.' {
		// Number ends with a dot?
		tok.Type = token.Illegal
		return tok, Error{Message: "expected a digit to follow the dot in a float literal", Location: tok.SourceRange}
	}
	if isFloat {
		tok.Type = token.FloatLiteral
	}
	return tok, nil
}

func (l *Lexer) readString() (token.Token, error) {
	start := l.position
	column := l.column
	l.readChar()
	escaping := false
	for {
		l.readChar()
		if l.character == 0 {
			break
		}
		if l.character == '\\' {
			escaping = true
			continue
		}
		if escaping {
			if l.character == 'n' || l.character == 't' || l.character == '"' || l.character == '\\' {
				escaping = false
				continue
			}
			tok := l.newTokenWithRange(token.Illegal, start, l.position-start, l.line, column)
			return tok, Error{Message: fmt.Sprintf("encountered an invalid string escape sequence: \\%s", string(l.character)), Location: tok.SourceRange}
		}
		if l.character == '"' {
			break
		}
	}
	tok := l.newTokenWithRange(token.StringLiteral, start, l.position-start, l.line, column)
	if l.character == 0 {
		tok.Type = token.Illegal
		return tok, Error{Message: "reached end of file while reading string literal", Location: tok.SourceRange}
	}
	l.readChar()
	return tok, nil
}

func (l *Lexer) readComment() (token.Token, error) {
	tok := token.Token{
		Type: token.Illegal,
		SourceRange: source.Range{
			File:       l.file,
			ByteOffset: l.position,
			Line:       l.line,
			Column:     l.column,
		},
	}
	if l.character == '{' {
		// Doc comment
		for l.character != 0 && l.character != '}' {
			l.readChar()
		}

		if l.character == 0 {
			tok.SourceRange.Length = l.position - tok.SourceRange.ByteOffset
			return tok, Error{Message: "reached end of file while reading doc comment", Location: tok.SourceRange}
		}
		l.readChar()
		tok.Type = token.DocComment
		tok.SourceRange.Length = l.position - tok.SourceRange.ByteOffset
		return tok, nil
	}
	l.readChar()
	if l.character == '/' {
		// Block comment
		l.readChar()
		for {
			if l.character == 0 {
				break
			}
			if l.character == '/' {
				l.readChar()
				if l.character == ';' {
					break
				}
			}
			l.readChar()
		}

		if l.character == 0 {
			tok.SourceRange.Length = l.position - tok.SourceRange.ByteOffset
			return tok, Error{Message: "reached end of file while reading block comment", Location: tok.SourceRange}
		}
		l.readChar()
		tok.Type = token.BlockComment
		tok.SourceRange.Length = l.position - tok.SourceRange.ByteOffset
		return tok, nil
	}
	// Line comment
	for l.character != 0 && l.character != '\n' {
		l.readChar()
	}
	tok.Type = token.LineComment
	tok.SourceRange.Length = l.position - tok.SourceRange.ByteOffset
	return tok, nil
}

func (l *Lexer) skipWhitespace() {
	for l.character == ' ' || l.character == '\t' {
		l.readChar()
	}
}

func (l *Lexer) readChar() error {
	width := 1
	if l.character == '\n' {
		l.line++
		l.column = 0
	}
	if l.next >= len(l.file.Text) {
		l.character = 0
		l.column = 1
	} else {
		r, w := utf8.DecodeRune(l.file.Text[l.next:])
		if r == utf8.RuneError {
			return fmt.Errorf("encountered invalid UTF-8 at byte %d", l.next)
		}
		l.character = r
		width = w
		l.column++
	}
	l.position = l.next
	l.next += width
	return nil
}

func isLetter(char rune) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9'
}

func isHexDigit(char rune) bool {
	return '0' <= char && char <= '9' || 'a' <= char && char <= 'f' || 'A' <= char && char <= 'F'
}
