package parser

import (
	"github.com/ggql/ggql/ast"
)

type TypeCheckResult interface{}

type Equals struct{}

type NotEqualAndCantImplicitCast struct{}

type RightSideCasted struct {
	expr ast.Expression
}

type LeftSideCasted struct {
	expr ast.Expression
}

// nolint:gocyclo
func AreTypesEquals(scope *ast.Environment, lhs, rhs ast.Expression) TypeCheckResult {
	lhsType := lhs.ExprType(scope)
	rhsType := rhs.ExprType(scope)

	if lhsType == rhsType {
		return Equals{}
	}

	if lhsType.IsTime() && rhsType.IsText() && rhs.ExpressionKind() == ast.ExprString {
		expr, _ := rhs.(*ast.StringExpression)
		stringLiteralValue := expr.Value
		if !ast.IsValidTimeFormat(stringLiteralValue) {
			return NotEqualAndCantImplicitCast{}
		}

		return RightSideCasted{expr: &ast.StringExpression{
			Value:     stringLiteralValue,
			ValueType: ast.StringValueTime,
		}}
	}

	if lhsType.IsText() && rhsType.IsTime() && lhs.ExpressionKind() == ast.ExprString {
		expr, _ := lhs.(*ast.StringExpression)
		stringLiteralValue := expr.Value
		if !ast.IsValidTimeFormat(stringLiteralValue) {
			return NotEqualAndCantImplicitCast{}
		}

		return LeftSideCasted{expr: &ast.StringExpression{
			Value:     stringLiteralValue,
			ValueType: ast.StringValueTime,
		}}
	}

	if lhsType.IsDate() && rhsType.IsText() && rhs.ExpressionKind() == ast.ExprString {
		expr, _ := rhs.(*ast.StringExpression)
		stringLiteralValue := expr.Value
		if !ast.IsValidDateFormat(stringLiteralValue) {
			return NotEqualAndCantImplicitCast{}
		}

		return RightSideCasted{expr: &ast.StringExpression{
			Value:     stringLiteralValue,
			ValueType: ast.StringValueDate,
		}}
	}

	if lhsType.IsText() && rhsType.IsDate() && lhs.ExpressionKind() == ast.ExprString {
		expr, _ := lhs.(*ast.StringExpression)
		stringLiteralValue := expr.Value
		if !ast.IsValidDateFormat(stringLiteralValue) {
			return NotEqualAndCantImplicitCast{}
		}

		return LeftSideCasted{expr: &ast.StringExpression{
			Value:     stringLiteralValue,
			ValueType: ast.StringValueDate,
		}}
	}

	if lhsType.IsDateTime() && rhsType.IsText() && rhs.ExpressionKind() == ast.ExprString {
		expr, _ := rhs.(*ast.StringExpression)
		stringLiteralValue := expr.Value
		if !ast.IsValidDateTimeFormat(stringLiteralValue) {
			return NotEqualAndCantImplicitCast{}
		}

		return RightSideCasted{expr: &ast.StringExpression{
			Value:     stringLiteralValue,
			ValueType: ast.StringValueDateTime,
		}}
	}

	return NotEqualAndCantImplicitCast{}
}
