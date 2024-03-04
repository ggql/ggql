package parser

import (
	"fmt"

	"github.com/ggql/ggql/ast"
)

type ParserContext struct {
	Aggregations        map[string]ast.AggregateValue
	SelectedFields      []string
	HiddenSelections    []string
	GeneratedFieldCount int32
	IsSingleValueQuery  bool
	HasGroupByStatement bool
}

func NewParserContext() *ParserContext {
	return &ParserContext{
		Aggregations:        make(map[string]ast.AggregateValue),
		SelectedFields:      make([]string, 0),
		HiddenSelections:    make([]string, 0),
		GeneratedFieldCount: 0,
		IsSingleValueQuery:  false,
		HasGroupByStatement: false,
	}
}

func (p *ParserContext) GenerateColumnName() string {
	p.GeneratedFieldCount++
	return fmt.Sprintf("column_%d", p.GeneratedFieldCount)
}
