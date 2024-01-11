package ast

import (
	"reflect"
)

type StatementKind int

const (
	Select StatementKind = iota
	Where
	Having
	Limit
	Offset
	OrderBy
	GroupBy
	AggregateFunction
	GlobalVariable
)

type Statement interface {
	GetStatementKind() StatementKind
	AsAny() reflect.Value
}

type Query struct {
	Select                    *GQLQuery
	GlobalVariableDeclaration *GlobalVariableStatement
}

type GQLQuery struct {
	Statements             map[string]Statement
	HasAggregationFunction bool
	HasGroupByStatement    bool
	HiddenSelections       []string
}

type SelectStatement struct {
	TableName    string
	FieldsNames  []string
	FieldsValues []Expression
	AliasTable   map[string]string
	IsDistinct   bool
}

func (s *SelectStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *SelectStatement) GetStatementKind() StatementKind {
	return Select
}

type WhereStatement struct {
	Condition Expression
}

func (s *WhereStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *WhereStatement) GetStatementKind() StatementKind {
	return Where
}

type HavingStatement struct {
	Condition Expression
}

func (s *HavingStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *HavingStatement) GetStatementKind() StatementKind {
	return Having
}

type LimitStatement struct {
	Count int
}

func (s *LimitStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *LimitStatement) GetStatementKind() StatementKind {
	return Limit
}

type OffsetStatement struct {
	Count int
}

func (s *OffsetStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *OffsetStatement) GetStatementKind() StatementKind {
	return Offset
}

type SortingOrder int

const (
	Ascending SortingOrder = iota
	Descending
)

type OrderByStatement struct {
	Arguments     []Expression
	SortingOrders []SortingOrder
}

func (s *OrderByStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *OrderByStatement) GetStatementKind() StatementKind {
	return OrderBy
}

type GroupByStatement struct {
	FieldName string
}

func (s *GroupByStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *GroupByStatement) GetStatementKind() StatementKind {
	return GroupBy
}

type AggregateFunctions struct {
	FunctionName string
	Argument     string
}

type AggregationFunctionsStatement struct {
	Aggregations map[string]*AggregateFunctions
}

func (s *AggregationFunctionsStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *AggregationFunctionsStatement) GetStatementKind() StatementKind {
	return AggregateFunction
}

type GlobalVariableStatement struct {
	Name  string
	Value Expression
}

func (s *GlobalVariableStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *GlobalVariableStatement) GetStatementKind() StatementKind {
	return GlobalVariable
}
