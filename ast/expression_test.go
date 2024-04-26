package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpressionIsConst(t *testing.T) {
	t.Skip("Skipping TestExpressionIsConst.")
}

func TestAssignmentExpressionKind(t *testing.T) {
	t.Skip("Skipping TestAssignmentExpressionKind.")
}

func TestAssignmentExpressionExprType(t *testing.T) {
	expr := &AssignmentExpression{
		Symbol: "",
		Value: &StringExpression{
			Value:     "",
			ValueType: StringValueText,
		},
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsText())
}

func TestStringExpressionKind(t *testing.T) {
	t.Skip("Skipping TestStringExpressionKind.")
}

func TestStringExpressionExprType(t *testing.T) {
	expr := &StringExpression{
		Value:     "",
		ValueType: StringValueText,
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsText())
}

func TestSymbolExpressionKind(t *testing.T) {
	t.Skip("Skipping TestSymbolExpressionKind.")
}

func TestSymbolExpressionExprType(t *testing.T) {
	expr := &SymbolExpression{
		Value: "field1",
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	scope.Scopes["field1"] = Text{}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsText())

	expr = &SymbolExpression{
		Value: "title",
	}

	ret = expr.ExprType(scope)
	assert.Equal(t, true, ret.IsText())

	expr = &SymbolExpression{
		Value: "invalid",
	}

	ret = expr.ExprType(scope)
	assert.Equal(t, true, ret.IsUndefined())
}

func TestGlobalVariableExpressionKind(t *testing.T) {
	t.Skip("Skipping TestGlobalVariableExpressionKind.")
}

func TestGlobalVariableExpressionExprType(t *testing.T) {
	expr := &GlobalVariableExpression{
		Name: "field1",
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	scope.Scopes["field1"] = Text{}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsText())

	expr = &GlobalVariableExpression{
		Name: "invalid",
	}

	ret = expr.ExprType(scope)
	assert.Equal(t, true, ret.IsUndefined())
}

func TestNumberExpressionKind(t *testing.T) {
	t.Skip("Skipping TestNumberExpressionKind.")
}

func TestNumberExpressionExprType(t *testing.T) {
	expr := &NumberExpression{
		Value: TextValue{"field"},
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsText())
}

func TestBooleanExpressionKind(t *testing.T) {
	t.Skip("Skipping TestBooleanExpressionKind.")
}

func TestBooleanExpressionExprType(t *testing.T) {
	expr := &BooleanExpression{
		IsTrue: false,
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsBool())
}

func TestPrefixUnaryExpressionKind(t *testing.T) {
	t.Skip("Skipping TestPrefixUnaryExpressionKind.")
}

func TestPrefixUnaryExpressionExprType(t *testing.T) {
	expr := &PrefixUnary{
		Right: &NumberExpression{Value: NullValue{}},
		Op:    POMinus,
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsInt())

	expr = &PrefixUnary{
		Right: &NumberExpression{Value: NullValue{}},
		Op:    POBang,
	}

	ret = expr.ExprType(scope)
	assert.Equal(t, true, ret.IsBool())
}

func TestArithmeticExpressionKind(t *testing.T) {
	t.Skip("Skipping TestArithmeticExpressionKind.")
}

func TestArithmeticExpressionExprType(t *testing.T) {
	expr := &ArithmeticExpression{
		Left:     &NumberExpression{Value: IntegerValue{1}},
		Operator: AOPlus,
		Right:    &NumberExpression{Value: IntegerValue{}},
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsInt())

	expr = &ArithmeticExpression{
		Left:     &NumberExpression{Value: IntegerValue{1}},
		Operator: AOPlus,
		Right:    &NumberExpression{Value: FloatValue{1.0}},
	}

	ret = expr.ExprType(scope)
	assert.Equal(t, true, ret.IsFloat())
}

func TestComparisonExpressionKind(t *testing.T) {
	t.Skip("Skipping TestComparisonExpressionKind.")
}

func TestComparisonExpressionExprType(t *testing.T) {
	expr := &ComparisonExpression{
		Left:     &NumberExpression{Value: IntegerValue{1}},
		Operator: CONullSafeEqual,
		Right:    &NumberExpression{Value: IntegerValue{1}},
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsInt())

	expr = &ComparisonExpression{
		Left:     &NumberExpression{Value: IntegerValue{1}},
		Operator: CONotEqual,
		Right:    &NumberExpression{Value: IntegerValue{1}},
	}

	ret = expr.ExprType(scope)
	assert.Equal(t, true, ret.IsBool())
}

func TestLikeExpressionKind(t *testing.T) {
	t.Skip("Skipping TestLikeExpressionKind.")
}

func TestLikeExpressionExprType(t *testing.T) {
	expr := &LikeExpression{
		Input:   &NumberExpression{Value: IntegerValue{1}},
		Pattern: &NumberExpression{Value: IntegerValue{1}},
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsBool())
}

func TestGlobExpressionKind(t *testing.T) {
	t.Skip("Skipping TestGlobExpressionKind.")
}

func TestGlobExpressionExprType(t *testing.T) {
	expr := &GlobExpression{
		Input:   &NumberExpression{Value: IntegerValue{1}},
		Pattern: &NumberExpression{Value: IntegerValue{1}},
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsBool())
}

func TestLogicalExpressionKind(t *testing.T) {
	t.Skip("Skipping TestLogicalExpressionKind.")
}

func TestLogicalExpressionExprType(t *testing.T) {
	expr := &LogicalExpression{
		Left:     &NumberExpression{Value: IntegerValue{1}},
		Operator: LOOr,
		Right:    &NumberExpression{Value: IntegerValue{1}},
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsBool())
}

func TestBitwiseExpressionKind(t *testing.T) {
	t.Skip("Skipping TestBitwiseExpressionKind.")
}

func TestBitwiseExpressionExprType(t *testing.T) {
	expr := &BitwiseExpression{
		Left:     &NumberExpression{Value: IntegerValue{1}},
		Operator: BOOr,
		Right:    &NumberExpression{Value: IntegerValue{1}},
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsInt())
}

func TestCallExpressionKind(t *testing.T) {
	t.Skip("Skipping TestCallExpressionKind.")
}

func TestCallExpressionExprType(t *testing.T) {
	expr := &CallExpression{
		FunctionName: "lower",
		Arguments: []Expression{
			&NumberExpression{Value: IntegerValue{1}},
		},
		IsAggregation: false,
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsText())
}

func TestBetweenExpressionKind(t *testing.T) {
	t.Skip("Skipping TestBetweenExpressionKind.")
}

func TestBetweenExpressionExprType(t *testing.T) {
	expr := &BetweenExpression{
		Value:      &NumberExpression{Value: IntegerValue{1}},
		RangeStart: &NumberExpression{Value: IntegerValue{1}},
		RangeEnd:   &NumberExpression{Value: IntegerValue{1}},
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsBool())
}

func TestCaseExpressionKind(t *testing.T) {
	t.Skip("Skipping TestCaseExpressionKind.")
}

func TestCaseExpressionExprType(t *testing.T) {
	expr := &CaseExpression{
		Conditions:   []Expression{},
		Values:       []Expression{},
		DefaultValue: nil,
		ValuesType:   Text{},
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsText())
}

func TestInExpressionKind(t *testing.T) {
	t.Skip("Skipping TestInExpressionKind.")
}

func TestInExpressionExprType(t *testing.T) {
	expr := &InExpression{
		Argument:      &NumberExpression{Value: IntegerValue{1}},
		Values:        []Expression{},
		ValuesType:    Text{},
		HasNotKeyword: false,
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsText())
}

func TestIsNullExpressionKind(t *testing.T) {
	t.Skip("Skipping TestIsNullExpressionKind.")
}

func TestIsNullExpressionExprType(t *testing.T) {
	expr := &IsNullExpression{
		Argument: &NumberExpression{Value: IntegerValue{1}},
		HasNot:   false,
	}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsBool())
}

func TestNullExpressionKind(t *testing.T) {
	t.Skip("Skipping TestNullExpressionKind.")
}

func TestNullExpressionExprType(t *testing.T) {
	expr := &NullExpression{}

	scope := &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
	}

	ret := expr.ExprType(scope)
	assert.Equal(t, true, ret.IsNull())
}
