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

	env.Define("field1", Text{})
	assert.Equal(t, true, env.Scopes["field1"].IsText())
}

func TestDefineGlobal(t *testing.T) {
	env := Environment{
		Globals:      map[string]Value{},
		GlobalsTypes: map[string]DataType{},
		Scopes:       map[string]DataType{},
	}

	env.DefineGlobal("field1", Text{})
	assert.Equal(t, true, env.GlobalsTypes["field1"].IsText())
}

func TestContains(t *testing.T) {
	env := Environment{
		Globals:      map[string]Value{},
		GlobalsTypes: map[string]DataType{},
		Scopes:       map[string]DataType{},
	}

	env.Define("field1", Text{})
	env.DefineGlobal("field2", Integer{})

	assert.Equal(t, true, env.Contains("field1"))
	assert.Equal(t, true, env.Contains("field2"))
	assert.Equal(t, false, env.Contains("invalid"))
}

func TestResolveType(t *testing.T) {
	env := Environment{
		Globals:      map[string]Value{},
		GlobalsTypes: map[string]DataType{},
		Scopes:       map[string]DataType{},
	}

	env.Define("field1", Text{})
	env.DefineGlobal("@field2", Integer{})

	ret, err := env.ResolveType("field1")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, ret.IsText())

	ret, err = env.ResolveType("@field2")
	assert.Equal(t, nil, err)
	assert.Equal(t, true, ret.IsInt())

	ret, err = env.ResolveType("invalid")
	assert.NotEqual(t, nil, err)
}

func TestClearSession(t *testing.T) {
	env := Environment{
		Globals:      map[string]Value{},
		GlobalsTypes: map[string]DataType{},
		Scopes:       map[string]DataType{},
	}

	env.Define("field1", Text{})

	env.ClearSession()
	assert.Equal(t, 0, len(env.Scopes))
}
