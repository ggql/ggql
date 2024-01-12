package ast

import (
	"strings"
	"sync"
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

func (e *Environment) Define(str string, dataType DataType) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.Scopes[str] = dataType
}

func (e *Environment) DefineGlobal(str string, dataType DataType) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.GlobalsTypes[str] = dataType
}

func (e *Environment) Contains(str string) bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	_, bGlobalsTypesExist := e.GlobalsTypes[str]
	_, bScopesExist := e.Scopes[str]

	return bScopesExist || bGlobalsTypesExist
}

func (e *Environment) ResolveType(str string) interface{} {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if strings.HasPrefix(str, "@") {
		if value, bGlobalsTypesExist := e.GlobalsTypes[str]; bGlobalsTypesExist {
			return value
		}
	} else {
		if value, bScopesExist := e.Scopes[str]; bScopesExist {
			return value
		}
	}

	return nil
}

func (e *Environment) ClearSession() {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.Scopes = make(map[string]DataType)
}
