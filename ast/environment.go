package ast

import (
	"strings"
	"sync"
)

var TABLES_FIELDS_NAMES = map[string][]string{
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
	mutex        sync.RWMutex
}

func NewEnvironment() *Environment {
	return &Environment{
		Globals:      make(map[string]Value),
		GlobalsTypes: make(map[string]DataType),
		Scopes:       make(map[string]DataType),
		mutex:        sync.RWMutex{},
	}
}

func (env *Environment) Define(str string, dataType DataType) {
	env.mutex.Lock()
	defer env.mutex.Unlock()

	env.Scopes[str] = dataType
}

func (env *Environment) DefineGlobal(str string, dataType DataType) {
	env.mutex.Lock()
	defer env.mutex.Unlock()

	env.GlobalsTypes[str] = dataType
}

func (env *Environment) Contains(str string) bool {
	env.mutex.Lock()
	defer env.mutex.Unlock()

	_, bGlobalsTypesExist := env.GlobalsTypes[str]
	_, bScopesExist := env.Scopes[str]

	return (bScopesExist || bGlobalsTypesExist)
}

func (env *Environment) ResolveType(str string) interface{} {
	env.mutex.Lock()
	defer env.mutex.Unlock()

	if strings.HasPrefix(str, "@") {
		if value, bGlobalsTypesExist := env.GlobalsTypes[str]; bGlobalsTypesExist {
			return value
		}
	} else {
		if value, bScopesExist := env.Scopes[str]; bScopesExist {
			return value
		}
	}

	return nil
}

func (env *Environment) ClearSession() {
	env.mutex.Lock()
	defer env.mutex.Unlock()

	env.Scopes = make(map[string]DataType)
}
