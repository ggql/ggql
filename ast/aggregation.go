package ast

import (
	"strings"
)

type Aggregation func(string, []string, *Group) Value

type AggregationPrototype struct {
	Parameter DataType
	Result    DataType
}

var Aggregations = map[string]Aggregation{
	"max":   aggregationMax,
	"min":   aggregationMin,
	"sum":   aggregationSum,
	"avg":   aggregationAverage,
	"count": aggregationCount,
}

var AggregationsProtos = map[string]AggregationPrototype{
	"max": {
		Parameter: Variant{
			Integer{},
			Float{},
			Text{},
			Date{},
			Time{},
			DateTime{},
		},
		Result: Integer{},
	},
	"min": {
		Parameter: Variant{
			Integer{},
			Float{},
			Text{},
			Date{},
			Time{},
			DateTime{},
		},
		Result: Integer{},
	},
	"sum":   {Parameter: Integer{}, Result: Integer{}},
	"avg":   {Parameter: Integer{}, Result: Integer{}},
	"count": {Parameter: Any{}, Result: Integer{}},
}

func aggregationMax(fieldName string, titles []string, objects *Group) Value {
	var columnIndex int

	for i, v := range titles {
		if strings.Compare(v, fieldName) == 0 {
			columnIndex = i
			break
		}
	}

	maxValue := objects.Rows[0].Values[columnIndex]

	for _, row := range objects.Rows {
		fieldValue := row.Values[columnIndex]
		if ret := maxValue.Compare(fieldValue); ret == Greater {
			maxValue = fieldValue
		}
	}

	return maxValue
}

func aggregationMin(fieldName string, titles []string, objects *Group) Value {
	var columnIndex int

	for i, v := range titles {
		if strings.Compare(v, fieldName) == 0 {
			columnIndex = i
			break
		}
	}

	minValue := objects.Rows[0].Values[columnIndex]

	for _, row := range objects.Rows {
		fieldValue := row.Values[columnIndex]
		if ret := minValue.Compare(fieldValue); ret == Less {
			minValue = fieldValue
		}
	}

	return minValue
}

func aggregationSum(fieldName string, titles []string, objects *Group) Value {
	var sum int64
	var columnIndex int

	for i, v := range titles {
		if strings.Compare(v, fieldName) == 0 {
			columnIndex = i
			break
		}
	}

	for _, row := range objects.Rows {
		fieldValue := row.Values[columnIndex]
		sum += fieldValue.AsInt()
	}

	return IntegerValue{Value: sum}
}

func aggregationAverage(fieldName string, titles []string, objects *Group) Value {
	var sum int64
	var columnIndex int

	for i, v := range titles {
		if strings.Compare(v, fieldName) == 0 {
			columnIndex = i
			break
		}
	}

	count := int64(len(objects.Rows))

	for _, row := range objects.Rows {
		fieldValue := row.Values[columnIndex]
		sum += fieldValue.AsInt()
	}

	avg := sum / count

	return IntegerValue{Value: avg}
}

func aggregationCount(_ string, _ []string, objects *Group) Value {
	return IntegerValue{int64(objects.Len())}
}
