package parser

import (
	"fmt"
	"math"
	"strconv"

	"github.com/ggql/ggql/ast"
)

func ParserGql(tokens []Token, env *ast.Environment) (ast.Query, Diagnostic) {
	position := 0
	first_token := tokens[position]
	var query_result ast.Query
	var err Diagnostic
	switch first_token.Kind {
	case Set:
		query_result, err = ParseSetQuery(env, &tokens, &position)
	case Select:
		query_result, err = ParseSelectQuery(env, &tokens, &position)
	default:
		return ast.Query{}, UnExpectedStatementError(&tokens, &position)
	}

	// 消耗可选的 `;` 在有效语句的末尾
	if position < len(tokens) {
		last_token := tokens[position]
		if last_token.Kind == Semicolon {
			position += 1
		}
	}

	// 检查有效语句后是否存在未预期的内容
	if position < len(tokens) {
		return ast.Query{}, UnExpectedContentAfterCorrectStatement(
			&first_token.Literal,
			&tokens,
			&position,
		)
	}

	return query_result, err
}

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

	if *position >= lentokens || IsAssignmentOperator(&(*tokens)[*position]) {
		return ast.Query{}, *NewError("Expect `=` or `:=` and Value after Variable name").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	// Consume `=` token
	*position += 1

	aggregations_count_before := len(context.Aggregations)
	value, _ := ParseExpression(&context, env, tokens, position)
	has_aggregations := len(context.Aggregations) != aggregations_count_before

	if has_aggregations {
		return ast.Query{}, *NewError("Aggregation value can't be assigned to global variable").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	env.DefineGlobal(name, value.ExprType(env))

	global_variable := ast.GlobalVariableStatement{
		Name:  name,
		Value: value,
	}

	return ast.Query{GlobalVariableDeclaration: &global_variable}, Diagnostic{}
}

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
			statement, _ := ParseSelectStatement(&context, env, tokens, position)
			statements["select"] = statement
			context.IsSingleValueQuery = len(context.Aggregations) != 0
		case Where:
			if _, ok := statements["where"]; !ok {
				return ast.Query{}, *NewError("You already used `WHERE` statement").AddNote("Can't use more than one `WHERE` statement in the same query").WithLocation(token.Location)
			}

			statement, _ := ParseWhereStatement(&context, env, tokens, position)
			statements["where"] = statement
		case Group:
			if _, ok := statements["group"]; !ok {
				return ast.Query{}, *NewError("You already used `GROUP BY` statement").AddNote("Can't use more than one `GROUP BY` statement in the same query").WithLocation(token.Location)
			}

			statement, _ := ParseGroupByStatement(&context, env, tokens, position)
			statements["group"] = statement
		case Having:
			if _, ok := statements["having"]; ok {
				return ast.Query{}, *NewError("You already used `HAVING` statement").AddNote("Can't use more than one `HAVING` statement in the same query").WithLocation(token.Location)
			}
			if _, ok := statements["group"]; !ok {
				return ast.Query{}, *NewError("You already used `GROUP BY` statement").AddNote("Can't use more than one `GROUP BY` statement in the same query").WithLocation(token.Location)
			}

			statement, _ := ParseHavingStatement(&context, env, tokens, position)
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

				if *position >= lentokens && (*tokens)[*position].Kind == Integer {
					return ast.Query{}, *NewError("Expects `OFFSET` amount as Integer value after `,`").AddHelp("Try to add constant Integer after comma").AddNote("`OFFSET` value must be a constant Integer").WithLocation(token.Location)
				}

				count_result, err := strconv.Atoi((*tokens)[*position].Literal)

				// Report clear error for Integer parsing
				if err != nil {
					return ast.Query{}, *NewError("`OFFSET` integer value is invalid").AddNote(fmt.Sprintf("`OFFSET` value must be between 0 and ", math.MaxInt)).WithLocation(token.Location)
				}

				*position += 1
				count := count_result
				statements["offset"] = &ast.OffsetStatement{Count: count}
			}
		case Offset:
			if _, ok := statements["select"]; !ok {
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
			break
		}
	}

	// If any aggregation function is used, add Aggregation Functions Node to the GQL Query
	if len(context.Aggregations) != 0 {
		aggregation_functions := &ast.AggregationsStatement{
			Aggregations: context.Aggregations,
		}
		statements["aggregation"] = aggregation_functions
	}

	// Remove all selected fields from hidden selection
	var hidden_selections []string
	for _, selection := range context.HiddenSelections {
		if !contains(context.SelectedFields, selection) {
			hidden_selections = append(hidden_selections, selection)
		}
	}
	return ast.Query{
		Select: &ast.GQLQuery{
			Statements:             statements,
			HasAggregationFunction: context.IsSingleValueQuery,
			HasGroupByStatement:    context.HasGroupByStatement,
			HiddenSelections:       hidden_selections,
		},
	}, Diagnostic{}
}

func ParseSelectStatement(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	*position += 1

	if *position >= len(*tokens) {
		return &ast.SelectStatement{}, *NewError("Incomplete input for select statement").AddHelp("Try select one or more values in the `SELECT` statement").AddNote("Select statements requires at least selecting one value").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	var table_name string
	var fields_names []string
	var fields_values []ast.Expression // 请替换为实际的 Expression 类型
	alias_table := make(map[string]string)
	is_select_all := false
	is_distinct := false

	// Check if select has distinct keyword after it
	if (*tokens)[*position].Kind == Distinct {
		is_distinct = true
		*position++
	}

	// Select all option
	if *position < len(*tokens) && (*tokens)[*position].Kind == Star {
		// Consume `*`
		*position++
		is_select_all = true
	} else {
		for *position < len(*tokens) && (*tokens)[*position].Kind != From {
			expression, _ := ParseExpression(context, env, tokens, position)
			expr_type := expression.ExprType(env)
			expression_name, _ := GetExpressionName(expression)

			field_name := expression_name
			// if expression_name != "" {
			// 	field_name := expression_name
			// } else {
			// 	field_name := context.GenerateColumnName()
			// }

			// Assert that each selected field is unique
			if contains(fields_names, field_name) {
				return nil, *NewError("Can't select the same field twice").WithLocation(GetSafeLocation(tokens, *position-1))
			}

			// Check for Field name alias
			if *position < len(*tokens) && (*tokens)[*position].Kind == As {
				// Consume `as` keyword
				*position += 1
				alias_name_token, err := ConsumeKind(*tokens, *position, Symbol)
				if err != nil {
					return nil, *NewError("Expect `identifier` as field alias name").WithLocation(GetSafeLocation(tokens, *position))
				}

				// Register alias name
				alias_name := alias_name_token.Literal
				if contains(context.SelectedFields, alias_name) || alias_table[alias_name] != "" {
					return nil, *NewError("You already have field with the same name").AddHelp("Try to use a new unique name for alias").WithLocation(GetSafeLocation(tokens, *position))
				}

				// Consume alias name
				*position += 1

				// Register alias name type
				env.Define(alias_name, expr_type)

				context.SelectedFields = append(context.SelectedFields, alias_name)
				alias_table[field_name] = alias_name
			}

			// Register field type
			env.Define(field_name, expr_type)

			fields_names = append(fields_names, field_name)
			context.SelectedFields = append(context.SelectedFields, field_name)
			fields_values = append(fields_values, expression)

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

		table_name_token, err := ConsumeKind(*tokens, *position, Symbol)
		if err != nil {
			return nil, *NewError("Expect `identifier` as a table name").AddNote("Table name must be an identifier").WithLocation(GetSafeLocation(tokens, *position))
		}

		// Consume table name
		*position += 1

		table_name = table_name_token.Literal

		if _, ok := ast.TablesFieldsNames[table_name]; ok {
			return nil, *NewError("Unresolved table name").AddHelp("Check the documentations to see available tables").WithLocation(GetSafeLocation(tokens, *position))
		}

		RegisterCurrentTableFieldsTypes(table_name, *env)
	}

	// Make sure `SELECT *` used with specific table
	if is_select_all && table_name == "" {
		return nil, *NewError("Expect `FROM` and table name after `SELECT *`").AddNote("Select all must be used with valid table name").WithLocation(GetSafeLocation(tokens, *position))
	}

	// Select input validations
	if !is_select_all && len(fields_names) == 0 {
		return nil, *NewError("Incomplete input for select statement").AddHelp("Try select one or more values in the `SELECT` statement").AddNote("Select statements requires at least selecting one value").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	// If it `select *` make all table fields selectable
	if is_select_all {
		SelectAllTableFields(
			table_name,
			context.SelectedFields,
			fields_names,
			fields_values,
		)
	}

	// Type check all selected fields has type registered in type table
	err := TypeCheckSelectedFields(env, table_name, &fields_names, tokens, *position)
	fmt.Println(err)

	return &ast.SelectStatement{
		TableName:    table_name,
		FieldsNames:  fields_names,
		FieldsValues: fields_values,
		AliasTable:   alias_table,
		IsDistinct:   is_distinct,
	}, Diagnostic{}
}

func ParseWhereStatement(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	*position++
	if *position >= len(*tokens) {
		return nil, *NewError("Expect expression after `WHERE` keyword").AddHelp("Try to add boolean expression after `WHERE` keyword").AddNote("`WHERE` statement expects expression as condition").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	aggregations_count_before := len(context.Aggregations)

	condition_location := (*tokens)[*position].Location
	condition, _ := ParseExpression(context, env, tokens, position)
	condition_type := condition.ExprType(env)
	if condition_type.Fmt() != "Boolean" {
		return nil, *NewError(fmt.Sprintf("Expect `WHERE` condition to be type %s but got %s", "Boolean", condition_type)).AddNote("`WHERE` statement condition must be Boolean").WithLocation(condition_location)
	}

	aggregations_count_after := len(context.Aggregations)
	if aggregations_count_before != aggregations_count_after {
		return nil, *NewError("Can't use Aggregation functions in `WHERE` statement").AddNote("Aggregation functions must be used after `GROUP BY` statement").AddNote("Aggregation functions evaluated after later after `GROUP BY` statement").WithLocation(condition_location)
	}

	return &ast.WhereStatement{Condition: condition}, Diagnostic{}
}

func ParseGroupByStatement(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	*position += 1
	if *position >= len(*tokens) || (*tokens)[*position].Kind != By {
		return nil, *NewError("Expect keyword `by` after keyword `group`").AddHelp("Try to use `BY` keyword after `GROUP").WithLocation(GetSafeLocation(tokens, *position-1))
	}
	*position += 1
	if *position >= len(*tokens) || (*tokens)[*position].Kind != Symbol {
		return nil, *NewError("Expect field name after `group by`").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	field_name := (*tokens)[*position].Literal
	*position += 1

	if !env.Contains(field_name) {
		return nil, *NewError("Current table not contains field with this name").AddHelp("Check the documentations to see available fields for each tables").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	context.HasGroupByStatement = true
	return &ast.GroupByStatement{FieldName: field_name}, Diagnostic{}
}

func ParseHavingStatement(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	*position += 1
	if *position >= len(*tokens) {
		return nil, *NewError("Expect expression after `HAVING` keyword").AddHelp("Try to add boolean expression after `HAVING` keyword").AddNote("`HAVING` statement expects expression as condition").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	condition_location := (*tokens)[*position].Location
	condition, _ := ParseExpression(context, env, tokens, position)
	condition_type := condition.ExprType(env)
	if condition_type.Fmt() != "Boolean" {
		return nil, *NewError(fmt.Sprintf("Expect `HAVING` condition to be type %s but got %s", "Boolean", condition_type)).AddNote("`HAVING` statement condition must be Boolean").WithLocation(condition_location)
	}

	return &ast.HavingStatement{Condition: condition}, Diagnostic{}
}

func ParseLimitStatement(tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	*position += 1
	if *position >= len(*tokens) || (*tokens)[*position].Kind != Integer {
		return nil, *NewError("Expect number after `LIMIT` keyword").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	count_result, err := strconv.Atoi((*tokens)[*position].Literal)

	// Report clear error for Integer parsing
	if err != nil {
		return &ast.OffsetStatement{}, *NewError("`OFFSET` integer value is invalid").AddNote(fmt.Sprintf("`LIMIT` value must be between 0 and ", math.MaxInt)).WithLocation(GetSafeLocation(tokens, *position))
	}

	*position += 1
	count := count_result
	return &ast.LimitStatement{Count: count}, Diagnostic{}
}

func ParseOffsetStatement(tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	*position += 1

	if *position >= len(*tokens) && (*tokens)[*position].Kind == Integer {
		return &ast.OffsetStatement{}, *NewError("Expect number after `OFFSET` keyword").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	count_result, err := strconv.Atoi((*tokens)[*position].Literal)

	// Report clear error for Integer parsing
	if err != nil {
		return &ast.OffsetStatement{}, *NewError("`LIMIT` integer value is invalid").AddNote(fmt.Sprintf("`OFFSET` value must be between 0 and ", math.MaxInt)).WithLocation(GetSafeLocation(tokens, *position))
	}

	*position += 1
	count := count_result
	return &ast.OffsetStatement{Count: count}, Diagnostic{}
}

func ParseOrderByStatement(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Statement, Diagnostic) {
	// Consume `ORDER` keyword
	*position += 1

	if *position >= len(*tokens) || (*tokens)[*position].Kind != By {
		return nil, *NewError("`Expect keyword `BY` after keyword `ORDER").AddHelp("Try to use `BY` keyword after `ORDER").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	// Consume `BY` keyword
	*position += 1

	var arguments []ast.Expression
	var sorting_orders []ast.SortingOrder

	for {
		argument, _ := ParseExpression(context, env, tokens, position)
		arguments = append(arguments, argument)

		order := ast.Ascending
		if *position < len(*tokens) && IsAscOrDesc(&(*tokens)[*position]) {
			if (*tokens)[*position].Kind == Descending {
				order = ast.Descending
			}

			// Consume `ASC or DESC` keyword
			*position += 1
		}

		sorting_orders = append(sorting_orders, order)
		if *position < len(*tokens) && (*tokens)[*position].Kind == Comma {
			// Consume `,` keyword
			*position += 1
		} else {
			break
		}
	}

	return &ast.OrderByStatement{
		Arguments:     arguments,
		SortingOrders: sorting_orders,
	}, Diagnostic{}
}

func ParseExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	aggregations_count_before := len(context.Aggregations)
	expression, _ := ParseAssignmentExpression(context, env, tokens, position)
	has_aggregations := (len(context.Aggregations) != aggregations_count_before)
	if has_aggregations {
		column_name := context.GenerateColumnName()
		env.Define(column_name, expression.ExprType(env))

		// Register the new aggregation generated field if the this expression is after group by
		if context.HasGroupByStatement && !contains(context.HiddenSelections, column_name) {
			context.HiddenSelections = append(context.HiddenSelections, column_name)
		}
		context.Aggregations[column_name] = ast.AggregateValue{Expression: &expression}
		return ast.SymbolExpression{
			Value: column_name,
		}, Diagnostic{}
	}
	return expression, Diagnostic{}
}

func ParseAssignmentExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, _ := ParseIsNullExpression(context, env, tokens, position)
	if *position < len(*tokens) && (*tokens)[*position].Kind == ColonEqual {
		if expression.Kind() != ast.ExpressionKind(GlobalVariable) {
			return nil, *NewError("Assignment expressions expect global variable name before `:=`").WithLocation((*tokens)[*position].Location)
		}

		expr := expression.(*ast.GlobalVariableExpression)
		variable_name := expr.Name

		// Consume `:=` operator
		*position += 1

		value, _ := ParseIsNullExpression(context, env, tokens, position)
		env.DefineGlobal(variable_name, value.ExprType(env))

		return &ast.AssignmentExpression{
			Symbol: variable_name,
			Value:  value,
		}, Diagnostic{}
	}

	return expression, Diagnostic{}
}

// // getExpressionName 函数翻译
// func getExpressionName(expression Expression) (string, error) {
// 	switch expr := expression.asAny().(type) {
// 	case *SymbolExpression:
// 		return expr.value, nil
// 	case *GlobalVariableExpression:
// 		return expr.name, nil
// 	default:
// 		return "", fmt.Errorf("Invalid expression type")
// 	}
// }

/// Remove last token if it semicolon, because it's optional
func ConsumeOptionalSemicolonIfExists(tokens *[]Token) {
	if len(*tokens) == 0 {
		return
	}
	lastToken := (*tokens)[len(*tokens)-1]
	if lastToken.Kind == Semicolon {
		*tokens = (*tokens)[:len(*tokens)-1]
	}
}

func ParseIsNullExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, _ := ParseInExpression(context, env, tokens, position)
	if *position < len(*tokens) && (*tokens)[*position].Kind == Is {
		is_location := (*tokens)[*position].Location

		// Consume `IS` keyword
		*position += 1

		has_not_keyword := false
		if *position < len(*tokens) && (*tokens)[*position].Kind == Not {
			// Consume `NOT` keyword
			*position++
			has_not_keyword = true
		}

		if *position < len(*tokens) && (*tokens)[*position].Kind == Null {
			// Consume `Null` keyword
			*position += 1

			return &ast.IsNullExpression{
				Argument: expression,
				HasNot:   has_not_keyword,
			}, Diagnostic{}
		}

		return &ast.IsNullExpression{}, *NewError("Expects `NULL` Keyword after `IS` or `IS NOT`").WithLocation(is_location)
	}
	return expression, Diagnostic{}
}

func ParseInExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, _ := ParseBetweenExpression(context, env, tokens, position)

	has_not_keyword := false
	if *position < len(*tokens) && (*tokens)[*position].Kind == Not {
		has_not_keyword = true
	}

	if *position < len(*tokens) && (*tokens)[*position].Kind == In {
		in_location := (*tokens)[*position].Location

		// Consume `IN` keyword
		*position += 1

		if _, err := ConsumeKind(*tokens, *position, LeftParen); err != nil {
			return nil, *NewError("Expects values between `(` and `)` after `IN` keyword").WithLocation(in_location)
		}

		values, _ := ParseArgumentsExpressions(context, env, tokens, position)
		if len(values) == 0 {
			return ast.BooleanExpression{IsTrue: has_not_keyword}, Diagnostic{}
		}

		values_type_result := CheckAllValuesAreSameType(env, values)
		if values_type_result == nil {
			return nil, *NewError("Expects values between `(` and `)` to have the same type").WithLocation(in_location)
		}

		// Check that argument and values has the same type
		values_type := values_type_result
		if values_type.Fmt() != "Any" && expression.ExprType(env) != values_type {
			return nil, *NewError("Argument and Values of In Expression must have the same type").WithLocation(in_location)
		}

		return &ast.InExpression{
			Argument:      expression,
			Values:        values,
			ValuesType:    values_type,
			HasNotKeyword: has_not_keyword,
		}, Diagnostic{}
	}

	if has_not_keyword {
		return nil, *NewError("Expects `IN` expression after this `NOT` keyword").AddHelp("Try to use `IN` expression after NOT keyword").AddHelp("Try to remove `NOT` keyword").AddNote("Expect to see `NOT` then `IN` keyword with a list of values").WithLocation(GetSafeLocation(tokens, *position-1))
	}

	return expression, Diagnostic{}
}

func ParseBetweenExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, _ := ParseLogicalOrExpression(context, env, tokens, position)

	if *position < len(*tokens) && (*tokens)[*position].Kind == Between {
		between_location := (*tokens)[*position].Location

		// Consume `BETWEEN` keyword
		*position += 1

		if *position >= len(*tokens) {
			return nil, *NewError("`BETWEEN` keyword expects two range after it").WithLocation(between_location)
		}

		argument_type := expression.ExprType(env)
		range_start, _ := ParseLogicalOrExpression(context, env, tokens, position)

		if *position >= len(*tokens) || (*tokens)[*position].Kind != DotDot {
			return nil, *NewError("Expect `..` after `BETWEEN` range start").WithLocation(between_location)
		}

		// Consume `..` token
		*position += 1
		range_end, _ := ParseLogicalOrExpression(context, env, tokens, position)

		if argument_type != range_start.ExprType(env) || argument_type != range_end.ExprType(env) {
			return nil, *NewError(fmt.Sprintf("Expect `BETWEEN` argument, range start and end to has same type but got %s, %s and %s", argument_type, range_start.ExprType(env), range_end.ExprType(env))).AddHelp("Try to make sure all of them has same type").WithLocation(between_location)
		}

		return &ast.BetweenExpression{
			Value:      expression,
			RangeStart: range_start,
			RangeEnd:   range_end,
		}, Diagnostic{}
	}

	return expression, Diagnostic{}
}

func ParseLogicalOrExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseLogicalAndExpression(context, env, tokens, position)
	if err.message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression
	for *position < len(*tokens) && (*tokens)[*position].Kind != LogicalOr {
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

		lhs = ast.LogicalExpression{
			Left:     lhs,
			Operator: ast.Or,
			Right:    rhs,
		}
	}

	return lhs, Diagnostic{}
}

func ParseLogicalAndExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseBitwiseOrExpression(context, env, tokens, position)
	if err.message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression
	for *position < len(*tokens) && (*tokens)[*position].Kind != LogicalOr {
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

		lhs = ast.LogicalExpression{
			Left:     lhs,
			Operator: ast.And,
			Right:    rhs,
		}
	}

	return lhs, Diagnostic{}
}
func ParseBitwiseOrExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseLogicalXorExpression(context, env, tokens, position)
	if err.message != "" || *position >= len(*tokens) {
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

		rhs, _ := ParseLogicalXorExpression(context, env, tokens, position)
		if rhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		return &ast.BitwiseExpression{
			Left:     lhs,
			Operator: ast.BitwiseOperator(ast.Or),
			Right:    rhs,
		}, Diagnostic{}
	}

	return lhs, Diagnostic{}
}

func ParseLogicalXorExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseBitwiseAndExpression(context, env, tokens, position)
	if err.message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression
	for *position < len(*tokens) && (*tokens)[*position].Kind != LogicalXor {
		*position += 1
		if lhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position-2].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		rhs, _ := ParseBitwiseAndExpression(context, env, tokens, position)
		if rhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		lhs = ast.LogicalExpression{
			Left:     lhs,
			Operator: ast.Xor,
			Right:    rhs,
		}
	}

	return lhs, Diagnostic{}
}

func ParseBitwiseAndExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseEqualityExpression(context, env, tokens, position)
	if err.message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression

	if *position < len(*tokens) && (*tokens)[*position].Kind != BitwiseAnd {
		*position += 1
		if lhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position-2].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		rhs, _ := ParseEqualityExpression(context, env, tokens, position)
		if rhs.ExprType(env).Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				(*tokens)[*position].Location,
				ast.Boolean{},
				lhs.ExprType(env),
			)
		}

		lhs = ast.BitwiseExpression{
			Left:     lhs,
			Operator: ast.BOAnd,
			Right:    rhs,
		}
	}

	return lhs, Diagnostic{}
}

func ParseEqualityExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseComparisonExpression(context, env, tokens, position)
	if err.message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression

	operator := &(*tokens)[*position]
	if operator.Kind == Equal || operator.Kind == BangEqual {
		*position += 1
		var comparison_operator ast.ComparisonOperator
		if operator.Kind == Equal {
			comparison_operator = ast.COEqual
		} else {
			comparison_operator = ast.CONotEqual
		}

		rhs, _ := ParseComparisonExpression(context, env, tokens, position)

		switch AreTypesEquals(env, lhs, rhs) {
		case Equals{}:
			// do nothing
		case RightSideCasted{expr: expression}:
			rhs = expression
		case LeftSideCasted{expr: expression}:
			lhs = expression
		case NotEqualAndCantImplicitCast{}:
			lhs_type := lhs.ExprType(env)
			rhs_type := rhs.ExprType(env)
			diagnostic := *NewError(fmt.Sprintf(
				"Can't compare values of different types `%s` and `%s`",
				lhs_type,
				rhs_type,
			)).WithLocation(GetSafeLocation(tokens, *position-2))

			// Provides help messages if use compare null to non null value
			if lhs_type.IsNull() || rhs_type.IsNull() {
				return nil, *diagnostic.AddHelp("Try to use `IS NULL expr` expression").AddHelp("Try to use `ISNULL(expr)` function")
			}

			return nil, diagnostic
		default:
			return nil, *NewError("").WithLocation(GetSafeLocation(tokens, *position-2))
		}

		return ast.ComparisonExpression{
			Left:     lhs,
			Operator: comparison_operator,
			Right:    rhs,
		}, Diagnostic{}
	}

	return lhs, Diagnostic{}
}
func ParseComparisonExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseBitwiseShiftExpression(context, env, tokens, position)
	if err.message != "" || *position >= len(*tokens) {
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

		rhs, _ := ParseBitwiseShiftExpression(context, env, tokens, position)

		switch AreTypesEquals(env, lhs, rhs) {
		case Equals{}:
			// do nothing
		case RightSideCasted{expr: expression}:
			rhs = expression
		case LeftSideCasted{expr: expression}:
			lhs = expression
		case NotEqualAndCantImplicitCast{}:
			lhs_type := lhs.ExprType(env)
			rhs_type := rhs.ExprType(env)
			diagnostic := *NewError(fmt.Sprintf(
				"Can't compare values of different types `%s` and `%s`",
				lhs_type,
				rhs_type,
			)).WithLocation(GetSafeLocation(tokens, *position-2))

			// Provides help messages if use compare null to non null value
			if lhs_type.IsNull() || rhs_type.IsNull() {
				return nil, *diagnostic.AddHelp("Try to use `IS NULL expr` expression").AddHelp("Try to use `ISNULL(expr)` function")
			}

			return nil, diagnostic
		default:
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
func ParseBitwiseShiftExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	lhs, _ := ParseTermExpression(context, env, tokens, position)

	for *position < len(*tokens) && IsBitwiseShiftOperator(&(*tokens)[*position]) {
		operator := &(*tokens)[*position]
		*position += 1
		var bitwise_operator ast.BitwiseOperator
		if operator.Kind == BitwiseRightShift {
			bitwise_operator = ast.BORightShift
		} else {
			bitwise_operator = ast.BOLeftShift
		}

		rhs, _ := ParseTermExpression(context, env, tokens, position)

		// Make sure right and left hand side types are numbers
		if rhs.ExprType(env).Fmt() == "Integer" && rhs.ExprType(env) != lhs.ExprType(env) {
			return nil, *NewError(fmt.Sprintf(
				"Bitwise operators require number types but got `%s` and `%s`",
				lhs.ExprType(env),
				rhs.ExprType(env),
			)).WithLocation(GetSafeLocation(tokens, *position-2))
		}

		lhs = &ast.BitwiseExpression{
			Left:     lhs,
			Operator: bitwise_operator,
			Right:    rhs,
		}
	}

	return lhs, Diagnostic{}
}

func ParseTermExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	lhs, _ := ParseFactorExpression(context, env, tokens, position)

	for *position < len(*tokens) && IsTermOperator(&(*tokens)[*position]) {
		operator := &(*tokens)[*position]
		*position += 1
		var math_operator ast.ArithmeticOperator
		if operator.Kind == Plus {
			math_operator = ast.AOPlus
		} else {
			math_operator = ast.AOMinus
		}

		rhs, _ := ParseFactorExpression(context, env, tokens, position)

		lhs_type := lhs.ExprType(env)
		rhs_type := rhs.ExprType(env)

		// Make sure right and left hand side types are numbers
		if lhs_type.IsNumber() && rhs_type.IsNumber() {
			lhs = &ast.ArithmeticExpression{
				Left:     lhs,
				Operator: math_operator,
				Right:    rhs,
			}

			continue
		}

		if math_operator == ast.ArithmeticOperator(Plus) {
			return nil, *NewError(fmt.Sprintf(
				"Math operators `+` both sides to be number types but got `%s` and `%s`",
				lhs_type,
				rhs_type,
			)).AddHelp("You can use `CONCAT(Any, Any, ...Any)` function to concatenate values with different types").WithLocation(operator.Location)
		}

		return nil, *NewError(fmt.Sprintf(
			"Math operators require number types but got `%s` and `%s`",
			lhs_type,
			rhs_type,
		)).WithLocation(operator.Location)
	}

	return lhs, Diagnostic{}
}

func ParseFactorExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseLikeExpression(context, env, tokens, position)
	if err.message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression
	for *position < len(*tokens) && IsFactorOperator(&(*tokens)[*position]) {
		operator := &(*tokens)[*position]
		*position += 1

		var factor_operator ast.ArithmeticOperator
		switch operator.Kind {
		case Star:
			factor_operator = ast.AOStar
		case Slash:
			factor_operator = ast.AOSlash
		default:
			factor_operator = ast.AOModulus
		}

		rhs, _ := ParseLikeExpression(context, env, tokens, position)

		lhs_type := lhs.ExprType(env)
		rhs_type := rhs.ExprType(env)

		// Make sure right and left hand side types are numbers
		if lhs_type.IsNumber() && rhs_type.IsNumber() {
			lhs = &ast.ArithmeticExpression{
				Left:     lhs,
				Operator: factor_operator,
				Right:    rhs,
			}
			continue
		}

		return nil, *NewError(fmt.Sprintf(
			"Math operators require number types but got `%s` and `%s`",
			lhs_type,
			rhs_type,
		)).WithLocation(GetSafeLocation(tokens, *position-2))
	}

	return lhs, Diagnostic{}
}

func ParseLikeExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, err := ParseGlobExpression(context, env, tokens, position)
	if err.message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression
	if (*tokens)[*position].Kind == Like {
		location := (*tokens)[*position].Location
		*position += 1

		if !lhs.ExprType(env).IsText() {
			return nil, *NewError(fmt.Sprintf("Expect `LIKE` left hand side to be `TEXT` but got %s", lhs.ExprType(env))).WithLocation(location)
		}

		pattern, _ := ParseGlobExpression(context, env, tokens, position)
		if !pattern.ExprType(env).IsText() {
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
	if err.message != "" || *position >= len(*tokens) {
		return expression, err
	}

	lhs := expression
	if (*tokens)[*position].Kind == Glob {
		location := (*tokens)[*position].Location
		*position += 1

		if !lhs.ExprType(env).IsText() {
			return nil, *NewError(fmt.Sprintf("Expect `GLOB` left hand side to be `TEXT` but got %s", lhs.ExprType(env))).WithLocation(location)
		}

		pattern, _ := ParseUnaryExpression(context, env, tokens, position)
		if !pattern.ExprType(env).IsText() {
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
			op = ast.Bang
		} else {
			op = ast.Minus
		}

		*position += 1

		rhs, _ := ParseUnaryExpression(context, env, tokens, position)
		rhs_type := rhs.ExprType(env)
		if op == ast.Bang && rhs_type.Fmt() != "Boolean" {
			return nil, TypeMismatchError(
				GetSafeLocation(tokens, *position-1),
				ast.Boolean{},
				rhs_type,
			)
		} else if op == ast.Minus && rhs_type.Fmt() != "Integer" {
			return nil, TypeMismatchError(
				GetSafeLocation(tokens, *position-1),
				ast.Integer{},
				rhs_type,
			)
		}

		return &ast.PrefixUnary{Right: rhs, Op: op}, Diagnostic{}
	}

	return ParseFunctionCallExpression(context, env, tokens, position)
}

func ParseFunctionCallExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	expression, _ := ParsePrimaryExpression(context, env, tokens, position)

	if *position < len(*tokens) && (*tokens)[*position].Kind == LeftParen {
		symbol_expression := expression.(*ast.SymbolExpression)
		function_name_location := GetSafeLocation(tokens, *position)

		// Make sure function name is SymbolExpression
		if symbol_expression == nil {
			return nil, *NewError("Function name must be an identifier").WithLocation(function_name_location)
		}

		function_name := symbol_expression.Value

		// Check if this function is a Standard library functions
		if _, ok := ast.Functions[function_name]; ok {
			arguments, _ := ParseArgumentsExpressions(context, env, tokens, position)
			prototype := ast.Prototypes[function_name]
			parameters := prototype.Parameters
			return_type := prototype.Result

			_, err := CheckFunctionCallArguments(
				env,
				&arguments,
				&parameters,
				function_name,
				function_name_location,
			)
			if err.message != "" {
				return nil, err
			}

			// Register function name with return type
			env.Define(function_name, return_type)

			return &ast.CallExpression{
				FunctionName:  function_name,
				Arguments:     arguments,
				IsAggregation: false,
			}, Diagnostic{}
		}

		// Check if this function is an Aggregation functions
		if _, ok := ast.Aggregations[function_name]; ok {
			arguments, _ := ParseArgumentsExpressions(context, env, tokens, position)
			prototype := ast.AggregationsProtos[function_name]
			parameters := []ast.DataType{prototype.Parameter}
			return_type := prototype.Result

			_, err := CheckFunctionCallArguments(
				env,
				&arguments,
				&parameters,
				function_name,
				function_name_location,
			)
			if err.message != "" {
				return nil, err
			}

			argument_result, argueerr := GetExpressionName(arguments[0])
			if argueerr.Error() != "" {
				return nil, *NewError("Invalid Aggregation function argument").AddHelp("Try to use field name as Aggregation function argument").AddNote("Aggregation function accept field name as argument").WithLocation(function_name_location)
			}

			argument := argument_result
			column_name := context.GenerateColumnName()

			context.HiddenSelections = append(context.HiddenSelections, column_name)

			// Register aggregation generated name with return type
			env.Define(column_name, return_type)

			context.Aggregations[column_name] = ast.AggregateValue{
				Function: struct {
					First  string
					Second string
				}{
					First:  function_name,
					Second: argument,
				},
			}

			return &ast.SymbolExpression{
				Value: column_name,
			}, Diagnostic{}
		}

		// Report that this function name is not standard or aggregation
		return nil, *NewError("No such function name").AddHelp(fmt.Sprintf("Function `%s` is not an Aggregation or Standard library function name", function_name)).WithLocation(function_name_location)
	}
	return expression, Diagnostic{}
}

func ParseArgumentsExpressions(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) ([]ast.Expression, Diagnostic) {
	var arguments []ast.Expression
	if _, err := ConsumeKind(*tokens, *position, LeftParen); err == nil {
		*position += 1

		for (*tokens)[*position].Kind != RightParen {
			argument, _ := ParseExpression(context, env, tokens, position)

			argument_literal, err := GetExpressionName(argument)
			if err == nil {
				literal := argument_literal
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

func ParsePrimaryExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	if *position >= len(*tokens) {
		return nil, UnExpectedExpressionError(tokens, position)
	}

	switch (*tokens)[*position].Kind {
	case String:
		*position += 1
		return ast.StringExpression{
			Value:     (*tokens)[*position-1].Literal,
			ValueType: ast.StringValueText,
		}, Diagnostic{}
	case Symbol:
		value := (*tokens)[*position-1].Literal
		*position += 1
		if !contains(context.SelectedFields, value) {
			context.HiddenSelections = append(context.HiddenSelections, value)
		}
		return ast.SymbolExpression{Value: value}, Diagnostic{}
	case GlobalVariable:
		name := (*tokens)[*position-1].Literal
		*position += 1
		return ast.GlobalVariableExpression{Name: name}, Diagnostic{}
	case Integer:
		if integer, err := strconv.ParseInt((*tokens)[*position-1].Literal, 10, 64); err == nil {
			*position += 1
			value := ast.IntegerValue{Value: integer}
			return ast.NumberExpression{Value: value}, Diagnostic{}
		}

		return nil, *NewError("Too big Integer value").AddHelp(fmt.Sprintf("Integer value must be between %s and %s", math.MinInt64, math.MaxInt64)).WithLocation((*tokens)[*position].Location)
	case Float:
		if float, err := strconv.ParseFloat((*tokens)[*position-1].Literal, 64); err == nil {
			*position += 1
			value := ast.FloatValue{Value: float}
			return ast.NumberExpression{Value: value}, Diagnostic{}
		}

		return nil, *NewError("Too big Float value").AddHelp("Try to use smaller value").AddNote(fmt.Sprintf("Float value must be between %s and %s", -math.MaxFloat64, math.MaxFloat64)).WithLocation((*tokens)[*position].Location)
	case True:
		*position += 1
		return ast.BooleanExpression{IsTrue: true}, Diagnostic{}
	case False:
		*position += 1
		return ast.BooleanExpression{IsTrue: false}, Diagnostic{}
	case Null:
		*position += 1
		return ast.NullExpression{}, Diagnostic{}
	case LeftParen:
		return ParseGroupExpression(context, env, tokens, position)
	case Case:
		return ParseCaseExpression(context, env, tokens, position)
	default:
		return nil, UnExpectedExpressionError(tokens, position)
	}
}

func ParseGroupExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	*position += 1
	expression, _ := ParseExpression(context, env, tokens, position)
	if (*tokens)[*position].Kind != RightParen {
		return nil, *NewError("Expect `)` to end group expression").WithLocation(GetSafeLocation(tokens, *position)).AddHelp("Try to add ')' at the end of group expression")
	}
	*position += 1
	return expression, Diagnostic{}
}

func ParseCaseExpression(context *ParserContext, env *ast.Environment, tokens *[]Token, position *int) (ast.Expression, Diagnostic) {
	conditions := []ast.Expression{}
	values := []ast.Expression{}
	var default_value ast.Expression = nil

	// Consume `case` keyword
	case_location := (*tokens)[*position].Location
	*position += 1

	has_else_branch := false

	for *position < len(*tokens) && (*tokens)[*position].Kind != End {
		// Else branch
		if (*tokens)[*position].Kind == Else {
			if has_else_branch {
				return nil, *NewError("This `CASE` expression already has else branch").AddNote("`CASE` expression can has only one `ELSE` branch").WithLocation(GetSafeLocation(tokens, *position))
			}

			// Consume `ELSE` keyword
			*position += 1

			default_value_expr, _ := ParseExpression(context, env, tokens, position)
			default_value = default_value_expr
			has_else_branch = true
			continue
		}

		// Check if current token kind is `WHEN` keyword
		if _, err := ConsumeKind(*tokens, *position, When); err != nil {
			return nil, *NewError("Expect `when` before case condition").AddHelp("Try to add `WHEN` keyword before any condition").WithLocation(GetSafeLocation(tokens, *position))
		}

		// Consume when keyword
		*position += 1

		condition, _ := ParseExpression(context, env, tokens, position)
		if condition.ExprType(env).Fmt() != "Boolean" {
			return nil, *NewError("Case condition must be a boolean type").WithLocation(GetSafeLocation(tokens, *position))
		}

		conditions = append(conditions, condition)

		if _, err := ConsumeKind(*tokens, *position, Then); err != nil {
			return nil, *NewError("Expect `THEN` after case condition").WithLocation(GetSafeLocation(tokens, *position))
		}

		// Consume then keyword
		*position += 1

		expression, _ := ParseExpression(context, env, tokens, position)
		values = append(values, expression)
	}

	// Make sure case expression has at least else branch
	if len(conditions) == 0 && !has_else_branch {
		return nil, *NewError("Case expression must has at least else branch").WithLocation(GetSafeLocation(tokens, *position))
	}

	// Make sure case expression end with END keyword
	if *position >= len(*tokens) || (*tokens)[*position].Kind != End {
		return nil, *NewError("Expect `END` after case branches").WithLocation(GetSafeLocation(tokens, *position))
	}

	// Consume end
	*position += 1

	// Make sure this case expression has else branch
	if !has_else_branch {
		return nil, *NewError("Case expression must has else branch").WithLocation(GetSafeLocation(tokens, *position))
	}

	// Assert that all values has the same type
	values_type := values[0].ExprType(env)
	for i, value := range values[1:] {
		if values_type != value.ExprType(env) {
			return nil, *NewError(fmt.Sprintf("Case value in branch %s has different type than the last branch", i+1)).AddNote("All values in `CASE` expression must has the same Type").WithLocation(case_location)
		}
	}

	return ast.CaseExpression{
		Conditions:   conditions,
		Values:       values,
		DefaultValue: default_value,
		ValuesType:   values_type,
	}, Diagnostic{}
}

func CheckFunctionCallArguments(
	env *ast.Environment,
	arguments *[]ast.Expression,
	parameters *[]ast.DataType,
	functionname string,
	location Location,
) (ast.Expression, Diagnostic) {
	parameters_len := len(*parameters)
	arguments_len := len(*arguments)

	has_optional_parameter := false
	has_varargs_parameter := false
	if len(*parameters) != 0 {
		last_parameter := (*parameters)[len(*parameters)-1]
		has_optional_parameter = last_parameter.IsOptional()
		has_varargs_parameter = last_parameter.IsVarargs()
	}

	// Has Optional parameter type at the end
	if has_optional_parameter {
		// If function last parameter is optional make sure it at least has
		if arguments_len < parameters_len-1 {
			return nil, *NewError(fmt.Sprintf("Function `%s` expects at least `%s` arguments but got `%s`", functionname, parameters_len-1, arguments_len)).WithLocation(location)
		}

		// Make sure function with optional parameter not called with too much arguments
		if arguments_len > parameters_len {
			return nil, *NewError(fmt.Sprintf("Function `%s` expects at most `%s` arguments but got `%s`", functionname, parameters_len, arguments_len)).WithLocation(location)
		}
	} else if has_varargs_parameter {
		// If function last parameter is optional make sure it at least has
		if arguments_len < parameters_len-1 {
			return nil, *NewError(fmt.Sprintf("Function `%s` expects at least `%s` arguments but got `%s`", functionname, parameters_len-1, arguments_len)).WithLocation(location)
		}
	} else if arguments_len != parameters_len {
		return nil, *NewError(fmt.Sprintf("Function `%s` expects `%s` arguments but got `%s`", functionname, parameters_len, arguments_len)).WithLocation(location)
	}

	last_required_parameter_index := parameters_len
	if has_optional_parameter || has_varargs_parameter {
		last_required_parameter_index -= 1
	}

	// Check each argument vs parameter type
	for index := 0; index < last_required_parameter_index; index++ {
		parameter_type := (*parameters)[index]
		argument := (*arguments)[index]

		switch IsExpressionTypeEquals(env, argument, parameter_type) {
		case Equals{}:
			// do nothing
		case RightSideCasted{expr: argument}:
			(*arguments)[index] = argument
		case LeftSideCasted{expr: argument}:
			(*arguments)[index] = argument
		case NotEqualAndCantImplicitCast{}:
			argument_type := argument.ExprType(env)
			return nil, *NewError(fmt.Sprintf("Function `%s` argument number %s with type `%s` don't match expected type `%s`", functionname, index, argument_type, parameter_type)).WithLocation(location)
		case Error{}:
			return nil, *NewError("")
		}
	}

	if has_optional_parameter || has_varargs_parameter {
		last_parameter_type := (*parameters)[last_required_parameter_index]

		for index := last_required_parameter_index; index < arguments_len; index++ {
			argument := (*arguments)[index]
			switch IsExpressionTypeEquals(env, argument, last_parameter_type) {
			case Equals{}:
				// do nothing
			case RightSideCasted{expr: argument}:
				(*arguments)[index] = argument
			case LeftSideCasted{expr: argument}:
				(*arguments)[index] = argument
			case NotEqualAndCantImplicitCast{}:
				argument_type := (*arguments)[index]
				return nil, *NewError(fmt.Sprintf("Function `%s` argument number %s with type `%s` don't match expected type `%s`", functionname, index, argument_type, last_parameter_type)).WithLocation(location)
			case Error{}:
				return nil, *NewError("")
			}
		}
	}

	return nil, Diagnostic{}
}

func TypeCheckSelectedFields(
	env *ast.Environment,
	tableName string,
	fieldNames *[]string,
	tokens *[]Token,
	position int,
) Diagnostic {
	for _, fieldName := range *fieldNames {
		if dataType, err := env.ResolveType(fieldName); err != nil {
			if dataType.Fmt() == "Undefined" {
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

func UnExpectedContentAfterCorrectStatement(statementname *string, tokens *[]Token, position *int) Diagnostic {
	error_message := fmt.Sprintf("Unexpected content after the end of `%s` statement", statementname)
	location_of_extra_content := Location{
		Start: (*tokens)[*position].Location.Start,
		End:   (*tokens)[len(*tokens)-1].Location.End,
	}

	return *NewError(error_message).AddHelp("Try to check if statement keyword is missing").AddHelp("Try remove un expected extra content").WithLocation(location_of_extra_content)
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

func RegisterCurrentTableFieldsTypes(table_name string, symbol_table ast.Environment) {
	table_fields_names := ast.TablesFieldsNames[table_name]
	for _, field_name := range table_fields_names {
		field_type := ast.TablesFieldsTypes[field_name]
		symbol_table.Define(field_name, field_type)
	}
}

func SelectAllTableFields(table_name string, selected_fields, fields_names []string, fields_values []ast.Expression) {
	if table_fields, ok := ast.TablesFieldsNames[table_name]; ok {
		for _, field := range table_fields {
			if !contains(fields_names, field) {
				fields_names = append(fields_names, field)
				selected_fields = append(selected_fields, field)
				literal_expr := &ast.SymbolExpression{
					Value: field,
				}
				fields_values = append(fields_values, literal_expr)
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

func TypeMismatchError(location Location, expected ast.DataType, actual ast.DataType) Diagnostic {
	return *NewError(fmt.Sprintf("Type mismatch expected `%s`, got `%s`", expected, actual)).WithLocation(location)
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
