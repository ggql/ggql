package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ggql/ggql/ast"
)

func TestIsExpressionTypeEquals(t *testing.T) {
	// Cast equal
	scope := &ast.Environment{}
	expr := &ast.StringExpression{
		Value:     "name",
		ValueType: ast.StringValueText,
	}
	text := ast.Text{}

	result := IsExpressionTypeEquals(scope, expr, &text)
	assert.Equal(t, Equals{}, result)

	// Cast DataType::Text to DataType::Time
	scope = &ast.Environment{}
	expr = &ast.StringExpression{
		Value:     "12:36:31",
		ValueType: ast.StringValueText,
	}
	time := ast.Time{}

	result = IsExpressionTypeEquals(scope, expr, &time)
	assert.Equal(t, RightSideCasted{}, result)

	// Cast DataType::Text to DataType::Date
	scope = &ast.Environment{}
	expr = &ast.StringExpression{
		Value:     "2024-01-10",
		ValueType: ast.StringValueText,
	}
	date := ast.Date{}

	result = IsExpressionTypeEquals(scope, expr, &date)
	assert.Equal(t, RightSideCasted{}, result)

	// Cast DataType::Text to DataType::DateTime
	scope = &ast.Environment{}
	expr = &ast.StringExpression{
		Value:     "2024-01-10 12:36:31",
		ValueType: ast.StringValueText,
	}
	dateTime := ast.DateTime{}

	result = IsExpressionTypeEquals(scope, expr, &dateTime)
	assert.Equal(t, RightSideCasted{}, result)

	// Cast not equal
	scope = &ast.Environment{}
	expr = &ast.StringExpression{
		Value:     "invalid",
		ValueType: ast.StringValueText,
	}
	integer := ast.Integer{}

	result = IsExpressionTypeEquals(scope, expr, &integer)
	assert.Equal(t, NotEqualAndCantImplicitCast{}, result)
}

func TestAreTypesEquals(t *testing.T) {
	// Cast equal
	scope := &ast.Environment{}
	lhs := &ast.StringExpression{
		Value:     "name",
		ValueType: ast.StringValueText,
	}
	rhs := &ast.StringExpression{
		Value:     "name",
		ValueType: ast.StringValueText,
	}

	result := AreTypesEquals(scope, lhs, rhs)
	assert.Equal(t, Equals{}, result)

	// Cast DataType::Text to DataType::Time for rhs
	scope = &ast.Environment{}
	lhs = &ast.StringExpression{
		Value:     "name",
		ValueType: ast.StringValueTime,
	}
	rhs = &ast.StringExpression{
		Value:     "12:36:31",
		ValueType: ast.StringValueText,
	}

	result = AreTypesEquals(scope, lhs, rhs)
	assert.Equal(t, RightSideCasted{}, result)

	// Cast DataType::Text to DataType::Time for lhs
	scope = &ast.Environment{}
	lhs = &ast.StringExpression{
		Value:     "12:36:31",
		ValueType: ast.StringValueText,
	}
	rhs = &ast.StringExpression{
		Value:     "name",
		ValueType: ast.StringValueTime,
	}

	result = AreTypesEquals(scope, lhs, rhs)
	assert.Equal(t, LeftSideCasted{}, result)

	// Cast DataType::Text to DataType::Date for rhs
	scope = &ast.Environment{}
	lhs = &ast.StringExpression{
		Value:     "name",
		ValueType: ast.StringValueDate,
	}
	rhs = &ast.StringExpression{
		Value:     "2024-01-10",
		ValueType: ast.StringValueText,
	}

	result = AreTypesEquals(scope, lhs, rhs)
	assert.Equal(t, RightSideCasted{}, result)

	// Cast DataType::Text to DataType::Date for lhs
	scope = &ast.Environment{}
	lhs = &ast.StringExpression{
		Value:     "2024-01-10",
		ValueType: ast.StringValueText,
	}
	rhs = &ast.StringExpression{
		Value:     "name",
		ValueType: ast.StringValueDate,
	}

	result = AreTypesEquals(scope, lhs, rhs)
	assert.Equal(t, LeftSideCasted{}, result)

	// Cast DataType::Text to DataType::DateTime for rhs
	scope = &ast.Environment{}
	lhs = &ast.StringExpression{
		Value:     "name",
		ValueType: ast.StringValueDateTime,
	}
	rhs = &ast.StringExpression{
		Value:     "2024-01-10 12:36:31",
		ValueType: ast.StringValueText,
	}

	result = AreTypesEquals(scope, lhs, rhs)
	assert.Equal(t, RightSideCasted{}, result)

	// Cast DataType::Text to DataType::DateTime for lhs
	scope = &ast.Environment{}
	lhs = &ast.StringExpression{
		Value:     "2024-01-10 12:36:31",
		ValueType: ast.StringValueText,
	}
	rhs = &ast.StringExpression{
		Value:     "name",
		ValueType: ast.StringValueDateTime,
	}

	result = AreTypesEquals(scope, lhs, rhs)
	assert.Equal(t, LeftSideCasted{}, result)

	// Cast not equal
	scope = &ast.Environment{}
	lhsNumer := &ast.NumberExpression{
		Value: ast.IntegerValue{Value: 1},
	}
	rhsNumer := &ast.NumberExpression{
		Value: ast.FloatValue{Value: 1.0},
	}

	result = AreTypesEquals(scope, lhsNumer, rhsNumer)
	assert.Equal(t, NotEqualAndCantImplicitCast{}, result)
}

func TestCheckAllValuesAreSameType(t *testing.T) {
	var arguments []ast.Expression

	// Check null type
	env := &ast.Environment{}

	result := CheckAllValuesAreSameType(env, arguments)
	assert.Equal(t, ast.Any{}, result)

	// Check different type
	env = &ast.Environment{}
	arg1 := &ast.StringExpression{
		Value:     "name",
		ValueType: ast.StringValueText,
	}
	argInteger := &ast.NumberExpression{
		Value: ast.IntegerValue{Value: 1},
	}
	arguments = []ast.Expression{arg1, argInteger}

	result = CheckAllValuesAreSameType(env, arguments)
	assert.Equal(t, nil, result)

	// Check the same type
	env = &ast.Environment{}
	arg1 = &ast.StringExpression{
		Value:     "name",
		ValueType: ast.StringValueText,
	}
	argText := &ast.StringExpression{
		Value:     "name",
		ValueType: ast.StringValueText,
	}
	arguments = []ast.Expression{arg1, argText}

	result = CheckAllValuesAreSameType(env, arguments)
	assert.NotEqual(t, nil, result)
}
