package engine

import (
	"testing"

	"github.com/ggql/ggql/ast"
	"github.com/ggql/ggql/parser"
	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"
)

const querystr = "SELECT * FROM commits"

func TestEvaluate(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}
	if err := TestNewRepo(path); err != nil {
		t.Fatal("failed to create repo:", err)
	}
	repo, err := git.PlainOpen(path)
	if err != nil {
		t.Fatal("failed to open repo")
	}
	repos := []*git.Repository{repo}

	tokens, errtoken := parser.Tokenize(querystr)
	if errtoken.Message != "" {
		t.Fatal("failed to tokenize")
	}
	query, err2 := parser.ParserGql(tokens, &env)
	if err2.Message != "" {
		t.Fatal("failed to parser")
	}

	_, errret := Evaluate(&env, repos, query)
	if errret != nil {
		if err := TestDeleteRepo(path); err != nil {
			t.Fatal("failed to delete repo:", err)
		}
	}

	if err := TestDeleteRepo(path); err != nil {
		t.Fatal("failed to delete repo:", err)
	}
}

func TestEvaluateSelectQuery(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	if err := TestNewRepo(path); err != nil {
		t.Fatal("failed to create repo:", err)
	}
	repo, err := git.PlainOpen(path)
	if err != nil {
		t.Fatal("failed to open repo")
	}
	repos := []*git.Repository{repo}

	tokens, errtoken := parser.Tokenize(querystr)
	if errtoken.Message != "" {
		t.Fatal("failed to tokenize")
	}
	query, err2 := parser.ParserGql(tokens, &env)
	if err2.Message != "" {
		t.Fatal("failed to parser")
	}

	switch query {
	case ast.Query{Select: query.Select}:
		_, errret := Evaluate(&env, repos, query)
		if errret != nil {
			if err := TestDeleteRepo(path); err != nil {
				t.Fatal("failed to delete repo:", err)
			}
		}
	default:
		if err := TestDeleteRepo(path); err != nil {
			t.Fatal("failed to delete repo:", err)
		}
	}
	if err := TestDeleteRepo(path); err != nil {
		t.Fatal("failed to delete repo:", err)
	}
}

func TestApplyDistinctOnObjectsGroup(t *testing.T) {
	object := ast.GitQLObject{
		Titles: []string{"title1", "title2"},
		Groups: []ast.Group{
			{Rows: []ast.Row{
				{Values: []ast.Value{
					ast.IntegerValue{Value: 1}, ast.IntegerValue{Value: 2},
				}},
				{Values: []ast.Value{
					ast.IntegerValue{Value: 3}, ast.IntegerValue{Value: 4},
				}},
			}},
		},
	}

	var selections []string
	ApplyDistinctOnObjectsGroup(&object, selections)
	assert.Equal(t, len(object.Groups[0].Rows), 2)

	object2 := ast.GitQLObject{
		Titles: []string{"title1", "title2"},
		Groups: []ast.Group{
			{Rows: []ast.Row{
				{Values: []ast.Value{
					ast.IntegerValue{Value: 1}, ast.IntegerValue{Value: 2},
				}},
				{Values: []ast.Value{
					ast.IntegerValue{Value: 1}, ast.IntegerValue{Value: 2},
				}},
			}},
		},
	}
	var selections2 []string
	ApplyDistinctOnObjectsGroup(&object2, selections2)
	assert.Equal(t, len(object2.Groups[0].Rows), 1)
}
