package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ggql/ggql/ast"
)

func TestRenderObjects(t *testing.T) {
	object := ast.GitQLObject{
		Titles: []string{"title1", "title2"},
		Groups: []ast.Group{
			{
				Rows: []ast.Row{
					{
						Values: []ast.Value{ast.IntegerValue{Value: 1}, ast.IntegerValue{Value: 1}},
					},
				},
			},
			{
				Rows: []ast.Row{
					{
						Values: []ast.Value{ast.TextValue{Value: "hello"}, ast.TextValue{Value: "world"}},
					},
				},
			},
		},
	}

	err := RenderObjects(&object, []string{"item"}, false, 1)
	assert.Equal(t, nil, err)
}

func TestPrintGroupAsTable(t *testing.T) {
	titles := []string{"title1", "title2"}
	tableHeaders := []string{}
	rows := []ast.Row{
		{
			Values: []ast.Value{ast.TextValue{Value: "hello"}, ast.TextValue{Value: "world"}},
		},
	}

	for _, item := range titles {
		tableHeaders = append(tableHeaders, item)
	}

	err := printGroupAsTable(titles, tableHeaders, rows)
	assert.Equal(t, nil, err)
}

func TestHandlePaginationInput(t *testing.T) {
	t.Skip("Skipping TestHandlePaginationInput.")
}
