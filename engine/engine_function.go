package engine

import (
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"

	"github.com/ggql/ggql/ast"
)

func SelectGQLObjects(
	env *ast.Environment,
	repo *git.Repository,
	table string,
	fieldsNames []string,
	titles []string,
	fieldsValues []ast.Expression,
) (*ast.Group, error) {
	switch table {
	case "refs":
		return selectReferences(env, repo, fieldsNames, titles, fieldsValues)
	case "commits":
		return selectCommits(env, repo, fieldsNames, titles, fieldsValues)
	case "branches":
		return selectBranches(env, repo, fieldsNames, titles, fieldsValues)
	case "diffs":
		return selectDiffs(env, repo, fieldsNames, titles, fieldsValues)
	case "tags":
		return selectTags(env, repo, fieldsNames, titles, fieldsValues)
	default:
		return selectValues(env, titles, fieldsValues)
	}
}

// nolint:goconst
func selectReferences(
	env *ast.Environment,
	repo *git.Repository,
	fieldsNames []string,
	titles []string,
	fieldsValues []ast.Expression,
) (*ast.Group, error) {
	var rows []ast.Row

	storer, _ := repo.Storer.(*filesystem.Storage)
	repoPath := storer.Filesystem().Root()

	gitReferences, err := repo.References()
	if err != nil {
		return &ast.Group{Rows: rows}, nil
	}

	namesLen := int64(len(fieldsNames))
	valuesLen := int64(len(fieldsValues))
	padding := namesLen - valuesLen

	_ = gitReferences.ForEach(func(ref *plumbing.Reference) error {
		var values []ast.Value
		for index := int64(0); index < namesLen; index++ {
			fieldName := fieldsNames[index]
			if index-padding >= 0 {
				value := fieldsValues[index-padding]
				if _, ok := value.(*ast.SymbolExpression); !ok {
					evaluated, _ := EvaluateExpression(env, value, titles, values)
					values = append(values, evaluated)
					continue
				}
			}
			switch fieldName {
			case "name":
				name := ref.Name().Short()
				values = append(values, ast.TextValue{Value: name})
			case "full_name":
				fullName := ref.Name().String()
				values = append(values, ast.TextValue{Value: fullName})
			case "type":
				var nameType string
				name := ref.Name()
				if name.IsBranch() {
					nameType = "branch"
				} else if name.IsRemote() {
					nameType = "remote"
				} else if name.IsTag() {
					nameType = "tag"
				} else if name.IsNote() {
					nameType = "note"
				} else {
					nameType = "other"
				}
				value := ast.TextValue{Value: nameType}
				values = append(values, value)
			case "repo":
				value := ast.TextValue{Value: repoPath}
				values = append(values, value)
			default:
				value := ast.NullValue{}
				values = append(values, value)
			}
		}
		row := ast.Row{Values: values}
		rows = append(rows, row)
		return nil
	})

	return &ast.Group{Rows: rows}, nil
}

// nolint:goconst
func selectCommits(
	env *ast.Environment,
	repo *git.Repository,
	fieldsNames []string,
	titles []string,
	fieldsValues []ast.Expression,
) (*ast.Group, error) {
	var rows []ast.Row

	storer, _ := repo.Storer.(*filesystem.Storage)
	repoPath := storer.Filesystem().Root()

	namesLen := int64(len(fieldsNames))
	valuesLen := int64(len(fieldsValues))
	padding := namesLen - valuesLen

	commitObjects, err := repo.CommitObjects()
	if err != nil {
		return &ast.Group{Rows: rows}, nil
	}

	_ = commitObjects.ForEach(func(commit *object.Commit) error {
		var values []ast.Value
		for index := int64(0); index < namesLen; index++ {
			fieldName := fieldsNames[index]
			if index-padding >= 0 {
				value := fieldsValues[index-padding]
				if _, ok := value.(*ast.SymbolExpression); !ok {
					evaluated, _ := EvaluateExpression(env, value, titles, values)
					values = append(values, evaluated)
					continue
				}
			}
			switch fieldName {
			case "commit_id":
				commitID := commit.ID().String()
				values = append(values, ast.TextValue{Value: commitID})
			case "name":
				name := commit.Author.Name
				values = append(values, ast.TextValue{Value: name})
			case "email":
				email := commit.Author.Email
				values = append(values, ast.TextValue{Value: email})
			case "title":
				summary := strings.Split(commit.Message, "\n\n")[0]
				values = append(values, ast.TextValue{Value: summary})
			case "message":
				message := commit.Message
				values = append(values, ast.TextValue{Value: message})
			case "datetime":
				timeStamp := commit.Author.When.Unix()
				values = append(values, ast.DateTimeValue{Value: timeStamp})
			case "repo":
				value := ast.TextValue{Value: repoPath}
				values = append(values, value)
			default:
				value := ast.NullValue{}
				values = append(values, value)
			}
		}
		row := ast.Row{Values: values}
		rows = append(rows, row)
		return nil
	})

	return &ast.Group{Rows: rows}, nil
}

func selectBranches(
	env *ast.Environment,
	repo *git.Repository,
	fieldsNames []string,
	titles []string,
	fieldsValues []ast.Expression,
) (*ast.Group, error) {
	helper := func(ref *plumbing.Reference) int64 {
		var count int64
		commit, err := repo.CommitObject(ref.Hash())
		if err != nil {
			return -1
		}
		commitIter := object.NewCommitIterCTime(commit, nil, nil)
		if err = commitIter.ForEach(func(c *object.Commit) error {
			count++
			return nil
		}); err != nil {
			return -1
		}
		return count
	}

	var rows []ast.Row

	storer, _ := repo.Storer.(*filesystem.Storage)
	repoPath := storer.Filesystem().Root()

	localAndRemoteBranches, _ := repo.References()
	headRef, _ := repo.Head()

	namesLen := int64(len(fieldsNames))
	valuesLen := int64(len(fieldsValues))
	padding := namesLen - valuesLen

	_ = localAndRemoteBranches.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.InvalidReference {
			return nil
		}
		if !ref.Name().IsBranch() && !ref.Name().IsRemote() {
			return nil
		}
		var values []ast.Value
		for index := int64(0); index < namesLen; index++ {
			fieldName := fieldsNames[index]
			if index-padding >= 0 {
				value := fieldsValues[index-padding]
				if _, ok := value.(*ast.SymbolExpression); !ok {
					evaluated, _ := EvaluateExpression(env, value, titles, values)
					values = append(values, evaluated)
					continue
				}
			}
			switch fieldName {
			case "name":
				branchName := ref.Name().String()
				values = append(values, ast.TextValue{Value: branchName})
			case "commit_count":
				commitCount := helper(ref)
				values = append(values, ast.IntegerValue{Value: commitCount})
			case "is_head":
				isHead := ref.Hash().String() == headRef.Hash().String()
				values = append(values, ast.BooleanValue{Value: isHead})
			case "is_remote":
				isRemote := ref.Name().IsRemote()
				values = append(values, ast.BooleanValue{Value: isRemote})
			case "repo":
				values = append(values, ast.TextValue{Value: repoPath})
			default:
				value := ast.NullValue{}
				values = append(values, value)
			}
		}
		row := ast.Row{Values: values}
		rows = append(rows, row)
		return nil
	})

	return &ast.Group{Rows: rows}, nil
}

// nolint:goconst,gocyclo
func selectDiffs(
	env *ast.Environment,
	repo *git.Repository,
	fieldsNames []string,
	titles []string,
	fieldsValues []ast.Expression,
) (*ast.Group, error) {
	var rows []ast.Row

	storer, _ := repo.Storer.(*filesystem.Storage)
	repoPath := storer.Filesystem().Root()

	namesLen := int64(len(fieldsNames))
	valuesLen := int64(len(fieldsValues))
	padding := namesLen - valuesLen

	commitObjects, err := repo.CommitObjects()
	if err != nil {
		return &ast.Group{Rows: rows}, nil
	}

	_ = commitObjects.ForEach(func(commit *object.Commit) error {
		var values []ast.Value
		for index := int64(0); index < namesLen; index++ {
			fieldName := fieldsNames[index]
			if index-padding >= 0 {
				value := fieldsValues[index-padding]
				if _, ok := value.(*ast.SymbolExpression); !ok {
					evaluated, _ := EvaluateExpression(env, value, titles, values)
					values = append(values, evaluated)
					continue
				}
			}
			switch fieldName {
			case "commit_id":
				commitID := commit.ID().String()
				values = append(values, ast.TextValue{Value: commitID})
			case "name":
				name := commit.Author.Name
				values = append(values, ast.TextValue{Value: name})
			case "email":
				email := commit.Author.Email
				values = append(values, ast.TextValue{Value: email})
			case "repo":
				values = append(values, ast.TextValue{Value: repoPath})
			case "insertions", "deletions", "files_changed":
				var insertions, deletions, filesChanged int64
				current := commit
				previous := commit.Parents()
				selectInsertionsOrDeletions := fieldName == "insertions" || fieldName == "deletions"
				_ = previous.ForEach(func(commit *object.Commit) error {
					patch, _ := commit.Patch(current)
					for _, stat := range patch.Stats() {
						filesChanged += 1
						if selectInsertionsOrDeletions {
							insertions += int64(stat.Addition)
							deletions += int64(stat.Deletion)
						}
					}
					return nil
				})
				if fieldName == "insertions" {
					values = append(values, ast.IntegerValue{Value: insertions})
				} else if fieldName == "deletions" {
					values = append(values, ast.IntegerValue{Value: deletions})
				} else if fieldName == "files_changed" {
					values = append(values, ast.IntegerValue{Value: filesChanged})
				}
			default:
				value := ast.NullValue{}
				values = append(values, value)
			}
		}
		row := ast.Row{Values: values}
		rows = append(rows, row)
		return nil
	})

	return &ast.Group{Rows: rows}, nil
}

func selectTags(
	env *ast.Environment,
	repo *git.Repository,
	fieldsNames []string,
	titles []string,
	fieldsValues []ast.Expression,
) (*ast.Group, error) {
	var rows []ast.Row

	tags, err := repo.Tags()
	if err != nil {
		return &ast.Group{Rows: rows}, nil
	}

	storer, _ := repo.Storer.(*filesystem.Storage)
	repoPath := storer.Filesystem().Root()

	namesLen := int64(len(fieldsNames))
	valuesLen := int64(len(fieldsValues))
	padding := namesLen - valuesLen

	_ = tags.ForEach(func(ref *plumbing.Reference) error {
		var values []ast.Value
		for index := int64(0); index < namesLen; index++ {
			fieldName := fieldsNames[index]
			if index-padding >= 0 {
				value := fieldsValues[index-padding]
				if _, ok := value.(*ast.SymbolExpression); !ok {
					evaluated, _ := EvaluateExpression(env, value, titles, values)
					values = append(values, evaluated)
					continue
				}
			}
			switch fieldName {
			case "name":
				if ref.Name().IsTag() {
					tagName := ref.Name().Short()
					values = append(values, ast.TextValue{Value: tagName})
				}
			case "repo":
				value := ast.TextValue{Value: repoPath}
				values = append(values, value)
			default:
				value := ast.NullValue{}
				values = append(values, value)
			}
		}
		row := ast.Row{Values: values}
		rows = append(rows, row)
		return nil
	})

	return &ast.Group{Rows: rows}, nil
}

func selectValues(
	env *ast.Environment,
	titles []string,
	fieldsValues []ast.Expression,
) (*ast.Group, error) {
	var group ast.Group
	var values []ast.Value

	for _, value := range fieldsValues {
		evaluated, _ := EvaluateExpression(env, value, titles, values)
		values = append(values, evaluated)
	}

	row := ast.Row{Values: values}
	group.Rows = append(group.Rows, row)

	return &group, nil
}

func GetColumnName(aliasTable map[string]string, name string) string {
	if columnName, ok := aliasTable[name]; ok {
		return columnName
	}

	return name
}
