package parser

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ggql/ggql/ast"
)

func ParserGql(tokens []Token, env *ast.Environment) (ast.Query, Diagnostic) {
	position := 0
	firstToken := tokens[position]
	var queryResult ast.Query
	var err Diagnostic
	switch firstToken.Kind {
	case Set:
		queryResult, err = ParseSetQuery(env, &tokens, &position)
	case Select:
		queryResult, err = ParseSelectQuery(env, &tokens, &position)
	default:
		err = UnExpectedStatementError(&tokens, &position)
	}

	if position < len(tokens) {
		lastToken := tokens[position]
		if lastToken.Kind == Semicolon {
			position += 1
		}
	}

	if position < len(tokens) {
		err = UnExpectedContentAfterCorrectStatement(
			&firstToken.Literal,
			&tokens,
			&position,
		)
	}

	return queryResult, err
}

// nolint:lll
func ParseSetQuery(env *ast.Environment, tokens *[]Token, position *int) (ast.Query, Diagnostic) {
	lentokens := len(*tokens)
	context := ParserContext{}
	// Consume Set keyword
	*position += 1

	if *position >= lentokens || (*tokens)[*position].Kind != GlobalVariable {
		return ast.Query{}, *NewError("Expect Global variable name start with `@` after `SET` keyword").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	name := (*tokens)[*position].Literal

	// Consume variable name
	*position += 1

	if *position >= lentokens || !IsAssignmentOperator(&(*tokens)[*position]) {
		return ast.Query{}, *NewError("Expect `=` or `:=` and Value after Variable name").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	// Consume `=` token
	*position += 1

	aggregationsCountBefore := len(context.Aggregations)
	value, err := ParseExpression(&context, env, tokens, position)
	if err.Message != "" {
		return ast.Query{}, err
	}
	hasAggregations := len(context.Aggregations) != aggregationsCountBefore

	if hasAggregations {
		return ast.Query{}, *NewError("Aggregation value can't be assigned to global variable").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	env.DefineGlobal(name, value.ExprType(env))

	globalVariable := ast.GlobalVariableStatement{
		Name:  name,
		Value: value,
	}

	return ast.Query{GlobalVariableDeclaration: &globalVariable}, Diagnostic{}
}

// nolint:funlen,gocyclo,lll
func ParseSelectQuery(env *ast.Environment, tokens *[]Token, position *int) (ast.Query, Diagnostic) {
	lentokens := len(*tokens)
	context := ParserContext{}
	statements := make(map[string]ast.Statement)

	for *position < lentokens {
		token := &(*tokens)[*position]
		switch token.Kind {
		case Select:
			if _, ok := statements["select"]; ok {
				return ast.Query{}, *NewError("You already used `SELECT` statement").AddNote("Can't use more than one `SELECT` statement in the same query").WithLocation(token.Location)
			}
			statement, err := ParseSelectStatement(&context, env, tokens, position)
			if err.Message != "" {
				return ast.Query{}, err
			}
			statements["select"] = statement
			context.IsSingleValueQuery = len(context.Aggregations) != 0
		case Where:
			if _, ok := statements["where"]; ok {
				return ast.Query{}, *NewError("You already used `WHERE` statement").AddNote("Can't use more than one `WHERE` statement in the same query").WithLocation(token.Location)
			}
			statement, err := ParseWhereStatement(&context, env, tokens, position)
			if err.Message != "" {
				return ast.Query{}, err
			}
			statements["where"] = statement
		case Group:
			if _, ok := statements["group"]; ok {
				return ast.Query{}, *NewError("You already used `GROUP BY` statement").AddNote("Can't use more than one `GROUP BY` statement in the same query").WithLocation(token.Location)
			}
			statement, err := ParseGroupByStatement(&context, env, tokens, position)
			if err.Message != "" {
				return ast.Query{}, err
			}
			statements["group"] = statement
		case Having:
			if _, ok := statements["having"]; ok {
				return ast.Query{}, *NewError("You already used `HAVING` statement").AddNote("Can't use more than one `HAVING` statement in the same query").WithLocation(token.Location)
			}
			if _, ok := statements["group"]; !ok {
				return ast.Query{}, *NewError("`HAVING` must be used after `GROUP BY` statement").AddNote("`HAVING` statement must be used in a query that has `GROUP BY` statement").WithLocation(token.Location)
			}
			statement, err := ParseHavingStatement(&context, env, tokens, position)
			if err.Message != "" {
				return ast.Query{}, err
			}
			statements["having"] = statement
		case Limit:
			if _, ok := statements["limit"]; ok {
				return ast.Query{}, *NewError("You already used `LIMIT` statement").AddNote("Can't use more than one `LIMIT` statement in the same query").WithLocation(token.Location)
			}
			statement, _ := ParseLimitStatement(tokens, position)
			statements["limit"] = statement
			if *position < lentokens && (*tokens)[*position].Kind == Comma {
				if _, ok := statements["offset"]; ok {
					return ast.Query{}, *NewError("You already used `OFFSET` statement").AddNote("Can't use more than one `OFFSET` statement in the same query").WithLocation(token.Location)
				}
				*position += 1
				if *position >= lentokens || (*tokens)[*position].Kind != Integer {
					return ast.Query{}, *NewError("Expects `OFFSET` amount as Integer value after `,`").AddHelp("Try to add constant Integer after comma").AddNote("`OFFSET` value must be a constant Integer").WithLocation(token.Location)
				}
				count, err := strconv.Atoi((*tokens)[*position].Literal)
				// Report clear error for Integer parsing
				if err != nil {
					return ast.Query{}, *NewError("`OFFSET` integer value is invalid").AddNote(fmt.Sprintf("`OFFSET` value must be between 0 and %d", math.MaxInt)).WithLocation(token.Location)
				}
				*position += 1
				statements["offset"] = &ast.OffsetStatement{Count: count}
			}
		case Offset:
			if _, ok := statements["offset"]; ok {
				return ast.Query{}, *NewError("You already used `OFFSET` statement").AddHelp("Can't use more than one `OFFSET` statement in the same query").WithLocation(token.Location)
			}
			statement, _ := ParseOffsetStatement(tokens, position)
			statements["offset"] = statement
		case Order:
			if _, ok := statements["order"]; ok {
				return ast.Query{}, *NewError("You already used `ORDER BY` statement").AddHelp("Can't use more than one `ORDER BY` statement in the same query").WithLocation(token.Location)
			}
			statement, _ := ParseOrderByStatement(&context, env, tokens, position)
			statements["order"] = statement
		default:
			*position += 1
			break
		}
	}

	// If any aggregation function is used, add Aggregation Functions Node to the GQL Query
	if len(context.Aggregations) != 0 {
		aggregationFunctions := &ast.AggregationsStatement{
			Aggregations: context.Aggregations,
		}
		statements["aggregation"] = aggregationFunctions
	}

	// Remove all selected fields from hidden selection
	hiddenSelections := make([]string, 0)
	for _, selection := range context.HiddenSelections {
		if !contains(context.SelectedFields, selection) {
			hiddenSelections = append(hiddenSelections, selection)
		}
	}

	return ast.Query{
		Select: &ast.GQLQuery{
			Statements:             statements,
			HasAggregationFunction: context.IsSingleValueQuery,
			HasGroupByStatement:    context.HasGroupByStatement,
			HiddenSelections:       hiddenSelections,
		},
	}, Diagnostic{}
}

// nolint:funlen,gocyclo,lll
func ParseSelectStatement(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	*position += 1

	if *position >= len(*tokens) {
		return &ast.SelectStatement{}, *NewError("Incomplete input for select statement").AddHelp("Try select one or more values in the `SELECT` statement").AddNote("Select statements requires at least selecting one value").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	var tableName string
	var fieldsNames []string
	var fieldsValues []ast.Expression
	aliasTable := make(map[string]string)
	isSelectAll := false
	isDistinct := false

	// Check if select has distinct keyword after it
	if (*tokens)[*position].Kind == Distinct {
		isDistinct = true
		*position++
	}

	// Select all option
	if *position < len(*tokens) && (*tokens)[*position].Kind == Star {
		// Consume `*`
		*position++
		isSelectAll = true
	} else {
		for *position < len(*tokens) && (*tokens)[*position].Kind != From {
			expression, err := ParseExpression(context, env, tokens, position)
			if err.Message != "" {
				return nil, err
			}
			exprType := expression.ExprType(env)
			expressionName, _ := GetExpressionName(expression)

			var fieldName string
			if expressionName != "" {
				fieldName = expressionName
			} else {
				fieldName = context.GenerateColumnName()
			}

			// Assert that each selected field is unique
			if contains(fieldsNames, fieldName) {
				return nil, *NewError("Can't select the same field twice").WithLocation(GetSafeLocation(tokens, *position-1))
			}

			// Check for Field name alias
			if *position < len(*tokens) && (*tokens)[*position].Kind == As {
				// Consume `as` keyword
				*position += 1
				aliasNameToken, err := ConsumeKind(*tokens, *position, Symbol)
				if err != nil {
					return nil, *NewError("Expect `identifier` as field alias name").WithLocation(GetSafeLocation(tokens, *position))
				}

				// Register alias name
				aliasName := aliasNameToken.Literal
				if contains(context.SelectedFields, aliasName) || aliasTable[aliasName] != "" {
					return nil, *NewError("You already have field with the same name").AddHelp("Try to use a new unique name for alias").WithLocation(GetSafeLocation(tokens, *position))
				}

				// Consume alias name
				*position += 1

				// Register alias name type
				env.Define(aliasName, exprType)

				context.SelectedFields = append(context.SelectedFields, aliasName)
				aliasTable[fieldName] = aliasName
			}

			// Register field type
			env.Define(fieldName, exprType)

			fieldsNames = append(fieldsNames, fieldName)
			context.SelectedFields = append(context.SelectedFields, fieldName)
			fieldsValues = append(fieldsValues, expression)

			// Consume `,` or break
			if *position < len(*tokens) && (*tokens)[*position].Kind == Comma {
				*position += 1
			} else {
				break
			}
		}
	}

	// Parse optional Form statement
	if *position < len(*tokens) && (*tokens)[*position].Kind == From {
		// Consume `from` keyword
		*position += 1

		tableNameToken, err := ConsumeKind(*tokens, *position, Symbol)
		if err != nil {
			return nil, *NewError("Expect `identifier` as a table name").AddNote("Table name must be an identifier").WithLocation(GetSafeLocation(tokens, *position))
		}

		// Consume table name
		*position += 1

		tableName = tableNameToken.Literal

		if _, ok := ast.TablesFieldsNames[tableName]; !ok {
			return nil, *NewError("Unresolved table name").AddHelp("Check the documentations to see available tables").WithLocation(GetSafeLocation(tokens, *position))
		}

		RegisterCurrentTableFieldsTypes(tableName, *env)
	}

	// Make sure `SELECT *` used with specific table
	if isSelectAll && tableName == "" {
		return nil, *NewError("Expect `FROM` and table name after `SELECT *`").AddNote("Select all must be used with valid table name").WithLocation(GetSafeLocation(tokens, *position))
	}

	// Select input validations
	if !isSelectAll && len(fieldsNames) == 0 {
		return nil, *NewError("Incomplete input for select statement").AddHelp("Try select one or more values in the `SELECT` statement").AddNote("Select statements requires at least selecting one value").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	// If it `select *` make all table fields selectable
	if isSelectAll {
		SelectAllTableFields(
			tableName,
			context.SelectedFields,
			fieldsNames,
			fieldsValues,
		)
	}

	// Type check all selected fields has type registered in type table
	err := TypeCheckSelectedFields(env, tableName, &fieldsNames, tokens, *position)
	if err.Message != "" {
		return nil, err
	}

	return &ast.SelectStatement{
		TableName:    tableName,
		FieldsNames:  fieldsNames,
		FieldsValues: fieldsValues,
		AliasTable:   aliasTable,
		IsDistinct:   isDistinct,
	}, Diagnostic{}
}

// nolint:goconst,lll
func ParseWhereStatement(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	*position++
	if *position >= len(*tokens) {
		return nil, *NewError("Expect expression after `WHERE` keyword").AddHelp("Try to add boolean expression after `WHERE` keyword").AddNote("`WHERE` statement expects expression as condition").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	aggregationsCountBefore := len(context.Aggregations)

	conditionLocation := (*tokens)[*position].Location
	condition, _ := ParseExpression(context, env, tokens, position)
	conditionType := condition.ExprType(env)
	if conditionType.Fmt() != "Boolean" {
		return nil, *NewError(fmt.Sprintf("Expect `WHERE` condition to be type %s but got %s", "Boolean", conditionType)).AddNote("`WHERE` statement condition must be Boolean").WithLocation(conditionLocation)
	}

	aggregationsCountAfter := len(context.Aggregations)
	if aggregationsCountBefore != aggregationsCountAfter {
		return nil, *NewError("Can't use Aggregation functions in `WHERE` statement").AddNote("Aggregation functions must be used after `GROUP BY` statement").AddNote("Aggregation functions evaluated after later after `GROUP BY` statement").WithLocation(conditionLocation)
	}

	return &ast.WhereStatement{Condition: condition}, Diagnostic{}
}

// nolint:lll
func ParseGroupByStatement(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	*position += 1
	if *position >= len(*tokens) || (*tokens)[*position].Kind != By {
		return nil, *NewError("Expect keyword `by` after keyword `group`").AddHelp("Try to use `BY` keyword after `GROUP").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	*position += 1
	if *position >= len(*tokens) || (*tokens)[*position].Kind != Symbol {
		return nil, *NewError("Expect field name after `group by`").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	fieldName := (*tokens)[*position].Literal
	*position += 1

	if !env.Contains(fieldName) {
		return nil, *NewError("Current table not contains field with this name").AddHelp("Check the documentations to see available fields for each tables").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	context.HasGroupByStatement = true
	return &ast.GroupByStatement{FieldName: fieldName}, Diagnostic{}
}

// nolint:lll
func ParseHavingStatement(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	*position += 1
	if *position >= len(*tokens) {
		return nil, *NewError("Expect expression after `HAVING` keyword").AddHelp("Try to add boolean expression after `HAVING` keyword").AddNote("`HAVING` statement expects expression as condition").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	conditionLocation := (*tokens)[*position].Location
	condition, _ := ParseExpression(context, env, tokens, position)
	conditionType := condition.ExprType(env)
	if conditionType.Fmt() != "Boolean" {
		return nil, *NewError(fmt.Sprintf("Expect `HAVING` condition to be type %s but got %s", "Boolean", conditionType)).AddNote("`HAVING` statement condition must be Boolean").WithLocation(conditionLocation)
	}

	return &ast.HavingStatement{Condition: condition}, Diagnostic{}
}

// nolint:lll
func ParseLimitStatement(tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	*position += 1
	if *position >= len(*tokens) || (*tokens)[*position].Kind != Integer {
		return nil, *NewError("Expect number after `LIMIT` keyword").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	count, err := strconv.Atoi((*tokens)[*position].Literal)

	// Report clear error for Integer parsing
	if err != nil {
		return &ast.OffsetStatement{}, *NewError("`OFFSET` integer value is invalid").AddNote(fmt.Sprintf("`LIMIT` value must be between 0 and %d", math.MaxInt)).WithLocation(GetSafeLocation(tokens, *position))
	}

	*position += 1
	return &ast.LimitStatement{Count: count}, Diagnostic{}
}

// nolint:lll
func ParseOffsetStatement(tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	*position += 1
	if *position >= len(*tokens) || (*tokens)[*position].Kind != Integer {
		return &ast.OffsetStatement{}, *NewError("Expect number after `OFFSET` keyword").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	count, err := strconv.Atoi((*tokens)[*position].Literal)

	// Report clear error for Integer parsing
	if err != nil {
		return &ast.OffsetStatement{}, *NewError("`LIMIT` integer value is invalid").AddNote(fmt.Sprintf("`OFFSET` value must be between 0 and %d", math.MaxInt)).WithLocation(GetSafeLocation(tokens, *position))
	}

	*position += 1
	return &ast.OffsetStatement{Count: count}, Diagnostic{}
}

// nolint:lll
func ParseOrderByStatement(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	// Consume `ORDER` keyword
	*position += 1

	if *position >= len(*tokens) || (*tokens)[*position].Kind != By {
		return nil, *NewError("`Expect keyword `BY` after keyword `ORDER").AddHelp("Try to use `BY` keyword after `ORDER").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	// Consume `BY` keyword
	*position += 1

	var arguments []ast.Expression
	var sortingOrders []ast.SortingOrder

	for {
		argument, err := ParseExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}
		arguments = append(arguments, argument)

		order := ast.Ascending
		if *position < len(*tokens) && IsAscOrDesc(&(*tokens)[*position]) {
			if (*tokens)[*position].Kind == Descending {
				order = ast.Descending
			}

			// Consume `ASC or DESC` keyword
			*position += 1
		}

		sortingOrders = append(sortingOrders, order)
		if *position < len(*tokens) && (*tokens)[*position].Kind == Comma {
			// Consume `,` keyword
			*position += 1
		} else {
			break
		}
	}

	return &ast.OrderByStatement{
		Arguments:     arguments,
		SortingOrders: sortingOrders,
	}, Diagnostic{}
}

func ParseExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	aggregationsCountBefore := len(context.Aggregations)
	expression, err := ParseAssignmentExpression(context, env, tokens, position)
	if err.Message != "" {
		return nil, err
	}
	hasAggregations := len(context.Aggregations) != aggregationsCountBefore

	if hasAggregations {
		columnName := context.GenerateColumnName()
		env.Define(columnName, expression.ExprType(env))

		// Register the new aggregation generated field if the this expression is after group by
		if context.HasGroupByStatement && !contains(context.HiddenSelections, columnName) {
			context.HiddenSelections = append(context.HiddenSelections, columnName)
		}

		context.Aggregations[columnName] = ast.AggregateValue{Expression: expression}

		return &ast.SymbolExpression{
			Value: columnName,
		}, Diagnostic{}
	}

	return expression, Diagnostic{}
}

func ParseAssignmentExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseIsNullExpression(context, env, tokens, position)
	if err.Message != "" {
		return nil, err
	}
	if *position < len(*tokens) && (*tokens)[*position].Kind == ColonEqual {
		if expression.Kind() != ast.ExprGlobalVariable {
			return nil, *NewError("Assignment expressions expect global variable name before `:=`").WithLocation((*tokens)[*position].Location)
		}

		expr := expression.(*ast.GlobalVariableExpression)
		variableName := expr.Name

		// Consume `:=` operator
		*position += 1

		value, err := ParseIsNullExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}
		env.DefineGlobal(variableName, value.ExprType(env))

		return &ast.AssignmentExpression{
			Symbol: variableName,
			Value:  value,
		}, Diagnostic{}
	}

	return expression, Diagnostic{}
}

func ParseIsNullExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseInExpression(context, env, tokens, position)
	if err.Message != "" {
		return nil, err
	}
	if *position < len(*tokens) && (*tokens)[*position].Kind == Is {
		isLocation := (*tokens)[*position].Location

		// Consume `IS` keyword
		*position += 1

		hasNotKeyword := false
		if *position < len(*tokens) && (*tokens)[*position].Kind == Not {
			// Consume `NOT` keyword
			*position++
			hasNotKeyword = true
		}

		if *position < len(*tokens) && (*tokens)[*position].Kind == Null {
			// Consume `Null` keyword
			*position += 1

			return &ast.IsNullExpression{
				Argument: expression,
				HasNot:   hasNotKeyword,
			}, Diagnostic{}
		}

		return &ast.IsNullExpression{}, *NewError("Expects `NULL` Keyword after `IS` or `IS NOT`").WithLocation(isLocation)
	}

	return expression, Diagnostic{}
}

// nolint:lll
func ParseInExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseBetweenExpression(context, env, tokens, position)
	if err.Message != "" {
		return nil, err
	}

	hasNotKeyword := false
	if *position < len(*tokens) && (*tokens)[*position].Kind == Not {
		*position += 1
		hasNotKeyword = true
	}

	if *position < len(*tokens) && (*tokens)[*position].Kind == In {
		inLocation := (*tokens)[*position].Location

		// Consume `IN` keyword
		*position += 1

		if _, err := ConsumeKind(*tokens, *position, LeftParen); err != nil {
			return nil, *NewError("Expects values between `(` and `)` after `IN` keyword").WithLocation(inLocation)
		}

		values, err := ParseArgumentsExpressions(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}

		if len(values) == 0 {
			return &ast.BooleanExpression{IsTrue: hasNotKeyword}, Diagnostic{}
		}

		valuesTypeResult := CheckAllValuesAreSameType(env, values)
		if valuesTypeResult == nil {
			return nil, *NewError("Expects values between `(` and `)` to have the same type").WithLocation(inLocation)
		}

		// Check that argument and values has the same type
		valuesType := valuesTypeResult
		if valuesType.Fmt() != "Any" && expression.ExprType(env) != valuesType {
			return nil, *NewError("Argument and Values of In Expression must have the same type").WithLocation(inLocation)
		}

		return &ast.InExpression{
			Argument:      expression,
			Values:        values,
			ValuesType:    valuesType,
			HasNotKeyword: hasNotKeyword,
		}, Diagnostic{}
	}

	if hasNotKeyword {
		return nil, *NewError("Expects `IN` expression after this `NOT` keyword").AddHelp("Try to use `IN` expression after NOT keyword").AddHelp("Try to remove `NOT` keyword").AddNote("Expect to see `NOT` then `IN` keyword with a list of values").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	return expression, Diagnostic{}
}

// nolint:lll
func ParseBetweenExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseLogicalOrExpression(context, env, tokens, position)
	if err.Message != "" {
		return nil, err
	}

	if *position < len(*tokens) && (*tokens)[*position].Kind == Between {
		betweenLocation := (*tokens)[*position].Location

		// Consume `BETWEEN` keyword
		*position += 1

		if *position >= len(*tokens) {
			return nil, *NewError("`BETWEEN` keyword expects two range after it").WithLocation(betweenLocation)
		}

		argumentType := expression.ExprType(env)
		rangeStart, err := ParseLogicalOrExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}

		if *position >= len(*tokens) || (*tokens)[*position].Kind != DotDot {
			return nil, *NewError("Expect `..` after `BETWEEN` range start").WithLocation(betweenLocation)
		}

		// Consume `..` token
		*position += 1
		rangeEnd, err := ParseLogicalOrExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}

		if argumentType != rangeStart.ExprType(env) || argumentType != rangeEnd.ExprType(env) {
			return nil, *NewError(fmt.Sprintf("Expect `BETWEEN` argument, range start and end to has same type but got %s, %s and %s", argumentType, rangeStart.ExprType(env), rangeEnd.ExprType(env))).AddHelp("Try to make sure all of them has same type").WithLocation(betweenLocation)
		}

		return &ast.BetweenExpression{
			Value:      expression,
			RangeStart: rangeStart,
			RangeEnd:   rangeEnd,
		}, Diagnostic{}
	}

	return expression, Diagnostic{}
}

func ParseLogicalOrExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseLogicalAndExpression(context, env, tokens, position)
	if err.Message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression
	for *position < len(*tokens) && (*tokens)[*position].Kind == LogicalOr {
		*position += 1
		if lhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position-2].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		rhs, _ := ParseLogicalAndExpression(context, env, tokens, position)
		if rhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		lhs = &ast.LogicalExpression{
			Left:     lhs,
			Operator: ast.LOOr,
			Right:    rhs,
		}
	}

	return lhs, Diagnostic{}
}

func ParseLogicalAndExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseBitwiseOrExpression(context, env, tokens, position)
	if err.Message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression
	for *position < len(*tokens) && (*tokens)[*position].Kind == LogicalAnd {
		*position += 1
		if lhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position-2].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		rhs, _ := ParseBitwiseOrExpression(context, env, tokens, position)
		if rhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		lhs = &ast.LogicalExpression{
			Left:     lhs,
			Operator: ast.LOAnd,
			Right:    rhs,
		}
	}

	return lhs, Diagnostic{}
}

func ParseBitwiseOrExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseLogicalXorExpression(context, env, tokens, position)
	if err.Message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression

	operator := &(*tokens)[*position]
	if operator.Kind == BitwiseOr {
		*position += 1

		if lhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position-2].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		rhs, err := ParseLogicalXorExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}
		if rhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		return &ast.BitwiseExpression{
			Left:     lhs,
			Operator: ast.BOOr,
			Right:    rhs,
		}, Diagnostic{}
	}

	return lhs, Diagnostic{}
}

func ParseLogicalXorExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseBitwiseAndExpression(context, env, tokens, position)
	if err.Message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression
	for *position < len(*tokens) && (*tokens)[*position].Kind == LogicalXor {
		*position += 1
		if lhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position-2].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		rhs, err := ParseBitwiseAndExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}
		if rhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		lhs = &ast.LogicalExpression{
			Left:     lhs,
			Operator: ast.LOXor,
			Right:    rhs,
		}
	}

	return lhs, Diagnostic{}
}

func ParseBitwiseAndExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseEqualityExpression(context, env, tokens, position)
	if err.Message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression

	if *position < len(*tokens) && (*tokens)[*position].Kind == BitwiseAnd {
		*position += 1
		if lhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position-2].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		rhs, err := ParseEqualityExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}
		if rhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		lhs = &ast.BitwiseExpression{
			Left:     lhs,
			Operator: ast.BOAnd,
			Right:    rhs,
		}
	}

	return lhs, Diagnostic{}
}

// nolint:gomnd
func ParseEqualityExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseComparisonExpression(context, env, tokens, position)
	if err.Message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression

	operator := &(*tokens)[*position]
	if operator.Kind == Equal || operator.Kind == BangEqual {
		*position += 1
		var comparisonOperator ast.ComparisonOperator
		if operator.Kind == Equal {
			comparisonOperator = ast.COEqual
		} else {
			comparisonOperator = ast.CONotEqual
		}

		rhs, err := ParseComparisonExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}

		switch AreTypesEquals(env, lhs, rhs) {
		case Equals{}:
			// do nothing
		case RightSideCasted{expr: expression}:
			rhs = expression
		case LeftSideCasted{expr: expression}:
			lhs = expression
		case NotEqualAndCantImplicitCast{}:
			lhsType := lhs.ExprType(env)
			rhsType := rhs.ExprType(env)
			diagnostic := *NewError(fmt.Sprintf(
				"Can't compare values of different types `%s` and `%s`",
				lhsType,
				rhsType,
			)).WithLocation(GetSafeLocation(tokens, *position-2))

			// Provides help messages if use compare null to non null value
			if lhsType.IsNull() || rhsType.IsNull() {
				return nil, *diagnostic.AddHelp("Try to use `IS NULL expr` expression").AddHelp("Try to use `ISNULL(expr)` function")
			}

			return nil, diagnostic
		case Error{}:
			return nil, *NewError("").WithLocation(GetSafeLocation(tokens, *position-2))
		}

		return &ast.ComparisonExpression{
			Left:     lhs,
			Operator: comparisonOperator,
			Right:    rhs,
		}, Diagnostic{}
	}

	return lhs, Diagnostic{}
}

// nolint:gocyclo,gomnd
func ParseComparisonExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseBitwiseShiftExpression(context, env, tokens, position)
	if err.Message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression
	if IsComparisonOperator(&(*tokens)[*position]) {
		operator := &(*tokens)[*position]
		*position += 1
		var comparisonOperator ast.ComparisonOperator
		switch operator.Kind {
		case Greater:
			comparisonOperator = ast.COGreater
		case GreaterEqual:
			comparisonOperator = ast.COGreaterEqual
		case Less:
			comparisonOperator = ast.COLess
		case LessEqual:
			comparisonOperator = ast.COLessEqual
		default:
			comparisonOperator = ast.CONullSafeEqual
		}

		rhs, err := ParseBitwiseShiftExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}

		switch AreTypesEquals(env, lhs, rhs) {
		case Equals{}:
			// do nothing
		case RightSideCasted{expr: expression}:
			rhs = expression
		case LeftSideCasted{expr: expression}:
			lhs = expression
		case NotEqualAndCantImplicitCast{}:
			lhsType := lhs.ExprType(env)
			rhsType := rhs.ExprType(env)
			diagnostic := *NewError(fmt.Sprintf(
				"Can't compare values of different types `%s` and `%s`",
				lhsType,
				rhsType,
			)).WithLocation(GetSafeLocation(tokens, *position-2))

			// Provides help messages if use compare null to non null value
			if lhsType.IsNull() || rhsType.IsNull() {
				return nil, *diagnostic.AddHelp("Try to use `IS NULL expr` expression").AddHelp("Try to use `ISNULL(expr)` function")
			}

			return nil, diagnostic
		case Error{}:
			return nil, *NewError("").WithLocation(GetSafeLocation(tokens, *position-2))
		}

		return &ast.ComparisonExpression{
			Left:     lhs,
			Operator: comparisonOperator,
			Right:    rhs,
		}, Diagnostic{}
	}

	return lhs, Diagnostic{}
}

// nolint:goconst,gomnd,lll
func ParseBitwiseShiftExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	lhs, err := ParseTermExpression(context, env, tokens, position)
	if err.Message != "" {
		return nil, err
	}

	for *position < len(*tokens) && IsBitwiseShiftOperator(&(*tokens)[*position]) {
		operator := &(*tokens)[*position]
		*position += 1
		var bitwiseOperator ast.BitwiseOperator
		if operator.Kind == BitwiseRightShift {
			bitwiseOperator = ast.BORightShift
		} else {
			bitwiseOperator = ast.BOLeftShift
		}

		rhs, err := ParseTermExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}

		// Make sure right and left hand side types are numbers
		if rhs.ExprType(env).IsInt() && rhs.ExprType(env) != lhs.ExprType(env) {
			return nil, *NewError(fmt.Sprintf(
				"Bitwise operators require number types but got `%s` and `%s`",
				lhs.ExprType(env),
				rhs.ExprType(env),
			)).WithLocation(GetSafeLocation(tokens, *position-2))
		}

		lhs = &ast.BitwiseExpression{
			Left:     lhs,
			Operator: bitwiseOperator,
			Right:    rhs,
		}
	}

	return lhs, Diagnostic{}
}

func ParseTermExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	lhs, err := ParseFactorExpression(context, env, tokens, position)
	if err.Message != "" {
		return nil, err
	}

	for *position < len(*tokens) && IsTermOperator(&(*tokens)[*position]) {
		operator := &(*tokens)[*position]
		*position += 1
		var mathOperator ast.ArithmeticOperator
		if operator.Kind == Plus {
			mathOperator = ast.AOPlus
		} else {
			mathOperator = ast.AOMinus
		}

		rhs, err := ParseFactorExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}

		lhsType := lhs.ExprType(env)
		rhsType := rhs.ExprType(env)

		// Make sure right and left hand side types are numbers
		if lhsType.IsNumber() && rhsType.IsNumber() {
			lhs = &ast.ArithmeticExpression{
				Left:     lhs,
				Operator: mathOperator,
				Right:    rhs,
			}

			continue
		}

		if mathOperator == ast.AOPlus {
			return nil, *NewError(fmt.Sprintf(
				"Math operators `+` both sides to be number types but got `%s` and `%s`",
				lhsType,
				rhsType,
			)).AddHelp("You can use `CONCAT(Any, Any, ...Any)` function to concatenate values with different types").WithLocation(operator.Location)
		}

		return nil, *NewError(fmt.Sprintf(
			"Math operators require number types but got `%s` and `%s`",
			lhsType,
			rhsType,
		)).WithLocation(operator.Location)
	}

	return lhs, Diagnostic{}
}

// nolint:gomnd
func ParseFactorExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseLikeExpression(context, env, tokens, position)
	if err.Message != "" || *position >= len(*tokens) {
		return expression, err
	}
	lhs := expression
	for *position < len(*tokens) && IsFactorOperator(&(*tokens)[*position]) {
		operator := &(*tokens)[*position]
		*position += 1

		var factorOperator ast.ArithmeticOperator
		switch operator.Kind {
		case Star:
			factorOperator = ast.AOStar
		case Slash:
			factorOperator = ast.AOSlash
		default:
			factorOperator = ast.AOModulus
		}

		rhs, err := ParseLikeExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}

		lhsType := lhs.ExprType(env)
		rhsType := rhs.ExprType(env)

		// Make sure right and left hand side types are numbers
		if lhsType.IsNumber() && rhsType.IsNumber() {
			lhs = &ast.ArithmeticExpression{
				Left:     lhs,
				Operator: factorOperator,
				Right:    rhs,
			}
			continue
		}

		return nil, *NewError(fmt.Sprintf(
			"Math operators require number types but got `%s` and `%s`",
			lhsType,
			rhsType,
		)).WithLocation(GetSafeLocation(tokens, *position-2))
	}

	return lhs, Diagnostic{}
}

func ParseLikeExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseGlobExpression(context, env, tokens, position)
	if err.Message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression
	if (*tokens)[*position].Kind == Like {
		location := (*tokens)[*position].Location
		*position += 1

		if !lhs.ExprType(env).IsText() {
			return nil, *NewError(fmt.Sprintf("Expect `LIKE` left hand side to be `TEXT` but got %s", lhs.ExprType(env))).WithLocation(location)
		}

		pattern, err := ParseGlobExpression(context, env, tokens, position)
		if err.Message != "" || !pattern.ExprType(env).IsText() {
			return nil, *NewError(fmt.Sprintf("Expect `LIKE` right hand side to be `TEXT` but got %s", lhs.ExprType(env))).WithLocation(location)
		}

		return &ast.LikeExpression{
			Input:   lhs,
			Pattern: pattern,
		}, Diagnostic{}
	}

	return lhs, Diagnostic{}
}

func ParseGlobExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseUnaryExpression(context, env, tokens, position)
	if err.Message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression
	if (*tokens)[*position].Kind == Glob {
		location := (*tokens)[*position].Location
		*position += 1

		if !lhs.ExprType(env).IsText() {
			return nil, *NewError(fmt.Sprintf("Expect `GLOB` left hand side to be `TEXT` but got %s", lhs.ExprType(env))).WithLocation(location)
		}

		pattern, err := ParseUnaryExpression(context, env, tokens, position)
		if err.Message != "" || !pattern.ExprType(env).IsText() {
			return nil, *NewError(fmt.Sprintf("Expect `GLOB` right hand side to be `TEXT` but got %s", pattern.ExprType(env))).WithLocation(location)
		}

		return &ast.GlobExpression{
			Input:   lhs,
			Pattern: pattern,
		}, Diagnostic{}
	}

	return lhs, Diagnostic{}
}

func ParseUnaryExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	if *position < len(*tokens) && IsPrefixUnaryOperator(&(*tokens)[*position]) {
		var op ast.PrefixUnaryOperator
		if (*tokens)[*position].Kind == Bang {
			op = ast.POBang
		} else {
			op = ast.POMinus
		}

		*position += 1

		rhs, err := ParseUnaryExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}
		rhsType := rhs.ExprType(env)
		if op == ast.POBang && rhsType.Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				GetSafeLocation(tokens, *position-1),
				ast.Boolean{},
				rhsType,
			)
		}

		if op == ast.POMinus && rhsType.Fmt() != "Integer" {
			return nil, TypeMismatchError(
				GetSafeLocation(tokens, *position-1),
				ast.Integer{},
				rhsType,
			)
		}

		return &ast.PrefixUnary{Right: rhs, Op: op}, Diagnostic{}
	}
	return ParseFunctionCallExpression(context, env, tokens, position)
}

// nolint:funlen,lll
func ParseFunctionCallExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParsePrimaryExpression(context, env, tokens, position)
	if err.Message != "" {
		return nil, err
	}

	if *position < len(*tokens) && (*tokens)[*position].Kind == LeftParen {
		symbolExpression, ok := expression.(*ast.SymbolExpression)
		functionNameLocation := GetSafeLocation(tokens, *position)

		// Make sure function name is SymbolExpression
		if !ok {
			return nil, *NewError("Function name must be an identifier").WithLocation(functionNameLocation)
		}

		functionName := symbolExpression.Value

		// Check if this function is a Standard library functions
		if _, ok := ast.Functions[functionName]; ok {
			arguments, err := ParseArgumentsExpressions(context, env, tokens, position)
			if err.Message != "" {
				return nil, err
			}

			prototype := ast.Prototypes[functionName]
			parameters := prototype.Parameters
			return_type := prototype.Result

			_, err = CheckFunctionCallArguments(
				env,
				&arguments,
				&parameters,
				functionName,
				functionNameLocation,
			)
			if err.Message != "" {
				return nil, err
			}

			// Register function name with return type
			env.Define(functionName, return_type)

			return &ast.CallExpression{
				FunctionName:  functionName,
				Arguments:     arguments,
				IsAggregation: false,
			}, Diagnostic{}
		}

		// Check if this function is an Aggregation functions
		if _, ok := ast.Aggregations[functionName]; ok {
			arguments, err := ParseArgumentsExpressions(context, env, tokens, position)
			if err.Message != "" {
				return nil, err
			}

			prototype := ast.AggregationsProtos[functionName]
			parameters := []ast.DataType{prototype.Parameter}
			returnType := prototype.Result

			_, err = CheckFunctionCallArguments(
				env,
				&arguments,
				&parameters,
				functionName,
				functionNameLocation,
			)
			if err.Message != "" {
				return nil, err
			}

			argumentResult, argueErr := GetExpressionName(arguments[0])
			if argueErr.Error() != "" {
				return nil, *NewError("Invalid Aggregation function argument").AddHelp("Try to use field name as Aggregation function argument").AddNote("Aggregation function accept field name as argument").WithLocation(functionNameLocation)
			}

			argument := argumentResult
			columnName := context.GenerateColumnName()

			context.HiddenSelections = append(context.HiddenSelections, columnName)

			// Register aggregation generated name with return type
			env.Define(columnName, returnType)

			context.Aggregations[columnName] = ast.AggregateValue{
				Function: struct {
					Name string
					Arg  string
				}{
					Name: functionName,
					Arg:  argument,
				},
			}

			return &ast.SymbolExpression{
				Value: columnName,
			}, Diagnostic{}
		}

		// Report that this function name is not standard or aggregation
		return nil, *NewError("No such function name").AddHelp(fmt.Sprintf("Function `%s` is not an Aggregation or Standard library function name", functionName)).WithLocation(functionNameLocation)
	}
	return expression, Diagnostic{}
}

// nolint:lll
func ParseArgumentsExpressions(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) ([]ast.Expression, Diagnostic) {
	var arguments []ast.Expression
	if _, err := ConsumeKind(*tokens, *position, LeftParen); err == nil {
		*position += 1

		for (*tokens)[*position].Kind != RightParen {
			argument, err := ParseExpression(context, env, tokens, position)
			if err.Message != "" {
				return nil, err
			}

			argumentLiteral, erral := GetExpressionName(argument)
			if erral == nil {
				literal := argumentLiteral
				context.HiddenSelections = append(context.HiddenSelections, literal)
			}

			arguments = append(arguments, argument)

			if (*tokens)[*position].Kind == Comma {
				*position += 1
			} else {
				break
			}
		}

		if _, err := ConsumeKind(*tokens, *position, RightParen); err != nil {
			return nil, *NewError("Expect `)` after function call arguments").AddHelp("Try to add ')' at the end of function call, after arguments").WithLocation(GetSafeLocation(tokens, *position))
		}

		*position += 1
	}
	return arguments, Diagnostic{}
}

// nolint:lll
func ParsePrimaryExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	if *position >= len(*tokens) {
		return nil, UnExpectedExpressionError(tokens, position)
	}
	switch (*tokens)[*position].Kind {
	case String:
		*position += 1
		return &ast.StringExpression{
			Value:     (*tokens)[*position-1].Literal,
			ValueType: ast.StringValueText,
		}, Diagnostic{}
	case Symbol:
		value := (*tokens)[*position].Literal
		*position += 1
		if !contains(context.SelectedFields, value) {
			context.HiddenSelections = append(context.HiddenSelections, value)
		}
		return &ast.SymbolExpression{Value: value}, Diagnostic{}
	case GlobalVariable:
		name := (*tokens)[*position].Literal
		*position += 1
		return &ast.GlobalVariableExpression{Name: name}, Diagnostic{}
	case Integer:
		if integer, err := strconv.ParseInt((*tokens)[*position].Literal, 10, 64); err == nil {
			*position += 1
			value := ast.IntegerValue{Value: integer}
			return &ast.NumberExpression{Value: value}, Diagnostic{}
		}
		return nil, *NewError("Too big Integer value").AddHelp(fmt.Sprintf("Integer value must be between %d and %d", math.MinInt64, math.MaxInt64)).WithLocation((*tokens)[*position].Location)
	case Float:
		if float, err := strconv.ParseFloat((*tokens)[*position].Literal, 64); err == nil {
			*position += 1
			value := ast.FloatValue{Value: float}
			return &ast.NumberExpression{Value: value}, Diagnostic{}
		}
		return nil, *NewError("Too big Float value").AddHelp("Try to use smaller value").AddNote(fmt.Sprintf("Float value must be between %f and %f", -math.MaxFloat64, math.MaxFloat64)).WithLocation((*tokens)[*position].Location)
	case True:
		*position += 1
		return &ast.BooleanExpression{IsTrue: true}, Diagnostic{}
	case False:
		*position += 1
		return &ast.BooleanExpression{IsTrue: false}, Diagnostic{}
	case Null:
		*position += 1
		return &ast.NullExpression{}, Diagnostic{}
	case LeftParen:
		return ParseGroupExpression(context, env, tokens, position)
	case Case:
		return ParseCaseExpression(context, env, tokens, position)
	default:
		return nil, UnExpectedExpressionError(tokens, position)
	}
}

// nolint:lll
func ParseGroupExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	*position += 1

	expression, err := ParseExpression(context, env, tokens, position)
	if err.Message != "" {
		return nil, err
	}

	if (*tokens)[*position].Kind != RightParen {
		return nil, *NewError("Expect `)` to end group expression").WithLocation(GetSafeLocation(tokens, *position)).AddHelp("Try to add ')' at the end of group expression")
	}

	*position += 1
	return expression, Diagnostic{}
}

// nolint:funlen,gocyclo,lll
func ParseCaseExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	var conditions []ast.Expression
	var values []ast.Expression
	var defaultValue ast.Expression

	// Consume `case` keyword
	caseLocation := (*tokens)[*position].Location
	*position += 1

	hasElseBranch := false

	for *position < len(*tokens) && (*tokens)[*position].Kind != End {
		// Else branch
		if (*tokens)[*position].Kind == Else {
			if hasElseBranch {
				return nil, *NewError("This `CASE` expression already has else branch").AddNote("`CASE` expression can has only one `ELSE` branch").WithLocation(GetSafeLocation(tokens, *position))
			}

			// Consume `ELSE` keyword
			*position += 1

			defaultValueExpr, err := ParseExpression(context, env, tokens, position)
			if err.Message != "" {
				return nil, err
			}
			defaultValue = defaultValueExpr
			hasElseBranch = true
			continue
		}

		// Check if current token kind is `WHEN` keyword
		if _, err := ConsumeKind(*tokens, *position, When); err != nil {
			return nil, *NewError("Expect `when` before case condition").AddHelp("Try to add `WHEN` keyword before any condition").WithLocation(GetSafeLocation(tokens, *position))
		}

		// Consume when keyword
		*position += 1

		condition, err := ParseExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}
		if condition.ExprType(env).Fmt() != "Boolean" {
			return nil, *NewError("Case condition must be a boolean type").WithLocation(GetSafeLocation(tokens, *position))
		}

		conditions = append(conditions, condition)

		if _, err := ConsumeKind(*tokens, *position, Then); err != nil {
			return nil, *NewError("Expect `THEN` after case condition").WithLocation(GetSafeLocation(tokens, *position))
		}

		// Consume then keyword
		*position += 1

		expression, err := ParseExpression(context, env, tokens, position)
		if err.Message != "" {
			return nil, err
		}
		values = append(values, expression)
	}

	// Make sure case expression has at least else branch
	if len(conditions) == 0 && !hasElseBranch {
		return nil, *NewError("Case expression must has at least else branch").WithLocation(GetSafeLocation(tokens, *position))
	}

	// Make sure case expression end with END keyword
	if *position >= len(*tokens) || (*tokens)[*position].Kind != End {
		return nil, *NewError("Expect `END` after case branches").WithLocation(GetSafeLocation(tokens, *position))
	}

	// Consume end
	*position += 1

	// Make sure this case expression has else branch
	if !hasElseBranch {
		return nil, *NewError("Case expression must has else branch").WithLocation(GetSafeLocation(tokens, *position))
	}

	// Assert that all values has the same type
	valuesType := values[0].ExprType(env)
	for i, value := range values[1:] {
		if valuesType != value.ExprType(env) {
			return nil, *NewError(fmt.Sprintf("Case value in branch %d has different type than the last branch", i+1)).AddNote("All values in `CASE` expression must has the same Type").WithLocation(caseLocation)
		}
	}

	return &ast.CaseExpression{
		Conditions:   conditions,
		Values:       values,
		DefaultValue: defaultValue,
		ValuesType:   valuesType,
	}, Diagnostic{}
}

// nolint:gocyclo,lll
func CheckFunctionCallArguments(env *ast.Environment, arguments *[]ast.Expression, parameters *[]ast.DataType, functionName string, location Location) (ast.Expression, Diagnostic) {
	parametersLen := len(*parameters)
	argumentsLen := len(*arguments)

	hasOptionalParameter := false
	hasVarargsParameter := false
	if len(*parameters) != 0 {
		lastParameter := (*parameters)[len(*parameters)-1]
		hasOptionalParameter = lastParameter.IsOptional()
		hasVarargsParameter = lastParameter.IsVarargs()
	}

	// Has Optional parameter type at the end
	if hasOptionalParameter {
		// If function last parameter is optional make sure it at least has
		if argumentsLen < parametersLen-1 {
			return nil, *NewError(fmt.Sprintf("Function `%s` expects at least `%d` arguments but got `%d`", functionName, parametersLen-1, argumentsLen)).WithLocation(location)
		}
		// Make sure function with optional parameter not called with too much arguments
		if argumentsLen > parametersLen {
			return nil, *NewError(fmt.Sprintf("Function `%s` expects at most `%d` arguments but got `%d`", functionName, parametersLen, argumentsLen)).WithLocation(location)
		}
	} else if hasVarargsParameter {
		// If function last parameter is optional make sure it at least has
		if argumentsLen < parametersLen-1 {
			return nil, *NewError(fmt.Sprintf("Function `%s` expects at least `%d` arguments but got `%d`", functionName, parametersLen-1, argumentsLen)).WithLocation(location)
		}
	} else if argumentsLen != parametersLen {
		return nil, *NewError(fmt.Sprintf("Function `%s` expects `%d` arguments but got `%d`", functionName, parametersLen, argumentsLen)).WithLocation(location)
	}

	lastRequiredParameterIndex := parametersLen
	if hasOptionalParameter || hasVarargsParameter {
		lastRequiredParameterIndex -= 1
	}

	// Check each argument vs parameter type
	for index := 0; index < lastRequiredParameterIndex; index++ {
		parameterType := (*parameters)[index]
		argument := (*arguments)[index]

		switch IsExpressionTypeEquals(env, argument, parameterType) {
		case Equals{}:
			// do nothing
		case RightSideCasted{expr: argument}:
			// ================================= ???
			(*arguments)[index] = argument
		case LeftSideCasted{expr: argument}:
			(*arguments)[index] = argument
		case NotEqualAndCantImplicitCast{}:
			argumentType := argument.ExprType(env)
			return nil, *NewError(fmt.Sprintf("Function `%s` argument number %d with type `%s` don't match expected type `%s`", functionName, index, argumentType, parameterType)).WithLocation(location)
		case Error{}:
			return nil, *NewError("")
		}
	}

	if hasOptionalParameter || hasVarargsParameter {
		lastParameterType := (*parameters)[lastRequiredParameterIndex]

		for index := lastRequiredParameterIndex; index < argumentsLen; index++ {
			argument := (*arguments)[index]
			switch IsExpressionTypeEquals(env, argument, lastParameterType) {
			case Equals{}:
				// do nothing
			case RightSideCasted{expr: argument}:
				// ================================= ???
				(*arguments)[index] = argument
			case LeftSideCasted{expr: argument}:
				(*arguments)[index] = argument
			case NotEqualAndCantImplicitCast{}:
				argumentType := (*arguments)[index].ExprType(env)
				if !lastParameterType.Equal(argumentType) {
					return nil, *NewError(fmt.Sprintf("Function `%s` argument number %d with type `%s` don't match expected type `%s`", functionName, index, argumentType, lastParameterType)).WithLocation(location)
				}
			case Error{}:
				return nil, *NewError("")
			}
		}
	}

	return nil, Diagnostic{}
}

func TypeCheckSelectedFields(env *ast.Environment, tableName string, fieldNames *[]string, tokens *[]Token, position int) Diagnostic {
	for _, fieldName := range *fieldNames {
		if dataType, err := env.ResolveType(fieldName); err != nil {
			if dataType.IsUndefined() {
				return *NewError(fmt.Sprintf("No field with name `%s`", fieldName)).WithLocation(GetSafeLocation(tokens, position))
			}
			continue
		}
		return *NewError(fmt.Sprintf("Table %s has no field with name %s", tableName, fieldName)).WithLocation(GetSafeLocation(tokens, position))
	}

	return Diagnostic{}
}

func UnExpectedStatementError(tokens *[]Token, position *int) Diagnostic {
	token := (*tokens)[*position]
	location := token.Location

	if location.Start == 0 {
		return *NewError("Unexpected statement").AddHelp("Expect query to start with `SELECT` or `SET` keyword").WithLocation(location)
	}

	return *NewError("Unexpected statement").WithLocation(location)
}

// nolint:gocyclo
func UnExpectedExpressionError(tokens *[]Token, position *int) Diagnostic {
	location := GetSafeLocation(tokens, *position)

	if *position == 0 || *position >= len(*tokens) {
		return *NewError("Can't complete parsing this expression").WithLocation(location)
	}

	current := &(*tokens)[*position]
	previous := &(*tokens)[*position-1]

	if current.Kind == Ascending || current.Kind == Descending {
		return *NewError("`ASC` and `DESC` must be used in `ORDER BY` statement").WithLocation(location)
	}

	if previous.Kind == Equal && current.Kind == Equal {
		return *NewError("Unexpected `==`, Just use `=` to check equality").AddHelp("Try to remove the extra `=`").WithLocation(location)
	}

	if previous.Kind == Greater && current.Kind == Equal {
		return *NewError("Unexpected `> =`, do you mean `>=`?").AddHelp("Try to remove space between `> =`").WithLocation(location)
	}

	if previous.Kind == Less && current.Kind == Equal {
		return *NewError("Unexpected `< =`, do you mean `<=`?").AddHelp("Try to remove space between `< =`").WithLocation(location)
	}

	if previous.Kind == Greater && current.Kind == Greater {
		return *NewError("Unexpected `> >`, do you mean `>>`?").AddHelp("Try to remove space between `> >`").WithLocation(location)
	}

	if previous.Kind == Less && current.Kind == Less {
		return *NewError("Unexpected `< <`, do you mean `<<`?").AddHelp("Try to remove space between `< <`").WithLocation(location)
	}

	if previous.Kind == Less && current.Kind == Greater {
		return *NewError("Unexpected `< >`, do you mean `<>`?").AddHelp("Try to remove space between `< >`").WithLocation(location)
	}

	return *NewError("Can't complete parsing this expression").WithLocation(location)
}

// nolint:lll
func UnExpectedContentAfterCorrectStatement(statementName *string, tokens *[]Token, position *int) Diagnostic {
	errorMessage := fmt.Sprintf("Unexpected content after the end of `%s` statement", strings.ToUpper(*statementName))
	locationOfExtraContent := Location{
		Start: (*tokens)[*position].Location.Start,
		End:   (*tokens)[len(*tokens)-1].Location.End,
	}

	return *NewError(errorMessage).AddHelp("Try to check if statement keyword is missing").AddHelp("Try remove un expected extra content").WithLocation(locationOfExtraContent)
}

func GetExpressionName(expression ast.Expression) (string, error) {
	if symbol, ok := expression.(*ast.SymbolExpression); ok {
		return symbol.Value, nil
	}
	if variable, ok := expression.(*ast.GlobalVariableExpression); ok {
		return variable.Name, nil
	}
	return "", fmt.Errorf("unsupported expression type")
}

func RegisterCurrentTableFieldsTypes(tableName string, symbolTable ast.Environment) {
	tableFieldsNames := ast.TablesFieldsNames[tableName]
	for _, fieldName := range tableFieldsNames {
		fieldType := ast.TablesFieldsTypes[fieldName]
		symbolTable.Define(fieldName, fieldType)
	}
}

func SelectAllTableFields(tableName string, selectedFields, fieldsNames []string, fieldsValues []ast.Expression) {
	if tableFields, ok := ast.TablesFieldsNames[tableName]; ok {
		for _, field := range tableFields {
			if !contains(fieldsNames, field) {
				fieldsNames = append(fieldsNames, field)
				selectedFields = append(selectedFields, field)
				literalExpr := &ast.SymbolExpression{
					Value: field,
				}
				fieldsValues = append(fieldsValues, literalExpr)
			}
		}
	}
}

func ConsumeKind(tokens []Token, position int, kind TokenKind) (*Token, error) {
	if position < len(tokens) && tokens[position].Kind == kind {
		return &tokens[position], nil
	}
	return nil, fmt.Errorf("error")
}

func GetSafeLocation(tokens *[]Token, position int) Location {
	if position < len(*tokens) {
		return (*tokens)[position].Location
	}
	return (*tokens)[len(*tokens)-1].Location
}

func IsAssignmentOperator(token *Token) bool {
	return token.Kind == Equal || token.Kind == ColonEqual
}

func IsTermOperator(token *Token) bool {
	return token.Kind == Plus || token.Kind == Minus
}

func IsBitwiseShiftOperator(token *Token) bool {
	return token.Kind == BitwiseLeftShift || token.Kind == BitwiseRightShift
}

func IsPrefixUnaryOperator(token *Token) bool {
	return token.Kind == Bang || token.Kind == Minus
}

func IsComparisonOperator(token *Token) bool {
	return token.Kind == Greater ||
		token.Kind == GreaterEqual ||
		token.Kind == Less ||
		token.Kind == LessEqual ||
		token.Kind == NullSafeEqual
}

func IsFactorOperator(token *Token) bool {
	return token.Kind == Star ||
		token.Kind == Slash ||
		token.Kind == Percentage
}

func IsAscOrDesc(token *Token) bool {
	return token.Kind == Ascending || token.Kind == Descending
}

func TypeMismatchError(location Location, expected, actual ast.DataType) Diagnostic {
	return *NewError(fmt.Sprintf("Type mismatch expected `%s`, got `%s`", expected.Fmt(), actual.Fmt())).WithLocation(location)
}

func contains(arr []string, items ...string) bool {
	for _, item := range items {
		for _, a := range arr {
			if item == a {
				return true
			}
		}
	}

	return false
}
