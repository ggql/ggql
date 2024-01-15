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
	assert.Equal(t, true, env.Scopes["key1"].isText())
}

func TestDefineGlobal(t *testing.T) {
	env := Environment{
		Globals:      map[string]Value{},
		GlobalsTypes: map[string]DataType{},
		Scopes:       map[string]DataType{},
	}

	env.DefineGlobal("key1", Text)
	assert.Equal(t, true, env.GlobalsTypes["key1"].isText())
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
	assert.Equal(t, true, ret.isText())

	env.DefineGlobal("@key2", Text)
	ret, err = env.ResolveType("@key2")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, ret.isText())
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
