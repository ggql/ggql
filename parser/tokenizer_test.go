package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen,goconst,gocritic
func TestTokenize(t *testing.T) {
	// Symbol: NAME
	script := "NAME"
	tokens, err := Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 4, tokens[0].Location.End)
	assert.Equal(t, "name", tokens[0].Literal)
	assert.Equal(t, Symbol, tokens[0].Kind)

	// GlobalVariable: @NAME
	script = "@NAME"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 5, tokens[0].Location.End)
	assert.Equal(t, "@name", tokens[0].Literal)
	assert.Equal(t, GlobalVariable, tokens[0].Kind)

	// Integer: 0x01
	script = "0x01"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 2, tokens[0].Location.Start)
	assert.Equal(t, 4, tokens[0].Location.End)
	assert.Equal(t, "1", tokens[0].Literal)
	assert.Equal(t, Integer, tokens[0].Kind)

	// Integer: 0b01
	script = "0b01"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 2, tokens[0].Location.Start)
	assert.Equal(t, 4, tokens[0].Location.End)
	assert.Equal(t, "1", tokens[0].Literal)
	assert.Equal(t, Integer, tokens[0].Kind)

	// Integer: 0o01
	script = "0o01"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 2, tokens[0].Location.Start)
	assert.Equal(t, 4, tokens[0].Location.End)
	assert.Equal(t, "1", tokens[0].Literal)
	assert.Equal(t, Integer, tokens[0].Kind)

	// Integer: 1
	script = "1"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 1, tokens[0].Location.End)
	assert.Equal(t, "1", tokens[0].Literal)
	assert.Equal(t, Integer, tokens[0].Kind)

	// Float: 0.1
	script = "0.1"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 3, tokens[0].Location.End)
	assert.Equal(t, "0.1", tokens[0].Literal)
	assert.Equal(t, Float, tokens[0].Kind)

	// String: "name"
	script = "\"name\""
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 6, tokens[0].Location.End)
	assert.Equal(t, "name", tokens[0].Literal)
	assert.Equal(t, String, tokens[0].Kind)

	// Symbol: `name`
	script = "`name`"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 6, tokens[0].Location.End)
	assert.Equal(t, "name", tokens[0].Literal)
	assert.Equal(t, Symbol, tokens[0].Kind)

	// Plus: +
	script = "+"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "+", tokens[0].Literal)
	assert.Equal(t, Plus, tokens[0].Kind)

	// Minus: -
	script = "-"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "-", tokens[0].Literal)
	assert.Equal(t, Minus, tokens[0].Kind)

	// Star: *
	script = "*"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "*", tokens[0].Literal)
	assert.Equal(t, Star, tokens[0].Kind)

	// Slash: /
	script = "/"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "/", tokens[0].Literal)
	assert.Equal(t, Slash, tokens[0].Kind)

	// Percentage: %
	script = "%"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "%", tokens[0].Literal)
	assert.Equal(t, Percentage, tokens[0].Kind)

	// BitwiseOr: |
	script = "|"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "|", tokens[0].Literal)
	assert.Equal(t, BitwiseOr, tokens[0].Kind)

	// LogicalOr: ||
	script = "||"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "||", tokens[0].Literal)
	assert.Equal(t, LogicalOr, tokens[0].Kind)

	// BitwiseAnd: &
	script = "&"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "&", tokens[0].Literal)
	assert.Equal(t, BitwiseAnd, tokens[0].Kind)

	// LogicalAnd: &&
	script = "&&"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "&&", tokens[0].Literal)
	assert.Equal(t, LogicalAnd, tokens[0].Kind)

	// LogicalXor: ^
	script = "^"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "^", tokens[0].Literal)
	assert.Equal(t, LogicalXor, tokens[0].Kind)

	// Comma: ,
	script = ","
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, ",", tokens[0].Literal)
	assert.Equal(t, Comma, tokens[0].Kind)

	// Dot: .
	script = "."
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, ".", tokens[0].Literal)
	assert.Equal(t, Dot, tokens[0].Kind)

	// DotDot: ..
	script = ".."
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "..", tokens[0].Literal)
	assert.Equal(t, DotDot, tokens[0].Kind)

	// Greater: >
	script = ">"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, ">", tokens[0].Literal)
	assert.Equal(t, Greater, tokens[0].Kind)

	// GreaterEqual: >=
	script = ">="
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, ">=", tokens[0].Literal)
	assert.Equal(t, GreaterEqual, tokens[0].Kind)

	// BitwiseRightShift: >>
	script = ">>"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, ">>", tokens[0].Literal)
	assert.Equal(t, BitwiseRightShift, tokens[0].Kind)

	// Less: <
	script = "<"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "<", tokens[0].Literal)
	assert.Equal(t, Less, tokens[0].Kind)

	// NulllSafeEqual: <=>
	script = "<=>"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "<=>", tokens[0].Literal)
	assert.Equal(t, NullSafeEqual, tokens[0].Kind)

	// LessEqual: <=
	script = "<="
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "<=", tokens[0].Literal)
	assert.Equal(t, LessEqual, tokens[0].Kind)

	// BitwiseLeftShift: <<
	script = "<<"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "<<", tokens[0].Literal)
	assert.Equal(t, BitwiseLeftShift, tokens[0].Kind)

	// BangEqual: <>
	script = "<>"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "<>", tokens[0].Literal)
	assert.Equal(t, BangEqual, tokens[0].Kind)

	// Equal: =
	script = "="
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "=", tokens[0].Literal)
	assert.Equal(t, Equal, tokens[0].Kind)

	// ColonEqual: :
	script = ":"
	_, err = Tokenize(script)
	assert.Equal(t, "Expect `=` after `:`", err.Message())

	// ColonEqual: :=
	script = ":="
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, ":=", tokens[0].Literal)
	assert.Equal(t, ColonEqual, tokens[0].Kind)

	// Bang: !
	script = "!"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "!", tokens[0].Literal)
	assert.Equal(t, Bang, tokens[0].Kind)

	// BangEqual: !=
	script = "!="
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "!=", tokens[0].Literal)
	assert.Equal(t, BangEqual, tokens[0].Kind)

	// LeftParen: (
	script = "("
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, "(", tokens[0].Literal)
	assert.Equal(t, LeftParen, tokens[0].Kind)

	// RightParen: )
	script = ")"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, ")", tokens[0].Literal)
	assert.Equal(t, RightParen, tokens[0].Kind)

	// Semicolon: ;
	script = ";"
	tokens, err = Tokenize(script)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 1, len(tokens))
	assert.Equal(t, 0, tokens[0].Location.Start)
	assert.Equal(t, 0, tokens[0].Location.End)
	assert.Equal(t, ";", tokens[0].Literal)
	assert.Equal(t, Semicolon, tokens[0].Kind)

	// Invalid: ?
	script = "?"
	_, err = Tokenize(script)
	assert.Equal(t, "Unexpected character", err.Message())
}

func TestConsumeGlobalVariableName(t *testing.T) {
	// Invalid: @_
	chars := []rune{'@', '_'}
	start := 0
	pos := 0
	_, err := consumeGlobalVariableName(chars, &pos, &start)
	assert.Equal(t, "Global variable name must start with alphabetic character", err.Message())

	// GlobalVariable: @N
	chars = []rune{'@', 'N'}
	start = 0
	pos = 0
	token, err := consumeGlobalVariableName(chars, &pos, &start)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 0, token.Location.Start)
	assert.Equal(t, 2, token.Location.End)
	assert.Equal(t, "@n", token.Literal)
	assert.Equal(t, GlobalVariable, token.Kind)
}

func TestConsumeIdentifier(t *testing.T) {
	// Set: SET
	chars := []rune{'S', 'E', 'T'}
	start := 0
	pos := 0
	token := consumeIdentifier(chars, &pos, &start)
	assert.Equal(t, 0, token.Location.Start)
	assert.Equal(t, 3, token.Location.End)
	assert.Equal(t, "set", token.Literal)
	assert.Equal(t, Set, token.Kind)
}

func TestConsumeNumber(t *testing.T) {
	// Integer: 1
	chars := []rune{'1'}
	start := 0
	pos := 0
	token, err := consumeNumber(chars, &pos, &start)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 0, token.Location.Start)
	assert.Equal(t, 1, token.Location.End)
	assert.Equal(t, "1", token.Literal)
	assert.Equal(t, Integer, token.Kind)

	// Integer: 1_0
	chars = []rune{'1', '_', '0'}
	start = 0
	pos = 0
	token, err = consumeNumber(chars, &pos, &start)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 0, token.Location.Start)
	assert.Equal(t, 3, token.Location.End)
	assert.Equal(t, "10", token.Literal)
	assert.Equal(t, Integer, token.Kind)

	// Float: 1.0
	chars = []rune{'1', '.', '0'}
	start = 0
	pos = 0
	token, err = consumeNumber(chars, &pos, &start)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 0, token.Location.Start)
	assert.Equal(t, 3, token.Location.End)
	assert.Equal(t, "1.0", token.Literal)
	assert.Equal(t, Float, token.Kind)

	// Integer: 1_0.0
	chars = []rune{'1', '_', '0', '.', '0'}
	start = 0
	pos = 0
	token, err = consumeNumber(chars, &pos, &start)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 0, token.Location.Start)
	assert.Equal(t, 5, token.Location.End)
	assert.Equal(t, "10.0", token.Literal)
	assert.Equal(t, Float, token.Kind)
}

func TestConsumeBackticksIdentifier(t *testing.T) {
	// Symbol: `N
	chars := []rune{'`', 'N'}
	start := 0
	pos := 0
	_, err := consumeBackticksIdentifier(chars, &pos, &start)
	assert.Equal(t, "Unterminated backticks", err.Message())

	// Symbol: `N`
	chars = []rune{'`', 'N', '`'}
	start = 0
	pos = 0
	token, err := consumeBackticksIdentifier(chars, &pos, &start)
	assert.Equal(t, 0, token.Location.Start)
	assert.Equal(t, 3, token.Location.End)
	assert.Equal(t, "N", token.Literal)
	assert.Equal(t, Symbol, token.Kind)
}

func TestConsumeBinaryNumber(t *testing.T) {
	// Integer: 2
	chars := []rune{'2'}
	start := 0
	pos := 0
	_, err := consumeBinaryNumber(chars, &pos, &start)
	assert.Equal(t, "Missing digits after the integer base prefix", err.Message())

	// Integer: 010
	chars = []rune{'0', '1', '0'}
	start = 0
	pos = 0
	token, err := consumeBinaryNumber(chars, &pos, &start)
	assert.Equal(t, 0, token.Location.Start)
	assert.Equal(t, 3, token.Location.End)
	assert.Equal(t, "2", token.Literal)
	assert.Equal(t, Integer, token.Kind)
}

func TestConsumeOctalNumber(t *testing.T) {
	// Integer: 8
	chars := []rune{'8'}
	start := 0
	pos := 0
	_, err := consumeOctalNumber(chars, &pos, &start)
	assert.Equal(t, "Invalid octal number", err.Message())

	// Integer: 0_7
	chars = []rune{'0', '_', '7'}
	start = 0
	pos = 0
	token, err := consumeOctalNumber(chars, &pos, &start)
	assert.Equal(t, 0, token.Location.Start)
	assert.Equal(t, 3, token.Location.End)
	assert.Equal(t, "7", token.Literal)
	assert.Equal(t, Integer, token.Kind)
}

func TestConsumeHexNumber(t *testing.T) {
	// Integer: G
	chars := []rune{'G'}
	start := 0
	pos := 0
	_, err := consumeHexNumber(chars, &pos, &start)
	assert.Equal(t, "Missing digits after the integer base prefix", err.Message())

	// Integer: 01EF
	chars = []rune{'0', '1', 'E', 'F'}
	start = 0
	pos = 0
	token, err := consumeHexNumber(chars, &pos, &start)
	assert.Equal(t, 0, token.Location.Start)
	assert.Equal(t, 4, token.Location.End)
	assert.Equal(t, "495", token.Literal)
	assert.Equal(t, Integer, token.Kind)
}

func TestConsumeString(t *testing.T) {
	// String: "N
	chars := []rune{'"', 'N'}
	start := 0
	pos := 0
	_, err := consumeString(chars, &pos, &start)
	assert.Equal(t, "Unterminated double quote string", err.Message())

	// String: "N"
	chars = []rune{'"', 'N', '"'}
	start = 0
	pos = 0
	token, err := consumeString(chars, &pos, &start)
	assert.Equal(t, 0, token.Location.Start)
	assert.Equal(t, 3, token.Location.End)
	assert.Equal(t, "N", token.Literal)
	assert.Equal(t, String, token.Kind)
}

func TestIgnoreSingleLineComment(t *testing.T) {
	// Comment: "-- N\n"
	chars := []rune{'-', '-', ' ', 'N', '\n'}
	pos := 0
	ignoreSingleLineComment(chars, &pos)
	assert.Equal(t, 5, pos)
}

// nolint:gocritic
func TestIgnoreCStyleComment(t *testing.T) {
	// Comment: /*N
	chars := []rune{'/', '*', 'N'}
	pos := 0
	err := ignoreCStyleComment(chars, &pos)
	assert.Equal(t, "C Style comment must end with */", err.Message())

	// Comment: /*N*/
	chars = []rune{'/', '*', 'N', '*', '/'}
	pos = 0
	err = ignoreCStyleComment(chars, &pos)
	assert.Equal(t, "", err.Message())
	assert.Equal(t, 5, pos)
}

func TestResolveSymbolKind(t *testing.T) {
	// Set: SET
	literal := "SET"
	kind := resolveSymbolKind(literal)
	assert.Equal(t, Set, kind)

	// Symbol: NAME
	literal = "NAME"
	kind = resolveSymbolKind(literal)
	assert.Equal(t, Symbol, kind)
}
