package engine

import (
	"errors"
	"sort"

	"github.com/go-git/go-git/v5"

	"github.com/ggql/ggql/ast"
)

func ExecuteStatement(
	env *ast.Environment,
	statement ast.Statement,
	repo *git.Repository,
	gitqlObject *ast.GitQLObject,
	aliasTable map[string]string,
	hiddenSelection []string,
) error {
	switch statement.Kind() {
	case ast.Select:
		selectStatement := statement.(*ast.SelectStatement)
		for _, alias := range selectStatement.AliasTable {
			aliasTable[alias] = alias
		}
		return executeSelectStatement(env, selectStatement, repo, gitqlObject, hiddenSelection)
	case ast.Where:
		whereStatement := statement.(*ast.WhereStatement)
		return executeWhereStatement(env, whereStatement, gitqlObject)
	case ast.Having:
		havingStatement := statement.(*ast.HavingStatement)
		return executeHavingStatement(env, havingStatement, gitqlObject)
	case ast.Limit:
		limitStatement := statement.(*ast.LimitStatement)
		return executeLimitStatement(limitStatement, gitqlObject)
	case ast.Offset:
		offsetStatement := statement.(*ast.OffsetStatement)
		return executeOffsetStatement(offsetStatement, gitqlObject)
	case ast.OrderBy:
		orderByStatement := statement.(*ast.OrderByStatement)
		return executeOrderByStatement(env, orderByStatement, gitqlObject)
	case ast.GroupBy:
		groupByStatement := statement.(*ast.GroupByStatement)
		return executeGroupByStatement(groupByStatement, gitqlObject)
	case ast.AggregateFunction:
		aggregationsStatement := statement.(*ast.AggregationsStatement)
		return executeAggregationFunctionStatement(env, aggregationsStatement, gitqlObject, aliasTable)
	case ast.GlobalVariable:
		globalVariableStatement := statement.(*ast.GlobalVariableStatement)
		return executeGlobalVariableStatement(env, globalVariableStatement)
	default:
		return errors.New("unknown statement kind")
	}
}

func executeSelectStatement(
	env *ast.Environment,
	statement *ast.SelectStatement,
	repo *git.Repository,
	gitqlObject *ast.GitQLObject,
	hiddenSelections []string,
) error {
	fieldsNames := statement.FieldsNames
	if statement.TableName != "" {
		for _, hidden := range hiddenSelections {
			if !contains(fieldsNames, hidden) {
				fieldsNames = append(fieldsNames, hidden)
			}
		}
	}

	for _, fieldName := range fieldsNames {
		gitqlObject.Titles = append(gitqlObject.Titles, GetColumnName(statement.AliasTable, fieldName))
	}

	objects, err := SelectGQLObjects(env, repo, statement.TableName, fieldsNames, gitqlObject.Titles, statement.FieldsValues)
	if err != nil {
		return err
	}

	if gitqlObject.IsEmpty() {
		gitqlObject.Groups = append(gitqlObject.Groups, *objects)
	} else {
		gitqlObject.Groups[0].Rows = append(gitqlObject.Groups[0].Rows, objects.Rows...)
	}

	return nil
}

func executeWhereStatement(
	env *ast.Environment,
	statement *ast.WhereStatement,
	gitqlObject *ast.GitQLObject,
) error {
	if gitqlObject.IsEmpty() {
		return nil
	}

	filteredGroup := ast.Group{}
	firstGroup := gitqlObject.Groups[0].Rows
	for _, object := range firstGroup {
		evalResult, err := EvaluateExpression(env, statement.Condition, gitqlObject.Titles, object.Values)
		if err != nil {
			return err
		}

		if evalResult.AsBool() {
			filteredGroup.Rows = append(filteredGroup.Rows, ast.Row{Values: object.Values})
		}
	}

	gitqlObject.Groups = gitqlObject.Groups[1:]
	gitqlObject.Groups = append(gitqlObject.Groups, filteredGroup)

	return nil
}

func executeHavingStatement(
	env *ast.Environment,
	statement *ast.HavingStatement,
	gitqlObject *ast.GitQLObject,
) error {
	if gitqlObject.IsEmpty() {
		return nil
	}

	if len(gitqlObject.Groups) > 1 {
		gitqlObject.Flat()
	}

	filteredGroup := ast.Group{}
	firstGroup := gitqlObject.Groups[0].Rows
	for _, object := range firstGroup {
		evalResult, err := EvaluateExpression(env, statement.Condition, gitqlObject.Titles, object.Values)
		if err != nil {
			return err
		}

		if evalResult.AsBool() {
			filteredGroup.Rows = append(filteredGroup.Rows, ast.Row{Values: object.Values})
		}
	}

	gitqlObject.Groups = gitqlObject.Groups[1:]
	gitqlObject.Groups = append(gitqlObject.Groups, filteredGroup)

	return nil
}

func executeLimitStatement(
	statement *ast.LimitStatement,
	gitqlObject *ast.GitQLObject,
) error {
	if gitqlObject.IsEmpty() {
		return nil
	}

	if gitqlObject.Len() > 1 {
		gitqlObject.Flat()
	}

	mainGroup := &gitqlObject.Groups[0]
	if statement.Count <= mainGroup.Len() {
		mainGroup.Rows = mainGroup.Rows[:statement.Count]
	}

	return nil
}

func executeOffsetStatement(
	statement *ast.OffsetStatement,
	gitqlObject *ast.GitQLObject,
) error {
	if gitqlObject.IsEmpty() {
		return nil
	}

	if len(gitqlObject.Groups) > 1 {
		gitqlObject.Flat()
	}

	mainGroup := &gitqlObject.Groups[0]
	if statement.Count <= len(mainGroup.Rows) {
		mainGroup.Rows = mainGroup.Rows[statement.Count:]
	} else {
		mainGroup.Rows = nil
	}

	return nil
}

func executeOrderByStatement(
	env *ast.Environment,
	statement *ast.OrderByStatement,
	gitqlObject *ast.GitQLObject,
) error {
	if gitqlObject.IsEmpty() {
		return nil
	}

	if gitqlObject.Len() > 1 {
		gitqlObject.Flat()
	}

	mainGroup := &gitqlObject.Groups[0]
	if mainGroup.IsEmpty() {
		return nil
	}

	sort.SliceStable(mainGroup.Rows, func(i, j int) bool {
		for idx := range statement.Arguments {
			argument := statement.Arguments[idx]
			if argument.IsConst() {
				continue
			}

			first, err := EvaluateExpression(env, argument, gitqlObject.Titles, mainGroup.Rows[i].Values)
			if err != nil {
				return false
			}

			other, err := EvaluateExpression(env, argument, gitqlObject.Titles, mainGroup.Rows[j].Values)
			if err != nil {
				return false
			}

			currentOrdering := compare(first, other)

			if currentOrdering == 0 {
				continue
			}

			if statement.SortingOrders[idx] == ast.Descending {
				return currentOrdering > 0
			} else {
				return currentOrdering < 0
			}
		}

		return false
	})

	return nil
}

func executeGroupByStatement(
	statement *ast.GroupByStatement,
	gitqlObject *ast.GitQLObject,
) error {
	if gitqlObject.IsEmpty() {
		return nil
	}

	mainGroup := gitqlObject.Groups[0]
	if mainGroup.IsEmpty() {
		return nil
	}

	groupsMap := make(map[string]int)
	nextGroupIndex := 0

	for _, object := range mainGroup.Rows {
		fieldIndex := indexOf(gitqlObject.Titles, statement.FieldName)
		fieldValue := object.Values[fieldIndex].Fmt()
		if _, ok := groupsMap[fieldValue]; !ok {
			groupsMap[fieldValue] = nextGroupIndex
			nextGroupIndex++
			gitqlObject.Groups = append(gitqlObject.Groups, ast.Group{Rows: []ast.Row{object}})
		} else {
			index := groupsMap[fieldValue]
			targetGroup := &gitqlObject.Groups[index]
			targetGroup.Rows = append(targetGroup.Rows, object)
		}
	}

	return nil
}

func executeAggregationFunctionStatement(
	env *ast.Environment,
	statement *ast.AggregationsStatement,
	gitqlObject *ast.GitQLObject,
	aliasTable map[string]string,
) error {
	if len(statement.Aggregations) == 0 {
		return nil
	}

	groupsCount := gitqlObject.Len()

	for _, group := range gitqlObject.Groups {
		if group.IsEmpty() {
			continue
		}
		// Resolve all aggregations functions first
		for _, aggregation := range statement.Aggregations {
			fn := aggregation.Function
			if fn.Name != "" && fn.Arg != "" {
				// Get alias name if exists or column name by default
				resultColumnName := aggregation.Function.Arg
				columnName := GetColumnName(aliasTable, resultColumnName)
				columnIndex := indexOf(gitqlObject.Titles, columnName)
				aggregationFunction := ast.Aggregations[fn.Name]
				result := aggregationFunction(fn.Arg, gitqlObject.Titles, &group)
				for _, object := range group.Rows {
					if columnIndex < len(object.Values) {
						object.Values[columnIndex] = result
					} else {
						object.Values = append(object.Values, result)
					}
				}
			}
		}
		// Resolve aggregations expressions
		for _, aggregation := range statement.Aggregations {
			if expr := aggregation.Expression; expr != nil {
				resultColumnName := aggregation.Function.Arg
				columnName := GetColumnName(aliasTable, resultColumnName)
				columnIndex := indexOf(gitqlObject.Titles, columnName)

				for _, object := range group.Rows {
					result, _ := EvaluateExpression(env, expr, gitqlObject.Titles, object.Values)
					if columnIndex < len(object.Values) {
						object.Values[columnIndex] = result
					} else {
						object.Values = append(object.Values, result)
					}
				}
			}
		}
		if groupsCount > 1 {
			group.Rows = group.Rows[:1]
		}
	}

	return nil
}

func executeGlobalVariableStatement(
	env *ast.Environment,
	statement *ast.GlobalVariableStatement,
) error {
	value, err := EvaluateExpression(env, statement.Value, []string{}, nil)
	if err != nil {
		return err
	}

	env.Globals[statement.Name] = value
	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func indexOf(slice []string, item string) int {
	for i, s := range slice {
		if s == item {
			return i
		}
	}
	return -1
}

func compare(a, b interface{}) int {
	switch a := a.(type) {
	case int:
		b := b.(int)
		if a < b {
			return -1
		} else if a > b {
			return 1
		} else {
			return 0
		}
	case float64:
		b := b.(float64)
		if a < b {
			return -1
		} else if a > b {
			return 1
		} else {
			return 0
		}
	case string:
		b := b.(string)
		if sort.StringsAreSorted([]string{a, b}) {
			return 1
		} else {
			return 0
		}
	default:
		return 0
	}
}
