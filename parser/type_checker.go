package parser

import (
	"fmt"

	"github.com/ggql/ggql/ast"
)

type TypeCheckResult interface{}

type Equals struct{}

type NotEqualAndCantImplicitCast struct{}

type Error struct {
	diag *Diagnostic
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
				diag: NewError(fmt.Sprintf("Can't compare Time and Text %s because it can't be implicitly casted to Time", stringLiteralValue)).
					AddHelp("A valid Time format must match `HH:MM:SS` or `HH:MM:SS.SSS`").
					AddHelp("You can use `MAKETIME(hour, minute, second)` function to create date value"),
			}
		}

		return RightSideCasted{
			expr: &ast.StringExpression{
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
				diag: NewError(fmt.Sprintf("Can't compare Date and Text %s because it can't be implicitly casted to Date", stringLiteralValue)).
					AddHelp("A valid Date format must match `YYYY-MM-DD`").
					AddHelp("You can use `MAKEDATE(year, dayOfYear)` function to a create date value"),
			}
		}

		return RightSideCasted{
			expr: &ast.StringExpression{
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
				diag: NewError(fmt.Sprintf("Can't compare DateTime and Text %s because it can't be implicitly casted to DateTime", stringLiteralValue)).
					AddHelp("A valid DateTime format must match `YYYY-MM-DD HH:MM:SS` or `YYYY-MM-DD HH:MM:SS.SSS`"),
			}
		}

		return RightSideCasted{
			expr: &ast.StringExpression{
				Value:     stringLiteralValue,
				ValueType: ast.StringValueDateTime,
			},
		}
	}

	return NotEqualAndCantImplicitCast{}
}

// nolint:funlen,gocyclo
func AreTypesEquals(scope *ast.Environment, lhs, rhs ast.Expression) TypeCheckResult {
	lhsType := lhs.ExprType(scope)
	rhsType := rhs.ExprType(scope)

	if lhsType == rhsType {
		return Equals{}
	}

	if lhsType.IsTime() && rhsType.IsText() && rhs.Kind() == ast.ExprString {
		expr := rhs.AsAny().(ast.StringExpression)
		stringLiteralValue := expr.Value
		if !ast.IsValidTimeFormat(stringLiteralValue) {
			return Error{
				diag: NewError(fmt.Sprintf("Can't compare Time and Text %s because it can't be implicitly casted to Time", stringLiteralValue)).
					AddHelp("A valid Time format must match `HH:MM:SS` or `HH:MM:SS.SSS`").
					AddHelp("You can use `MAKETIME(hour, minute, second)` function to create date value"),
			}
		}

		return RightSideCasted{
			expr: &ast.StringExpression{
				Value:     stringLiteralValue,
				ValueType: ast.StringValueTime,
			},
		}
	}

	if lhsType.IsText() && rhsType.IsTime() && lhs.Kind() == ast.ExprString {
		expr := lhs.AsAny().(ast.StringExpression)
		stringLiteralValue := expr.Value
		if !ast.IsValidTimeFormat(stringLiteralValue) {
			return Error{
				diag: NewError(fmt.Sprintf("Can't compare Text %s and Time because it can't be implicitly casted to Time", stringLiteralValue)).
					AddHelp("A valid Time format must match `HH:MM:SS` or `HH:MM:SS.SSS`").
					AddHelp("You can use `MAKETIME(hour, minute, second)` function to a create date value"),
			}
		}

		return LeftSideCasted{
			expr: &ast.StringExpression{
				Value:     stringLiteralValue,
				ValueType: ast.StringValueTime,
			},
		}
	}

	if lhsType.IsDate() && rhsType.IsText() && rhs.Kind() == ast.ExprString {
		expr := rhs.AsAny().(ast.StringExpression)
		stringLiteralValue := expr.Value
		if !ast.IsValidDateFormat(stringLiteralValue) {
			return Error{
				diag: NewError(fmt.Sprintf("Can't compare Date and Text %s because Text can't be implicitly casted to Date", stringLiteralValue)).
					AddHelp("A valid Date format should be matching `YYYY-MM-DD`").
					AddHelp("You can use `MAKEDATE(year, dayOfYear)` function to a create date value"),
			}
		}

		return RightSideCasted{
			expr: &ast.StringExpression{
				Value:     stringLiteralValue,
				ValueType: ast.StringValueDate,
			},
		}
	}

	if lhsType.IsText() && rhsType.IsDate() && lhs.Kind() == ast.ExprString {
		expr := lhs.AsAny().(ast.StringExpression)
		stringLiteralValue := expr.Value
		if !ast.IsValidDateFormat(stringLiteralValue) {
			return Error{
				diag: NewError(fmt.Sprintf("Can't compare Text %s and Date because Text can't be implicitly casted to Date", stringLiteralValue)).
					AddHelp("A valid Date format should be matching `YYYY-MM-DD`").
					AddHelp("You can use `MAKEDATE(year, dayOfYear)` function to a create date value"),
			}
		}

		return LeftSideCasted{
			expr: &ast.StringExpression{
				Value:     stringLiteralValue,
				ValueType: ast.StringValueDate,
			},
		}
	}

	if lhsType.IsDateTime() && rhsType.IsText() && rhs.Kind() == ast.ExprString {
		expr := rhs.AsAny().(ast.StringExpression)
		stringLiteralValue := expr.Value
		if !ast.IsValidDateTimeFormat(stringLiteralValue) {
			return Error{
				diag: NewError(fmt.Sprintf("Can't compare DateTime and Text %s because it can't be implicitly casted to DateTime", stringLiteralValue)).
					AddHelp("A valid DateTime format must match `YYYY-MM-DD HH:MM:SS` or `YYYY-MM-DD HH:MM:SS.SSS`"),
			}
		}

		return RightSideCasted{
			expr: &ast.StringExpression{
				Value:     stringLiteralValue,
				ValueType: ast.StringValueDateTime,
			},
		}
	}

	if lhsType.IsText() && rhsType.IsDateTime() && lhs.Kind() == ast.ExprString {
		expr := lhs.AsAny().(ast.StringExpression)
		stringLiteralValue := expr.Value
		if !ast.IsValidDateTimeFormat(stringLiteralValue) {
			return Error{
				diag: NewError(fmt.Sprintf("Can't compare Text %s and DateTime because it can't be implicitly casted to DateTime", stringLiteralValue)).
					AddHelp("A valid DateTime format must match `YYYY-MM-DD HH:MM:SS` or `YYYY-MM-DD HH:MM:SS.SSS`"),
			}
		}

		return LeftSideCasted{
			expr: &ast.StringExpression{
				Value:     stringLiteralValue,
				ValueType: ast.StringValueDateTime,
			},
		}
	}

	return NotEqualAndCantImplicitCast{}
}

func CheckAllValuesAreSameType(env *ast.Environment, arguments []ast.Expression) ast.DataType {
	argumentsCount := len(arguments)
	if argumentsCount == 0 {
		return ast.Any{}
	}

	dataType := arguments[0].ExprType(env)

	for _, argument := range arguments[1:] {
		exprType := argument.ExprType(env)
		if dataType != exprType {
			return nil
		}
	}

	return dataType
}
