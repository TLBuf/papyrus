// Package lexer implements the Papyrus lexer.
package lexer

import (
	"fmt"
	"unicode/utf8"

	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/token"
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

// lexer provides the ability to lex a Papyrus script.
type lexer struct {
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

// Lex processes the given file into a stream of tokens or an [Error]
// if the input file could not be lexed into a token stream.
func Lex(file *source.File) (TokenStream, error) {
	stream := TokenStream{}
	l := &lexer{
		file:            file,
		line:            1,
		column:          0,
		lineStartOffset: 0,
		lineEndOffset:   -1,
	}
	if err := l.readChar(); err != nil {
		return stream, fmt.Errorf("failed to initialize lexer: %v", err)
	}
	for {
		t, err := l.NextToken()
		if err != nil {
			return stream, err
		}
		stream.tokens = append(stream.tokens, t)
		if t.Kind == token.EOF {
			return stream, nil
		}
	}

}

// NextToken scans the input for the next [token.Token].
//
// Returns an [Error] if the input could not be lexed as a token.
func (l *lexer) NextToken() (token.Token, error) {
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
		err = fmt.Errorf("lexer in unknown lexing mode: %d", l.mode)
	}
	if err != nil || tok.Kind == token.Newline || tok.Kind == token.EOF {
		return tok, err
	}
	if l.terminal.Kind != token.Newline {
		tok.Location.PostambleLength = l.lineEndOffset - tok.Location.ByteOffset - tok.Location.Length
	}
	return tok, err
}

func (l *lexer) nextToken() (token.Token, error) {
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
		if l.character != '=' {
			return l.newTokenAt(token.Assign, start), nil
		}
		tok = l.newTokenFrom(token.Equal, start)
	case '+':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character != '=' {
			return l.newTokenAt(token.Plus, start), nil
		}
		tok = l.newTokenFrom(token.AssignAdd, start)
	case '-':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character != '=' {
			return l.newTokenAt(token.Minus, start), nil
		}
		tok = l.newTokenFrom(token.AssignSubtract, start)
	case '*':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character != '=' {
			return l.newTokenAt(token.Multiply, start), nil
		}
		tok = l.newTokenFrom(token.AssignMultiply, start)
	case '/':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character != '=' {
			return l.newTokenAt(token.Divide, start), nil
		}
		tok = l.newTokenFrom(token.AssignDivide, start)
	case '%':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character != '=' {
			return l.newTokenAt(token.Modulo, start), nil
		}
		tok = l.newTokenFrom(token.AssignModulo, start)
	case '!':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character != '=' {
			return l.newTokenAt(token.LogicalNot, start), nil
		}
		tok = l.newTokenFrom(token.NotEqual, start)
	case '>':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character != '=' {
			return l.newTokenAt(token.Greater, start), nil
		}
		tok = l.newTokenFrom(token.GreaterOrEqual, start)
	case '<':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character != '=' {
			return l.newTokenAt(token.Less, start), nil
		}
		tok = l.newTokenFrom(token.LessOrEqual, start)
	case '|':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character != '|' {
			return l.newTokenAt(token.Illegal, start), newError(tok.Location, "'|' is not a valid operator")
		}
		tok = l.newTokenFrom(token.LogicalOr, start)
	case '&':
		start := l.here()
		if err := l.readChar(); err != nil {
			return l.newToken(token.Illegal), err
		}
		if l.character != '&' {
			return l.newTokenAt(token.Illegal, start), newError(tok.Location, "'&' is not a valid operator")
		}
		tok = l.newTokenFrom(token.LogicalAnd, start)
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
		if l.character != '/' {
			l.mode = commentLine
			return l.newTokenAt(token.Semicolon, start), nil
		}
		l.mode = commentBlock
		tok = l.newTokenFrom(token.BlockCommentOpen, start)
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

func (l *lexer) newToken(t token.Kind) token.Token {
	return token.Token{
		Kind:     t,
		Location: l.here(),
	}
}

func (l *lexer) newTokenAt(t token.Kind, at source.Location) token.Token {
	return token.Token{
		Kind:     t,
		Location: at,
	}
}

func (l *lexer) newTokenFrom(t token.Kind, from source.Location) token.Token {
	return token.Token{
		Kind:     t,
		Location: source.Span(from, l.here()),
	}
}

func (l *lexer) here() source.Location {
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

func (l *lexer) nextByteLocation() source.Location {
	return source.Location{
		File:           l.file,
		ByteOffset:     l.position,
		Length:         1,
		StartLine:      l.line,
		StartColumn:    l.column + 1,
		EndLine:        l.line,
		EndColumn:      l.column + 1,
		PreambleLength: l.position - l.lineStartOffset + 1,
	}
}

func (l *lexer) readIdentifier() (token.Token, error) {
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
	tok := l.newTokenAt(token.LookupIdentifier(string(loc.Text())), loc)
	if l.character == '[' {
		next, err := l.peek()
		if err != nil {
			return l.newTokenAt(token.Illegal, l.nextByteLocation()), err
		}
		if next == ']' {
			// Array type (rather than index or declaration).
			if err := l.readChar(); err != nil {
				return l.newToken(token.Illegal), err
			}
			kind := tok.Kind
			switch kind {
			case token.Bool:
				kind = token.BoolArray
			case token.Int:
				kind = token.IntArray
			case token.Float:
				kind = token.FloatArray
			case token.String:
				kind = token.StringArray
			default:
				kind = token.ObjectArray
			}
			tok = l.newTokenFrom(kind, tok.Location)
			if err := l.readChar(); err != nil {
				return l.newToken(token.Illegal), err
			}
		}
	}
	return tok, nil
}

func (l *lexer) readNumber() (token.Token, error) {
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

func (l *lexer) readString() (token.Token, error) {
	start := l.here()
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
func (l *lexer) commentLine() (token.Token, error) {
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
func (l *lexer) commentBlock() (token.Token, error) {
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
func (l *lexer) commentDoc() (token.Token, error) {
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

func (l *lexer) skipWhitespace() error {
	for l.character == ' ' || l.character == '\t' {
		if err := l.readChar(); err != nil {
			return err
		}
	}
	return nil
}

func (l *lexer) readChar() error {
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
			return newError(l.nextByteLocation(), "encountered invalid UTF-8 at byte %d", l.next)
		}
		l.character = r
		width = w
		l.column++
	}
	l.position = l.next
	l.next += width
	return nil
}

func (l *lexer) peek() (rune, error) {
	r, _ := utf8.DecodeRune(l.file.Text[l.next:])
	if r == utf8.RuneError {
		return 0, newError(l.nextByteLocation(), "encountered invalid UTF-8 at byte %d", l.next)
	}
	return r, nil
}

func (l *lexer) findNextNewlineOffset() int {
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
