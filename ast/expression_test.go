package ast

import "testing"

func TestStringExpressionExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	strExpr := &StringExpression{
		Value:     "value",
		ValueType: StringValueText,
	}

	if kind := strExpr.ExpressionKind(); kind != ExprString {
		t.Errorf("Expected stringExpression %v, got %v", ExprString, kind)
	}

	dataType := strExpr.ExprType(testEnvironment)
	if dataType != Text {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Text)
	}

	anyValue := strExpr.AsAny()
	if anyValue != strExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, strExpr)
	}
}

func TestSymbolExpressionExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	testEnvironment.Scopes["someSymbol"] = Text

	sblExpr := &SymbolExpression{
		Value: "someSymbol",
	}

	if kind := sblExpr.ExpressionKind(); kind != ExprSymbol {
		t.Errorf("Expected SymbolExpression %v, got %v", ExprSymbol, kind)
	}

	if dataType := sblExpr.ExprType(testEnvironment); dataType != Text {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Text)
	}

	if anyValue := sblExpr.AsAny(); anyValue != sblExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, sblExpr)
	}
}

func TestGlobalVariableExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	gvExpr := &GlobalVariableExpression{
		Name: "globalname",
	}

	if kind3 := gvExpr.ExpressionKind(); kind3 != ExprGlobalVariable {
		t.Errorf("Expected GlobalVariableExpression %v, got %v", ExprGlobalVariable, kind3)
	}

	testEnvironment.Scopes["globalname"] = Text
	dataType3 := gvExpr.ExprType(testEnvironment)
	if dataType3 != Text {
		t.Errorf("ExprType() returned %v, expected %v", dataType3, Text)
	}

	if anyValue3 := gvExpr.AsAny(); anyValue3 != gvExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue3, gvExpr)
	}
}

func TestNumberExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	testEnvironment.Scopes["someNumber"] = Text

	numberExpr := &NumberExpression{
		Value: TextValue{},
	}

	if kind := numberExpr.ExpressionKind(); kind != ExprNumber {
		t.Errorf("Expected NumberExpression %v, got %v", ExprNumber, kind)
	}

	if dataType := numberExpr.ExprType(testEnvironment); dataType != Text {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Text)
	}

	if anyValue := numberExpr.AsAny(); anyValue != numberExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, numberExpr)
	}
}

func TestBooleanExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	booleanExpr := &BooleanExpression{
		IsTrue: true,
	}

	if kind := booleanExpr.ExpressionKind(); kind != ExprBoolean {
		t.Errorf("Expected BooleanExpression %v, got %v", ExprBoolean, kind)
	}

	if dataType := booleanExpr.ExprType(testEnvironment); dataType != Boolean {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Boolean)
	}

	if anyValue := booleanExpr.AsAny(); anyValue != booleanExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, booleanExpr)
	}
}

func TestPrefixUnaryExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	prefixUnary := &PrefixUnary{
		Right: &PrefixUnary{},
		Op:    Bang,
	}

	if kind := prefixUnary.ExpressionKind(); kind != ExprPrefixUnary {
		t.Errorf("Expected PrefixUnary %v, got %v", ExprPrefixUnary, kind)
	}

	if dataType := prefixUnary.ExprType(testEnvironment); dataType != Boolean {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Boolean)
	}

	if anyValue := prefixUnary.AsAny(); anyValue != prefixUnary {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, prefixUnary)
	}
}

func TestArithmeticExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	arithmeticExpr := &ArithmeticExpression{
		Left: &InExpression{
			ValuesType: Integer,
		},
		Operator: AOPlus,
		Right: &InExpression{
			ValuesType: Integer,
		},
	}

	if kind := arithmeticExpr.ExpressionKind(); kind != ExprArithmetic {
		t.Errorf("Expected ArithmeticExpression %v, got %v", ExprArithmetic, kind)
	}

	if dataType := arithmeticExpr.ExprType(testEnvironment); dataType != Integer {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Integer)
	}

	if anyValue := arithmeticExpr.AsAny(); anyValue != arithmeticExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, arithmeticExpr)
	}
}

func TestComparisonExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	comparisonExpr := &ComparisonExpression{
		Left: &InExpression{
			ValuesType: Integer,
		},
		Operator: CONullSafeEqual,
		Right: &InExpression{
			ValuesType: Integer,
		},
	}

	if kind := comparisonExpr.ExpressionKind(); kind != ExprComparison {
		t.Errorf("Expected ComparisonExpression %v, got %v", ExprComparison, kind)
	}

	if dataType := comparisonExpr.ExprType(testEnvironment); dataType != Integer {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Integer)
	}

	if anyValue := comparisonExpr.AsAny(); anyValue != comparisonExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, comparisonExpr)
	}
}

func TestLikeExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	likeExpr := &LikeExpression{
		Input: &InExpression{
			ValuesType: Integer,
		},
		Pattern: &InExpression{
			ValuesType: Integer,
		},
	}

	if kind := likeExpr.ExpressionKind(); kind != ExprLike {
		t.Errorf("Expected LikeExpression %v, got %v", ExprLike, kind)
	}

	if dataType := likeExpr.ExprType(testEnvironment); dataType != Boolean {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Boolean)
	}

	if anyValue := likeExpr.AsAny(); anyValue != likeExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, likeExpr)
	}
}

func TestGlobExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	globExpr := &GlobExpression{
		Input: &InExpression{
			ValuesType: Integer,
		},
		Pattern: &InExpression{
			ValuesType: Integer,
		},
	}

	if kind := globExpr.ExpressionKind(); kind != ExprGlob {
		t.Errorf("Expected GlobExpression %v, got %v", ExprGlob, kind)
	}

	if dataType := globExpr.ExprType(testEnvironment); dataType != Boolean {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Boolean)
	}

	if anyValue := globExpr.AsAny(); anyValue != globExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, globExpr)
	}
}

func TestLogicalExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	logicalExpr := &LogicalExpression{
		Left: &InExpression{
			ValuesType: Integer,
		},
		Operator: Or,
		Right: &InExpression{
			ValuesType: Integer,
		},
	}

	if kind := logicalExpr.ExpressionKind(); kind != ExprLogical {
		t.Errorf("Expected LogicalExpression %v, got %v", ExprLogical, kind)
	}

	if dataType := logicalExpr.ExprType(testEnvironment); dataType != Boolean {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Boolean)
	}

	if anyValue := logicalExpr.AsAny(); anyValue != logicalExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, logicalExpr)
	}
}

func TestBitwiseExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	bitwiseExpr := &BitwiseExpression{
		Left: &InExpression{
			ValuesType: Integer,
		},
		Operator: BOOr,
		Right: &InExpression{
			ValuesType: Integer,
		},
	}

	if kind := bitwiseExpr.ExpressionKind(); kind != ExprBitwise {
		t.Errorf("Expected BitwiseExpression %v, got %v", ExprBitwise, kind)
	}

	if dataType := bitwiseExpr.ExprType(testEnvironment); dataType != Integer {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Integer)
	}

	if anyValue := bitwiseExpr.AsAny(); anyValue != bitwiseExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, bitwiseExpr)
	}
}

func TestCallExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	callExpr := &CallExpression{
		FunctionName: "upper",
	}

	if kind := callExpr.ExpressionKind(); kind != ExprCall {
		t.Errorf("Expected CallExpression %v, got %v", ExprCall, kind)
	}

	if dataType := callExpr.ExprType(testEnvironment); dataType != Text {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Text)
	}

	if anyValue := callExpr.AsAny(); anyValue != callExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, callExpr)
	}
}

func TestBetweenExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	betweenExpr := &BetweenExpression{
		Value: &InExpression{
			ValuesType: Integer,
		},
		RangeStart: &InExpression{
			ValuesType: Integer,
		},
		RangeEnd: &InExpression{
			ValuesType: Integer,
		},
	}

	if kind := betweenExpr.ExpressionKind(); kind != ExprBetween {
		t.Errorf("Expected BetweenExpression %v, got %v", ExprBetween, kind)
	}

	if dataType := betweenExpr.ExprType(testEnvironment); dataType != Boolean {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Boolean)
	}

	if anyValue := betweenExpr.AsAny(); anyValue != betweenExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, betweenExpr)
	}
}

func TestCaseExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	caseExpr := &CaseExpression{
		Conditions:   []Expression{&InExpression{ValuesType: Integer}},
		Values:       []Expression{&InExpression{ValuesType: Integer}},
		DefaultValue: &InExpression{ValuesType: Integer},
		ValuesType:   Boolean,
	}

	if kind := caseExpr.ExpressionKind(); kind != ExprCase {
		t.Errorf("Expected CaseExpression %v, got %v", ExprCase, kind)
	}

	if dataType := caseExpr.ExprType(testEnvironment); dataType != Boolean {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Boolean)
	}

	if anyValue := caseExpr.AsAny(); anyValue != caseExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, caseExpr)
	}
}

func TestInExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	inExpr := &InExpression{
		Argument:   &InExpression{ValuesType: Integer},
		Values:     []Expression{&InExpression{ValuesType: Integer}},
		ValuesType: Boolean,
	}

	if kind := inExpr.ExpressionKind(); kind != ExprIn {
		t.Errorf("Expected InExpression %v, got %v", ExprIn, kind)
	}

	if dataType := inExpr.ExprType(testEnvironment); dataType != Boolean {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Boolean)
	}

	if anyValue := inExpr.AsAny(); anyValue != inExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, inExpr)
	}
}

func TestIsNullExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	isnullExpr := &IsNullExpression{
		Argument: &InExpression{ValuesType: Integer},
		HasNot:   true,
	}

	if kind := isnullExpr.ExpressionKind(); kind != ExprIsNull {
		t.Errorf("Expected IsNullExpression %v, got %v", ExprIsNull, kind)
	}

	if dataType := isnullExpr.ExprType(testEnvironment); dataType != Boolean {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Boolean)
	}

	if anyValue := isnullExpr.AsAny(); anyValue != isnullExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, isnullExpr)
	}
}

func TestNullExpressionKind(t *testing.T) {
	testEnvironment := &Environment{
		Scopes: make(map[string]DataType),
	}

	nullExpr := &NullExpression{}

	if kind := nullExpr.ExpressionKind(); kind != ExprNull {
		t.Errorf("Expected NullExpression %v, got %v", ExprNull, kind)
	}

	if dataType := nullExpr.ExprType(testEnvironment); dataType != Null {
		t.Errorf("ExprType() returned %v, expected %v", dataType, Null)
	}

	if anyValue := nullExpr.AsAny(); anyValue != nullExpr {
		t.Errorf("AsAny() returned %v, expected %v", anyValue, nullExpr)
	}
}
