package main

import (
    "errors"
    "fmt"
    "strings"
)

type Environment struct {
    // implementation details
}

type Group struct {
    Rows []Row
}

type Row struct {
    Values []Value
}

type Value interface {
    // implementation details
}

type Expression interface {
    // implementation details
}

type SymbolExpression struct {
    // implementation details
}

type Repository struct {
    // implementation details
}

type Category int

const (
    LocalBranch Category = iota
    RemoteBranch
    Tag
    Note
)

func SelectGQLObjects(
    env *Environment,
    repo *Repository,
    table string,
    fieldsNames []string,
    titles []string,
    fieldsValues []Expression,
) (*Group, error) {
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

func selectReferences(
    env *Environment,
    repo *Repository,
    fieldsNames []string,
    titles []string,
    fieldsValues []Expression,
) (*Group, error) {
    repoPath := repo.Path()

    var rows []Row

    gitReferences, err := repo.References()
    if err != nil {
        return &Group{Rows: rows}, nil
    }

    references := gitReferences.All()
    namesLen := int64(len(fieldsNames))
    valuesLen := int64(len(fieldsValues))
    padding := namesLen - valuesLen

    for _, reference := range references {
        var values []Value

        for index := int64(0); index < namesLen; index++ {
            fieldName := fieldsNames[index]

            if index-padding >= 0 {
                value := fieldsValues[index-padding]
                if _, ok := value.(*SymbolExpression); !ok {
                    evaluated, err := evaluateExpression(env, value, titles, values)
                    if err != nil {
                        return nil, err
                    }
                    values = append(values, evaluated)
                    continue
                }
            }

            switch fieldName {
            case "name":
                name := reference.Name().CategoryAndShortName().ShortName
                values = append(values, TextValue(name))
            case "full_name":
                fullName := reference.Name().AsBstr().String()
                values = append(values, TextValue(fullName))
            case "type":
                category := reference.Name().Category()
                var valueType string
                switch category {
                case LocalBranch:
                    valueType = "branch"
                case RemoteBranch:
                    valueType = "remote"
                case Tag:
                    valueType = "tag"
                case Note:
                    valueType = "note"
                default:
                    valueType = "other"
                }
                values = append(values, TextValue(valueType))
            case "repo":
                values = append(values, TextValue(repoPath))
            default:
                values = append(values, NullValue())
            }
        }

        row := Row{Values: values}
        rows = append(rows, row)
    }

    return &Group{Rows: rows}, nil
}

func selectCommits(
    env *Environment,
    repo *Repository,
    fieldsNames []string,
    titles []string,
    fieldsValues []Expression,
) (*Group, error) {
    repoPath := repo.Path()

    var rows []Row

    headID, err := repo.HeadID()
    if err != nil {
        return &Group{Rows: rows}, nil
    }

    revwalk := headID.Ancestors().All()

    namesLen := int64(len(fieldsNames))
    valuesLen := int64(len(fieldsValues))
    padding := namesLen - valuesLen

    for commitInfo := range revwalk {
        commit := commitInfo.ID().Object().Commit()

        var values []Value

        for index := int64(0); index < namesLen; index++ {
            fieldName := fieldsNames[index]

            if index-padding >= 0 {
                value := fieldsValues[index-padding]
                if _, ok := value.(*SymbolExpression); !ok {
                    evaluated, err := evaluateExpression(env, value, titles, values)
                    if err != nil {
                        return nil, err
                    }
                    values = append(values, evaluated)
                    continue
                }
            }

            switch fieldName {
            case "commit_id":
                commitID := commitInfo.ID().String()
                values = append(values, TextValue(commitID))
            case "name":
                name := commit.Author().Name
                values = append(values, TextValue(name))
            case "email":
                email := commit.Author().Email
                values = append(values, TextValue(email))
            case "title":
                summary := commit.Message().Summary()
                values = append(values, TextValue(summary))
            case "message":
                message := commit.Message().String()
                values = append(values, TextValue(message))
            case "datetime":
                timeStamp := commitInfo.CommitTime()
                values = append(values, DateTimeValue(timeStamp))
            case "repo":
                values = append(values, TextValue(repoPath))
            default:
                values = append(values, NullValue())
            }
        }

        row := Row{Values: values}
        rows = append(rows, row)
    }

    return &Group{Rows: rows}, nil
}

func selectBranches(
    env *Environment,
    repo *Repository,
    fieldsNames []string,
    titles []string,
    fieldsValues []Expression,
) (*Group, error) {
    var rows []Row

    repoPath := repo.Path()
    platform := repo.References()
    localBranches := platform.LocalBranches()
    remoteBranches := platform.RemoteBranches()
    localAndRemoteBranches := append(localBranches, remoteBranches...)
    headRef, err := repo.HeadRef()
    if err != nil {
        return &Group{Rows: rows}, nil
    }
    if headRef == nil {
        return &Group{Rows: rows}, nil
    }

    namesLen := int64(len(fieldsNames))
    valuesLen := int64(len(fieldsValues))
    padding := namesLen - valuesLen

    for _, branch := range localAndRemoteBranches {
        var values []Value

        for index := int64(0); index < namesLen; index++ {
            fieldName := fieldsNames[index]

            if index-padding >= 0 {
                value := fieldsValues[index-padding]
                if _, ok := value.(*SymbolExpression); !ok {
                    evaluated, err := evaluateExpression(env, value, titles, values)
                    if err != nil {
                        return nil, err
                    }
                    values = append(values, evaluated)
                    continue
                }
            }

            switch fieldName {
            case "name":
                branchName := branch.Name().AsBstr().String()
                values = append(values, TextValue(branchName))
            case "commit_count":
                commitCount := -1
                if id := branch.TryID(); id != nil {
                    if revwalk, err := id.Ancestors().All(); err == nil {
                        commitCount = revwalk.Count()
                    }
                }
                values = append(values, IntegerValue(int64(commitCount)))
            case "is_head":
                isHead := branch.Inner() == headRef.Inner()
                values = append(values, BooleanValue(isHead))
            case "is_remote":
                isRemote := branch.Name().Category() == RemoteBranch
                values = append(values, BooleanValue(isRemote))
            case "repo":
                values = append(values, TextValue(repoPath))
            default:
                values = append(values, NullValue())
            }
        }

        row := Row{Values: values}
        rows = append(rows, row)
    }

    return &Group{Rows: rows}, nil
}

func selectDiffs(
    env *Environment,
    repo *Repository,
    fieldsNames []string,
    titles []string,
    fieldsValues []Expression,
) (*Group, error) {
    var rows []Row

    repo.ObjectCacheSizeIfUnset(4 * 1024 * 1024)
    revwalk := repo.HeadID().Ancestors().All()
    repoPath := repo.Path()

    namesLen := int64(len(fieldsNames))
    valuesLen := int64(len(fieldsValues))
    padding := namesLen - valuesLen

    for commitInfo := range revwalk {
        commit := commitInfo.ID().Object().Commit()

        var values []Value

        for index := int64(0); index < namesLen; index++ {
            fieldName := fieldsNames[index]

            if index-padding >= 0 {
                value := fieldsValues[index-padding]
                if _, ok := value.(*SymbolExpression); !ok {
                    evaluated, err := evaluateExpression(env, value, titles, values)
                    if err != nil {
                        return nil, err
                    }
                    values = append(values, evaluated)
                    continue
                }
            }

            switch fieldName {
            case "commit_id":
                values = append(values, TextValue(commitInfo.ID().String()))
            case "name":
                name := commit.Author().Name
                values = append(values, TextValue(name))
            case "email":
                email := commit.Author().Email
                values = append(values, TextValue(email))
            case "repo":
                values = append(values, TextValue(repoPath))
            case "insertions", "deletions", "files_changed":
                current := commit.Tree()
                previous := commitInfo.ParentIDs().Next().Object().Commit().Tree()

                selectInsertionsOrDeletions := fieldName == "insertions" || fieldName == "deletions"

                var insertions, deletions, filesChanged int64

                previous.Changes().ForEachToObtainTreeWithCache(
                    current,
                    func(change gix.ObjectTreeDiffChange) error {
                        filesChanged += int64(change.Event.EntryMode().IsNoTree())
                        if selectInsertionsOrDeletions {
                            platform, err := change.Diff()
                            if err == nil {
                                if counts, err := platform.LineCounts(); err == nil {
                                    insertions += int64(counts.Insertions)
                                    deletions += int64(counts.Removals)
                                }
                            }
                        }
                        return nil
                    },
                )

                if fieldName == "insertions" {
                    values = append(values, IntegerValue(insertions))
                } else if fieldName == "deletions" {
                    values = append(values, IntegerValue(deletions))
                } else if fieldName == "files_changed" {
                    values = append(values, IntegerValue(filesChanged))
                }
            default:
                values = append(values, NullValue())
            }
        }

        row := Row{Values: values}
        rows = append(rows, row)
    }

    return &Group{Rows: rows}, nil
}

func selectTags(
    env *Environment,
    repo *Repository,
    fieldsNames []string,
    titles []string,
    fieldsValues []Expression,
) (*Group, error) {
    var rows []Row

    platform := repo.References()
    tagNames := platform.Tags()
    repoPath := repo.Path()

    namesLen := int64(len(fieldsNames))
    valuesLen := int64(len(fieldsValues))
    padding := namesLen - valuesLen

    for _, tagRef := range tagNames {
        var values []Value

        for index := int64(0); index < namesLen; index++ {
            fieldName := fieldsNames[index]

            if index-padding >= 0 {
                value := fieldsValues[index-padding]
                if _, ok := value.(*SymbolExpression); !ok {
                    evaluated, err := evaluateExpression(env, value, titles, values)
                    if err != nil {
                        return nil, err
                    }
                    values = append(values, evaluated)
                    continue
                }
            }

            switch fieldName {
            case "name":
                tagName := tagRef.Name().CategoryAndShortName().ShortName
                values = append(values, TextValue(tagName))
            case "repo":
                values = append(values, TextValue(repoPath))
            default:
                values = append(values, NullValue())
            }
        }

        row := Row{Values: values}
        rows = append(rows, row)
    }

    return &Group{Rows: rows}, nil
}

func selectValues(
    env *Environment,
    titles []string,
    fieldsValues []Expression,
) (*Group, error) {
    var group Group
    var values []Value

    for _, value := range fieldsValues {
        evaluated, err := evaluateExpression(env, value, titles, values)
        if err != nil {
            return nil, err
        }
        values = append(values, evaluated)
    }

    row := Row{Values: values}
    group.Rows = append(group.Rows, row)
    return &group, nil
}

func GetColumnName(aliasTable map[string]string, name string) string {
    if columnName, ok := aliasTable[name]; ok {
        return columnName
    }
    return name
}

func evaluateExpression(
    env *Environment,
    expression Expression,
    titles []string,
    values []Value,
) (Value, error) {
    // implementation details
    return nil, errors.New("evaluation error")
}

func TextValue(value string) Value {
    // implementation details
    return nil
}

func NullValue() Value {
    // implementation details
    return nil
}

func IntegerValue(value int64) Value {
    // implementation details
    return nil
}

func BooleanValue(value bool) Value {
    // implementation details
    return nil
}

func DateTimeValue(value int64) Value {
    // implementation details
    return nil
}

func main() {
    // main function code
}
