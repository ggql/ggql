package parser

import (
	"unicode"
)

type TokenKind int

const (
	Set TokenKind = iota
	Select
	Distinct
	From
	Group
	Where
	Having
	Limit
	Offset
	Order
	By
	In
	Is
	Not
	Like
	Glob

	Case
	When
	Then
	Else
	End

	Between
	DotDot

	Greater
	GreaterEqual
	Less
	LessEqual
	Equal
	Bang
	BangEqual
	NullSafeEqual

	As

	LeftParen
	RightParen

	LogicalOr
	LogicalAnd
	LogicalXor

	BitwiseOr
	BitwiseAnd
	BitwiseRightShift
	BitwiseLeftShift

	Symbol
	GlobalVariable
	Integer
	Float
	String

	True
	False
	Null

	ColonEqual

	Plus
	Minus
	Star
	Slash
	Percentage

	Comma
	Dot
	Semicolon

	Ascending
	Descending
)

type Location struct {
	Start int
	End   int
}

type Token struct {
	Location Location
	Kind     TokenKind
	Literal  string
}

// nolint:funlen,gocyclo
func tokenize(script string) ([]Token, GQLError) {
	var tokens []Token

	characters := []rune(script)
	position := 0
	columnStart := 0
	length := len(characters)

	for position < length {
		columnStart = position
		char := characters[position]

		// Symbol
		if unicode.IsLetter(char) {
			identifier := consumeIdentifier(characters, &position, &columnStart)
			tokens = append(tokens, identifier)
			continue
		}

		// Global Variable Symbol
		if char == '@' {
			identifier := consumeGlobalVariableName(characters, &position, &columnStart)
			tokens = append(tokens, identifier)
			continue
		}

		// Number
		if unicode.IsDigit(char) {
			if char == '0' && position+1 < length {
				if characters[position+1] == 'x' {
					position += 2
					columnStart += 2
					result, err := consumeHexNumber(characters, &position, &columnStart)
					if err.Message == "" {
						tokens = append(tokens, result)
					}
					continue
				}

				if characters[position+1] == 'b' {
					position += 2
					columnStart += 2
					result, err := consumeBinaryNumber(characters, &position, &columnStart)
					if err.Message == "" {
						tokens = append(tokens, result)
					}
					continue
				}

				if characters[position+1] == 'o' {
					position += 2
					columnStart += 2
					result, err := consumeOctalNumber(characters, &position, &columnStart)
					if err.Message == "" {
						tokens = append(tokens, result)
					}
					continue
				}
			}

			number, err := consumeNumber(characters, &position, &columnStart)
			if err.Message == "" {
				tokens = append(tokens, number)
			}
			continue
		}

		// String literal
		if char == '"' {
			result, err := consumeString(characters, &position, &columnStart)
			if err.Message == "" {
				tokens = append(tokens, result)
			}
			continue
		}

		// Plus
		// TODO: FIXME

		// Minus
		// TODO: FIXME

		// Star
		// TODO: FIXME

		// Slash
		// TODO: FIXME

		// Percentage
		// TODO: FIXME

		// Or
		// TODO: FIXME

		// And
		// TODO: FIXME

		// Xor
		// TODO: FIXME

		// Comma
		// TODO: FIXME

		// Dot or Range (DotDot)
		// TODO: FIXME

		// Greater or GreaterEqual
		// TODO: FIXME

		// Less, LessEqual or NULL-safe equal
		// TODO: FIXME

		// Equal
		// TODO: FIXME

		// Colon Equal
		// TODO: FIXME

		// Bang or Bang Equal
		// TODO: FIXME

		// Left Paren
		// TODO: FIXME

		// Right Paren
		// TODO: FIXME

		// Semicolon
		// TODO: FIXME

		// Characters to ignoring
		// TODO: FIXME
	}

	return tokens, GQLError{}
}

func consumeGlobalVariableName(chars []rune, pos, start *int) Token {
	// TODO: FIXME
	return Token{}
}

func consumeIdentifier(chars []rune, pos, start *int) Token {
	// TODO: FIXME
	return Token{}
}

func consumeNumber(chars []rune, pos, start *int) (Token, GQLError) {
	// TODO: FIXME
	return Token{}, GQLError{}
}

func consumeBinaryNumber(chars []rune, pos, start *int) (Token, GQLError) {
	// TODO: FIXME
	return Token{}, GQLError{}
}

func consumeOctalNumber(chars []rune, pos, start *int) (Token, GQLError) {
	// TODO: FIXME
	return Token{}, GQLError{}
}

func consumeHexNumber(chars []rune, pos, start *int) (Token, GQLError) {
	// TODO: FIXME
	return Token{}, GQLError{}
}

func consumeString(chars []rune, pos, start *int) (Token, GQLError) {
	// TODO: FIXME
	return Token{}, GQLError{}
}

func ignoreSingleLineComment(chars []rune, pos *int) {
	// TODO: FIXME
}

func ignoreCStyleComment(chars []rune, pos *int) GQLError {
	// TODO: FIXME
	return GQLError{}
}

func resolveSymbolKind(literal string) TokenKind {
	// TODO: FIXME
	return Set
}
