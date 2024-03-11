package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregationMax(t *testing.T) {
	titles := []string{"field1", "field2"}

	objects := Group{
		Rows: []Row{
			{Values: []Value{IntegerValue{1}, IntegerValue{2}}},
			{Values: []Value{IntegerValue{3}, IntegerValue{4}}},
			{Values: []Value{IntegerValue{5}, IntegerValue{6}}},
		},
	}

	ret := aggregationMax("field1", titles, &objects)
	assert.Equal(t, int64(5), ret.AsInt())
}

func TestAggregationMin(t *testing.T) {
	titles := []string{"field1", "field2"}

	objects := Group{
		Rows: []Row{
			{Values: []Value{IntegerValue{1}, IntegerValue{2}}},
			{Values: []Value{IntegerValue{3}, IntegerValue{4}}},
			{Values: []Value{IntegerValue{5}, IntegerValue{6}}},
		},
	}

	ret := aggregationMin("field1", titles, &objects)
	assert.Equal(t, int64(1), ret.AsInt())
}

func TestAggregationSum(t *testing.T) {
	titles := []string{"field1", "field2"}

	objects := Group{
		Rows: []Row{
			{Values: []Value{IntegerValue{1}, IntegerValue{2}}},
			{Values: []Value{IntegerValue{3}, IntegerValue{4}}},
			{Values: []Value{IntegerValue{5}, IntegerValue{6}}},
		},
	}

	ret := aggregationSum("field1", titles, &objects)
	assert.Equal(t, int64(9), ret.AsInt())
}

func TestAggregationAverage(t *testing.T) {
	titles := []string{"field1", "field2"}

	objects := Group{
		Rows: []Row{
			{Values: []Value{IntegerValue{1}, IntegerValue{2}}},
			{Values: []Value{IntegerValue{3}, IntegerValue{4}}},
			{Values: []Value{IntegerValue{5}, IntegerValue{6}}},
		},
	}

	ret := aggregationAverage("field1", titles, &objects)
	assert.Equal(t, int64(3), ret.AsInt())
}

func TestAggregationCount(t *testing.T) {
	titles := []string{"field1", "field2"}

	objects := Group{
		Rows: []Row{
			{Values: []Value{IntegerValue{1}, IntegerValue{2}}},
			{Values: []Value{IntegerValue{3}, IntegerValue{4}}},
			{Values: []Value{IntegerValue{5}, IntegerValue{6}}},
		},
	}

	ret := aggregationCount("field1", titles, &objects)
	assert.Equal(t, int64(3), ret.AsInt())
}
