package parser

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ggql/ggql/ast"
)

//nolint:gocritic
func TestParseGql(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	// Test: SET @name = value
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Set,
			Literal: "SET",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    GlobalVariable,
			Literal: "@name",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Equal,
			Literal: "=",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    String,
			Literal: "value",
		},
	}

	_, err := ParserGql(tokens, &env)

	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseSetQuery(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	// Test: SET @invalid
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Set,
			Literal: "SET",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    GlobalVariable,
			Literal: "@one",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Equal,
			Literal: "=",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    Integer,
			Literal: "1",
		},
	}
	position := 0

	ret, err := ParseSetQuery(&env, &tokens, &position)
	fmt.Println(ret)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseSelectQuery(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	// Test: SELECT * FROM branches WHERE is_head = "true"
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Select,
			Literal: "SELECT",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Star,
			Literal: "*",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    From,
			Literal: "FROM",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    Symbol,
			Literal: "branches",
		},
		{
			Location: Location{
				Start: 5,
				End:   6,
			},
			Kind:    Where,
			Literal: "WHERE",
		},
		{
			Location: Location{
				Start: 6,
				End:   7,
			},
			Kind:    Symbol,
			Literal: "is_head",
		},
		{
			Location: Location{
				Start: 7,
				End:   8,
			},
			Kind:    True,
			Literal: "false",
		},
	}

	position := 0

	_, err := ParseSelectQuery(&env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseSelectStatement(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: SELECT
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Select,
			Literal: "SELECT",
		},
	}
	position := 1

	_, err := ParseSelectStatement(&context, &env, &tokens, &position)
	if err.Message == "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseWhereStatement(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: WHERE
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Where,
			Literal: "WHERE",
		},
	}
	position := 0

	_, err := ParseWhereStatement(&context, &env, &tokens, &position)
	if err.Message == "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseGroupByStatement(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: WHERE
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Group,
			Literal: "GROUP",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    By,
			Literal: "BY",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Symbol,
			Literal: "name",
		},
	}
	env.DefineGlobal("name", ast.Text{})
	position := 0

	_, err := ParseGroupByStatement(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseHavingStatement(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: Having is_head = "true"
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Having,
			Literal: "HAVING",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Symbol,
			Literal: "is_head",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Equal,
			Literal: "=",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    True,
			Literal: "true",
		},
	}

	position := 0

	_, err := ParseHavingStatement(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseLimitStatement(t *testing.T) {
	// Test: LIMIT 1
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Limit,
			Literal: "LIMIT",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Integer,
			Literal: "1",
		},
	}

	position := 0

	_, err := ParseLimitStatement(&tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseOffsetStatement(t *testing.T) {
	// Test: OFFSET 1
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Offset,
			Literal: "OFFSET",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Integer,
			Literal: "1",
		},
	}

	position := 0

	_, err := ParseOffsetStatement(&tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseOrderByStatement(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: ORDER BY name
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Order,
			Literal: "ORDER",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    By,
			Literal: "BY",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Symbol,
			Literal: "name",
		},
	}

	position := 0

	_, err := ParseOrderByStatement(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: commit_count > -1
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Greater,
			Literal: ">",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "-1",
		},
	}

	position := 0

	_, err := ParseExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseAssignmentExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: commit_count := 1
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    GlobalVariable,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    ColonEqual,
			Literal: ":=",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "1",
		},
	}

	position := 0

	_, err := ParseAssignmentExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseIsNullExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: 1 IS NULL
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Integer,
			Literal: "1",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Is,
			Literal: "IS",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Null,
			Literal: "NULL",
		},
	}

	position := 0

	_, err := ParseIsNullExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseInExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: "One" IN ("One", "Two")
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    String,
			Literal: "One",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    In,
			Literal: "IN",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    LeftParen,
			Literal: "(",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    String,
			Literal: "One",
		},
		{
			Location: Location{
				Start: 5,
				End:   6,
			},
			Kind:    Comma,
			Literal: ",",
		},
		{
			Location: Location{
				Start: 6,
				End:   7,
			},
			Kind:    String,
			Literal: "Two",
		},
		{
			Location: Location{
				Start: 7,
				End:   8,
			},
			Kind:    RightParen,
			Literal: ")",
		},
	}

	position := 0

	_, err := ParseInExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseBetweenExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: commit_count BETWEEN 2 .. 30000
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Between,
			Literal: "BETWEEN",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "2",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    DotDot,
			Literal: "..",
		},
		{
			Location: Location{
				Start: 5,
				End:   6,
			},
			Kind:    Integer,
			Literal: "30000",
		},
	}

	position := 0

	_, err := ParseBetweenExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseLogicalOrExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: commit_count > 0 || commit_count < 0
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Greater,
			Literal: ">",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "0",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    LogicalOr,
			Literal: "||",
		},
		{
			Location: Location{
				Start: 5,
				End:   6,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 6,
				End:   7,
			},
			Kind:    Less,
			Literal: "<",
		},
		{
			Location: Location{
				Start: 7,
				End:   8,
			},
			Kind:    Integer,
			Literal: "0",
		},
	}

	position := 0

	_, err := ParseLogicalOrExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseLogicalAndExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: commit_count > 0 && commit_count < 0
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Greater,
			Literal: ">",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "0",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    LogicalAnd,
			Literal: "&&",
		},
		{
			Location: Location{
				Start: 5,
				End:   6,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 6,
				End:   7,
			},
			Kind:    Less,
			Literal: "<",
		},
		{
			Location: Location{
				Start: 7,
				End:   8,
			},
			Kind:    Integer,
			Literal: "0",
		},
	}

	position := 0

	_, err := ParseLogicalAndExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:goconst,gocritic
func TestParseBitwiseOrExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: commit_count > 0 | commit_count < 0
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Greater,
			Literal: ">",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "0",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    BitwiseOr,
			Literal: "|",
		},
		{
			Location: Location{
				Start: 5,
				End:   6,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 6,
				End:   7,
			},
			Kind:    Less,
			Literal: "<",
		},
		{
			Location: Location{
				Start: 7,
				End:   8,
			},
			Kind:    Integer,
			Literal: "0",
		},
	}

	position := 0

	_, err := ParseBitwiseOrExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseLogicalXorExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: commit_count > 0 ^ commit_count < 0
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Greater,
			Literal: ">",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "0",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    LogicalXor,
			Literal: "^",
		},
		{
			Location: Location{
				Start: 5,
				End:   6,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 6,
				End:   7,
			},
			Kind:    Less,
			Literal: "<",
		},
		{
			Location: Location{
				Start: 7,
				End:   8,
			},
			Kind:    Integer,
			Literal: "0",
		},
	}

	position := 0

	_, err := ParseLogicalXorExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseBitwiseAndExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: commit_count > 0 & commit_count < 0
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Greater,
			Literal: ">",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "0",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    BitwiseAnd,
			Literal: "&",
		},
		{
			Location: Location{
				Start: 5,
				End:   6,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 6,
				End:   7,
			},
			Kind:    Less,
			Literal: "<",
		},
		{
			Location: Location{
				Start: 7,
				End:   8,
			},
			Kind:    Integer,
			Literal: "0",
		},
	}

	position := 0

	_, err := ParseBitwiseAndExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseEqualityExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: commit_count = 0
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Equal,
			Literal: "=",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "0",
		},
	}

	position := 0

	_, err := ParseEqualityExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseComparisonExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: commit_count > 0
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Greater,
			Literal: ">",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "0",
		},
	}

	position := 0

	_, err := ParseComparisonExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseBitwiseShiftExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: commit_count << 1
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "commit_count",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    BitwiseLeftShift,
			Literal: "<<",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "1",
		},
	}

	position := 0

	_, err := ParseBitwiseShiftExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseTermExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: commit_count > 0
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Integer,
			Literal: "1",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Plus,
			Literal: "+",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "1",
		},
	}

	position := 0

	_, err := ParseTermExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseFactorExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: 1 * 2
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Integer,
			Literal: "1",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Star,
			Literal: "*",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Integer,
			Literal: "2",
		},
	}

	position := 0
	_, err := ParseFactorExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseLikeExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: "10 usd" LIKE "[0-9]* usd"
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    String,
			Literal: "10 usd",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Like,
			Literal: "LIKE",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    String,
			Literal: "[0-9]* usd",
		},
	}

	position := 0

	_, err := ParseLikeExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseGLobExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: "Git Query Language" GLOB "Git*"
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    String,
			Literal: "Git Query Language",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Glob,
			Literal: "GLOB",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    String,
			Literal: "Git*",
		},
	}

	position := 0

	_, err := ParseGlobExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseUnaryExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: !1
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Bang,
			Literal: "!",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Integer,
			Literal: "1",
		},
	}

	position := 0

	_, err := ParseUnaryExpression(&context, &env, &tokens, &position)
	if err.Message == "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseFunctionCallExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: lower(name)
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "lower",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    LeftParen,
			Literal: "(",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    Symbol,
			Literal: "name",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    RightParen,
			Literal: ")",
		},
	}

	position := 0

	_, err := ParseFunctionCallExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseArgumentsExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: (name)
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    LeftParen,
			Literal: "(",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Symbol,
			Literal: "name",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    RightParen,
			Literal: ")",
		},
	}

	position := 0

	_, err := ParseArgumentsExpressions(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParsePrimaryExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: CASE WHEN isRemote THEN 1 ELSE 0 END
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Case,
			Literal: "CASE",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    When,
			Literal: "WHEN",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    True,
			Literal: "isRemote",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    Then,
			Literal: "THEN",
		},
		{
			Location: Location{
				Start: 5,
				End:   6,
			},
			Kind:    Integer,
			Literal: "1",
		},
		{
			Location: Location{
				Start: 6,
				End:   7,
			},
			Kind:    Else,
			Literal: "ELSE",
		},
		{
			Location: Location{
				Start: 7,
				End:   8,
			},
			Kind:    Integer,
			Literal: "0",
		},
		{
			Location: Location{
				Start: 8,
				End:   9,
			},
			Kind:    End,
			Literal: "END",
		},
	}

	position := 0

	_, err := ParsePrimaryExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseGroupExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: ("One")
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    LeftParen,
			Literal: "(",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    String,
			Literal: "One",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    RightParen,
			Literal: ")",
		},
	}

	position := 0

	_, err := ParseGroupExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestParseCaseExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	context := ParserContext{}

	// Test: CASE WHEN isRemote THEN 1 ELSE 0 END
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Case,
			Literal: "CASE",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    When,
			Literal: "WHEN",
		},
		{
			Location: Location{
				Start: 3,
				End:   4,
			},
			Kind:    True,
			Literal: "isRemote",
		},
		{
			Location: Location{
				Start: 4,
				End:   5,
			},
			Kind:    Then,
			Literal: "THEN",
		},
		{
			Location: Location{
				Start: 5,
				End:   6,
			},
			Kind:    Integer,
			Literal: "1",
		},
		{
			Location: Location{
				Start: 6,
				End:   7,
			},
			Kind:    Else,
			Literal: "ELSE",
		},
		{
			Location: Location{
				Start: 7,
				End:   8,
			},
			Kind:    Integer,
			Literal: "0",
		},
		{
			Location: Location{
				Start: 8,
				End:   9,
			},
			Kind:    End,
			Literal: "END",
		},
	}

	position := 0

	_, err := ParseCaseExpression(&context, &env, &tokens, &position)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestCheckFunctionCallExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	// Test: lower(name)
	arguments := []ast.Expression{&ast.SymbolExpression{
		Value: "name",
	}}
	parameters := []ast.DataType{ast.Text{}}
	function_name := "lower"
	location := Location{
		Start: 1,
		End:   2,
	}

	_, err := CheckFunctionCallArguments(&env, &arguments, &parameters, function_name, location)
	if err.Message != "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:goconst,gocritic,gomnd
func TestTypeCheckSelectedFields(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	// Test: invalid
	table_name := "invalid"
	fields_names := []string{"invalid"}
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "invalid",
		},
	}

	position := 0

	for k := range env.Scopes {
		delete(env.Scopes, k)
	}
	env.Scopes["commit_id"] = ast.Text{}

	err := TypeCheckSelectedFields(&env, table_name, &fields_names, &tokens, position)
	if err.Message == "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

//nolint:gocritic
func TestUnExpectedStatementError(t *testing.T) {
	// Test: start == 0
	tokens := []Token{
		{
			Location: Location{
				Start: 0,
				End:   0,
			},
			Kind:    Symbol,
			Literal: "select",
		},
	}
	position := 0

	err := UnExpectedStatementError(&tokens, &position)
	assert.Equal(t, "Unexpected statement", err.Message)
}

//nolint:gocritic
func TestUnExpectedExpressionError(t *testing.T) {
	// Test: position == 0
	tokens := []Token{
		{
			Location: Location{
				Start: 0,
				End:   0,
			},
			Kind:    Symbol,
			Literal: "select",
		},
	}
	position := 0

	err := UnExpectedExpressionError(&tokens, &position)
	assert.Equal(t, "Can't complete parsing this expression", err.Message)

	// Test: current.kind == =
	tokens2 := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Equal,
			Literal: "==",
		},
		{
			Location: Location{
				Start: 2,
				End:   3,
			},
			Kind:    Equal,
			Literal: "==",
		},
	}

	position2 := 1

	err2 := UnExpectedExpressionError(&tokens2, &position2)
	assert.Equal(t, "Unexpected `==`, Just use `=` to check equality", err2.Message)
}

//nolint:gocritic
func TestUnExpectedContentAfterCorrectStatement(t *testing.T) {
	// Test: invalid
	statement_name := "invalid"
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "invalid",
		},
	}
	position := 0

	err := UnExpectedContentAfterCorrectStatement(&statement_name, &tokens, &position)
	assert.Equal(t, "Unexpected content after the end of `INVALID` statement", err.Message)
}

//nolint:gocritic
func TestGetExpressionName(t *testing.T) {
	// Test: symbol
	expression := ast.SymbolExpression{
		Value: "symbol",
	}
	statement, _ := GetExpressionName(&expression)
	if statement == "" {
		assert.Equal(t, "symbol", statement)
	}
}

//nolint:goconst,gocritic
func TestRegisterCurrentTableFieldsTypes(t *testing.T) {
	// Test: commits
	table_name := "commits"
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	RegisterCurrentTableFieldsTypes(table_name, env)
	assert.Equal(t, ast.Text{}, env.Scopes["commit_id"])
}

//nolint:gocritic
func TestSelectAllTableFields(t *testing.T) {
	var selectedFields []string
	var fieldsNames []string
	var fieldsValues []ast.Expression

	tableName := "branches"
	SelectAllTableFields(tableName, &selectedFields, &fieldsNames, &fieldsValues)

	assert.Equal(t, len(selectedFields), 5)
	assert.Equal(t, len(fieldsNames), 5)
	assert.Equal(t, len(fieldsValues), 5)
}

func TestConsumeKind(t *testing.T) {
	// position = 1
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "select",
		},
	}
	position := 1
	kind := Symbol

	_, err := ConsumeKind(tokens, position, kind)
	if err.Error() == "" {
		t.Errorf("ParserGql failed with error: %v", err)
	}
}

func TestGetSafeLocation(t *testing.T) {
	tokens := []Token{
		{
			Location: Location{
				Start: 1,
				End:   2,
			},
			Kind:    Symbol,
			Literal: "select",
		},
	}

	// position = 0
	position := 0

	location := GetSafeLocation(&tokens, position)
	assert.Equal(t, 1, location.Start)
	assert.Equal(t, 2, location.End)
}

//nolint:gocritic
func TestIsAssignmentOperator(t *testing.T) {
	// Test: kind = TokenKind::Symbol
	tokens := Token{
		Location: Location{
			Start: 1,
			End:   2,
		},
		Kind:    Symbol,
		Literal: "select",
	}

	status := IsAssignmentOperator(&tokens)
	assert.Equal(t, false, status)
}

//nolint:gocritic
func TestIsTermOperator(t *testing.T) {
	// Test: kind = TokenKind::Symbol
	tokens := Token{
		Location: Location{
			Start: 1,
			End:   2,
		},
		Kind:    Symbol,
		Literal: "select",
	}

	status := IsTermOperator(&tokens)
	assert.Equal(t, false, status)
}

//nolint:gocritic
func TestIsBitwiseShiftOperator(t *testing.T) {
	// Test: kind = TokenKind::Symbol
	tokens := Token{
		Location: Location{
			Start: 1,
			End:   2,
		},
		Kind:    Symbol,
		Literal: "select",
	}

	status := IsBitwiseShiftOperator(&tokens)
	assert.Equal(t, false, status)
}

//nolint:gocritic
func TestIsPrefixUnaryOperator(t *testing.T) {
	// Test: kind = TokenKind::Symbol
	tokens := Token{
		Location: Location{
			Start: 1,
			End:   2,
		},
		Kind:    Symbol,
		Literal: "select",
	}

	status := IsPrefixUnaryOperator(&tokens)
	assert.Equal(t, false, status)
}

//nolint:gocritic
func TestIsComparisonOperator(t *testing.T) {
	// Test: kind = TokenKind::Symbol
	tokens := Token{
		Location: Location{
			Start: 1,
			End:   2,
		},
		Kind:    Symbol,
		Literal: "select",
	}
	status := IsComparisonOperator(&tokens)
	assert.Equal(t, false, status)
}

//nolint:gocritic
func TestIsFactorOperator(t *testing.T) {
	// Test: kind = TokenKind::Symbol
	tokens := Token{
		Location: Location{
			Start: 1,
			End:   2,
		},
		Kind:    Symbol,
		Literal: "select",
	}
	status := IsFactorOperator(&tokens)
	assert.Equal(t, false, status)
}

//nolint:gocritic
func TestIsAscOrDesc(t *testing.T) {
	// Test: kind = TokenKind::Symbol
	tokens := Token{
		Location: Location{
			Start: 1,
			End:   2,
		},
		Kind:    Symbol,
		Literal: "select",
	}
	status := IsAscOrDesc(&tokens)
	assert.Equal(t, false, status)
}

func TestTypeMismatchError(t *testing.T) {
	location := Location{
		Start: 1,
		End:   2,
	}
	expected := ast.Text{}
	actual := ast.Integer{}

	status := TypeMismatchError(location, expected, actual)
	assert.Equal(t, "Type mismatch expected `Text`, got `Integer`", status.Message)
}
