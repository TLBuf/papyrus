// Package lexer implements the Papyrus lexer.
package lexer

import (
	"fmt"
	"unicode/utf8"

	"github.com/TLBuf/papyrus/pkg/token"
)

// Error defines an error raised by the lexer.
type Error struct {
	msg string
	// OffendingText is the segment of input text that caused an error.
	OffendingText []byte
	// ByteOffset is the byte offset of the offending text in the input.
	ByteOffset int
}

// Error implments the error interface.
func (e Error) Error() string {
	return e.msg
}

type Lexer struct {
	text      []byte
	position  int
	next      int
	character rune
}

// New returns a [*Lexer] initialized for the given text.
func New(text []byte) *Lexer {
	l := &Lexer{text: text}
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
		tok = token.Token{Type: token.EOF, ByteOffset: l.position}
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
		l.readChar()
		if l.character == '\n' {
			tok = token.Token{Type: token.Newline, Text: l.text[l.position-1 : l.next], ByteOffset: l.position - 1}
		}
		text := l.text[l.position-1 : l.position]
		return token.Token{Type: token.Illegal, Text: text, ByteOffset: l.position - 1},
			Error{msg: "expected a newline after carriage return", OffendingText: text, ByteOffset: l.position - 1}
	case '\\':
		l.readChar()
		tok, err := l.NextToken()
		if err != nil {
			return tok, err
		}
		if tok.Type != token.Newline {
			return token.Token{Type: token.Illegal, Text: tok.Text, ByteOffset: tok.ByteOffset},
				Error{msg: "expected a newline immediately after '/'", OffendingText: tok.Text, ByteOffset: tok.ByteOffset}
		}
		return l.NextToken()
	case '=':
		l.readChar()
		if l.character == '=' {
			tok = token.Token{Type: token.Equal, Text: l.text[l.position-1 : l.next], ByteOffset: l.position - 1}
		}
		return token.Token{Type: token.Assign, Text: l.text[l.position-1 : l.position], ByteOffset: l.position - 1}, nil
	case '+':
		l.readChar()
		if l.character == '=' {
			tok = token.Token{Type: token.AssignAdd, Text: l.text[l.position-1 : l.next], ByteOffset: l.position - 1}
		}
		return token.Token{Type: token.Add, Text: l.text[l.position-1 : l.position], ByteOffset: l.position - 1}, nil
	case '-':
		l.readChar()
		if l.character == '=' {
			tok = token.Token{Type: token.AssignSubtract, Text: l.text[l.position-1 : l.next], ByteOffset: l.position - 1}
		}
		return token.Token{Type: token.Subtract, Text: l.text[l.position-1 : l.position], ByteOffset: l.position - 1}, nil
	case '*':
		l.readChar()
		if l.character == '=' {
			tok = token.Token{Type: token.AssignMultiply, Text: l.text[l.position-1 : l.next], ByteOffset: l.position - 1}
		}
		return token.Token{Type: token.Multiply, Text: l.text[l.position-1 : l.position], ByteOffset: l.position - 1}, nil
	case '/':
		l.readChar()
		if l.character == '=' {
			tok = token.Token{Type: token.AssignDivide, Text: l.text[l.position-1 : l.next], ByteOffset: l.position - 1}
		}
		return token.Token{Type: token.Divide, Text: l.text[l.position-1 : l.position], ByteOffset: l.position - 1}, nil
	case '%':
		l.readChar()
		if l.character == '=' {
			tok = token.Token{Type: token.AssignModulo, Text: l.text[l.position-1 : l.next], ByteOffset: l.position - 1}
		}
		return token.Token{Type: token.Modulo, Text: l.text[l.position-1 : l.position], ByteOffset: l.position - 1}, nil
	case '!':
		l.readChar()
		if l.character == '=' {
			tok = token.Token{Type: token.NotEqual, Text: l.text[l.position-1 : l.next], ByteOffset: l.position - 1}
		}
		return token.Token{Type: token.LogicalNot, Text: l.text[l.position-1 : l.position], ByteOffset: l.position - 1}, nil
	case '>':
		l.readChar()
		if l.character == '=' {
			tok = token.Token{Type: token.GreaterOrEqual, Text: l.text[l.position-1 : l.next], ByteOffset: l.position - 1}
		}
		return token.Token{Type: token.Greater, Text: l.text[l.position-1 : l.position], ByteOffset: l.position - 1}, nil
	case '<':
		l.readChar()
		if l.character == '=' {
			tok = token.Token{Type: token.LessOrEqual, Text: l.text[l.position-1 : l.next], ByteOffset: l.position - 1}
		}
		return token.Token{Type: token.Less, Text: l.text[l.position-1 : l.position], ByteOffset: l.position - 1}, nil
	case '|':
		l.readChar()
		if l.character == '|' {
			tok = token.Token{Type: token.LogicalOr, Text: l.text[l.position-1 : l.next], ByteOffset: l.position - 1}
		}
		text := l.text[l.position-1 : l.position]
		return token.Token{Type: token.Illegal, Text: text, ByteOffset: l.position - 1},
			Error{msg: "'|' is not a valid operator", OffendingText: text, ByteOffset: l.position - 1}
	case '&':
		l.readChar()
		if l.character == '&' {
			tok = token.Token{Type: token.LogicalOr, Text: l.text[l.position-1 : l.next], ByteOffset: l.position - 1}
		}
		text := l.text[l.position-1 : l.position]
		return token.Token{Type: token.Illegal, Text: text, ByteOffset: l.position - 1},
			Error{msg: "'&' is not a valid operator", OffendingText: text, ByteOffset: l.position - 1}
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
			return tok, Error{msg: "failed to lex any token", OffendingText: tok.Text, ByteOffset: tok.ByteOffset}
		}
	}
	l.readChar()
	return tok, nil
}

func (l *Lexer) newToken(t token.Type) token.Token {
	return token.Token{Type: t, Text: l.text[l.position:l.next], ByteOffset: l.position}
}

func (l *Lexer) readIdentifier() token.Token {
	start := l.position
	l.readChar()
	for isLetter(l.character) || isDigit(l.character) {
		l.readChar()
	}
	text := l.text[start:l.position]
	return token.Token{Type: token.LookupIdentifier(string(text)), Text: text, ByteOffset: start}
}

func (l *Lexer) readNumber() (token.Token, error) {
	start := l.position
	first := l.character
	l.readChar()
	if first == '0' && (l.character == 'x' || l.character == 'X') {
		// Hex Int
		l.readChar()
		for isHexDigit(l.character) {
			l.readChar()
		}
		text := l.text[start:l.position]
		if text[len(text)-1] == 'x' || text[len(text)-1] == 'X' {
			return token.Token{Type: token.Illegal, Text: text, ByteOffset: start},
				Error{msg: fmt.Sprintf("expected a digit to follow the %s in a hex int literal", string(text[len(text)-1])), OffendingText: text, ByteOffset: start}
		}
		return token.Token{Type: token.IntLiteral, Text: text, ByteOffset: start}, nil
	}
	isFloat := false
	for isDigit(l.character) || l.character == '.' {
		if l.character == '.' {
			isFloat = true
		}
		l.readChar()
	}
	text := l.text[start:l.position]
	if text[len(text)-1] == '.' {
		// Number ends with a dot?
		return token.Token{Type: token.Illegal, Text: text, ByteOffset: start},
			Error{msg: "expected a digit to follow the dot in a float literal", OffendingText: text, ByteOffset: start}
	}
	if isFloat {
		return token.Token{Type: token.FloatLiteral, Text: text, ByteOffset: start}, nil
	}
	return token.Token{Type: token.IntLiteral, Text: text, ByteOffset: start}, nil
}

func (l *Lexer) readString() (token.Token, error) {
	start := l.position
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
			text := l.text[start:l.position]
			return token.Token{Type: token.Illegal, Text: text, ByteOffset: start},
				Error{msg: fmt.Sprintf("encountered an invalid string escape sequence: \\%s", string(l.character)), OffendingText: text, ByteOffset: start}
		}
		if l.character == '"' {
			break
		}
	}
	text := l.text[start:l.position]
	if l.character == 0 {
		return token.Token{Type: token.Illegal, Text: text, ByteOffset: start},
			Error{msg: "reached end of file while reading string literal", OffendingText: text, ByteOffset: start}
	}
	l.readChar()
	return token.Token{Type: token.StringLiteral, Text: text, ByteOffset: start}, nil
}

func (l *Lexer) readComment() (token.Token, error) {
	start := l.position
	if l.character == '{' {
		// Doc comment
		for l.character != 0 && l.character != '}' {
			l.readChar()
		}
		text := l.text[start:l.position]
		if l.character == 0 {
			return token.Token{Type: token.Illegal, Text: text, ByteOffset: start},
				Error{msg: "reached end of file while reading doc comment", OffendingText: text, ByteOffset: start}
		}
		l.readChar()
		return token.Token{Type: token.DocComment, Text: l.text[start:l.position], ByteOffset: start}, nil
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
		text := l.text[start:l.position]
		if l.character == 0 {
			return token.Token{Type: token.Illegal, Text: text, ByteOffset: start},
				Error{msg: "reached end of file while reading block comment", OffendingText: text, ByteOffset: start}
		}
		l.readChar()
		return token.Token{Type: token.BlockComment, Text: l.text[start:l.position], ByteOffset: start}, nil
	}
	// Line comment
	for l.character != 0 && l.character != '\n' {
		l.readChar()
	}
	text := l.text[start:l.position]
	if l.character == 0 {
		return token.Token{Type: token.LineComment, Text: text, ByteOffset: start}, nil
	}
	return token.Token{Type: token.LineComment, Text: text, ByteOffset: start}, nil
}

func (l *Lexer) skipWhitespace() {
	for l.character == ' ' || l.character == '\t' {
		l.readChar()
	}
}

func (l *Lexer) readChar() error {
	width := 1
	if l.next >= len(l.text) {
		l.character = 0
	} else {
		r, w := utf8.DecodeRune(l.text[l.next:])
		if r == utf8.RuneError {
			return fmt.Errorf("encountered invalid UTF-8 at byte %d", l.next)
		}
		l.character = r
		width = w
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
