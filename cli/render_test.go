package cli

import (
	"testing"

	"github.com/pterm/pterm"
	"github.com/stretchr/testify/assert"

	"github.com/ggql/ggql/ast"
)

func TestRenderObjects(t *testing.T) {
	var group []ast.GQLObject
	var groups [][]ast.GQLObject

	hiddenSelections := []string{"item"}

	err := RenderObjects(&groups, hiddenSelections, false, 1)
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(groups))

	group = append(group, ast.GQLObject{
		Attributes: map[string]ast.Value{},
	})

	groups = groups[:0]
	groups = append(groups, group)

	err = RenderObjects(&groups, hiddenSelections, false, 1)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(groups))

	groups = groups[:0]
	groups = append(groups, group, group)

	err = RenderObjects(&groups, hiddenSelections, false, 1)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(groups))
}

func TestPrintGroupAsTable(t *testing.T) {
	var group []ast.GQLObject
	var tableHeaders []string
	var titles []string

	titles = append(titles, "title1", "title2")

	for _, item := range titles {
		tableHeaders = append(tableHeaders, pterm.Green(item))
	}

	group = append(group, ast.GQLObject{
		Attributes: map[string]ast.Value{
			"title1": "value1",
			"title2": "value2",
		},
	})

	err := printGroupAsTable(titles, tableHeaders, group)
	assert.Equal(t, nil, err)
}
