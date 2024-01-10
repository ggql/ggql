package ast

import (
	"math"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type Ordering int

const (
	Less    = -1
	Equal   = 0
	Greater = 1
)

type Value interface {
	DataType() DataType
	Equals(Value) bool
	Compare(Value) (Ordering, error)
	Plus(Value) (Value, error)
	Minus(Value) (Value, error)
	Mul(Value) (Value, error)
	Div(Value) (Value, error)
	Modulus(Value) (Value, error)
	Literal() string
	AsInt() int64
	AsFloat() float64
	AsText() string
	AsBool() bool
	AsDateTime() int64
	AsDate() int64
	AsTime() string
}

// IntegerValue implementation

type IntegerValue struct {
	value int64
}

func (v IntegerValue) DataType() DataType {
	return Integer
}

func (v IntegerValue) Equals(other Value) bool {
	if other.DataType() != Integer {
		return false
	}

	return v.AsInt() == other.AsInt()
}

func (v IntegerValue) Compare(other Value) (Ordering, error) {
	helper := func(a, b int64) Ordering {
		if a < b {
			return Less
		} else if a > b {
			return Greater
		}
		return Equal
	}

	if other.DataType() != Integer {
		return Less, errors.New("invalid data type")
	}

	return helper(v.AsInt(), other.AsInt()), nil
}

func (v IntegerValue) Plus(other Value) (Value, error) {
	if other.DataType() != Integer && other.DataType() != Float {
		return nil, errors.New("invalid data type")
	}

	if other.DataType() == Float {
		return FloatValue{float64(v.AsInt()) + other.AsFloat()}, nil
	}

	return IntegerValue{v.AsInt() + other.AsInt()}, nil
}

func (v IntegerValue) Minus(other Value) (Value, error) {
	if other.DataType() != Integer && other.DataType() != Float {
		return nil, errors.New("invalid data type")
	}

	if other.DataType() == Float {
		return FloatValue{float64(v.AsInt()) - other.AsFloat()}, nil
	}

	return IntegerValue{v.AsInt() - other.AsInt()}, nil
}

func (v IntegerValue) Mul(other Value) (Value, error) {
	if other.DataType() != Integer && other.DataType() != Float {
		return nil, errors.New("invalid data type")
	}

	if other.DataType() == Float {
		return FloatValue{float64(v.AsInt()) * other.AsFloat()}, nil
	}

	lhs := v.AsInt()
	rhs := other.AsInt()

	ret := lhs * rhs
	if ret/lhs != rhs {
		return nil, errors.Errorf(
			"Attempt to compute `%d * %d`, which would overflow",
			lhs, rhs,
		)
	}

	return IntegerValue{ret}, nil
}

func (v IntegerValue) Div(other Value) (Value, error) {
	if other.DataType() != Integer && other.DataType() != Float {
		return nil, errors.New("invalid data type")
	}

	if other.DataType() == Integer {
		if other.AsInt() == 0 {
			return nil, errors.Errorf("Attempt to divide %s by zero", v.Literal())
		}
	}

	if other.DataType() == Float {
		return FloatValue{float64(v.AsInt()) / other.AsFloat()}, nil
	}

	return IntegerValue{v.AsInt() / other.AsInt()}, nil
}

func (v IntegerValue) Modulus(other Value) (Value, error) {
	if other.DataType() != Integer && other.DataType() != Float {
		return nil, errors.New("invalid data type")
	}

	if other.DataType() == Integer {
		if other.AsInt() == 0 {
			return nil, errors.Errorf("Attempt to calculate the remainder of %s with a divisor of zero", v.Literal())
		}
	}

	if other.DataType() == Float {
		return FloatValue{math.Mod(float64(v.AsInt()), other.AsFloat())}, nil
	}

	return IntegerValue{v.AsInt() % other.AsInt()}, nil
}

func (v IntegerValue) Literal() string {
	return v.AsText()
}

func (v IntegerValue) AsInt() int64 {
	return v.value
}

func (v IntegerValue) AsFloat() float64 {
	return float64(v.value)
}

func (v IntegerValue) AsText() string {
	return strconv.FormatInt(v.AsInt(), 10)
}

func (v IntegerValue) AsBool() bool {
	return false
}

func (v IntegerValue) AsDateTime() int64 {
	return 0
}

func (v IntegerValue) AsDate() int64 {
	return 0
}

func (v IntegerValue) AsTime() string {
	return ""
}

// FloatValue implementation

type FloatValue struct {
	value float64
}

func (v FloatValue) DataType() DataType {
	return Float
}

func (v FloatValue) Equals(other Value) bool {
	if other.DataType() != Float {
		return false
	}

	return v.AsFloat() == other.AsFloat()
}

func (v FloatValue) Compare(other Value) (Ordering, error) {
	helper := func(a, b float64) Ordering {
		if a < b {
			return Less
		} else if a > b {
			return Greater
		}
		return Equal
	}

	if other.DataType() != Float {
		return Less, errors.New("invalid data type")
	}

	return helper(v.AsFloat(), other.AsFloat()), nil
}

func (v FloatValue) Plus(other Value) (Value, error) {
	if other.DataType() != Float && other.DataType() != Integer {
		return nil, errors.New("invalid data type")
	}

	if other.DataType() == Integer {
		return FloatValue{v.AsFloat() + float64(other.AsInt())}, nil
	}

	return FloatValue{v.AsFloat() + other.AsFloat()}, nil
}

func (v FloatValue) Minus(other Value) (Value, error) {
	if other.DataType() != Float && other.DataType() != Integer {
		return nil, errors.New("invalid data type")
	}

	if other.DataType() == Integer {
		return FloatValue{v.AsFloat() - float64(other.AsInt())}, nil
	}

	return FloatValue{v.AsFloat() - other.AsFloat()}, nil
}

func (v FloatValue) Mul(other Value) (Value, error) {
	if other.DataType() != Float && other.DataType() != Integer {
		return nil, errors.New("invalid data type")
	}

	if other.DataType() == Integer {
		return FloatValue{v.AsFloat() * float64(other.AsInt())}, nil
	}

	return FloatValue{v.AsFloat() * other.AsFloat()}, nil
}

func (v FloatValue) Div(other Value) (Value, error) {
	if other.DataType() != Float && other.DataType() != Integer {
		return nil, errors.New("invalid data type")
	}

	if other.DataType() == Integer {
		if other.AsInt() == 0 {
			return nil, errors.Errorf("Attempt to divide %s by zero", v.Literal())
		}
		return FloatValue{v.AsFloat() / float64(other.AsInt())}, nil
	}

	return FloatValue{v.AsFloat() / other.AsFloat()}, nil
}

func (v FloatValue) Modulus(other Value) (Value, error) {
	if other.DataType() != Integer && other.DataType() != Float {
		return nil, errors.New("invalid data type")
	}

	if other.DataType() == Integer {
		if other.AsInt() == 0 {
			return nil, errors.Errorf("Attempt to calculate the remainder of %s with a divisor of zero", v.Literal())
		}
		return FloatValue{math.Mod(v.AsFloat(), float64(other.AsInt()))}, nil
	}

	return FloatValue{math.Mod(v.AsFloat(), other.AsFloat())}, nil
}

func (v FloatValue) Literal() string {
	return v.AsText()
}

func (v FloatValue) AsInt() int64 {
	return int64(v.value)
}

func (v FloatValue) AsFloat() float64 {
	return v.value
}

func (v FloatValue) AsText() string {
	return strconv.FormatFloat(v.AsFloat(), 'E', -1, 64)
}

func (v FloatValue) AsBool() bool {
	return false
}

func (v FloatValue) AsDateTime() int64 {
	return 0
}

func (v FloatValue) AsDate() int64 {
	return 0
}

func (v FloatValue) AsTime() string {
	return ""
}

// TextValue implementation

type TextValue struct {
	value string
}

func (v TextValue) DataType() DataType {
	return Text
}

func (v TextValue) Equals(other Value) bool {
	if other.DataType() != Text {
		return false
	}

	return v.AsText() == other.AsText()
}

func (v TextValue) Compare(other Value) (Ordering, error) {
	helper := func(a, b string) Ordering {
		ret := strings.Compare(other.AsText(), other.AsText())
		if ret == -1 {
			return Less
		} else if ret == 1 {
			return Greater
		}
		return Equal
	}

	if other.DataType() != Text {
		return Less, errors.New("invalid data type")
	}

	return helper(v.AsText(), other.AsText()), nil
}

func (v TextValue) Plus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v TextValue) Minus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v TextValue) Mul(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v TextValue) Div(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v TextValue) Modulus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v TextValue) Literal() string {
	return v.AsText()
}

func (v TextValue) AsInt() int64 {
	return 0
}

func (v TextValue) AsFloat() float64 {
	return 0
}

func (v TextValue) AsText() string {
	return v.value
}

func (v TextValue) AsBool() bool {
	return false
}

func (v TextValue) AsDateTime() int64 {
	return 0
}

func (v TextValue) AsDate() int64 {
	return 0
}

func (v TextValue) AsTime() string {
	return ""
}

// BooleanValue implementation

type BooleanValue struct {
	value bool
}

func (v BooleanValue) DataType() DataType {
	return Boolean
}

func (v BooleanValue) Equals(other Value) bool {
	if other.DataType() != Boolean {
		return false
	}

	return v.AsBool() == other.AsBool()
}

func (v BooleanValue) Compare(other Value) (Ordering, error) {
	if other.DataType() != Boolean {
		return Less, errors.New("invalid data type")
	}

	return Equal, nil
}

func (v BooleanValue) Plus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v BooleanValue) Minus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v BooleanValue) Mul(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v BooleanValue) Div(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v BooleanValue) Modulus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v BooleanValue) Literal() string {
	return v.AsText()
}

func (v BooleanValue) AsInt() int64 {
	return 0
}

func (v BooleanValue) AsFloat() float64 {
	return 0
}

func (v BooleanValue) AsText() string {
	return strconv.FormatBool(v.AsBool())
}

func (v BooleanValue) AsBool() bool {
	return v.value
}

func (v BooleanValue) AsDateTime() int64 {
	return 0
}

func (v BooleanValue) AsDate() int64 {
	return 0
}

func (v BooleanValue) AsTime() string {
	return ""
}

// DateTimeValue implementation

type DateTimeValue struct {
	value int64
}

func (v DateTimeValue) DataType() DataType {
	return DateTime
}

func (v DateTimeValue) Equals(other Value) bool {
	if other.DataType() != DateTime {
		return false
	}

	return v.AsDateTime() == other.AsDateTime()
}

func (v DateTimeValue) Compare(other Value) (Ordering, error) {
	helper := func(a, b int64) Ordering {
		if a < b {
			return Less
		} else if a > b {
			return Greater
		}
		return Equal
	}

	if other.DataType() != DateTime {
		return Less, errors.New("invalid data type")
	}

	return helper(v.AsDateTime(), other.AsDateTime()), nil
}

func (v DateTimeValue) Plus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v DateTimeValue) Minus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v DateTimeValue) Mul(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v DateTimeValue) Div(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v DateTimeValue) Modulus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v DateTimeValue) Literal() string {
	return v.AsText()
}

func (v DateTimeValue) AsInt() int64 {
	return 0
}

func (v DateTimeValue) AsFloat() float64 {
	return 0
}

func (v DateTimeValue) AsText() string {
	return strconv.FormatInt(v.AsDateTime(), 10)
}

func (v DateTimeValue) AsBool() bool {
	return false
}

func (v DateTimeValue) AsDateTime() int64 {
	return v.value
}

func (v DateTimeValue) AsDate() int64 {
	return 0
}

func (v DateTimeValue) AsTime() string {
	return ""
}

// DateValue implementation

type DateValue struct {
	value int64
}

func (v DateValue) DataType() DataType {
	return Date
}

func (v DateValue) Equals(other Value) bool {
	if other.DataType() != Date {
		return false
	}

	return v.AsDate() == other.AsDate()
}

func (v DateValue) Compare(other Value) (Ordering, error) {
	helper := func(a, b int64) Ordering {
		if a < b {
			return Less
		} else if a > b {
			return Greater
		}
		return Equal
	}

	if other.DataType() != Date {
		return Less, errors.New("invalid data type")
	}

	return helper(v.AsDate(), other.AsDate()), nil
}

func (v DateValue) Plus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v DateValue) Minus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v DateValue) Mul(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v DateValue) Div(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v DateValue) Modulus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v DateValue) Literal() string {
	return v.AsText()
}

func (v DateValue) AsInt() int64 {
	return 0
}

func (v DateValue) AsFloat() float64 {
	return 0
}

func (v DateValue) AsText() string {
	return strconv.FormatInt(v.AsDate(), 10)
}

func (v DateValue) AsBool() bool {
	return false
}

func (v DateValue) AsDateTime() int64 {
	return 0
}

func (v DateValue) AsDate() int64 {
	return v.value
}

func (v DateValue) AsTime() string {
	return ""
}

// TimeValue implementation

type TimeValue struct {
	value string
}

func (v TimeValue) DataType() DataType {
	return Time
}

func (v TimeValue) Equals(other Value) bool {
	if other.DataType() != Time {
		return false
	}

	return v.AsTime() == other.AsTime()
}

func (v TimeValue) Compare(other Value) (Ordering, error) {
	helper := func(a, b string) Ordering {
		ret := strings.Compare(other.AsText(), other.AsText())
		if ret == -1 {
			return Less
		} else if ret == 1 {
			return Greater
		}
		return Equal
	}

	if other.DataType() != Time {
		return Less, errors.New("invalid data type")
	}

	return helper(v.AsTime(), other.AsTime()), nil
}

func (v TimeValue) Plus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v TimeValue) Minus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v TimeValue) Mul(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v TimeValue) Div(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v TimeValue) Modulus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v TimeValue) Literal() string {
	return v.AsText()
}

func (v TimeValue) AsInt() int64 {
	return 0
}

func (v TimeValue) AsFloat() float64 {
	return 0
}

func (v TimeValue) AsText() string {
	return v.value
}

func (v TimeValue) AsBool() bool {
	return false
}

func (v TimeValue) AsDateTime() int64 {
	return 0
}

func (v TimeValue) AsDate() int64 {
	return 0
}

func (v TimeValue) AsTime() string {
	return v.value
}

// Undefined implementation

type UndefinedValue struct {
	value interface{}
}

func (v UndefinedValue) DataType() DataType {
	return Null
}

func (v UndefinedValue) Equals(other Value) bool {
	return true
}

func (v UndefinedValue) Compare(other Value) (Ordering, error) {
	if other.DataType() != Undefined {
		return Less, errors.New("invalid data type")
	}

	return Equal, nil
}

func (v UndefinedValue) Plus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v UndefinedValue) Minus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v UndefinedValue) Mul(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v UndefinedValue) Div(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v UndefinedValue) Modulus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v UndefinedValue) Literal() string {
	return v.AsText()
}

func (v UndefinedValue) AsInt() int64 {
	return 0
}

func (v UndefinedValue) AsFloat() float64 {
	return 0
}

func (v UndefinedValue) AsText() string {
	return UndefinedStr
}

func (v UndefinedValue) AsBool() bool {
	return false
}

func (v UndefinedValue) AsDateTime() int64 {
	return 0
}

func (v UndefinedValue) AsDate() int64 {
	return 0
}

func (v UndefinedValue) AsTime() string {
	return ""
}

// NullValue implementation

type NullValue struct {
	value interface{}
}

func (v NullValue) DataType() DataType {
	return Null
}

func (v NullValue) Equals(other Value) bool {
	return true
}

func (v NullValue) Compare(other Value) (Ordering, error) {
	if other.DataType() != Null {
		return Less, errors.New("invalid data type")
	}

	return Equal, nil
}

func (v NullValue) Plus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v NullValue) Minus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v NullValue) Mul(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v NullValue) Div(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v NullValue) Modulus(other Value) (Value, error) {
	return IntegerValue{0}, nil
}

func (v NullValue) Literal() string {
	return v.AsText()
}

func (v NullValue) AsInt() int64 {
	return 0
}

func (v NullValue) AsFloat() float64 {
	return 0
}

func (v NullValue) AsText() string {
	return NullStr
}

func (v NullValue) AsBool() bool {
	return false
}

func (v NullValue) AsDateTime() int64 {
	return 0
}

func (v NullValue) AsDate() int64 {
	return 0
}

func (v NullValue) AsTime() string {
	return ""
}
