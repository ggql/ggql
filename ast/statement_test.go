package ast

import (
	"reflect"
	"testing"
)

func TestSelectWhereHavingLimitStatement(t *testing.T) {
	selectStatement := &SelectStatement{
		TableName:    "users",
		FieldsNames:  []string{"id", "name"},
		FieldsValues: []Expression{},
		AliasTable:   map[string]string{"id": "user_id"},
		IsDistinct:   true,
	}

	if kind1 := selectStatement.GetStatementKind(); kind1 != Select {
		t.Errorf("Expected StatementKind %v, got %v", Select, kind1)
	}

	anyValue1 := selectStatement.AsAny()
	if anyValue1.Kind() != reflect.Ptr || anyValue1.Type() != reflect.TypeOf(&SelectStatement{}) {
		t.Error("AsAny did not return a pointer to SelectStatement")
	}

	if !selectStatement.IsDistinct {
		t.Error("Expected IsDistinct to be true, got false")
	}

	whereStatement := &WhereStatement{
		Condition: &StringExpression{},
	}

	if kind2 := whereStatement.GetStatementKind(); kind2 != Where {
		t.Errorf("Expected StatementKind %v, got %v", Select, kind2)
	}

	anyValue2 := whereStatement.AsAny()
	if anyValue2.Kind() != reflect.Ptr || anyValue2.Type() != reflect.TypeOf(&WhereStatement{}) {
		t.Error("AsAny did not return a pointer to WhereStatement")
	}

	havingStatement := &HavingStatement{
		Condition: &StringExpression{},
	}

	if kind3 := havingStatement.GetStatementKind(); kind3 != Having {
		t.Errorf("Expected StatementKind %v, got %v", Having, kind3)
	}

	anyValue3 := havingStatement.AsAny()
	if anyValue3.Kind() != reflect.Ptr || anyValue3.Type() != reflect.TypeOf(&HavingStatement{}) {
		t.Error("AsAny did not return a pointer to HavingStatement")
	}

	limitStatement := &LimitStatement{
		Count: 1,
	}

	if kind4 := limitStatement.GetStatementKind(); kind4 != Limit {
		t.Errorf("Expected StatementKind %v, got %v", Limit, kind4)
	}

	anyValue4 := limitStatement.AsAny()
	if anyValue4.Kind() != reflect.Ptr || anyValue4.Type() != reflect.TypeOf(&LimitStatement{}) {
		t.Error("AsAny did not return a pointer to LimitStatement")
	}
}

func TestOffsetOrderByGroupByStatement(t *testing.T) {
	offsetStatement := &OffsetStatement{
		Count: 1,
	}

	if kind5 := offsetStatement.GetStatementKind(); kind5 != Offset {
		t.Errorf("Expected StatementKind %v, got %v", Offset, kind5)
	}

	anyValue5 := offsetStatement.AsAny()
	if anyValue5.Kind() != reflect.Ptr || anyValue5.Type() != reflect.TypeOf(&OffsetStatement{}) {
		t.Error("AsAny did not return a pointer to OffsetStatement")
	}

	orderByStatement := &OrderByStatement{
		Arguments:     []Expression{},
		SortingOrders: []SortingOrder{1, 2},
	}

	if kind6 := orderByStatement.GetStatementKind(); kind6 != OrderBy {
		t.Errorf("Expected StatementKind %v, got %v", OrderBy, kind6)
	}

	anyValue6 := orderByStatement.AsAny()
	if anyValue6.Kind() != reflect.Ptr || anyValue6.Type() != reflect.TypeOf(&OrderByStatement{}) {
		t.Error("AsAny did not return a pointer to OrderByStatement")
	}

	groupByStatement := &GroupByStatement{
		FieldName: "field",
	}

	if kind7 := groupByStatement.GetStatementKind(); kind7 != GroupBy {
		t.Errorf("Expected StatementKind %v, got %v", GroupBy, kind7)
	}

	anyValue7 := groupByStatement.AsAny()
	if anyValue7.Kind() != reflect.Ptr || anyValue7.Type() != reflect.TypeOf(&GroupByStatement{}) {
		t.Error("AsAny did not return a pointer to GroupByStatement")
	}
}

func TestAggregationGlobalVariableStatement(t *testing.T) {
	aggregationFunctionsStatement := &AggregationFunctionsStatement{
		Aggregations: map[string]*AggregateFunctions{
			"a": {
				FunctionName: "functionname",
				Argument:     "argument",
			},
		},
	}

	if kind8 := aggregationFunctionsStatement.GetStatementKind(); kind8 != AggregateFunction {
		t.Errorf("Expected StatementKind %v, got %v", AggregateFunction, kind8)
	}

	anyValue8 := aggregationFunctionsStatement.AsAny()
	if anyValue8.Kind() != reflect.Ptr || anyValue8.Type() != reflect.TypeOf(&AggregationFunctionsStatement{}) {
		t.Error("AsAny did not return a pointer to GroupByStatement")
	}

	globalVariableStatement := &GlobalVariableStatement{
		Name:  "name",
		Value: &StringExpression{},
	}

	if kind9 := globalVariableStatement.GetStatementKind(); kind9 != GlobalVariable {
		t.Errorf("Expected StatementKind %v, got %v", GlobalVariable, kind9)
	}

	anyValue9 := globalVariableStatement.AsAny()
	if anyValue9.Kind() != reflect.Ptr || anyValue9.Type() != reflect.TypeOf(&GlobalVariableStatement{}) {
		t.Error("AsAny did not return a pointer to GlobalVariableStatement")
	}
}
