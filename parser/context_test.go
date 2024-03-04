package parser

import (
	"testing"

	"github.com/ggql/ggql/ast"
)

func TestGenerateColumnName(t *testing.T) {
	aggregationFunc := make(map[string]ast.AggregateValue)
	aggregationFunc["aggre"] = ast.AggregateValue{}
	parserContext := &ParserContext{
		Aggregations:        aggregationFunc,
		GeneratedFieldCount: 0,
	}

	columnName := parserContext.GenerateColumnName()

	if columnName != "column_1" {
		t.Error("Generated column name should be 'column_1'")
	}
	columnName2 := parserContext.GenerateColumnName()
	if columnName2 != "column_2" {
		t.Error("Generated column name should be 'column_2'")
	}
}
