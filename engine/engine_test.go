package engine

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/stretchr/testify/assert"

	"github.com/ggql/ggql/ast"
	"github.com/ggql/ggql/parser"
)

const (
	engineRepo = "ggql-engine-test.git"
	querystr   = "SELECT * FROM commits"
)

func newEngineRepo() *git.Repository {
	_, _ = git.PlainInit(engineRepo, true)
	repo, _ := git.PlainOpen(engineRepo)

	return repo
}

func deleteEngineRepo() {
	_ = os.RemoveAll(engineRepo)
}

func TestEvaluate(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	repo := newEngineRepo()
	defer deleteEngineRepo()

	repos := []*git.Repository{repo}

	tokens, errToken := parser.Tokenize(querystr)
	if errToken.Message != "" {
		t.Fatal("failed to tokenize")
	}
	query, errQuery := parser.ParserGql(tokens, &env)
	if errQuery.Message != "" {
		t.Fatal("failed to parser")
	}

	_, err := Evaluate(&env, repos, query)
	if err != nil {
		t.Fatal("failed to delete repo:", err)
	}
}

func TestEvaluateSelectQuery(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	repo := newEngineRepo()
	defer deleteEngineRepo()

	repos := []*git.Repository{repo}

	tokens, errToken := parser.Tokenize(querystr)
	if errToken.Message != "" {
		t.Fatal("failed to tokenize")
	}
	query, errQuery := parser.ParserGql(tokens, &env)
	if errQuery.Message != "" {
		t.Fatal("failed to parser")
	}

	switch query {
	case ast.Query{Select: query.Select}:
		_, err := Evaluate(&env, repos, query)
		if err != nil {
			t.Fatal("failed to delete repo:", err)
		}
	default:
		// PASS
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
