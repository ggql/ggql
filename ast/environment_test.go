package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefine(t *testing.T) {
	env := Environment{
		Globals:      map[string]Value{},
		GlobalsTypes: map[string]DataType{},
		Scopes:       map[string]DataType{},
	}

	env.Define("key1", Text)
	assert.Equal(t, true, env.Scopes["key1"].IsText())
}

func TestDefineGlobal(t *testing.T) {
	env := Environment{
		Globals:      map[string]Value{},
		GlobalsTypes: map[string]DataType{},
		Scopes:       map[string]DataType{},
	}

	env.DefineGlobal("key1", Text)
	assert.Equal(t, true, env.GlobalsTypes["key1"].IsText())
}

func TestContains(t *testing.T) {
	env := Environment{
		Globals:      map[string]Value{},
		GlobalsTypes: map[string]DataType{},
		Scopes:       map[string]DataType{},
	}

	env.Define("key1", Text)
	assert.Equal(t, true, env.Contains("key1"))

	env.DefineGlobal("key2", Text)
	assert.Equal(t, true, env.Contains("key2"))
}

func TestResolveType(t *testing.T) {
	env := Environment{
		Globals:      map[string]Value{},
		GlobalsTypes: map[string]DataType{},
		Scopes:       map[string]DataType{},
	}

	env.Define("key1", Text)
	ret, err := env.ResolveType("key1")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, ret.IsText())

	env.DefineGlobal("@key2", Text)
	ret, err = env.ResolveType("@key2")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, ret.IsText())
}

func TestClearSession(t *testing.T) {
	env := Environment{
		Globals:      map[string]Value{},
		GlobalsTypes: map[string]DataType{},
		Scopes:       map[string]DataType{},
	}

	env.Define("key1", Text)
	env.ClearSession()
	assert.Equal(t, 0, len(env.Scopes))
}
