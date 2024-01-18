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
func Tokenize(script string) ([]Token, GQLError) {
	var tokens []Token

	position := 0
	columnStart := 0

	characters := []rune(script)
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
					if err.Message != "" {
						return nil, err
					}
					tokens = append(tokens, result)
					continue
				}

				if characters[position+1] == 'b' {
					position += 2
					columnStart += 2
					result, err := consumeBinaryNumber(characters, &position, &columnStart)
					if err.Message != "" {
						return nil, err
					}
					tokens = append(tokens, result)
					continue
				}

				if characters[position+1] == 'o' {
					position += 2
					columnStart += 2
					result, err := consumeOctalNumber(characters, &position, &columnStart)
					if err.Message != "" {
						return nil, err
					}
					tokens = append(tokens, result)
					continue
				}
			}

			number, err := consumeNumber(characters, &position, &columnStart)
			if err.Message != "" {
				return nil, err
			}
			tokens = append(tokens, number)
			continue
		}

		// String literal
		if char == '"' {
			result, err := consumeString(characters, &position, &columnStart)
			if err.Message != "" {
				return nil, err
			}
			tokens = append(tokens, result)
			continue
		}

		// Plus
		if char == '+' {
			location := Location{Start: columnStart, End: position}
			token := Token{Location: location, Kind: Plus, Literal: "+"}
			tokens = append(tokens, token)
			position += 1
			continue
		}

		// Minus
		if char == '-' {
			// Ignore single line comment which from -- until the end of the current line
			if (position+1 < len(characters)) && (characters[position+1] == '-') {
				ignoreSingleLineComment(characters, &position)
				continue
			}
			location := Location{Start: columnStart, End: position}
			token := Token{Location: location, Kind: Minus, Literal: "-"}
			tokens = append(tokens, token)
			position += 1
			continue
		}

		// Star
		if char == '*' {
			location := Location{Start: columnStart, End: position}
			token := Token{Location: location, Kind: Star, Literal: "*"}
			tokens = append(tokens, token)
			position += 1
			continue
		}

		// Slash
		if char == '/' {
			// Ignore C style comment which from /* comment */
			if (position+1 < len(characters)) && (characters[position+1] == '*') {
				err := ignoreCStyleComment(characters, &position)
				if err.Message != "" {
					return nil, err
				}
				continue
			}
			location := Location{Start: columnStart, End: position}
			token := Token{Location: location, Kind: Slash, Literal: "/"}
			tokens = append(tokens, token)
			position += 1
			continue
		}

		// Percentage
		if char == '%' {
			location := Location{Start: columnStart, End: position}
			token := Token{Location: location, Kind: Percentage, Literal: "%"}
			tokens = append(tokens, token)
			position += 1
			continue
		}

		// Or
		if char == '|' {
			location := Location{Start: columnStart, End: position}
			position += 1
			kind := BitwiseOr
			literal := "|"
			if (position < length) && (characters[position] == '|') {
				position += 1
				kind = LogicalOr
				literal = "||"
			}
			token := Token{Location: location, Kind: kind, Literal: literal}
			tokens = append(tokens, token)
			continue
		}

		// And
		if char == '&' {
			location := Location{Start: columnStart, End: position}
			position += 1
			kind := BitwiseAnd
			literal := "&"
			if (position < length) && (characters[position] == '&') {
				position += 1
				kind = LogicalAnd
				literal = "&&"
			}
			token := Token{Location: location, Kind: kind, Literal: literal}
			tokens = append(tokens, token)
			continue
		}

		// Xor
		if char == '^' {
			location := Location{Start: columnStart, End: position}
			token := Token{Location: location, Kind: LogicalXor, Literal: "^"}
			tokens = append(tokens, token)
			position += 1
			continue
		}

		// Comma
		if char == ',' {
			location := Location{Start: columnStart, End: position}
			token := Token{Location: location, Kind: Comma, Literal: ","}
			tokens = append(tokens, token)
			position += 1
			continue
		}

		// Dot or Range (DotDot)
		if char == '.' {
			location := Location{Start: columnStart, End: position}
			position += 1
			kind := Dot
			literal := "."
			if (position < length) && (characters[position] == '.') {
				position += 1
				kind = DotDot
				literal = ".."
			}
			token := Token{Location: location, Kind: kind, Literal: literal}
			tokens = append(tokens, token)
			continue
		}

		// Greater or GreaterEqual
		if char == '>' {
			location := Location{Start: columnStart, End: position}
			position += 1
			kind := Greater
			literal := ">"
			if (position < length) && (characters[position] == '=') {
				position += 1
				kind = GreaterEqual
				literal = ">="
			} else if (position < length) && (characters[position] == '>') {
				position += 1
				kind = BitwiseRightShift
				literal = ">>"
			}
			token := Token{Location: location, Kind: kind, Literal: literal}
			tokens = append(tokens, token)
			continue
		}

		// Less, LessEqual or NULL-safe equal
		if char == '<' {
			location := Location{Start: columnStart, End: position}
			position += 1
			kind := Less
			literal := "<"
			if (position < length) && (characters[position] == '=') {
				position += 1
				if (position < length) && (characters[position] == '>') {
					position += 1
					kind = NullSafeEqual
					literal = "<=>"
				} else {
					kind = LessEqual
					literal = "<="
				}
			} else if (position < length) && (characters[position] == '<') {
				position += 1
				kind = BitwiseLeftShift
				literal = "<<"
			} else if (position < length) && (characters[position] == '>') {
				position += 1
				kind = BangEqual
				literal = "<>"
			}
			token := Token{Location: location, Kind: kind, Literal: literal}
			tokens = append(tokens, token)
			continue
		}

		// Equal
		if char == '=' {
			location := Location{Start: columnStart, End: position}
			token := Token{Location: location, Kind: Equal, Literal: "="}
			tokens = append(tokens, token)
			position += 1
			continue
		}

		// Colon Equal
		if char == ':' {
			if (position+1 < length) && (characters[position+1] == '=') {
				location := Location{Start: columnStart, End: position}
				token := Token{Location: location, Kind: ColonEqual, Literal: ":="}
				tokens = append(tokens, token)
				position += 2
				continue
			}
			return nil, GQLError{Message: "Expect `=` after `:`", Location: Location{Start: columnStart, End: position}}
		}

		// Bang or Bang Equal
		if char == '!' {
			location := Location{Start: columnStart, End: position}
			position += 1
			kind := Bang
			literal := "!"
			if (position < length) && (characters[position] == '=') {
				position += 1
				kind = BangEqual
				literal = "!="
			}
			token := Token{Location: location, Kind: kind, Literal: literal}
			tokens = append(tokens, token)
			continue
		}

		// Left Paren
		if char == '(' {
			location := Location{Start: columnStart, End: position}
			token := Token{Location: location, Kind: LeftParen, Literal: "("}
			tokens = append(tokens, token)
			position += 1
			continue
		}

		// Right Paren
		if char == ')' {
			location := Location{Start: columnStart, End: position}
			token := Token{Location: location, Kind: RightParen, Literal: ")"}
			tokens = append(tokens, token)
			position += 1
			continue
		}

		// Semicolon
		if char == ';' {
			location := Location{Start: columnStart, End: position}
			token := Token{Location: location, Kind: Semicolon, Literal: ";"}
			tokens = append(tokens, token)
			position += 1
			continue
		}

		// Characters to ignoring
		if char == ' ' || char == '\n' || char == '\t' {
			position += 1
			continue
		}

		return nil, GQLError{Message: "Unexpected character", Location: Location{Start: columnStart, End: position}}
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
