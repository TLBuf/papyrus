// Package lexer implements the Papyrus lexer.
package lexer

import (
	"fmt"
	"unicode/utf8"

	"github.com/TLBuf/papyrus/pkg/source"
	"github.com/TLBuf/papyrus/pkg/token"
)

// mode records what state the lexer is currently in (i.e. is it processing a
// block comment or just lexing normally).
type mode int

const (
	// Normal token lexing.
	normal mode = iota
	// Line comment lexing.
	commentLine
	// Block comment lexing.
	commentBlock
	// Doc comment lexing.
	commentDoc
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
	mode            mode
	terminal        token.Token
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
	var tok token.Token
	var err error
	switch l.mode {
	case normal:
		tok, err = l.nextToken()
	case commentLine:
		tok, err = l.commentLine()
	case commentBlock:
		tok, err = l.commentBlock()
	case commentDoc:
		tok, err = l.commentDoc()
	default:
		err = fmt.Errorf("Lexer in unknown lexing mode: %d", l.mode)
	}
	if err != nil || tok.Kind == token.Newline || tok.Kind == token.EOF {
		return tok, err
	}
	if l.terminal.Kind != token.Newline {
		tok.Location.PostambleLength = l.lineEndOffset - tok.Location.ByteOffset - tok.Location.Length
	}
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
			Kind:     token.EOF,
			Location: l.here(),
		}
		tok.Location.Length = 0
	case '(':
		tok = l.newToken(token.ParenthesisOpen)
	case ')':
		tok = l.newToken(token.ParenthesisClose)
	case '[':
		tok = l.newToken(token.BracketOpen)
	case ']':
		tok = l.newToken(token.BracketClose)
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
		if tok.Kind != token.Newline {
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
		return l.newTokenAt(token.Plus, start), nil
	case '-':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '=' {
			tok = l.newTokenFrom(token.AssignSubtract, start)
		}
		return l.newTokenAt(token.Minus, start), nil
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
			break
		}
		return l.newTokenAt(token.Illegal, start), newError(tok.Location, "'|' is not a valid operator")
	case '&':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '&' {
			tok = l.newTokenFrom(token.LogicalAnd, start)
			break
		}
		return l.newTokenAt(token.Illegal, start), newError(tok.Location, "'&' is not a valid operator")
	case '{':
		tok = l.newToken(token.BraceOpen)
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		l.mode = commentDoc
		return tok, nil
	case ';':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character == '/' {
			l.mode = commentBlock
			tok = l.newTokenFrom(token.BlockCommentOpen, start)
			break
		}
		l.mode = commentLine
		return l.newTokenAt(token.Semicolon, start), nil
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

func (l *Lexer) newToken(t token.Kind) token.Token {
	return token.Token{
		Kind:     t,
		Location: l.here(),
	}
}

func (l *Lexer) newTokenAt(t token.Kind, at source.Location) token.Token {
	return token.Token{
		Kind:     t,
		Location: at,
	}
}

func (l *Lexer) newTokenFrom(t token.Kind, from source.Location) token.Token {
	return token.Token{
		Kind:     t,
		Location: source.Span(from, l.here()),
	}
}

func (l *Lexer) here() source.Location {
	return source.Location{
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
			tok.Kind = token.Illegal
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
		tok.Kind = token.Illegal
		return tok, newError(tok.Location, "expected a digit to follow the dot in a float literal")
	}
	if isFloat {
		tok.Kind = token.FloatLiteral
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
		tok.Kind = token.Illegal
		return tok, newError(tok.Location, "reached end of file while reading string literal")
	}
	if err := l.readChar(); err != nil {
		return l.newToken(token.Illegal), err
	}
	return tok, nil
}

// commentLine handles lexing in the commentLine mode.
func (l *Lexer) commentLine() (token.Token, error) {
	if l.terminal.Kind == token.Newline || l.terminal.Kind == token.EOF {
		// Already hit the terminal token, clean up, and return it.
		terminal := l.terminal
		l.terminal = token.Token{}
		l.mode = normal
		return terminal, nil
	}
	start := l.here()
	var end source.Location
	for {
		terminal := l.here()
		if l.character != 0 && l.character != '\r' && l.character != '\n' {
			end = terminal
			if err := l.readChar(); err != nil {
				return l.newToken(token.Illegal), err
			}
			continue
		}
		if l.character == '\r' {
			// Maybe the end?
			if err := l.readChar(); err != nil {
				return l.newToken(token.Illegal), err
			}
		}
		if l.character == '\n' || l.character == 0 {
			// We've read as much conent as there is and hit the close token.
			comment := l.newTokenAt(token.Comment, source.Span(start, end))
			l.terminal = l.newTokenFrom(token.Newline, terminal)
			if l.character == 0 {
				l.terminal.Kind = token.EOF
				return comment, nil
			}
			if err := l.readChar(); err != nil {
				return l.newToken(token.Illegal), err
			}
			return comment, nil
		}
	}
}

// commentBlock handles lexing in the commentBlock mode.
func (l *Lexer) commentBlock() (token.Token, error) {
	if l.terminal.Kind == token.BlockCommentClose {
		// Already hit the terminal token, clean up, and return it.
		terminal := l.terminal
		l.terminal = token.Token{}
		l.mode = normal
		return terminal, nil
	}
	start := l.here()
	var end source.Location
	for l.character != 0 {
		if l.character == '/' {
			terminal := l.here()
			// Maybe the end?
			if err := l.readChar(); err != nil {
				return l.newToken(token.Illegal), err
			}
			if l.character == 0 {
				break // Unexpected EOF
			}
			if l.character == ';' {
				// We've read as much conent as there is and hit the close token.
				comment := l.newTokenAt(token.Comment, source.Span(start, end))
				l.terminal = l.newTokenFrom(token.BlockCommentClose, terminal)
				if err := l.readChar(); err != nil {
					return l.newToken(token.Illegal), err
				}
				// Do NOT flip the mode yet, we need to come back for the terminal.
				return comment, nil
			}
			// False alarm, keep reading.
		}
		end = l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
	}
	l.mode = normal
	tok := l.newTokenFrom(token.Illegal, start)
	return tok, newError(tok.Location, "reached end of file while reading block comment")
}

// commentDoc handles lexing in the commentDoc mode.
func (l *Lexer) commentDoc() (token.Token, error) {
	if l.terminal.Kind == token.BraceClose {
		// Already hit the terminal token, clean up, and return it.
		terminal := l.terminal
		l.terminal = token.Token{}
		l.mode = normal
		return terminal, nil
	}
	start := l.here()
	var end source.Location
	for l.character != 0 && l.character != '}' {
		end = l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
	}
	if l.character == 0 {
		l.mode = normal
		tok := l.newTokenFrom(token.Illegal, start)
		return tok, newError(tok.Location, "reached end of file while reading doc comment")
	}
	// We've read as much conent as there is and hit the close token.
	comment := l.newTokenAt(token.Comment, source.Span(start, end))
	l.terminal = l.newToken(token.BraceClose)
	if err := l.readChar(); err != nil {
		return l.newToken(token.Illegal), err
	}
	// Do NOT flip the mode yet, we need to come back for the terminal.
	return comment, nil
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
