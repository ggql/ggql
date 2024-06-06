package engine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"

	"github.com/ggql/ggql/ast"
)

var gqlCommandsInOrder = []string{
	"select",
	"where",
	"group",
	"aggregation",
	"having",
	"order",
	"offset",
	"limit",
}

type EvaluationResult struct {
	SelectedGroups struct {
		Obj ast.GitQLObject
		Str []string
	}
	SetGlobalVariable bool
}

func Evaluate(env *ast.Environment, repos []*git.Repository, query ast.Query) (EvaluationResult, error) {
	if query.Select != nil {
		return EvaluateSelectQuery(env, repos, *query.Select)
	}

	if query.GlobalVariableDeclaration != nil {
		err := executeGlobalVariableStatement(env, query.GlobalVariableDeclaration)
		if err != nil {
			return EvaluationResult{}, err
		}
		return EvaluationResult{}, nil
	}

	return EvaluationResult{}, fmt.Errorf("unknown query type")
}

// nolint:gocyclo
func EvaluateSelectQuery(
	env *ast.Environment,
	repos []*git.Repository,
	query ast.GQLQuery,
) (EvaluationResult, error) {
	var gitqlObject ast.GitQLObject
	aliasTable := make(map[string]string)

	hiddenSelections := query.HiddenSelections
	statementsMap := query.Statements
	firstRepo := repos[0]

	for _, gqlCommand := range gqlCommandsInOrder {
		if statement, ok := statementsMap[gqlCommand]; ok {
			switch gqlCommand {
			case "select":
				selectStatement := statement.(*ast.SelectStatement)

				if selectStatement.TableName == "" {
					err := ExecuteStatement(env, statement, firstRepo, &gitqlObject, aliasTable, hiddenSelections)
					if err != nil {
						return EvaluationResult{}, err
					}

					if gitqlObject.IsEmpty() || gitqlObject.Groups[0].IsEmpty() {
						return EvaluationResult{SelectedGroups: struct {
							Obj ast.GitQLObject
							Str []string
						}{
							Obj: gitqlObject,
							Str: hiddenSelections,
						},
						}, nil
					}

					continue
				}

				for _, repo := range repos {
					err := ExecuteStatement(env, statement, repo, &gitqlObject, aliasTable, hiddenSelections)
					if err != nil {
						return EvaluationResult{SelectedGroups: struct {
							Obj ast.GitQLObject
							Str []string
						}{
							Obj: gitqlObject,
							Str: hiddenSelections,
						},
						}, err
					}
				}

				if gitqlObject.IsEmpty() || gitqlObject.Groups[0].IsEmpty() {
					return EvaluationResult{SelectedGroups: struct {
						Obj ast.GitQLObject
						Str []string
					}{
						Obj: gitqlObject,
						Str: hiddenSelections,
					},
					}, nil
				}

				if selectStatement.TableName != "" && selectStatement.IsDistinct {
					ApplyDistinctOnObjectsGroup(&gitqlObject, hiddenSelections)
				}

			default:
				err := ExecuteStatement(env, statement, firstRepo, &gitqlObject, aliasTable, hiddenSelections)
				if err != nil {
					return EvaluationResult{}, err
				}
			}
		}
	}

	if len(gitqlObject.Groups) > 1 {
		for _, group := range gitqlObject.Groups {
			if len(group.Rows) > 1 {
				group.Rows = group.Rows[:1]
			}
		}
	} else if len(gitqlObject.Groups) == 1 && !query.HasGroupByStatement && query.HasAggregationFunction {
		group := &gitqlObject.Groups[0]
		if len(group.Rows) > 1 {
			group.Rows = group.Rows[:1]
		}
	}

	return EvaluationResult{SelectedGroups: struct {
		Obj ast.GitQLObject
		Str []string
	}{
		Obj: gitqlObject,
		Str: hiddenSelections,
	},
	}, nil
}

func ApplyDistinctOnObjectsGroup(gitqlObject *ast.GitQLObject, hiddenSelections []string) {
	if gitqlObject.IsEmpty() {
		return
	}

	titles := make([]string, 0, len(gitqlObject.Titles))
	for _, title := range gitqlObject.Titles {
		if !contains(hiddenSelections, title) {
			titles = append(titles, title)
		}
	}

	titlesCount := len(titles)

	objects := gitqlObject.Groups[0].Rows
	newObjects := ast.Group{Rows: make([]ast.Row, 0)}
	valuesSet := make(map[string]bool)

	for _, object := range objects {
		rowValues := make([]string, titlesCount)
		for index := 0; index < len(titles); index++ {
			rowValues = append(rowValues, object.Values[index].Fmt())
		}

		hash := sha256.New()
		hash.Write([]byte(strings.Join(rowValues, "")))
		hashBytes := hash.Sum(nil)
		valuesHash := hex.EncodeToString(hashBytes)

		if valuesSet[valuesHash] {
			continue
		}

		valuesSet[valuesHash] = true
		newObjects.Rows = append(newObjects.Rows, ast.Row{Values: object.Values})
	}

	if len(objects) != newObjects.Len() {
		gitqlObject.Groups[0].Rows = newObjects.Rows
	}
}
