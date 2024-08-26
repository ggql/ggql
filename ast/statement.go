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
	Insert
)

type Statement interface {
	Kind() StatementKind
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

type Mutate struct {
	Insert                    *GQLMutate
	GlobalVariableDeclaration *GlobalVariableStatement
}

type GQLMutate struct {
	Statements             map[string]Statement
	HasAggregationFunction bool
	HasGroupByStatement    bool
}

type InsertStatement struct {
	TableName    string
	FieldsNames  []string
	FieldsValues []Expression
	AliasTable   map[string]string
	IsDistinct   bool
}

func (s *SelectStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *SelectStatement) Kind() StatementKind {
	return Select
}

func (i *InsertStatement) AsAny() reflect.Value {
	return reflect.ValueOf(i)
}

func (i *InsertStatement) Kind() StatementKind {
	return Insert
}

type WhereStatement struct {
	Condition Expression
}

func (s *WhereStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *WhereStatement) Kind() StatementKind {
	return Where
}

type HavingStatement struct {
	Condition Expression
}

func (s *HavingStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *HavingStatement) Kind() StatementKind {
	return Having
}

type LimitStatement struct {
	Count int
}

func (s *LimitStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *LimitStatement) Kind() StatementKind {
	return Limit
}

type OffsetStatement struct {
	Count int
}

func (s *OffsetStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *OffsetStatement) Kind() StatementKind {
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

func (s *OrderByStatement) Kind() StatementKind {
	return OrderBy
}

type GroupByStatement struct {
	FieldName string
}

func (s *GroupByStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *GroupByStatement) Kind() StatementKind {
	return GroupBy
}

type AggregateValue struct {
	Expression Expression
	Function   struct {
		Name string
		Arg  string
	}
}

type AggregationsStatement struct {
	Aggregations map[string]AggregateValue
}

func (s *AggregationsStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *AggregationsStatement) Kind() StatementKind {
	return AggregateFunction
}

type GlobalVariableStatement struct {
	Name  string
	Value Expression
}

func (s *GlobalVariableStatement) AsAny() reflect.Value {
	return reflect.ValueOf(s)
}

func (s *GlobalVariableStatement) Kind() StatementKind {
	return GlobalVariable
}
