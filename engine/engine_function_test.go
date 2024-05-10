package engine

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"

	"github.com/ggql/ggql/ast"
)

const (
	functionRepo = "ggql-engine-function-test.git"

	tableCommits = "commits"
)

func newFunctionRepo() *git.Repository {
	_, _ = git.PlainInit(functionRepo, true)
	repo, _ := git.PlainOpen(functionRepo)

	branch := config.Branch{
		Name: "master",
	}
	_ = repo.CreateBranch(&branch)

	name := "v0.0.1"
	hash := plumbing.Hash{}
	_, _ = repo.CreateTag(name, hash, nil)

	_, _ = repo.CreateRemote(&config.RemoteConfig{
		Name: "function-remote",
		URLs: []string{"https://github.com/function-project/" + functionRepo},
	})

	return repo
}

func deleteFunctionRepo() {
	_ = os.RemoveAll(functionRepo)
}

func TestSelectGQLObjects(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	repo := newFunctionRepo()
	defer deleteFunctionRepo()

	table := tableCommits
	fieldsNames := []string{"commit_id", "name", "email", "title", "message", "datetime", "repo"}

	var titles []string
	var fieldsValues []ast.Expression

	_, err := SelectGQLObjects(&env, repo, table, fieldsNames, titles, fieldsValues)
	assert.Equal(t, nil, err)
}

func TestSelectReferences(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	repo := newFunctionRepo()
	defer deleteFunctionRepo()

	fieldsNames := []string{"name"}
	titles := []string{"value"}
	fieldsValues := []ast.Expression{
		&ast.StringExpression{
			Value:     "value",
			ValueType: ast.StringValueText,
		},
	}

	group, err := selectReferences(&env, repo, fieldsNames, titles, fieldsValues)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(group.Rows))
	assert.Equal(t, 1, len(group.Rows[0].Values))
	assert.Equal(t, "value", group.Rows[0].Values[0].Fmt())
	assert.Equal(t, 1, len(group.Rows[1].Values))
	assert.Equal(t, "value", group.Rows[1].Values[0].Fmt())

	fieldsNames = []string{"name"}
	titles = []string{"value"}
	fieldsValues = []ast.Expression{
		&ast.SymbolExpression{
			Value: "value",
		},
	}

	group, err = selectReferences(&env, repo, fieldsNames, titles, fieldsValues)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(group.Rows))
	assert.Equal(t, 1, len(group.Rows[0].Values))
	assert.Equal(t, "v0.0.1", group.Rows[0].Values[0].Fmt())
	assert.Equal(t, 1, len(group.Rows[1].Values))
	assert.Equal(t, "HEAD", group.Rows[1].Values[0].Fmt())

	fieldsNames = []string{"full_name"}
	titles = []string{"value"}
	fieldsValues = []ast.Expression{
		&ast.SymbolExpression{
			Value: "value",
		},
	}

	group, err = selectReferences(&env, repo, fieldsNames, titles, fieldsValues)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(group.Rows))
	assert.Equal(t, 1, len(group.Rows[0].Values))
	assert.Equal(t, "refs/tags/v0.0.1", group.Rows[0].Values[0].Fmt())
	assert.Equal(t, 1, len(group.Rows[1].Values))
	assert.Equal(t, "HEAD", group.Rows[1].Values[0].Fmt())

	fieldsNames = []string{"type"}
	titles = []string{"value"}
	fieldsValues = []ast.Expression{
		&ast.SymbolExpression{
			Value: "value",
		},
	}

	group, err = selectReferences(&env, repo, fieldsNames, titles, fieldsValues)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(group.Rows))
	assert.Equal(t, 1, len(group.Rows[0].Values))
	assert.Equal(t, "tag", group.Rows[0].Values[0].Fmt())
	assert.Equal(t, 1, len(group.Rows[1].Values))
	assert.Equal(t, "other", group.Rows[1].Values[0].Fmt())
}

func TestSelectCommits(t *testing.T) {
}

func TestSelectBranches(t *testing.T) {
}

func TestSelectDiffs(t *testing.T) {
}

func TestSelectTags(t *testing.T) {
}

func TestSelectValues(t *testing.T) {
}

func TestGetColumnName(t *testing.T) {
}
