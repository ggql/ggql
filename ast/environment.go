package ast

import (
	"strings"

	"github.com/pkg/errors"
)

var TablesFieldsNames = map[string][]string{
	"refs":     {"name", "full_name", "type", "repo"},
	"commits":  {"commit_id", "title", "message", "name", "email", "datetime", "repo"},
	"branches": {"name", "commit_count", "is_head", "is_remote", "repo"},
	"diffs":    {"commit_id", "name", "email", "insertions", "deletions", "files_changed", "repo"},
	"tags":     {"name", "repo"},
}

type Environment struct {
	Globals      map[string]Value
	GlobalsTypes map[string]DataType
	Scopes       map[string]DataType
}

func (e *Environment) Define(str string, dataType DataType) {
	e.Scopes[str] = dataType
}

func (e *Environment) DefineGlobal(str string, dataType DataType) {
	e.GlobalsTypes[str] = dataType
}

func (e *Environment) Contains(str string) bool {
	_, globalsTypesOk := e.GlobalsTypes[str]
	_, scopesOk := e.Scopes[str]

	return scopesOk || globalsTypesOk
}

func (e *Environment) ResolveType(str string) (DataType, error) {
	if strings.HasPrefix(str, "@") {
		if val, ok := e.GlobalsTypes[str]; ok {
			return val, nil
		} else {
			return Undefined, errors.New("invalid data type")
		}
	}

	if val, ok := e.Scopes[str]; ok {
		return val, nil
	} else {
		return Undefined, errors.New("invalid data type")
	}
}

func (e *Environment) ClearSession() {
	e.Scopes = make(map[string]DataType)
}
