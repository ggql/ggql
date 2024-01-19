package ast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// IntegerValue implementation

func TestIntegerValueDataType(t *testing.T) {
	value := IntegerValue{1}
	ret := value.DataType()
	assert.Equal(t, true, ret.IsInt())
}

func TestIntegerValueEquals(t *testing.T) {
	value := IntegerValue{1}
	null := NullValue{}
	ret := value.Equals(null)
	assert.Equal(t, false, ret)

	value = IntegerValue{1}
	other := IntegerValue{1}
	ret = value.Equals(other)
	assert.Equal(t, true, ret)

	value = IntegerValue{1}
	other = IntegerValue{2}
	ret = value.Equals(other)
	assert.Equal(t, false, ret)
}

func TestIntegerValueCompare(t *testing.T) {
	value := IntegerValue{1}
	null := NullValue{}
	_, err := value.Compare(null)
	assert.NotEqual(t, nil, err)

	value = IntegerValue{1}
	other := IntegerValue{1}
	ret, err := value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Equal), ret)

	value = IntegerValue{1}
	other = IntegerValue{2}
	ret, err = value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Greater), ret)
}

func TestIntegerValuePlus(t *testing.T) {
	value := IntegerValue{1}
	null := NullValue{}
	_, err := value.Plus(null)
	assert.NotEqual(t, nil, err)

	value = IntegerValue{1}
	otherInt := IntegerValue{1}
	ret, err := value.Plus(otherInt)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(2), ret.AsInt())

	value = IntegerValue{1}
	otherFloat := FloatValue{1.0}
	ret, err = value.Plus(otherFloat)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(2.0), ret.AsFloat())
}

func TestIntegerValueMinus(t *testing.T) {
	value := IntegerValue{1}
	null := NullValue{}
	_, err := value.Minus(null)
	assert.NotEqual(t, nil, err)

	value = IntegerValue{1}
	otherInt := IntegerValue{1}
	ret, err := value.Minus(otherInt)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(0), ret.AsInt())

	value = IntegerValue{2}
	otherFloat := FloatValue{1.0}
	ret, err = value.Minus(otherFloat)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(1.0), ret.AsFloat())
}

func TestIntegerValueMul(t *testing.T) {
	value := IntegerValue{1}
	null := NullValue{}
	_, err := value.Mul(null)
	assert.NotEqual(t, nil, err)

	value = IntegerValue{1}
	otherInt := IntegerValue{2}
	ret, err := value.Mul(otherInt)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(2), ret.AsInt())

	value = IntegerValue{2}
	otherFloat := FloatValue{1.0}
	ret, err = value.Mul(otherFloat)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(2.0), ret.AsFloat())
}

func TestIntegerValueDiv(t *testing.T) {
	value := IntegerValue{1}
	null := NullValue{}
	_, err := value.Div(null)
	assert.NotEqual(t, nil, err)

	value = IntegerValue{1}
	otherZero := IntegerValue{0}
	_, err = value.Div(otherZero)
	assert.NotEqual(t, nil, err)

	value = IntegerValue{2}
	otherFloat := FloatValue{2.0}
	ret, err := value.Div(otherFloat)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(1.0), ret.AsFloat())

	value = IntegerValue{2}
	otherInt := IntegerValue{2}
	ret, err = value.Div(otherInt)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(1), ret.AsInt())
}

func TestIntegerValueModulus(t *testing.T) {
	value := IntegerValue{1}
	null := NullValue{}
	_, err := value.Modulus(null)
	assert.NotEqual(t, nil, err)

	value = IntegerValue{1}
	otherZero := IntegerValue{0}
	_, err = value.Modulus(otherZero)
	assert.NotEqual(t, nil, err)

	value = IntegerValue{5}
	otherFloat := FloatValue{3.0}
	ret, err := value.Modulus(otherFloat)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(2.0), ret.AsFloat())

	value = IntegerValue{5}
	otherInt := IntegerValue{3}
	ret, err = value.Modulus(otherInt)
	assert.Equal(t, nil, err)
	assert.Equal(t, int64(2), ret.AsInt())
}

func TestIntegerValueLiteral(t *testing.T) {
	value := IntegerValue{1}
	ret := value.Literal()
	fmt.Println("literal:", ret)
	assert.Equal(t, "1", ret)
}

func TestIntegerValueAsInt(t *testing.T) {
	value := IntegerValue{1}
	ret := value.AsInt()
	assert.Equal(t, int64(1), ret)
}

func TestIntegerValueAsFloat(t *testing.T) {
	value := IntegerValue{1}
	ret := value.AsFloat()
	assert.Equal(t, float64(1.0), ret)
}

func TestIntegerValueAsText(t *testing.T) {
	value := IntegerValue{1}
	ret := value.AsText()
	assert.Equal(t, "1", ret)
}

func TestIntegerValueAsBool(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestIntegerValueAsDateTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestIntegerValueAsDate(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestIntegerValueAsTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

// FloatValue implementation

func TestFloatValueDataType(t *testing.T) {
	value := FloatValue{1.0}
	ret := value.DataType()
	assert.Equal(t, true, ret.IsFloat())
}

func TestFloatValueEquals(t *testing.T) {
	value := FloatValue{1.0}
	other := FloatValue{1.0}
	ret := value.Equals(other)
	assert.Equal(t, true, ret)

	value = FloatValue{1.0}
	other = FloatValue{2.0}
	ret = value.Equals(other)
	assert.Equal(t, false, ret)
}

func TestFloatValueCompare(t *testing.T) {
	value := FloatValue{1.0}
	other := FloatValue{1.0}
	ret, err := value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Equal), ret)

	value = FloatValue{1.0}
	other = FloatValue{2.0}
	ret, err = value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Greater), ret)
}

func TestFloatValuePlus(t *testing.T) {
	value := FloatValue{1.0}
	null := NullValue{}
	_, err := value.Plus(null)
	assert.NotEqual(t, nil, err)

	value = FloatValue{1.0}
	otherInt := IntegerValue{1}
	ret, err := value.Plus(otherInt)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(2.0), ret.AsFloat())

	value = FloatValue{1.0}
	otherFloat := FloatValue{1.0}
	ret, err = value.Plus(otherFloat)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(2.0), ret.AsFloat())
}

func TestFloatValueMinus(t *testing.T) {
	value := FloatValue{1.0}
	null := NullValue{}
	_, err := value.Minus(null)
	assert.NotEqual(t, nil, err)

	value = FloatValue{1.0}
	otherInt := IntegerValue{1}
	ret, err := value.Minus(otherInt)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(0), ret.AsFloat())

	value = FloatValue{1.0}
	otherFloat := FloatValue{1.0}
	ret, err = value.Minus(otherFloat)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(0), ret.AsFloat())
}

func TestFloatValueMul(t *testing.T) {
	value := FloatValue{1.0}
	null := NullValue{}
	_, err := value.Mul(null)
	assert.NotEqual(t, nil, err)

	value = FloatValue{1.0}
	otherInt := IntegerValue{2}
	ret, err := value.Mul(otherInt)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(2.0), ret.AsFloat())

	value = FloatValue{1.0}
	otherFloat := FloatValue{2.0}
	ret, err = value.Mul(otherFloat)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(2.0), ret.AsFloat())
}

func TestFloatValueDiv(t *testing.T) {
	value := FloatValue{1}
	null := NullValue{}
	_, err := value.Div(null)
	assert.NotEqual(t, nil, err)

	value = FloatValue{2.0}
	otherZero := IntegerValue{0}
	_, err = value.Div(otherZero)
	assert.NotEqual(t, nil, err)

	value = FloatValue{2.0}
	otherFloat := IntegerValue{2}
	ret, err := value.Div(otherFloat)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(1.0), ret.AsFloat())

	value = FloatValue{2.0}
	otherInt := FloatValue{2.0}
	ret, err = value.Div(otherInt)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(1.0), ret.AsFloat())
}

func TestFloatValueModulus(t *testing.T) {
	value := FloatValue{1}
	null := NullValue{}
	_, err := value.Modulus(null)
	assert.NotEqual(t, nil, err)

	value = FloatValue{1.0}
	otherZero := IntegerValue{0}
	_, err = value.Modulus(otherZero)
	assert.NotEqual(t, nil, err)

	value = FloatValue{5.0}
	otherFloat := IntegerValue{3}
	ret, err := value.Modulus(otherFloat)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(2.0), ret.AsFloat())

	value = FloatValue{5.0}
	otherInt := FloatValue{3.0}
	ret, err = value.Modulus(otherInt)
	assert.Equal(t, nil, err)
	assert.Equal(t, float64(2.0), ret.AsFloat())
}

func TestFloatValueLiteral(t *testing.T) {
	value := FloatValue{1.0}
	ret := value.Literal()
	fmt.Println("literal:", ret)
	assert.Equal(t, "1", ret)

	value = FloatValue{1.1}
	ret = value.Literal()
	fmt.Println("literal:", ret)
	assert.Equal(t, "1.1", ret)
}

func TestFloatValueAsInt(t *testing.T) {
	value := FloatValue{1.0}
	ret := value.AsInt()
	assert.Equal(t, int64(1), ret)
}

func TestFloatValueAsFloat(t *testing.T) {
	value := FloatValue{1.0}
	ret := value.AsFloat()
	assert.Equal(t, float64(1.0), ret)
}

func TestFloatValueAsText(t *testing.T) {
	value := FloatValue{1.0}
	ret := value.AsText()
	assert.Equal(t, "1", ret)

	value = FloatValue{1.1}
	ret = value.AsText()
	assert.Equal(t, "1.1", ret)
}

func TestFloatValueAsBool(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestFloatValueAsDateTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestFloatValueAsDate(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestFloatValueAsTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

// TextValue implementation

func TestTextValueDataType(t *testing.T) {
	value := TextValue{"hello"}
	ret := value.DataType()
	assert.Equal(t, true, ret.IsText())
}

func TestTextValueEquals(t *testing.T) {
	value := TextValue{"hello"}
	other := TextValue{"hello"}
	ret := value.Equals(other)
	assert.Equal(t, true, ret)

	value = TextValue{"hello"}
	other = TextValue{"world"}
	ret = value.Equals(other)
	assert.Equal(t, false, ret)
}

func TestTextValueCompare(t *testing.T) {
	value := TextValue{"hello"}
	other := TextValue{"hello"}
	ret, err := value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Equal), ret)

	value = TextValue{"hello"}
	other = TextValue{"world"}
	ret, err = value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Greater), ret)
}

func TestTextValuePlus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTextValueMinus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTextValueMul(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTextValueDiv(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTextValueModulus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTextValueLiteral(t *testing.T) {
	value := TextValue{"hello"}
	ret := value.Literal()
	fmt.Println("literal:", ret)
	assert.Equal(t, "hello", ret)
}

func TestTextValueAsInt(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTextValueAsFloat(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTextValueAsText(t *testing.T) {
	value := TextValue{"hello"}
	ret := value.AsText()
	assert.Equal(t, "hello", ret)
}

func TestTextValueAsBool(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTextValueAsDateTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTextValueAsDate(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTextValueAsTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

// BooleanValue implementation

func TestBooleanValueDataType(t *testing.T) {
	value := BooleanValue{false}
	ret := value.DataType()
	assert.Equal(t, true, ret.IsType(Boolean))
}

func TestBooleanValueEquals(t *testing.T) {
	value := BooleanValue{true}
	other := BooleanValue{true}
	ret := value.Equals(other)
	assert.Equal(t, true, ret)

	value = BooleanValue{true}
	other = BooleanValue{false}
	ret = value.Equals(other)
	assert.Equal(t, false, ret)
}

func TestBooleanValueCompare(t *testing.T) {
	value := BooleanValue{true}
	other := BooleanValue{true}
	ret, err := value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Equal), ret)
}

func TestBooleanValuePlus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestBooleanValueMinus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestBooleanValueMul(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestBooleanValueDiv(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestBooleanValueModulus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestBooleanValueLiteral(t *testing.T) {
	value := BooleanValue{false}
	ret := value.Literal()
	fmt.Println("literal:", ret)
	assert.Equal(t, "false", ret)
}

func TestBooleanValueAsInt(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestBooleanValueAsFloat(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestBooleanValueAsText(t *testing.T) {
	value := BooleanValue{false}
	ret := value.AsText()
	assert.Equal(t, "false", ret)
}

func TestBooleanValueAsBool(t *testing.T) {
	value := BooleanValue{false}
	ret := value.AsBool()
	assert.Equal(t, false, ret)
}

func TestBooleanValueAsDateTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestBooleanValueAsDate(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestBooleanValueAsTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

// DateTimeValue implementation

func TestDateTimeValueDataType(t *testing.T) {
	value := DateTimeValue{1704890191}
	ret := value.DataType()
	assert.Equal(t, true, ret.IsDateTime())
}

func TestDateTimeValueEquals(t *testing.T) {
	value := DateTimeValue{1704890191}
	other := DateTimeValue{1704890191}
	ret := value.Equals(other)
	assert.Equal(t, true, ret)

	value = DateTimeValue{1704890191}
	other = DateTimeValue{1704890192}
	ret = value.Equals(other)
	assert.Equal(t, false, ret)
}

func TestDateTimeValueCompare(t *testing.T) {
	value := DateTimeValue{1704890191}
	other := DateTimeValue{1704890191}
	ret, err := value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Equal), ret)

	value = DateTimeValue{1704890191}
	other = DateTimeValue{1704890192}
	ret, err = value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Greater), ret)
}

func TestDateTimeValuePlus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateTimeValueMinus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateTimeValueMul(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateTimeValueDiv(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateTimeValueModulus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateTimeValueLiteral(t *testing.T) {
	value := DateTimeValue{1704890191}
	ret := value.Literal()
	fmt.Println("literal:", ret)
}

func TestDateTimeValueAsInt(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateTimeValueAsFloat(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateTimeValueAsText(t *testing.T) {
	value := DateTimeValue{1704890191}
	ret := value.AsText()
	assert.Equal(t, "1704890191", ret)
}

func TestDateTimeValueAsBool(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateTimeValueAsDateTime(t *testing.T) {
	value := DateTimeValue{1704890191}
	ret := value.AsDateTime()
	assert.Equal(t, int64(1704890191), ret)
}

func TestDateTimeValueAsDate(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateTimeValueAsTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

// DateValue implementation

func TestDateValueDataType(t *testing.T) {
	value := DateValue{1704890191}
	ret := value.DataType()
	assert.Equal(t, true, ret.IsDate())
}

func TestDateValueEquals(t *testing.T) {
	value := DateValue{1704890191}
	other := DateValue{1704890191}
	ret := value.Equals(other)
	assert.Equal(t, true, ret)

	value = DateValue{1704890191}
	other = DateValue{1704890192}
	ret = value.Equals(other)
	assert.Equal(t, false, ret)
}

func TestDateValueCompare(t *testing.T) {
	value := DateValue{1704890191}
	other := DateValue{1704890191}
	ret, err := value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Equal), ret)

	value = DateValue{1704890191}
	other = DateValue{1704890192}
	ret, err = value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Greater), ret)
}

func TestDateValuePlus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateValueMinus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateValueMul(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateValueDiv(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateValueModulus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateValueLiteral(t *testing.T) {
	value := DateValue{1704890191}
	ret := value.Literal()
	fmt.Println("literal:", ret)
}

func TestDateValueAsInt(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateValueAsFloat(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateValueAsText(t *testing.T) {
	value := DateValue{1704890191}
	ret := value.AsText()
	assert.Equal(t, "1704890191", ret)
}

func TestDateValueAsBool(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateValueAsDateTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestDateValueAsDate(t *testing.T) {
	value := DateValue{1704890191}
	ret := value.AsDate()
	assert.Equal(t, int64(1704890191), ret)
}

func TestDateValueAsTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

// TimeValue implementation

func TestTimeValueDataType(t *testing.T) {
	value := TimeValue{"12:36:31"}
	ret := value.DataType()
	assert.Equal(t, true, ret.IsTime())
}

func TestTimeValueEquals(t *testing.T) {
	value := TimeValue{"12:36:31"}
	other := TimeValue{"12:36:31"}
	ret := value.Equals(other)
	assert.Equal(t, true, ret)

	value = TimeValue{"12:36:31"}
	other = TimeValue{"12:36:32"}
	ret = value.Equals(other)
	assert.Equal(t, false, ret)
}

func TestTimeValueCompare(t *testing.T) {
	value := TimeValue{"12:36:31"}
	other := TimeValue{"12:36:31"}
	ret, err := value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Equal), ret)

	value = TimeValue{"12:36:31"}
	other = TimeValue{"12:36:32"}
	ret, err = value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Greater), ret)
}

func TestTimeValuePlus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTimeValueMinus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTimeValueMul(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTimeValueDiv(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTimeValueModulus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTimeValueLiteral(t *testing.T) {
	value := TimeValue{"12:36:31"}
	ret := value.Literal()
	fmt.Println("literal:", ret)
}

func TestTimeValueAsInt(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTimeValueAsFloat(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTimeValueAsText(t *testing.T) {
	value := TimeValue{"12:36:31"}
	ret := value.AsText()
	assert.Equal(t, "12:36:31", ret)
}

func TestTimeValueAsBool(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTimeValueAsDateTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTimeValueAsDate(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestTimeValueAsTime(t *testing.T) {
	value := TimeValue{"12:36:31"}
	ret := value.AsTime()
	assert.Equal(t, "12:36:31", ret)
}

// Undefined implementation

func TestUndefinedValueDataType(t *testing.T) {
	value := UndefinedValue{}
	ret := value.DataType()
	assert.Equal(t, true, ret.IsType(Undefined))
}

func TestUndefinedValueEquals(t *testing.T) {
	value := UndefinedValue{}
	other := UndefinedValue{}
	ret := value.Equals(other)
	assert.Equal(t, true, ret)
}

func TestUndefinedValueCompare(t *testing.T) {
	value := UndefinedValue{}
	other := UndefinedValue{}
	ret, err := value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Equal), ret)
}

func TestUndefinedValuePlus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestUndefinedValueMinus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestUndefinedValueMul(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestUndefinedValueDiv(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestUndefinedValueModulus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestUndefinedValueLiteral(t *testing.T) {
	value := UndefinedValue{}
	ret := value.Literal()
	fmt.Println("literal:", ret)
	assert.Equal(t, "Undefined", ret)
}

func TestUndefinedValueAsInt(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestUndefinedValueAsFloat(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestUndefinedValueAsText(t *testing.T) {
	value := UndefinedValue{}
	ret := value.AsText()
	assert.Equal(t, "Undefined", ret)
}

func TestUndefinedValueAsBool(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestUndefinedValueAsDateTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestUndefinedValueAsDate(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestUndefinedValueAsTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

// NullValue implementation

func TestNullValueDataType(t *testing.T) {
	value := NullValue{}
	ret := value.DataType()
	assert.Equal(t, true, ret.IsType(Null))
}

func TestNullValueEquals(t *testing.T) {
	value := NullValue{}
	other := NullValue{}
	ret := value.Equals(other)
	assert.Equal(t, true, ret)
}

func TestNullValueCompare(t *testing.T) {
	value := NullValue{}
	other := NullValue{}
	ret, err := value.Compare(other)
	assert.Equal(t, nil, err)
	assert.Equal(t, Ordering(Equal), ret)
}

func TestNullValuePlus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestNullValueMinus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestNullValueMul(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestNullValueDiv(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestNullValueModulus(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestNullValueLiteral(t *testing.T) {
	value := NullValue{}
	ret := value.Literal()
	fmt.Println("literal:", ret)
	assert.Equal(t, "Null", ret)
}

func TestNullValueAsInt(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestNullValueAsFloat(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestNullValueAsText(t *testing.T) {
	value := NullValue{}
	ret := value.AsText()
	assert.Equal(t, "Null", ret)
}

func TestNullValueAsBool(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestNullValueAsDateTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestNullValueAsDate(t *testing.T) {
	assert.Equal(t, nil, nil)
}

func TestNullValueAsTime(t *testing.T) {
	assert.Equal(t, nil, nil)
}
