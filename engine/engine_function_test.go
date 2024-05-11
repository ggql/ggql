package engine

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"

	"github.com/ggql/ggql/ast"
)

const (
	functionFile = "ggql-engine-function-test.txt"
	functionRepo = "ggql-engine-function-test.git"

	tableName = "refs"
)

func newFunctionRepo() *git.Repository {
	// Create a new repository
	_, _ = git.PlainInit(functionRepo, false)
	repo, _ := git.PlainOpen(functionRepo)

	tree, _ := repo.Worktree()

	// Create a new file
	filePath := filepath.Join(tree.Filesystem.Root(), functionFile)
	file, _ := os.Create(filePath)
	_, _ = file.WriteString("hello world")
	_ = file.Close()

	// Create a new commit
	_, _ = tree.Add(functionFile)
	commit, _ := tree.Commit("Adding "+functionFile, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "name",
			Email: "name@example.com",
			When:  time.Now(),
		},
	})

	_, _ = repo.CommitObject(commit)

	// Create a new tag
	ref, _ := repo.Head()
	_, _ = repo.CreateTag("v0.0.1", ref.Hash(), &git.CreateTagOptions{
		Message: "Create tag",
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

	table := tableName
	fieldsNames := []string{"name", "full_name", "type", "repo"}

	titles := []string{"title"}
	fieldsValues := []ast.Expression{
		&ast.StringExpression{
			Value:     "value",
			ValueType: ast.StringValueText,
		},
	}

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

	fieldsNames := []string{"name", "full_name", "type", "repo"}
	titles := []string{"title"}
	fieldsValues := []ast.Expression{
		&ast.SymbolExpression{
			Value: "value",
		},
	}

	group, err := selectReferences(&env, repo, fieldsNames, titles, fieldsValues)
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(group.Rows))
	assert.Equal(t, len(fieldsNames), len(group.Rows[0].Values))
}

func TestSelectCommits(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	repo := newFunctionRepo()
	defer deleteFunctionRepo()

	fieldsNames := []string{"commit_id", "name", "email", "title", "message", "datetime", "repo"}
	titles := []string{"title"}
	fieldsValues := []ast.Expression{
		&ast.SymbolExpression{
			Value: "value",
		},
	}

	group, err := selectCommits(&env, repo, fieldsNames, titles, fieldsValues)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(group.Rows))
	assert.Equal(t, len(fieldsNames), len(group.Rows[0].Values))
}

func TestSelectBranches(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	repo := newFunctionRepo()
	defer deleteFunctionRepo()

	fieldsNames := []string{"name", "commit_count", "is_head", "is_remote", "repo"}
	titles := []string{"title"}
	fieldsValues := []ast.Expression{
		&ast.SymbolExpression{
			Value: "value",
		},
	}

	group, err := selectBranches(&env, repo, fieldsNames, titles, fieldsValues)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(group.Rows))
	assert.Equal(t, len(fieldsNames), len(group.Rows[0].Values))
}

func TestSelectDiffs(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	repo := newFunctionRepo()
	defer deleteFunctionRepo()

	fieldsNames := []string{"commit_id", "name", "email", "repo", "insertions", "deletions", "files_changed"}
	titles := []string{"title"}
	fieldsValues := []ast.Expression{
		&ast.SymbolExpression{
			Value: "value",
		},
	}

	group, err := selectDiffs(&env, repo, fieldsNames, titles, fieldsValues)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(group.Rows))
	assert.Equal(t, len(fieldsNames), len(group.Rows[0].Values))
}

func TestSelectTags(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	repo := newFunctionRepo()
	defer deleteFunctionRepo()

	fieldsNames := []string{"name", "repo"}
	titles := []string{"title"}
	fieldsValues := []ast.Expression{
		&ast.SymbolExpression{
			Value: "value",
		},
	}

	group, err := selectTags(&env, repo, fieldsNames, titles, fieldsValues)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(group.Rows))
	assert.Equal(t, len(fieldsNames), len(group.Rows[0].Values))
}

func TestSelectValues(t *testing.T) {
	env := ast.Environment{
		Globals:      map[string]ast.Value{},
		GlobalsTypes: map[string]ast.DataType{},
		Scopes:       map[string]ast.DataType{},
	}

	titles := []string{"title"}
	fieldsValues := []ast.Expression{
		&ast.SymbolExpression{
			Value: "value",
		},
	}

	group, err := selectValues(&env, titles, fieldsValues)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(group.Rows))
}

func TestGetColumnName(t *testing.T) {
	aliasTable := map[string]string{
		"key": "value",
	}

	name := "key"
	ret := GetColumnName(aliasTable, name)
	assert.Equal(t, aliasTable[name], ret)

	name = "invalid"
	ret = GetColumnName(aliasTable, name)
	assert.Equal(t, name, ret)
}
