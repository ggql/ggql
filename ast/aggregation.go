package ast

type Aggregation func(string, []GQLObject) Value

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
	"max":   {Parameter: Integer, Result: Integer},
	"min":   {Parameter: Integer, Result: Integer},
	"sum":   {Parameter: Any, Result: Integer},
	"avg":   {Parameter: Any, Result: Integer},
	"count": {Parameter: Any, Result: Integer},
}

func aggregationMax(fieldName string, objects []GQLObject) Value {
	maxValue := objects[0].Attributes[fieldName]

	for _, object := range objects[1:] {
		fieldValue := object.Attributes[fieldName]
		ret, _ := maxValue.Compare(fieldValue)
		if ret == Greater {
			maxValue = fieldValue
		}
	}

	return maxValue
}

func aggregationMin(fieldName string, objects []GQLObject) Value {
	minValue := objects[0].Attributes[fieldName]

	for _, object := range objects[1:] {
		fieldValue := object.Attributes[fieldName]
		ret, _ := minValue.Compare(fieldValue)
		if ret == Less {
			minValue = fieldValue
		}
	}

	return minValue
}

func aggregationSum(fieldName string, objects []GQLObject) Value {
	var sum int64

	for _, object := range objects {
		fieldValue := object.Attributes[fieldName]
		sum += fieldValue.AsInt()
	}

	return IntegerValue{sum}
}

func aggregationAverage(fieldName string, objects []GQLObject) Value {
	var sum int64
	count := int64(len(objects))

	for _, object := range objects {
		fieldValue := object.Attributes[fieldName]
		sum += fieldValue.AsInt()
	}

	avg := sum / count

	return IntegerValue{avg}
}

func aggregationCount(_ string, objects []GQLObject) Value {
	return IntegerValue{int64(len(objects))}
}
