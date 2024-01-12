package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAggregationMax(t *testing.T) {
	objects := []GQLObject{
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{1},
			},
		},
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{2},
			},
		},
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{3},
			},
		},
	}

	ret := aggregationMax("field1", objects)
	assert.Equal(t, int64(3), ret.AsInt())
}

func TestAggregationMin(t *testing.T) {
	objects := []GQLObject{
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{1},
			},
		},
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{2},
			},
		},
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{3},
			},
		},
	}

	ret := aggregationMin("field1", objects)
	assert.Equal(t, int64(1), ret.AsInt())
}

func TestAggregationSum(t *testing.T) {
	objects := []GQLObject{
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{1},
			},
		},
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{2},
			},
		},
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{3},
			},
		},
	}

	ret := aggregationSum("field1", objects)
	assert.Equal(t, int64(6), ret.AsInt())
}

func TestAggregationAverage(t *testing.T) {
	objects := []GQLObject{
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{1},
			},
		},
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{2},
			},
		},
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{3},
			},
		},
	}

	ret := aggregationAverage("field1", objects)
	assert.Equal(t, int64(2), ret.AsInt())
}

func TestAggregationCount(t *testing.T) {
	objects := []GQLObject{
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{1},
			},
		},
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{2},
			},
		},
		{
			Attributes: map[string]Value{
				"field1": IntegerValue{3},
			},
		},
	}

	ret := aggregationCount("field1", objects)
	assert.Equal(t, int64(3), ret.AsInt())
}
