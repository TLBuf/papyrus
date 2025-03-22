// Package lexer implements the Papyrus lexer.
package lexer

import (
	"fmt"
	"unicode/utf8"

	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/TLBuf/papyrus/pkg/token"
)

// Lexer provides the ability to lex a Papyrus script.
type Lexer struct {
	file            *source.File
	position        int
	next            int
	character       rune
	column          int
	line            int
	lineStartOffset int
	lineEndOffset   int
}

// New returns a [*Lexer] initialized for the given text.
func New(file *source.File) (*Lexer, error) {
	l := &Lexer{
		file:            file,
		line:            1,
		column:          0,
		lineStartOffset: 0,
		lineEndOffset:   -1,
	}
	if err := l.readChar(); err != nil {
		return nil, fmt.Errorf("failed to initialize lexer: %v", err)
	}
	return l, nil
}

// NextToken scans the input for the next [token.Token].
//
// Returns an [Error] if the input could not be lexed as a token.
func (l *Lexer) NextToken() (token.Token, error) {
	tok, err := l.nextToken()
	if err != nil || tok.Type == token.Newline || tok.Type == token.EOF {
		return tok, err
	}
	tok.Location.PostambleLength = l.lineEndOffset - tok.Location.ByteOffset - tok.Location.Length
	return tok, err
}

func (l *Lexer) nextToken() (token.Token, error) {
	var tok token.Token
	if err := l.skipWhitespace(); err != nil {
		return l.newToken(token.Illegal), err
	}
	switch l.character {
	case 0:
		tok = token.Token{
			Type:     token.EOF,
			Location: l.here(),
		}
		tok.Location.Length = 0
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
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '\n' {
			tok = l.newTokenFrom(token.Newline, start)
		} else {
			return l.newTokenAt(token.Illegal, start), newError(tok.Location, "expected a newline after carriage return")
		}
	case '\\':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		tok, err := l.NextToken()
		if err != nil {
			return tok, err
		}
		if tok.Type != token.Newline {
			return l.newTokenAt(token.Illegal, start), newError(tok.Location, "expected a newline immediately after '/'")
		}
		return l.NextToken()
	case '=':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '=' {
			tok = l.newTokenFrom(token.Equal, start)
		}
		return l.newTokenAt(token.Assign, start), nil
	case '+':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '=' {
			tok = l.newTokenFrom(token.AssignAdd, start)
		}
		return l.newTokenAt(token.Add, start), nil
	case '-':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '=' {
			tok = l.newTokenFrom(token.AssignSubtract, start)
		}
		return l.newTokenAt(token.Subtract, start), nil
	case '*':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '=' {
			tok = l.newTokenFrom(token.AssignMultiply, start)
		}
		return l.newTokenAt(token.Multiply, start), nil
	case '/':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '=' {
			tok = l.newTokenFrom(token.AssignDivide, start)
		}
		return l.newTokenAt(token.Divide, start), nil
	case '%':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '=' {
			tok = l.newTokenFrom(token.AssignModulo, start)
		}
		return l.newTokenAt(token.Modulo, start), nil
	case '!':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '=' {
			tok = l.newTokenFrom(token.NotEqual, start)
		}
		return l.newTokenAt(token.LogicalNot, start), nil
	case '>':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '=' {
			tok = l.newTokenFrom(token.GreaterOrEqual, start)
		}
		return l.newTokenAt(token.Greater, start), nil
	case '<':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '=' {
			tok = l.newTokenFrom(token.LessOrEqual, start)
		}
		return l.newTokenAt(token.Less, start), nil
	case '|':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '|' {
			tok = l.newTokenFrom(token.LogicalOr, start)
		}
		return l.newTokenAt(token.Illegal, start), newError(tok.Location, "'|' is not a valid operator")
	case '&':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '&' {
			tok = l.newTokenFrom(token.LogicalAnd, start)
		}
		return l.newTokenAt(token.Illegal, start), newError(tok.Location, "'&' is not a valid operator")
	case '{', ';':
		return l.readComment()
	case '"':
		return l.readString()
	default:
		if isLetter(l.character) {
			return l.readIdentifier()
		} else if isDigit(l.character) {
			return l.readNumber()
		} else {
			tok = l.newToken(token.Illegal)
			if err := l.readChar(); err != nil {
				return l.newToken(token.Illegal), err
			}
			return tok, newError(tok.Location, "failed to lex any token")
		}
	}
	if err := l.readChar(); err != nil {
		return l.newToken(token.Illegal), err
	}
	return tok, nil
}

func (l *Lexer) newToken(t token.Type) token.Token {
	return token.Token{
		Type:     t,
		Location: l.here(),
	}
}

func (l *Lexer) newTokenAt(t token.Type, at source.Range) token.Token {
	return token.Token{
		Type:     t,
		Location: at,
	}
}

func (l *Lexer) newTokenFrom(t token.Type, from source.Range) token.Token {
	return token.Token{
		Type:     t,
		Location: source.Span(from, l.here()),
	}
}

func (l *Lexer) here() source.Range {
	return source.Range{
		File:           l.file,
		ByteOffset:     l.position,
		Length:         l.next - l.position,
		StartLine:      l.line,
		StartColumn:    l.column,
		EndLine:        l.line,
		EndColumn:      l.column,
		PreambleLength: l.position - l.lineStartOffset,
	}
}

func (l *Lexer) readIdentifier() (token.Token, error) {
	start := l.here()
	if err := l.readChar(); err != nil {
		return l.newToken(token.Illegal), err
	}
	for isLetter(l.character) || isDigit(l.character) {
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
	}
	loc := source.Span(start, l.here())
	loc.Length -= 1
	loc.EndColumn -= 1
	return l.newTokenAt(token.LookupIdentifier(string(loc.Text())), loc), nil
}

func (l *Lexer) readNumber() (token.Token, error) {
	start := l.here()
	first := l.character
	if err := l.readChar(); err != nil {
		return l.newToken(token.Illegal), err
	}
	end := start
	if first == '0' && (l.character == 'x' || l.character == 'X') {
		// Hex Int
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		for isHexDigit(l.character) {
			end = l.here()
			if err := l.readChar(); err != nil {
				return l.newToken(token.Illegal), err
			}
		}
		tok := l.newTokenAt(token.IntLiteral, source.Span(start, end))
		if l.file.Text[l.position-1] == 'x' || l.file.Text[l.position-1] == 'X' {
			tok.Type = token.Illegal
			return tok, newError(tok.Location, "expected a digit to follow the %s in a hex int literal", string(l.file.Text[l.position-1]))
		}
		return tok, nil
	}
	isFloat := false
	for isDigit(l.character) || l.character == '.' {
		if l.character == '.' {
			isFloat = true
		}
		end = l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
	}
	tok := l.newTokenAt(token.IntLiteral, source.Span(start, end))
	if l.file.Text[l.position-1] == '.' {
		// Number ends with a dot?
		tok.Type = token.Illegal
		return tok, newError(tok.Location, "expected a digit to follow the dot in a float literal")
	}
	if isFloat {
		tok.Type = token.FloatLiteral
	}
	return tok, nil
}

func (l *Lexer) readString() (token.Token, error) {
	start := l.here()
	if err := l.readChar(); err != nil {
		return l.newToken(token.Illegal), err
	}
	escaping := false
	for {
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
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
			tok := l.newTokenFrom(token.Illegal, start)
			return tok, newError(tok.Location, "encountered an invalid string escape sequence: \\%s", string(l.character))
		}
		if l.character == '"' {
			break
		}
	}
	tok := l.newTokenFrom(token.StringLiteral, start)
	if l.character == 0 {
		tok.Type = token.Illegal
		return tok, newError(tok.Location, "reached end of file while reading string literal")
	}
	if err := l.readChar(); err != nil {
		return l.newToken(token.Illegal), err
	}
	return tok, nil
}

func (l *Lexer) readComment() (token.Token, error) {
	start := l.here()
	if l.character == '{' {
		// Doc comment
		for l.character != 0 && l.character != '}' {
			if err := l.readChar(); err != nil {
				return l.newToken(token.Illegal), err
			}
		}
		if l.character == 0 {
			tok := l.newTokenFrom(token.Illegal, start)
			return tok, newError(tok.Location, "reached end of file while reading doc comment")
		}
		tok := l.newTokenFrom(token.DocComment, start)
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		return tok, nil
	}
	if err := l.readChar(); err != nil {
		return l.newToken(token.Illegal), err
	}
	if l.character == '/' {
		// Block comment
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		for {
			if l.character == 0 {
				break
			}
			if l.character == '/' {
				if err := l.readChar(); err != nil {
					return l.newToken(token.Illegal), err
				}
				if l.character == ';' {
					break
				}
			}
			if err := l.readChar(); err != nil {
				return l.newToken(token.Illegal), err
			}
		}
		if l.character == 0 {
			tok := l.newTokenFrom(token.Illegal, start)
			return tok, newError(tok.Location, "reached end of file while reading block comment")
		}
		tok := l.newTokenFrom(token.BlockComment, start)
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		return tok, nil
	}
	// Line comment
	var end source.Range
	for l.character != 0 && l.character != '\n' && l.character != '\r' {
		end = l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
	}
	return l.newTokenAt(token.LineComment, source.Span(start, end)), nil
}

func (l *Lexer) skipWhitespace() error {
	for l.character == ' ' || l.character == '\t' {
		if err := l.readChar(); err != nil {
			return err
		}
	}
	return nil
}

func (l *Lexer) readChar() error {
	width := 1
	if l.lineEndOffset < 0 {
		l.lineEndOffset = l.findNextNewlineOffset()
	}
	if l.character == '\n' {
		l.line++
		l.column = 0
		l.lineStartOffset = l.next
		l.lineEndOffset = l.findNextNewlineOffset()
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

func (l *Lexer) findNextNewlineOffset() int {
	for i := l.next; i < len(l.file.Text); i++ {
		b := l.file.Text[i]
		if b == '\n' || b == '\r' || b == 0 {
			return i
		}
	}
	return len(l.file.Text)
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
