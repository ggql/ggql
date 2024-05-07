package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ggql/ggql/ast"
)

func TestEvaluateExpression(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	assignmentExpression := ast.AssignmentExpression{
		Symbol: "=",
		Value: &ast.StringExpression{
			Value:     "value",
			ValueType: ast.StringValueText,
		},
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	_, err := EvaluateExpression(&env, &assignmentExpression, titles, object)
	assert.Equal(t, nil, err)

	nullExpression := ast.NullExpression{}

	titles = []string{"title"}
	object = []ast.Value{ast.NullValue{}}

	value, err := EvaluateExpression(&env, &nullExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.DataType().IsNull())
}

func TestEvaluateAssignment(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	assignmentExpression := ast.AssignmentExpression{
		Symbol: "=",
		Value: &ast.StringExpression{
			Value:     "value",
			ValueType: ast.StringValueText,
		},
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluateAssignment(&env, &assignmentExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, "value", value.AsText())
}

func TestEvaluateString(t *testing.T) {
	stringExpression := ast.StringExpression{
		Value:     "text",
		ValueType: ast.StringValueText,
	}

	value, err := EvaluateString(&stringExpression)
	assert.Equal(t, nil, err)
	assert.Equal(t, "text", value.AsText())

	stringExpression = ast.StringExpression{
		Value:     "12:36:31",
		ValueType: ast.StringValueTime,
	}

	value, err = EvaluateString(&stringExpression)
	assert.Equal(t, nil, err)
	assert.Equal(t, "12:36:31", value.AsText())

	stringExpression = ast.StringExpression{
		Value:     "2024-01-10",
		ValueType: ast.StringValueDate,
	}

	value, err = EvaluateString(&stringExpression)
	assert.Equal(t, nil, err)
	assert.Equal(t, ast.DateToTimeStamp("2024-01-10"), value.AsDate())

	stringExpression = ast.StringExpression{
		Value:     "2024-01-10 12:36:31",
		ValueType: ast.StringValueDateTime,
	}

	value, err = EvaluateString(&stringExpression)
	assert.Equal(t, nil, err)
	assert.Equal(t, ast.DateTimeToTimeStamp("2024-01-10 12:36:31"), value.AsDateTime())
}

func TestEvaluateSymbol(t *testing.T) {
	symbolExpression := ast.SymbolExpression{
		Value: "value",
	}

	titles := []string{"value"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluateSymbol(&symbolExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, "object", value.AsText())

	symbolExpression = ast.SymbolExpression{
		Value: "value",
	}

	titles = []string{"invalid"}
	object = []ast.Value{ast.TextValue{Value: "object"}}

	_, err = EvaluateSymbol(&symbolExpression, titles, object)
	assert.NotEqual(t, nil, err)
}

func TestEvaluateGlobalVariable(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	env.Globals["name"] = ast.TextValue{
		Value: "value",
	}

	globalVariableExpression := ast.GlobalVariableExpression{
		Name: "name",
	}

	value, err := EvaluateGlobalVariable(&env, &globalVariableExpression)
	assert.Equal(t, nil, err)
	assert.Equal(t, "value", value.AsText())

	globalVariableExpression = ast.GlobalVariableExpression{
		Name: "invalid",
	}

	_, err = EvaluateGlobalVariable(&env, &globalVariableExpression)
	assert.NotEqual(t, nil, err)
}

func TestEvaluateNumber(t *testing.T) {
	numberExpression := ast.NumberExpression{
		Value: ast.IntegerValue{
			Value: 1,
		},
	}

	value := EvaluateNumber(&numberExpression)
	assert.Equal(t, int64(1), value.AsInt())
}

func TestEvaluateBoolean(t *testing.T) {
	booleanExpression := ast.BooleanExpression{
		IsTrue: false,
	}

	value := EvaluateBoolean(&booleanExpression)
	assert.Equal(t, false, value.AsBool())
}

func TestEvaluatePrefixUnary(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	prefixUnary := ast.PrefixUnary{
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Op: ast.POMinus,
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluatePrefixUnary(&env, &prefixUnary, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(-1), value.AsInt())

	prefixUnary = ast.PrefixUnary{
		Right: &ast.NumberExpression{
			Value: ast.FloatValue{
				Value: 1.0,
			},
		},
		Op: ast.POMinus,
	}

	value, err = EvaluatePrefixUnary(&env, &prefixUnary, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(-1.0), value.AsFloat())

	prefixUnary = ast.PrefixUnary{
		Right: &ast.BooleanExpression{
			IsTrue: false,
		},
		Op: ast.POBang,
	}

	value, err = EvaluatePrefixUnary(&env, &prefixUnary, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())
}

func TestEvaluateArithmetic(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	arithmeticExpression := ast.ArithmeticExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Operator: ast.AOPlus,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluateArithmetic(&env, &arithmeticExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(2), value.AsInt())

	arithmeticExpression = ast.ArithmeticExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Operator: ast.AOMinus,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
	}

	value, err = EvaluateArithmetic(&env, &arithmeticExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(0), value.AsInt())

	arithmeticExpression = ast.ArithmeticExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 2,
			},
		},
		Operator: ast.AOStar,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
	}

	value, err = EvaluateArithmetic(&env, &arithmeticExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(2), value.AsInt())

	arithmeticExpression = ast.ArithmeticExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 2,
			},
		},
		Operator: ast.AOSlash,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
	}

	value, err = EvaluateArithmetic(&env, &arithmeticExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(2), value.AsInt())

	arithmeticExpression = ast.ArithmeticExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 2,
			},
		},
		Operator: ast.AOModulus,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
	}

	value, err = EvaluateArithmetic(&env, &arithmeticExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(0), value.AsInt())
}

// nolint:funlen
func TestEvaluateComparison(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	comparisonExpression := ast.ComparisonExpression{
		Left:     &ast.NullExpression{},
		Operator: ast.CONullSafeEqual,
		Right:    &ast.NullExpression{},
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluateComparison(&env, &comparisonExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1), value.AsInt())

	comparisonExpression = ast.ComparisonExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Operator: ast.CONullSafeEqual,
		Right:    &ast.NullExpression{},
	}

	value, err = EvaluateComparison(&env, &comparisonExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(0), value.AsInt())

	comparisonExpression = ast.ComparisonExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Operator: ast.CONullSafeEqual,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
	}

	value, err = EvaluateComparison(&env, &comparisonExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1), value.AsInt())

	comparisonExpression = ast.ComparisonExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Operator: ast.CONullSafeEqual,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 2,
			},
		},
	}

	value, err = EvaluateComparison(&env, &comparisonExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(0), value.AsInt())

	value, err = EvaluateComparison(&env, &comparisonExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(0), value.AsInt())

	comparisonExpression = ast.ComparisonExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 2,
			},
		},
		Operator: ast.COGreater,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
	}

	value, err = EvaluateComparison(&env, &comparisonExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())

	comparisonExpression = ast.ComparisonExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 2,
			},
		},
		Operator: ast.COGreaterEqual,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
	}

	value, err = EvaluateComparison(&env, &comparisonExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())

	comparisonExpression = ast.ComparisonExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Operator: ast.COLess,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 2,
			},
		},
	}

	value, err = EvaluateComparison(&env, &comparisonExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())

	comparisonExpression = ast.ComparisonExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Operator: ast.COLessEqual,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 2,
			},
		},
	}

	value, err = EvaluateComparison(&env, &comparisonExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())

	comparisonExpression = ast.ComparisonExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Operator: ast.COEqual,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
	}

	value, err = EvaluateComparison(&env, &comparisonExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())

	comparisonExpression = ast.ComparisonExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Operator: ast.CONullSafeEqual,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 2,
			},
		},
	}

	value, err = EvaluateComparison(&env, &comparisonExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, value.AsBool())
}

func TestEvaluateLike(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	likeExpression := ast.LikeExpression{
		Input: &ast.StringExpression{
			Value:     "10 usd",
			ValueType: ast.StringValueText,
		},
		Pattern: &ast.StringExpression{
			Value:     "[0-9]* usd",
			ValueType: ast.StringValueText,
		},
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluateLike(&env, &likeExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())

	likeExpression = ast.LikeExpression{
		Input: &ast.StringExpression{
			Value:     "10 usd",
			ValueType: ast.StringValueText,
		},
		Pattern: &ast.StringExpression{
			Value:     "1",
			ValueType: ast.StringValueText,
		},
	}

	value, err = EvaluateLike(&env, &likeExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, value.AsBool())
}

func TestEvaluateGlob(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	globExpression := ast.GlobExpression{
		Input: &ast.StringExpression{
			Value:     "Git Query Language",
			ValueType: ast.StringValueText,
		},
		Pattern: &ast.StringExpression{
			Value:     "Git*",
			ValueType: ast.StringValueText,
		},
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluateGlob(&env, &globExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())

	globExpression = ast.GlobExpression{
		Input: &ast.StringExpression{
			Value:     "Git Query Language",
			ValueType: ast.StringValueText,
		},
		Pattern: &ast.StringExpression{
			Value:     "1",
			ValueType: ast.StringValueText,
		},
	}

	value, err = EvaluateGlob(&env, &globExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, value.AsBool())
}

func TestEvaluateLogical(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	logicalExpression := ast.LogicalExpression{
		Left: &ast.BooleanExpression{
			IsTrue: false,
		},
		Operator: ast.LOAnd,
		Right: &ast.BooleanExpression{
			IsTrue: false,
		},
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluateLogical(&env, &logicalExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, value.AsBool())

	logicalExpression = ast.LogicalExpression{
		Left: &ast.BooleanExpression{
			IsTrue: false,
		},
		Operator: ast.LOOr,
		Right: &ast.BooleanExpression{
			IsTrue: true,
		},
	}

	value, err = EvaluateLogical(&env, &logicalExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())

	logicalExpression = ast.LogicalExpression{
		Left: &ast.BooleanExpression{
			IsTrue: false,
		},
		Operator: ast.LOXor,
		Right: &ast.BooleanExpression{
			IsTrue: true,
		},
	}

	value, err = EvaluateLogical(&env, &logicalExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())
}

func TestEvaluateBitwise(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	bitwiseExpression := ast.BitwiseExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Operator: ast.BOOr,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 0,
			},
		},
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluateBitwise(&env, &bitwiseExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1), value.AsInt())

	bitwiseExpression = ast.BitwiseExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Operator: ast.BOAnd,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 0,
			},
		},
	}

	value, err = EvaluateBitwise(&env, &bitwiseExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(0), value.AsInt())

	bitwiseExpression = ast.BitwiseExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 2,
			},
		},
		Operator: ast.BORightShift,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
	}

	value, err = EvaluateBitwise(&env, &bitwiseExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1), value.AsInt())

	bitwiseExpression = ast.BitwiseExpression{
		Left: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		Operator: ast.BOLeftShift,
		Right: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
	}

	value, err = EvaluateBitwise(&env, &bitwiseExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(2), value.AsInt())
}

func TestEvaluateCall(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	callExpression := ast.CallExpression{
		FunctionName: "lower",
		Arguments: []ast.Expression{
			&ast.StringExpression{
				Value:     "NAME",
				ValueType: ast.StringValueText,
			},
		},
		IsAggregation: false,
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluateCall(&env, &callExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, "name", value.AsText())
}

func TestEvaluateBetween(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	betweenExpression := ast.BetweenExpression{
		Value: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 0,
			},
		},
		RangeStart: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		RangeEnd: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 3,
			},
		},
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluateBetween(&env, &betweenExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, value.AsBool())

	betweenExpression = ast.BetweenExpression{
		Value: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		RangeStart: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		RangeEnd: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 3,
			},
		},
	}

	value, err = EvaluateBetween(&env, &betweenExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())

	betweenExpression = ast.BetweenExpression{
		Value: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 3,
			},
		},
		RangeStart: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		RangeEnd: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 3,
			},
		},
	}

	value, err = EvaluateBetween(&env, &betweenExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())

	betweenExpression = ast.BetweenExpression{
		Value: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 4,
			},
		},
		RangeStart: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		RangeEnd: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 3,
			},
		},
	}

	value, err = EvaluateBetween(&env, &betweenExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, value.AsBool())
}

func TestEvaluateCase(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	caseExpression := ast.CaseExpression{
		Conditions: []ast.Expression{
			&ast.StringExpression{
				Value:     "isRemote",
				ValueType: ast.StringValueText,
			},
		},
		Values: []ast.Expression{
			&ast.NumberExpression{
				Value: ast.IntegerValue{
					Value: 1,
				},
			},
		},
		DefaultValue: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 0,
			},
		},
		ValuesType: ast.Integer{},
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	_, err := EvaluateCase(&env, &caseExpression, titles, object)
	assert.Equal(t, nil, err)

	caseExpression = ast.CaseExpression{
		Conditions:   []ast.Expression{},
		Values:       []ast.Expression{},
		DefaultValue: nil,
		ValuesType:   ast.Integer{},
	}

	_, err = EvaluateCase(&env, &caseExpression, titles, object)
	assert.NotEqual(t, nil, err)
}

func TestEvaluateIn(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	inExpression := ast.InExpression{
		Argument: &ast.StringExpression{
			Value:     "One",
			ValueType: ast.StringValueText,
		},
		Values: []ast.Expression{
			&ast.StringExpression{
				Value:     "One",
				ValueType: ast.StringValueText,
			},
			&ast.StringExpression{
				Value:     "Two",
				ValueType: ast.StringValueText,
			},
		},
		ValuesType:    ast.Text{},
		HasNotKeyword: false,
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluateIn(&env, &inExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())
}

func TestEvaluateIsNull(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	isNullExpression := ast.IsNullExpression{
		Argument: &ast.NumberExpression{
			Value: ast.IntegerValue{
				Value: 1,
			},
		},
		HasNot: false,
	}

	titles := []string{"title"}
	object := []ast.Value{ast.TextValue{Value: "object"}}

	value, err := EvaluateIsNull(&env, &isNullExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, false, value.AsBool())

	isNullExpression = ast.IsNullExpression{
		Argument: &ast.NullExpression{},
		HasNot:   false,
	}

	value, err = EvaluateIsNull(&env, &isNullExpression, titles, object)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, value.AsBool())
}
