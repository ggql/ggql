package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlatGQLGroups(t *testing.T) {
	var group []GQLObject
	var groups [][]GQLObject

	err := FlatGQLGroups(&groups)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(groups))

	for _, item := range groups {
		assert.Equal(t, 0, len(item))
	}

	group = append(group, GQLObject{
		Attributes: map[string]Value{
			"key": "val",
		},
	})

	groups = groups[:0]
	groups = append(groups, group, group)
	assert.Equal(t, 2, len(groups))

	err = FlatGQLGroups(&groups)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(groups))

	for _, item := range groups {
		assert.Equal(t, 2, len(item))
	}
}
