package engine

import (
	"errors"
	"math"
	"regexp"
	"strings"

	"github.com/ggql/ggql/ast"
)

// nolint:funlen,gocyclo
func EvaluateExpression(env *ast.Environment, expression ast.Expression, titles []string, object []ast.Value) (ast.Value, error) {
	switch expression.Kind() {
	case ast.ExprAssignment:
		expr := expression.(*ast.AssignmentExpression)
		return EvaluateAssignment(env, expr, titles, object)
	case ast.ExprString:
		expr := expression.(*ast.StringExpression)
		return EvaluateString(expr)
	case ast.ExprSymbol:
		expr := expression.(*ast.SymbolExpression)
		return EvaluateSymbol(expr, titles, object)
	case ast.ExprGlobalVariable:
		expr := expression.(*ast.GlobalVariableExpression)
		return EvaluateGlobalVariable(env, expr)
	case ast.ExprNumber:
		expr := expression.(*ast.NumberExpression)
		return EvaluateNumber(expr), nil
	case ast.ExprBoolean:
		expr := expression.(*ast.BooleanExpression)
		return EvaluateBoolean(expr), nil
	case ast.ExprPrefixUnary:
		expr := expression.(*ast.PrefixUnary)
		return EvaluatePrefixUnary(env, expr, titles, object)
	case ast.ExprArithmetic:
		expr := expression.(*ast.ArithmeticExpression)
		return EvaluateArithmetic(env, expr, titles, object)
	case ast.ExprComparison:
		expr := expression.(*ast.ComparisonExpression)
		return EvaluateComparison(env, expr, titles, object)
	case ast.ExprLike:
		expr := expression.(*ast.LikeExpression)
		return EvaluateLike(env, expr, titles, object)
	case ast.ExprGlob:
		expr := expression.(*ast.GlobExpression)
		return EvaluateGlob(env, expr, titles, object)
	case ast.ExprLogical:
		expr := expression.(*ast.LogicalExpression)
		return EvaluateLogical(env, expr, titles, object)
	case ast.ExprBitwise:
		expr := expression.(*ast.BitwiseExpression)
		return EvaluateBitwise(env, expr, titles, object)
	case ast.ExprCall:
		expr := expression.(*ast.CallExpression)
		return EvaluateCall(env, expr, titles, object)
	case ast.ExprBetween:
		expr := expression.(*ast.BetweenExpression)
		return EvaluateBetween(env, expr, titles, object)
	case ast.ExprCase:
		expr := expression.(*ast.CaseExpression)
		return EvaluateCase(env, expr, titles, object)
	case ast.ExprIn:
		expr := expression.(*ast.InExpression)
		return EvaluateIn(env, expr, titles, object)
	case ast.ExprIsNull:
		expr := expression.(*ast.IsNullExpression)
		return EvaluateIsNull(env, expr, titles, object)
	case ast.ExprNull:
		return nil, nil
	default:
		return nil, errors.New("invalid expression kind")
	}
}

func EvaluateAssignment(env *ast.Environment, expr *ast.AssignmentExpression, titles []string, object []ast.Value) (ast.Value, error) {
	value, err := EvaluateExpression(env, expr.Value, titles, object)
	if err != nil {
		return nil, err
	}

	env.Globals[expr.Symbol] = value

	return value, nil
}

func EvaluateString(expr *ast.StringExpression) (ast.Value, error) {
	switch expr.ValueType {
	case ast.StringValueText:
		return ast.TextValue{Value: expr.Value}, nil
	case ast.StringValueTime:
		return ast.TimeValue{Value: expr.Value}, nil
	case ast.StringValueDate:
		return ast.DateValue{Value: ast.DateToTimeStamp(expr.Value)}, nil
	case ast.StringValueDateTime:
		return ast.DateValue{Value: ast.DateTimeToTimeStamp(expr.Value)}, nil
	default:
		return nil, errors.New("invalid string value type")
	}
}

func EvaluateSymbol(expr *ast.SymbolExpression, titles []string, object []ast.Value) (ast.Value, error) {
	for index, title := range titles {
		if expr.Value == title {
			return object[index], nil
		}
	}

	return nil, errors.New("invalid column name")
}

func EvaluateGlobalVariable(env *ast.Environment, expr *ast.GlobalVariableExpression) (ast.Value, error) {
	value, ok := env.Globals[expr.Name]
	if ok {
		return value, nil
	}

	return nil, errors.New("global variable not found")
}

func EvaluateNumber(expr *ast.NumberExpression) ast.Value {
	return expr.Value
}

func EvaluateBoolean(expr *ast.BooleanExpression) ast.Value {
	return ast.BooleanValue{Value: expr.IsTrue}
}

func EvaluatePrefixUnary(env *ast.Environment, expr *ast.PrefixUnary, titles []string, object []ast.Value) (ast.Value, error) {
	rhs, err := EvaluateExpression(env, expr.Right, titles, object)
	if err != nil {
		return nil, err
	}

	switch expr.Op {
	case ast.POMinus:
		if rhs.DataType().IsInt() {
			return ast.IntegerValue{Value: -rhs.AsInt()}, nil
		}
		return ast.FloatValue{Value: -rhs.AsFloat()}, nil
	case ast.POBang:
		return ast.BooleanValue{Value: !rhs.AsBool()}, nil
	default:
		return nil, errors.New("invalid prefix unary operator")
	}
}

func EvaluateArithmetic(env *ast.Environment, expr *ast.ArithmeticExpression, titles []string, object []ast.Value) (ast.Value, error) {
	lhs, err := EvaluateExpression(env, expr.Left, titles, object)
	if err != nil {
		return nil, err
	}

	rhs, err := EvaluateExpression(env, expr.Right, titles, object)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	case ast.AOPlus:
		return ast.FloatValue{Value: lhs.AsFloat() + rhs.AsFloat()}, nil
	case ast.AOMinus:
		return ast.FloatValue{Value: lhs.AsFloat() - rhs.AsFloat()}, nil
	case ast.AOStar:
		return ast.FloatValue{Value: lhs.AsFloat() * rhs.AsFloat()}, nil
	case ast.AOSlash:
		return ast.FloatValue{Value: lhs.AsFloat() / rhs.AsFloat()}, nil
	case ast.AOModulus:
		return ast.FloatValue{Value: math.Mod(lhs.AsFloat(), rhs.AsFloat())}, nil
	default:
		return nil, errors.New("invalid arithmetic operator")
	}
}

// nolint:gocyclo
func EvaluateComparison(env *ast.Environment, expr *ast.ComparisonExpression, titles []string, object []ast.Value) (ast.Value, error) {
	lhs, err := EvaluateExpression(env, expr.Left, titles, object)
	if err != nil {
		return nil, err
	}

	rhs, err := EvaluateExpression(env, expr.Right, titles, object)
	if err != nil {
		return nil, err
	}

	comparisonResult := lhs.Compare(rhs)

	if expr.Operator == ast.CONullSafeEqual {
		if lhs.DataType().IsNull() && rhs.DataType().IsNull() {
			return ast.IntegerValue{Value: 1}, nil
		} else if lhs.DataType().IsNull() || rhs.DataType().IsNull() {
			return ast.IntegerValue{Value: 0}, nil
		} else if comparisonResult == ast.Equal {
			return ast.IntegerValue{Value: 1}, nil
		} else {
			return ast.IntegerValue{Value: 0}, nil
		}
	}

	switch expr.Operator {
	case ast.COGreater:
		return ast.BooleanValue{Value: comparisonResult == ast.Greater}, nil
	case ast.COGreaterEqual:
		return ast.BooleanValue{Value: comparisonResult == ast.Greater || comparisonResult == ast.Equal}, nil
	case ast.COLess:
		return ast.BooleanValue{Value: comparisonResult == ast.Less}, nil
	case ast.COLessEqual:
		return ast.BooleanValue{Value: comparisonResult == ast.Less || comparisonResult == ast.Equal}, nil
	case ast.COEqual:
		return ast.BooleanValue{Value: comparisonResult == ast.Equal}, nil
	case ast.CONotEqual:
		return ast.BooleanValue{Value: comparisonResult != ast.Equal}, nil
	case ast.CONullSafeEqual:
		return ast.BooleanValue{Value: false}, nil
	default:
		return nil, errors.New("invalid comparison operator")
	}
}

func EvaluateLike(env *ast.Environment, expr *ast.LikeExpression, titles []string, object []ast.Value) (ast.Value, error) {
	rhs, err := EvaluateExpression(env, expr.Pattern, titles, object)
	if err != nil {
		return nil, err
	}

	pattern := "^" + strings.ToLower(rhs.AsText()) + "$"
	pattern = strings.ReplaceAll(pattern, "%", ".*")
	pattern = strings.ReplaceAll(pattern, "_", ".")

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	lhs, err := EvaluateExpression(env, expr.Input, titles, object)
	if err != nil {
		return nil, err
	}

	return ast.BooleanValue{Value: regex.MatchString(strings.ToLower(lhs.AsText()))}, nil
}

func EvaluateGlob(env *ast.Environment, expr *ast.GlobExpression, titles []string, object []ast.Value) (ast.Value, error) {
	rhs, err := EvaluateExpression(env, expr.Pattern, titles, object)
	if err != nil {
		return nil, err
	}

	pattern := "^" + regexp.QuoteMeta(rhs.AsText()) + "$"
	pattern = strings.ReplaceAll(pattern, ".", "\\.")
	pattern = strings.ReplaceAll(pattern, "*", ".*")
	pattern = strings.ReplaceAll(pattern, "?", ".")

	regex, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	lhs, err := EvaluateExpression(env, expr.Input, titles, object)
	if err != nil {
		return nil, err
	}

	return ast.BooleanValue{Value: regex.MatchString(lhs.AsText())}, nil
}

func EvaluateLogical(env *ast.Environment, expr *ast.LogicalExpression, titles []string, object []ast.Value) (ast.Value, error) {
	lhs, err := EvaluateExpression(env, expr.Left, titles, object)
	if err != nil {
		return nil, err
	}

	if expr.Operator == ast.LOAnd && !lhs.AsBool() {
		return ast.BooleanValue{Value: false}, nil
	}

	if expr.Operator == ast.LOOr && lhs.AsBool() {
		return ast.BooleanValue{Value: true}, nil
	}

	rhs, err := EvaluateExpression(env, expr.Right, titles, object)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	case ast.LOAnd:
		return ast.BooleanValue{Value: lhs.AsBool() && rhs.AsBool()}, nil
	case ast.LOOr:
		return ast.BooleanValue{Value: lhs.AsBool() || rhs.AsBool()}, nil
	case ast.LOXor:
		return ast.BooleanValue{Value: lhs.AsBool() != rhs.AsBool()}, nil
	default:
		return nil, errors.New("invalid logical operator")
	}
}

// nolint:gomnd
func EvaluateBitwise(env *ast.Environment, expr *ast.BitwiseExpression, titles []string, object []ast.Value) (ast.Value, error) {
	lhs, err := EvaluateExpression(env, expr.Left, titles, object)
	if err != nil {
		return nil, err
	}

	rhs, err := EvaluateExpression(env, expr.Right, titles, object)
	if err != nil {
		return nil, err
	}

	switch expr.Operator {
	case ast.BOOr:
		return ast.IntegerValue{Value: lhs.AsInt() | rhs.AsInt()}, nil
	case ast.BOAnd:
		return ast.IntegerValue{Value: lhs.AsInt() & rhs.AsInt()}, nil
	case ast.BORightShift:
		if rhs.AsInt() >= 64 {
			return nil, errors.New("attempt to shift right with overflow")
		}
		return ast.IntegerValue{Value: lhs.AsInt() >> rhs.AsInt()}, nil
	case ast.BOLeftShift:
		if rhs.AsInt() >= 64 {
			return nil, errors.New("attempt to shift left with overflow")
		}
		return ast.IntegerValue{Value: lhs.AsInt() << rhs.AsInt()}, nil
	default:
		return nil, errors.New("invalid bitwise operator")
	}
}

func EvaluateCall(env *ast.Environment, expr *ast.CallExpression, titles []string, object []ast.Value) (ast.Value, error) {
	function := ast.Functions[expr.FunctionName]
	if function == nil {
		return nil, errors.New("function not found")
	}

	arguments := make([]ast.Value, len(expr.Arguments))

	for i, arg := range expr.Arguments {
		value, err := EvaluateExpression(env, arg, titles, object)
		if err != nil {
			return nil, err
		}
		arguments[i] = value
	}

	return function(arguments), nil
}

func EvaluateBetween(env *ast.Environment, expr *ast.BetweenExpression, titles []string, object []ast.Value) (ast.Value, error) {
	value, err := EvaluateExpression(env, expr.Value, titles, object)
	if err != nil {
		return nil, err
	}

	rangeStart, err := EvaluateExpression(env, expr.RangeStart, titles, object)
	if err != nil {
		return nil, err
	}

	rangeEnd, err := EvaluateExpression(env, expr.RangeEnd, titles, object)
	if err != nil {
		return nil, err
	}

	retStart := value.Compare(rangeStart) == ast.Less || value.Compare(rangeStart) == ast.Equal
	retEnd := value.Compare(rangeEnd) == ast.Greater || value.Compare(rangeStart) == ast.Equal

	return ast.BooleanValue{Value: retStart && retEnd}, nil
}

func EvaluateCase(env *ast.Environment, expr *ast.CaseExpression, titles []string, object []ast.Value) (ast.Value, error) {
	conditions := expr.Conditions
	values := expr.Values

	for i := 0; i < len(conditions); i++ {
		condition, err := EvaluateExpression(env, conditions[i], titles, object)
		if err != nil {
			return nil, err
		}
		if condition.AsBool() {
			return EvaluateExpression(env, values[i], titles, object)
		}
	}

	if expr.DefaultValue == nil {
		return nil, errors.New("invalid case statement")
	}

	return EvaluateExpression(env, expr.DefaultValue, titles, object)
}

func EvaluateIn(env *ast.Environment, expr *ast.InExpression, titles []string, object []ast.Value) (ast.Value, error) {
	argument, err := EvaluateExpression(env, expr.Argument, titles, object)
	if err != nil {
		return nil, err
	}

	for _, valueExpr := range expr.Values {
		value, err := EvaluateExpression(env, valueExpr, titles, object)
		if err != nil {
			return nil, err
		}
		if argument.Equals(value) {
			return ast.BooleanValue{Value: !expr.HasNotKeyword}, nil
		}
	}

	return ast.BooleanValue{Value: expr.HasNotKeyword}, nil
}

func EvaluateIsNull(env *ast.Environment, expr *ast.IsNullExpression, titles []string, object []ast.Value) (ast.Value, error) {
	argument, err := EvaluateExpression(env, expr.Argument, titles, object)
	if err != nil {
		return nil, err
	}

	isNull := argument.DataType().IsNull()
	if expr.HasNot {
		return ast.BooleanValue{Value: !isNull}, nil
	}

	return ast.BooleanValue{Value: isNull}, nil
}
