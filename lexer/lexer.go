// Package lexer implements the Papyrus lexer.
package lexer

import (
	"bytes"
	"iter"
	"unicode/utf8"

	"github.com/TLBuf/papyrus/issue"
	"github.com/TLBuf/papyrus/source"
	"github.com/TLBuf/papyrus/token"
)

// Lex returns an iterator over the tokens lexed from a source file, the
// iterator returns an [issue.Issue] if lexing the next token failed and
// iteration is halted.
func Lex(file *source.File) iter.Seq2[token.Token, *issue.Issue] {
	l := &lexer{
		file:   file,
		length: file.Len(),
	}
	return func(yield func(token.Token, *issue.Issue) bool) {
		defer func() {
			if r := recover(); r != nil {
				_ = yield(token.Token{}, r.(*issue.Issue))
			}
		}()
		l.readChar()
		for {
			tok := l.nextToken()
			if !yield(tok, nil) {
				return
			}
			if tok.Kind == token.EOF {
				return
			}
		}
	}
}

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
	file      *source.File
	position  uint32
	next      uint32
	length    uint32
	character rune
	mode      mode
	terminal  token.Token
}

func (l *lexer) nextToken() token.Token {
	switch l.mode {
	case normal:
		return l.normal()
	case commentLine:
		return l.commentLine()
	case commentBlock:
		return l.commentBlock()
	case commentDoc:
		return l.commentDoc()
	}
	l.fail(intenalInvalidMode, l.here())
	return token.Token{} // Unreachable, fail panics.
}

func (l *lexer) normal() token.Token {
	var tok token.Token
	l.skipWhitespace()
	switch l.character {
	case 0:
		tok = token.Token{
			Kind:     token.EOF,
			Location: source.NewLocation(l.here().Start(), 0),
		}
	case '(':
		tok = l.newToken(token.ParenthesisOpen)
	case ')':
		tok = l.newToken(token.ParenthesisClose)
	case '[':
		start := l.here()
		l.readChar()
		if l.character != ']' {
			return l.newTokenAt(token.BracketOpen, start)
		}
		tok = l.newTokenFrom(token.ArrayType, start)
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
		l.readChar()
		if l.character != '\n' {
			l.fail(errorMissingNewlineCR, tok.Location)
		}
		tok = l.newTokenFrom(token.Newline, start)
	case '\\':
		l.readChar()
		tok := l.normal()
		if tok.Kind != token.Newline {
			l.fail(errorMissingNewlineTerm, tok.Location)
		}
		return l.normal()
	case '=':
		start := l.here()
		l.readChar()
		if l.character != '=' {
			return l.newTokenAt(token.Assign, start)
		}
		tok = l.newTokenFrom(token.Equal, start)
	case '+':
		start := l.here()
		l.readChar()
		if l.character != '=' {
			return l.newTokenAt(token.Plus, start)
		}
		tok = l.newTokenFrom(token.AssignAdd, start)
	case '-':
		// This section is a little complex since we have to determine whether the
		// '-' is a Minus token or should be folded into a number token.
		start := l.here()
		l.readChar()
		if bytes.EqualFold(l.file.Content()[l.position:l.position+3], []byte{'-', '0', 'x'}) {
			// Papyrus doesn't allow hex integers to be negative directly, emit Minus.
			return l.newTokenAt(token.Minus, start)
		}
		if isDigit(l.character) {
			// If we see a number next, just fold it into the number token.
			return l.readNumber(start)
		}
		if l.character != '=' {
			return l.newTokenAt(token.Minus, start)
		}
		tok = l.newTokenFrom(token.AssignSubtract, start)
	case '*':
		start := l.here()
		l.readChar()
		if l.character != '=' {
			return l.newTokenAt(token.Multiply, start)
		}
		tok = l.newTokenFrom(token.AssignMultiply, start)
	case '/':
		start := l.here()
		l.readChar()
		if l.character != '=' {
			return l.newTokenAt(token.Divide, start)
		}
		tok = l.newTokenFrom(token.AssignDivide, start)
	case '%':
		start := l.here()
		l.readChar()
		if l.character != '=' {
			return l.newTokenAt(token.Modulo, start)
		}
		tok = l.newTokenFrom(token.AssignModulo, start)
	case '!':
		start := l.here()
		l.readChar()
		if l.character != '=' {
			return l.newTokenAt(token.LogicalNot, start)
		}
		tok = l.newTokenFrom(token.NotEqual, start)
	case '>':
		start := l.here()
		l.readChar()
		if l.character != '=' {
			return l.newTokenAt(token.Greater, start)
		}
		tok = l.newTokenFrom(token.GreaterOrEqual, start)
	case '<':
		start := l.here()
		l.readChar()
		if l.character != '=' {
			return l.newTokenAt(token.Less, start)
		}
		tok = l.newTokenFrom(token.LessOrEqual, start)
	case '|':
		start := l.here()
		l.readChar()
		if l.character != '|' {
			l.fail(errorInvalidOpBitwiseOr, start)
		}
		tok = l.newTokenFrom(token.LogicalOr, start)
	case '&':
		start := l.here()
		l.readChar()
		if l.character != '&' {
			l.fail(errorInvalidOpBitwiseAnd, start)
		}
		tok = l.newTokenFrom(token.LogicalAnd, start)
	case '{':
		tok = l.newToken(token.BraceOpen)
		l.readChar()
		l.mode = commentDoc
		return tok
	case ';':
		start := l.here()
		l.readChar()
		if l.character != '/' {
			l.mode = commentLine
			return l.newTokenAt(token.Semicolon, start)
		}
		l.mode = commentBlock
		tok = l.newTokenFrom(token.BlockCommentOpen, start)
	case '"':
		return l.readString()
	default:
		switch {
		case isLetter(l.character):
			return l.readIdentifier()
		case isDigit(l.character):
			return l.readNumber(l.here())
		default:
			l.readChar()
			l.fail(errorUnknownToken, tok.Location)
		}
	}
	l.readChar()
	return tok
}

func (l *lexer) newToken(t token.Kind) token.Token {
	return token.Token{
		Kind:     t,
		Text:     l.file.Content()[l.position : l.position+1],
		Location: l.here(),
	}
}

func (l *lexer) newTokenAt(t token.Kind, at source.Location) token.Token {
	return token.Token{
		Kind:     t,
		Text:     l.file.Bytes(at),
		Location: at,
	}
}

func (l *lexer) newTokenFrom(t token.Kind, from source.Location) token.Token {
	return token.Token{
		Kind:     t,
		Text:     l.file.Bytes(from),
		Location: source.Span(from, l.here()),
	}
}

func (l *lexer) here() source.Location {
	return source.NewLocation(l.position, l.next-l.position)
}

func (l *lexer) nextByteLocation() source.Location {
	return source.NewLocation(l.next, 1)
}

func (l *lexer) readIdentifier() token.Token {
	start := l.here()
	l.readChar()
	end := start
	for isLetter(l.character) || isDigit(l.character) {
		end = l.here()
		l.readChar()
	}
	loc := source.Span(start, end)
	return l.newTokenAt(token.LookupIdentifier(string(l.file.Bytes(loc))), loc)
}

func (l *lexer) readNumber(start source.Location) token.Token {
	first := l.character
	end := l.here()
	l.readChar()
	if first == '0' && (l.character == 'x' || l.character == 'X') {
		// Hex Int
		l.readChar()
		for isHexDigit(l.character) {
			end = l.here()
			l.readChar()
		}
		tok := l.newTokenAt(token.IntLiteral, source.Span(start, end))
		if l.file.Content()[l.position-1] == 'x' || l.file.Content()[l.position-1] == 'X' {
			l.fail(errorInvalidIntTrailingX, tok.Location)
		}
		return tok
	}
	isFloat := false
	for isDigit(l.character) || l.character == '.' {
		if l.character == '.' {
			isFloat = true
		}
		end = l.here()
		l.readChar()
	}
	tok := l.newTokenAt(token.IntLiteral, source.Span(start, end))
	if l.file.Content()[l.position-1] == '.' {
		// Number ends with a dot?
		l.fail(errorInvalidFloatTrailingDot, tok.Location)
	}
	if isFloat {
		tok.Kind = token.FloatLiteral
	}
	return tok
}

func (l *lexer) readString() token.Token {
	start := l.here()
	escaping := false
	for {
		l.readChar()
		if l.character == 0 {
			break
		}
		if !escaping && l.character == '\\' {
			escaping = true
			continue
		}
		if escaping {
			if l.character == 'n' || l.character == 't' || l.character == '"' || l.character == '\\' {
				escaping = false
				continue
			}
			panic(
				issue.New(
					errorInvalidStringEscape,
					l.file,
					source.Span(start, l.here()),
				).WithDetail(
					`\%s`,
					string(l.character),
				),
			)
		}
		if l.character == '"' {
			break
		}
	}
	tok := l.newTokenFrom(token.StringLiteral, start)
	if l.character == 0 {
		tok.Kind = token.Illegal
		l.fail(errorUnclosedString, tok.Location)
	}
	l.readChar()
	return tok
}

// commentLine handles lexing in the commentLine mode.
func (l *lexer) commentLine() token.Token {
	if l.terminal.Kind == token.Newline || l.terminal.Kind == token.EOF {
		// Already hit the terminal token, clean up, and return it.
		terminal := l.terminal
		l.terminal = token.Token{}
		l.mode = normal
		return terminal
	}
	start := l.here()
	end := start
	for {
		terminal := l.here()
		if l.character != 0 && l.character != '\r' && l.character != '\n' {
			end = terminal
			l.readChar()
			continue
		}
		if l.character == '\r' {
			// Maybe the end?
			l.readChar()
		}
		if l.character == '\n' || l.character == 0 {
			// We've read as much conent as there is and hit the close token.
			comment := l.newTokenAt(token.Comment, source.Span(start, end))
			l.terminal = l.newTokenFrom(token.Newline, terminal)
			if l.character == 0 {
				l.terminal.Kind = token.EOF
				return comment
			}
			l.readChar()
			return comment
		}
	}
}

// commentBlock handles lexing in the commentBlock mode.
func (l *lexer) commentBlock() token.Token {
	if l.terminal.Kind == token.BlockCommentClose {
		// Already hit the terminal token, clean up, and return it.
		terminal := l.terminal
		l.terminal = token.Token{}
		l.mode = normal
		return terminal
	}
	start := l.here()
	var end source.Location
	for l.character != 0 {
		if l.character == '/' {
			terminal := l.here()
			// Maybe the end?
			l.readChar()
			if l.character == 0 {
				break // Unexpected EOF
			}
			if l.character == ';' {
				// We've read as much conent as there is and hit the close token.
				comment := l.newTokenAt(token.Comment, source.Span(start, end))
				l.terminal = l.newTokenFrom(token.BlockCommentClose, terminal)
				l.readChar()
				// Do NOT flip the mode yet, we need to come back for the terminal.
				return comment
			}
			// False alarm, keep reading.
		}
		end = l.here()
		l.readChar()
	}
	l.mode = normal
	l.fail(errorUnclosedBlockComment, source.Span(start, l.here()))
	return token.Token{} // Unreachable, fail panics.
}

// commentDoc handles lexing in the commentDoc mode.
func (l *lexer) commentDoc() token.Token {
	if l.terminal.Kind == token.BraceClose {
		// Already hit the terminal token, clean up, and return it.
		terminal := l.terminal
		l.terminal = token.Token{}
		l.mode = normal
		return terminal
	}
	start := l.here()
	var end source.Location
	for l.character != 0 && l.character != '}' {
		end = l.here()
		l.readChar()
	}
	if l.character == 0 {
		l.mode = normal
		l.fail(errorUnclosedDocComment, source.Span(start, l.here()))
	}
	// We've read as much conent as there is and hit the close token.
	comment := l.newTokenAt(token.Comment, source.Span(start, end))
	l.terminal = l.newToken(token.BraceClose)
	l.readChar()
	// Do NOT flip the mode yet, we need to come back for the terminal.
	return comment
}

func (l *lexer) skipWhitespace() {
	for l.character == ' ' || l.character == '\t' {
		l.readChar()
	}
}

func (l *lexer) readChar() {
	width := uint32(1)
	if l.next >= l.length {
		l.character = 0
	} else {
		r, w := utf8.DecodeRune(l.file.Content()[l.next:])
		if r == utf8.RuneError {
			l.fail(errorInvalidUTF8, l.nextByteLocation())
		}
		l.character = r
		width = uint32(w) // #nosec G115 -- DecodeRune only returns in range [0,4].
	}
	l.position = l.next
	l.next += width
}

func (l *lexer) fail(def *issue.Definition, loc source.Location) {
	panic(issue.New(def, l.file, loc))
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
