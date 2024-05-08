package engine

import (
	"fmt"
	"testing"

	"github.com/ggql/ggql/ast"
	"github.com/go-git/go-git/v5"
)

const path = "F:/tt/test.git"

func TestExecuteStatement(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	statement := &ast.SelectStatement{
		TableName:    "commits",
		FieldsNames:  []string{"commit_id", "title", "message", "name", "email", "datetime", "repo"},
		FieldsValues: []ast.Expression{},
		AliasTable:   make(map[string]string),
		IsDistinct:   false,
	}

	if err := TestNewRepo(path); err != nil {
		t.Fatal("failed to create repo:", err)
	}
	repo, err := git.PlainOpen(path)
	if err != nil {
		t.Fatal("failed to open repo")
	}

	var object ast.GitQLObject
	var table map[string]string
	selection := []string{}

	ret := ExecuteStatement(&env, statement, repo, &object, table, selection)
	if ret == nil {
		t.Log("execute statement succeeded")
	} else {
		t.Errorf("execute statement failed: %v", ret)
	}

	if err := TestDeleteRepo(path); err != nil {
		t.Fatal("failed to delete repo:", err)
	}
}

func TestExecuteSelectStatement(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	statement := &ast.SelectStatement{
		TableName:    "commits",
		FieldsNames:  []string{"commit_id", "title", "message", "name", "email", "datetime", "repo"},
		FieldsValues: []ast.Expression{},
		AliasTable:   make(map[string]string),
		IsDistinct:   false,
	}

	if err := TestNewRepo(path); err != nil {
		t.Fatal("failed to create repo:", err)
	}
	repo, err := git.PlainOpen(path)
	if err != nil {
		t.Fatal("failed to open repo")
	}

	var object ast.GitQLObject
	selections := []string{}

	ret := executeSelectStatement(&env, statement, repo, &object, selections)
	if ret == nil {
		t.Log("execute statement succeeded")
	} else {
		t.Errorf("execute statement failed: %v", ret)
	}

	if err := TestDeleteRepo(path); err != nil {
		t.Fatal("failed to delete repo:", err)
	}
}

func TestExecuteWhereStatement(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	statement := &ast.WhereStatement{
		Condition: &ast.NumberExpression{Value: ast.IntegerValue{Value: 1}},
	}

	var object ast.GitQLObject
	object.Titles = []string{"title1", "title2"}
	object.Groups = []ast.Group{
		{Rows: []ast.Row{
			{Values: []ast.Value{
				ast.IntegerValue{Value: 1}, ast.IntegerValue{Value: 2},
			}},
			{Values: []ast.Value{
				ast.IntegerValue{Value: 3}, ast.IntegerValue{Value: 4},
			}},
		}},
	}

	ret := executeWhereStatement(&env, statement, &object)
	if ret == nil {
		t.Log("execute statement succeeded")
	} else {
		t.Errorf("execute statement failed: %v", ret)
	}
}

func TestExecuteHavingStatement(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	statement := &ast.HavingStatement{
		Condition: &ast.NumberExpression{Value: ast.IntegerValue{Value: 1}},
	}
	var object ast.GitQLObject
	object.Titles = []string{"title1", "title2"}
	object.Groups = []ast.Group{
		{Rows: []ast.Row{
			{Values: []ast.Value{
				ast.IntegerValue{Value: 1}, ast.IntegerValue{Value: 2},
			}},
			{Values: []ast.Value{
				ast.IntegerValue{Value: 3}, ast.IntegerValue{Value: 4},
			}},
		}},
	}
	ret := executeHavingStatement(&env, statement, &object)
	if ret == nil {
		t.Log("execute statement succeeded")
	} else {
		t.Errorf("execute statement failed: %v", ret)
	}
}

func TestExecuteLimitStatement(t *testing.T) {
	statement := &ast.LimitStatement{
		Count: 0,
	}
	var object ast.GitQLObject
	object.Titles = []string{"title1", "title2"}
	object.Groups = []ast.Group{
		{Rows: []ast.Row{
			{Values: []ast.Value{
				ast.IntegerValue{Value: 1}, ast.IntegerValue{Value: 2},
			}},
			{Values: []ast.Value{
				ast.IntegerValue{Value: 3}, ast.IntegerValue{Value: 4},
			}},
		}},
	}
	ret := executeLimitStatement(statement, &object)
	if ret == nil {
		t.Log("execute statement succeeded")
	} else {
		t.Errorf("execute statement failed: %v", ret)
	}
}

func TestExecuteOffsetStatement(t *testing.T) {
	statement := &ast.OffsetStatement{
		Count: 0,
	}
	var object ast.GitQLObject
	object.Titles = []string{"title1", "title2"}
	object.Groups = []ast.Group{
		{Rows: []ast.Row{
			{Values: []ast.Value{
				ast.IntegerValue{Value: 1}, ast.IntegerValue{Value: 2},
			}},
			{Values: []ast.Value{
				ast.IntegerValue{Value: 3}, ast.IntegerValue{Value: 4},
			}},
		}},
	}
	ret := executeOffsetStatement(statement, &object)
	if ret == nil {
		t.Log("execute statement succeeded")
	} else {
		t.Errorf("execute statement failed: %v", ret)
	}
}

func TestExecuteOrderByStatement(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	statement := &ast.OrderByStatement{
		Arguments:     []ast.Expression{&ast.NumberExpression{Value: ast.IntegerValue{Value: 5}}},
		SortingOrders: []ast.SortingOrder{ast.Ascending},
	}
	var object ast.GitQLObject
	object.Titles = []string{"title1", "title2"}
	object.Groups = []ast.Group{
		{Rows: []ast.Row{
			{Values: []ast.Value{
				ast.IntegerValue{Value: 1}, ast.IntegerValue{Value: 2},
			}},
			{Values: []ast.Value{
				ast.IntegerValue{Value: 3}, ast.IntegerValue{Value: 4},
			}},
		}},
	}
	ret := executeOrderByStatement(&env, statement, &object)
	if ret == nil {
		t.Log("execute statement succeeded")
	} else {
		t.Errorf("execute statement failed: %v", ret)
	}
}

func TestExecuteGroupByStatement(t *testing.T) {
	statement := &ast.GroupByStatement{
		FieldName: "title1",
	}
	var object ast.GitQLObject
	object.Titles = []string{"title1", "title2"}
	object.Groups = []ast.Group{
		{Rows: []ast.Row{
			{Values: []ast.Value{
				ast.IntegerValue{Value: 1}, ast.IntegerValue{Value: 2},
			}},
			{Values: []ast.Value{
				ast.IntegerValue{Value: 3}, ast.IntegerValue{Value: 4},
			}},
		}},
	}
	ret := executeGroupByStatement(statement, &object)
	if ret == nil {
		t.Log("execute statement succeeded")
	} else {
		t.Errorf("execute statement failed: %v", ret)
	}
}

func TestExecuteAggregationFunctionStatement(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	statement := &ast.AggregationsStatement{
		Aggregations: make(map[string]ast.AggregateValue),
	}
	var a ast.AggregateValue
	a.Function.Name = "max"
	a.Function.Arg = "title3"

	statement.Aggregations["title"] = a
	var object ast.GitQLObject
	object.Titles = []string{"title1", "title2"}
	object.Groups = []ast.Group{
		{Rows: []ast.Row{
			{Values: []ast.Value{
				ast.IntegerValue{Value: 1}, ast.IntegerValue{Value: 2},
			}},
			{Values: []ast.Value{
				ast.IntegerValue{Value: 3}, ast.IntegerValue{Value: 4},
			}},
		}},
	}
	table := make(map[string]string)
	table["title"] = "title1"
	fmt.Println("5")
	ret := executeAggregationFunctionStatement(&env, statement, &object, table)
	fmt.Println("3")
	if ret == nil {
		t.Log("execute statement succeeded")
	} else {
		t.Errorf("execute statement failed: %v", ret)
	}
}

func TestExecuteGlobalVariableStatement(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	statement := &ast.GlobalVariableStatement{
		Name:  "name",
		Value: &ast.NumberExpression{Value: ast.IntegerValue{Value: 1}},
	}

	ret := executeGlobalVariableStatement(&env, statement)
	if ret == nil {
		t.Log("execute statement succeeded")
	} else {
		t.Errorf("execute statement failed: %v", ret)
	}
}
