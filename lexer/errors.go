package lexer

import (
	"github.com/TLBuf/papyrus/issue"
)

var (
	intenalInvalidMode = issue.NewInternal(
		"LEXR0001",
		`Lexer is in an invalid state.`,
	)
	errorUnknownToken = issue.NewError(
		"LEXR1001",
		`Encountered text that doesn't lex to any known token.`,
	)
	errorInvalidUTF8 = issue.NewError(
		"LEXR1002",
		`Encountered a byte that cannot be parsed as valid UTF-8.`,
	)
	errorUnclosedBlockComment = issue.NewError(
		"LEXR1003",
		`Reached the end of the file while reading a block comment; did you forget a closing '/;'?`,
	)
	errorUnclosedDocComment = issue.NewError(
		"LEXR1004",
		`Reached the end of the file while reading a doc comment; did you forget a closing '}'?`,
	)
	errorUnclosedString = issue.NewError(
		"LEXR1005",
		`Reached the end of the file while reading a string literal; did you forget a closing '"'?`,
	)
	errorInvalidStringEscape = issue.NewError(
		"LEXR1006",
		`Encountered an invalid string escape sequence; only '\n', '\t', '\\', and '\"' are allowed.`,
	)
	errorInvalidFloatTrailingDot = issue.NewError(
		"LEXR1007",
		`Expected a digit to follow the dot in a float literal; did you forget a '0' after the '.'?`,
	)
	errorInvalidIntTrailingX = issue.NewError(
		"LEXR1008",
		`Expected a digit to follow the 'x' in a hexadecimal integer literal; did you forget a '0' after the 'x'?`,
	)
	errorInvalidOpBitwiseAnd = issue.NewError(
		"LEXR1009",
		`Encountered '&', which is not a valid operator; did you mean '&&'?`,
	)
	errorInvalidOpBitwiseOr = issue.NewError(
		"LEXR1010",
		`Encountered '|', which is not a valid operator; did you mean '||'?`,
	)
	errorMissingNewlineTerm = issue.NewError(
		"LEXR1011",
		`Expected a newline immediately after '\'; did you mean '/'?`,
	)
	errorMissingNewlineCR = issue.NewError(
		"LEXR1012",
		`Expected a newline immediately after a carriage return (\r).`,
	)
)
