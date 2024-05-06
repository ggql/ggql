package parser

import (
	"strconv"
	"strings"
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

// nolint:funlen,goconst,gocyclo
func Tokenize(script string) ([]Token, *Diagnostic) {
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
			tokens = append(tokens, consumeIdentifier(characters, &position, &columnStart))
			continue
		}

		// Global Variable Symbol
		if char == '@' {
			result, err := consumeGlobalVariableName(characters, &position, &columnStart)
			if err.DiaMessage() != "" {
				return nil, err
			}
			tokens = append(tokens, result)
			continue
		}

		// Number
		if unicode.IsDigit(char) {
			if char == '0' && position+1 < length {
				if characters[position+1] == 'x' {
					position += 2
					columnStart += 2
					result, err := consumeHexNumber(characters, &position, &columnStart)
					if err.DiaMessage() != "" {
						return nil, err
					}
					tokens = append(tokens, result)
					continue
				}

				if characters[position+1] == 'b' {
					position += 2
					columnStart += 2
					result, err := consumeBinaryNumber(characters, &position, &columnStart)
					if err.DiaMessage() != "" {
						return nil, err
					}
					tokens = append(tokens, result)
					continue
				}

				if characters[position+1] == 'o' {
					position += 2
					columnStart += 2
					result, err := consumeOctalNumber(characters, &position, &columnStart)
					if err.DiaMessage() != "" {
						return nil, err
					}
					tokens = append(tokens, result)
					continue
				}
			}

			number, err := consumeNumber(characters, &position, &columnStart)
			if err.DiaMessage() != "" {
				return nil, err
			}
			tokens = append(tokens, number)
			continue
		}

		// String literal
		if char == '"' {
			result, err := consumeString(characters, &position, &columnStart)
			if err.DiaMessage() != "" {
				return nil, err
			}
			tokens = append(tokens, result)
			continue
		}

		// All chars between two backticks should be consumed as identifier
		if char == '`' {
			result, err := consumeBackticksIdentifier(characters, &position, &columnStart)
			if err.DiaMessage() != "" {
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
				if err.DiaMessage() != "" {
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
			return nil, NewError("Expect `=` after `:`").
				AddHelp("Only token that has `:` is `:=` so make sure you add `=` after `:`").
				WithLocationSpan(columnStart, position)
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

		return nil, NewError("Unexpected character").WithLocationSpan(columnStart, position)
	}

	return tokens, &Diagnostic{}
}

// nolint:lll
func consumeGlobalVariableName(chars []rune, pos, start *int) (Token, *Diagnostic) {
	// Consume `@`
	*pos += 1

	// Make sure first character is alphabetic
	if *pos < len(chars) && !unicode.IsLetter(chars[*pos]) {
		return Token{}, NewError("Global variable name must start with alphabetic character").AddHelp("Add at least one alphabetic character after @").WithLocationSpan(*start, *pos)
	}

	for *pos < len(chars) && (chars[*pos] == '_' || unicode.IsLetter(chars[*pos]) || unicode.IsDigit(chars[*pos])) {
		*pos += 1
	}

	// Identifier is case-insensitive by default, convert to lowercase to be easy to compare and lookup
	literal := string(chars[*start:*pos])
	buf := strings.ToLower(literal)

	location := Location{
		Start: *start,
		End:   *pos,
	}

	return Token{
		Location: location,
		Kind:     GlobalVariable,
		Literal:  buf,
	}, &Diagnostic{}
}

func consumeIdentifier(chars []rune, pos, start *int) Token {
	for *pos < len(chars) && (chars[*pos] == '_' || unicode.IsLetter(chars[*pos]) || unicode.IsDigit(chars[*pos])) {
		*pos++
	}

	// Identifier is being case-insensitive by default, convert to lowercase to be easy to compare and lookup
	literal := chars[*start:*pos]
	buf := strings.ToLower(string(literal))

	location := Location{
		Start: *start,
		End:   *pos,
	}

	return Token{
		Location: location,
		Kind:     resolveSymbolKind(buf),
		Literal:  buf,
	}
}

func consumeNumber(chars []rune, pos, start *int) (Token, *Diagnostic) {
	kind := Integer

	for *pos < len(chars) && (unicode.IsDigit(chars[*pos]) || chars[*pos] == '_') {
		*pos++
	}

	if *pos < len(chars) && chars[*pos] == '.' {
		*pos++
		kind = Float
		for *pos < len(chars) && (unicode.IsDigit(chars[*pos]) || chars[*pos] == '_') {
			*pos++
		}
	}

	literal := chars[*start:*pos]
	buf := string(literal)

	literalNum := strings.ReplaceAll(buf, "_", "")
	if _, err := strconv.ParseFloat(literalNum, 64); err != nil {
		return Token{}, NewError("invalid number")
	}

	location := Location{
		Start: *start,
		End:   *pos,
	}

	return Token{
		Location: location,
		Kind:     kind,
		Literal:  literalNum,
	}, &Diagnostic{}
}

func consumeBackticksIdentifier(chars []rune, pos, start *int) (Token, *Diagnostic) {
	*pos += 1

	for *pos < len(chars) && chars[*pos] != '`' {
		*pos += 1
	}

	if *pos >= len(chars) {
		return Token{}, NewError("Unterminated backticks").AddHelp("Add ` at the end of the identifier").WithLocationSpan(*start, *pos)
	}

	*pos += 1

	literal := chars[*start+1 : *pos-1]
	identifier := string(literal)

	location := struct {
		Start int
		End   int
	}{
		Start: *start,
		End:   *pos,
	}

	return Token{
		Location: location,
		Kind:     Symbol,
		Literal:  identifier,
	}, &Diagnostic{}
}

// nolint:lll
func consumeBinaryNumber(chars []rune, pos, start *int) (Token, *Diagnostic) {
	hasDigit := false

	for *pos < len(chars) && (chars[*pos] == '0' || chars[*pos] == '1' || chars[*pos] == '_') {
		*pos++
		hasDigit = true
	}

	if !hasDigit {
		return Token{}, NewError("Missing digits after the integer base prefix").AddHelp("Expect at least one binary digits after the prefix 0b").AddHelp("Binary digit mean 0 or 1").WithLocationSpan(*start, *pos)
	}

	literal := chars[*start:*pos]
	buf := string(literal)
	literalNum := strings.ReplaceAll(buf, "_", "")

	convertResult, err := strconv.ParseInt(literalNum, 2, 64)
	if err != nil {
		return Token{}, NewError("Invalid binary number").WithLocationSpan(*start, *pos)
	}

	location := struct {
		Start int
		End   int
	}{
		Start: *start,
		End:   *pos,
	}

	return Token{
		Location: location,
		Kind:     Integer,
		Literal:  strconv.FormatInt(convertResult, 10),
	}, &Diagnostic{}
}

// nolint:lll
func consumeOctalNumber(chars []rune, pos, start *int) (Token, *Diagnostic) {
	hasDigit := false

	for *pos < len(chars) && ((chars[*pos] >= '0' && chars[*pos] < '8') || chars[*pos] == '_') {
		*pos++
		hasDigit = true
	}

	if !hasDigit {
		return Token{}, NewError("Missing digits after the integer base prefix").AddHelp("Expect at least one octal digits after the prefix 0o").AddHelp("Octal digit mean 0 to 8 number").WithLocationSpan(*start, *pos)
	}

	literal := chars[*start:*pos]
	buf := string(literal)
	literalNum := strings.Replace(buf, "_", "", -1)

	convertResult, err := strconv.ParseInt(literalNum, 8, 64)
	if err != nil {
		return Token{}, NewError("Invalid octal number")
	}

	location := struct {
		Start int
		End   int
	}{
		Start: *start,
		End:   *pos,
	}

	return Token{
		Location: location,
		Kind:     Integer,
		Literal:  strconv.FormatInt(convertResult, 10),
	}, &Diagnostic{}
}

// nolint:lll
func consumeHexNumber(chars []rune, pos, start *int) (Token, *Diagnostic) {
	helper := func(r rune) bool {
		if _, err := strconv.ParseUint(string(r), 16, 64); err != nil {
			return false
		}
		return true
	}

	hasDigit := false

	for *pos < len(chars) && (helper(chars[*pos]) || chars[*pos] == '_') {
		*pos++
		hasDigit = true
	}

	if !hasDigit {
		return Token{}, NewError("Missing digits after the integer base prefix").AddHelp("Expect at least one hex digits after the prefix 0x").AddHelp("Hex digit mean 0 to 9 and a to f").WithLocationSpan(*start, *pos)
	}

	literal := chars[*start:*pos]
	buf := string(literal)
	literalNum := strings.ReplaceAll(buf, "_", "")

	convertResult, err := strconv.ParseInt(literalNum, 16, 64)
	if err != nil {
		return Token{}, NewError("Invalid hex decimal number")
	}

	location := struct {
		Start int
		End   int
	}{
		Start: *start,
		End:   *pos,
	}

	return Token{
		Location: location,
		Kind:     Integer,
		Literal:  strconv.FormatInt(convertResult, 10),
	}, &Diagnostic{}
}

// nolint:lll
func consumeString(chars []rune, pos, start *int) (Token, *Diagnostic) {
	*pos += 1

	for *pos < len(chars) && chars[*pos] != '"' {
		*pos += 1
	}

	if *pos >= len(chars) {
		return Token{}, NewError("Unterminated double quote string").AddHelp("Add \" at the end of the String literal").WithLocationSpan(*start, *pos)
	}

	*pos += 1

	literal := chars[*start+1 : *pos-1]
	stringLiteral := string(literal)

	location := struct {
		Start int
		End   int
	}{
		Start: *start,
		End:   *pos,
	}

	return Token{
		Location: location,
		Kind:     String,
		Literal:  stringLiteral,
	}, &Diagnostic{}
}

func ignoreSingleLineComment(chars []rune, pos *int) {
	*pos += 2

	for *pos < len(chars) && chars[*pos] != '\n' {
		*pos += 1
	}

	*pos += 1
}

func ignoreCStyleComment(chars []rune, pos *int) *Diagnostic {
	*pos += 2

	for *pos+1 < len(chars) && (chars[*pos] != '*' || chars[*pos+1] != '/') {
		*pos += 1
	}

	if *pos+2 > len(chars) {
		return NewError("C Style comment must end with */").AddHelp("Add */ at the end of C Style comments").WithLocationSpan(*pos, *pos)
	}

	*pos += 2

	return &Diagnostic{}
}

// nolint:funlen,gocyclo
func resolveSymbolKind(literal string) TokenKind {
	switch strings.ToLower(literal) {
	// Reserved keywords
	case "set":
		return Set
	case "select":
		return Select
	case "distinct":
		return Distinct
	case "from":
		return From
	case "group":
		return Group
	case "where":
		return Where
	case "having":
		return Having
	case "limit":
		return Limit
	case "offset":
		return Offset
	case "order":
		return Order
	case "by":
		return By
	case "case":
		return Case
	case "when":
		return When
	case "then":
		return Then
	case "else":
		return Else
	case "end":
		return End
	case "between":
		return Between
	case "in":
		return In
	case "is":
		return Is
	case "not":
		return Not
	case "like":
		return Like
	case "glob":
		return Glob
	// Logical Operators
	case "or":
		return LogicalOr
	case "and":
		return LogicalAnd
	case "xor":
		return LogicalXor
	// True, False and Null
	case "true":
		return True
	case "false":
		return False
	case "null":
		return Null
	case "as":
		return As
	// Order by DES and ASC
	case "asc":
		return Ascending
	case "desc":
		return Descending
	// Identifier
	default:
		return Symbol
	}
}
