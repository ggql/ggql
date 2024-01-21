package parser

import (
	"fmt"
	"github.com/ggql/ggql/ast"
)

type TypeCheckResult interface{}

type Equals struct{}

type NotEqualAndCantImplicitCast struct{}

type Error struct {
	err string
}

type RightSideCasted struct {
	expr ast.Expression
}

type LeftSideCasted struct {
	expr ast.Expression
}

func IsExpressionTypeEquals(scope *ast.Environment, expr ast.Expression, dataType ast.DataType) TypeCheckResult {
	exprType := expr.ExprType(scope)

	if exprType == dataType {
		return Equals{}
	}

	if dataType.IsTime() && exprType.IsText() && expr.Kind() == ast.ExprString {
		literal := expr.AsAny().(ast.StringExpression)
		stringLiteralValue := literal.Value
		if !ast.IsValidTimeFormat(stringLiteralValue) {
			return Error{
				err: fmt.Sprintf("Can't compare Time and Text %s because it can't be implicitly casted to Time\n%s\n%s\n",
					stringLiteralValue,
					"A valid Time format must match `HH:MM:SS` or `HH:MM:SS.SSS`",
					"You can use `MAKETIME(hour, minute, second)` function to create date value"),
			}
		}

		return RightSideCasted{
			expr: ast.StringExpression{
				Value:     stringLiteralValue,
				ValueType: ast.StringValueTime,
			},
		}
	}

	if dataType.IsDate() && exprType.IsText() && expr.Kind() == ast.ExprString {
		literal := expr.AsAny().(ast.StringExpression)
		stringLiteralValue := literal.Value
		if !ast.IsValidDateFormat(stringLiteralValue) {
			return Error{
				err: fmt.Sprintf("Can't compare Date and Text %s because it can't be implicitly casted to Date\n%s\n%s\n",
					stringLiteralValue,
					"A valid Date format must match `YYYY-MM-DD`",
					"You can use `MAKEDATE(year, dayOfYear)` function to a create date value"),
			}
		}

		return RightSideCasted{
			expr: ast.StringExpression{
				Value:     stringLiteralValue,
				ValueType: ast.StringValueDate,
			},
		}
	}

	if dataType.IsDateTime() && exprType.IsText() && expr.Kind() == ast.ExprString {
		literal := expr.AsAny().(ast.StringExpression)
		stringLiteralValue := literal.Value
		if !ast.IsValidDateTimeFormat(stringLiteralValue) {
			return Error{
				err: fmt.Sprintf("Can't compare DateTime and Text %s because it can't be implicitly casted to DateTime\n%s\n%s\n",
					stringLiteralValue,
					"A valid DateTime format must match `YYYY-MM-DD HH:MM:SS` or `YYYY-MM-DD HH:MM:SS.SSS`"),
			}
		}

		return RightSideCasted{
			expr: ast.StringExpression{
				Value:     stringLiteralValue,
				ValueType: ast.StringValueDateTime,
			},
		}
	}

	return NotEqualAndCantImplicitCast{}
}

func AreTypesEquals(scope *Environment, lhs Expression, rhs Expression) TypeCheckResult {
	lhsType := lhs.ExprType(scope)
	rhsType := rhs.ExprType(scope)

	if lhsType == rhsType {
		return Equals
	}

	if lhsType.IsTime() && rhsType.IsText() && rhs.Kind() == ExpressionKindString {
		expr := rhs.AsAny().(*StringExpression)
		stringLiteralValue := expr.Value
		if !IsValidTimeFormat(stringLiteralValue) {
			return Error
		}

		return RightSideCasted
	}

	if lhsType.IsText() && rhsType.IsTime() && lhs.Kind() == ExpressionKindString {
		expr := lhs.AsAny().(*StringExpression)
		stringLiteralValue := expr.Value
		if !IsValidTimeFormat(stringLiteralValue) {
			return Error
		}

		return LeftSideCasted
	}

	if lhsType.IsDate() && rhsType.IsText() && rhs.Kind() == ExpressionKindString {
		expr := rhs.AsAny().(*StringExpression)
		stringLiteralValue := expr.Value
		if !IsValidDateFormat(stringLiteralValue) {
			return Error
		}

		return RightSideCasted
	}

	if lhsType.IsText() && rhsType.IsDate() && lhs.Kind() == ExpressionKindString {
		expr := lhs.AsAny().(*StringExpression)
		stringLiteralValue := expr.Value
		if !IsValidDateFormat(stringLiteralValue) {
			return Error
		}

		return LeftSideCasted
	}

	if lhsType.IsDateTime() && rhsType.IsText() && rhs.Kind() == ExpressionKindString {
		expr := rhs.AsAny().(*StringExpression)
		stringLiteralValue := expr.Value
		if !IsValidDateTimeFormat(stringLiteralValue) {
			return Error
		}

		return RightSideCasted
	}

	if lhsType.IsText() && rhsType.IsDateTime() && lhs.Kind() == ExpressionKindString {
		expr := lhs.AsAny().(*StringExpression)
		stringLiteralValue := expr.Value
		if !IsValidDateTimeFormat(stringLiteralValue) {
			return Error
		}

		return LeftSideCasted
	}

	return NotEqualAndCantImplicitCast
}

func CheckAllValuesAreSameType(env *Environment, arguments []Expression) (DataType, error) {
	argumentsCount := len(arguments)
	if argumentsCount == 0 {
		return DataTypeAny, nil
	}

	dataType := arguments[0].ExprType(env)
	for _, argument := range arguments[1:] {
		exprType := argument.ExprType(env)
		if dataType != exprType {
			return 0, errors.New("not all values are of the same type")
		}
	}

	return dataType, nil
}
