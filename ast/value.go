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
	Fmt() string
	Equals(Value) bool
	Compare(Value) Ordering
	Plus(Value) (Value, error)
	Minus(Value) (Value, error)
	Mul(Value) (Value, error)
	Div(Value) (Value, error)
	Modulus(Value) (Value, error)
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
	Value int64
}

func (v IntegerValue) DataType() DataType {
	return Integer{}
}

func (v IntegerValue) Fmt() string {
	return v.AsText()
}

func (v IntegerValue) Equals(other Value) bool {
	if !other.DataType().IsInt() {
		return false
	}

	return v.AsInt() == other.AsInt()
}

func (v IntegerValue) Compare(other Value) Ordering {
	helper := func(a, b int64) Ordering {
		if a < b {
			return Less
		} else if a > b {
			return Greater
		}
		return Equal
	}

	if !other.DataType().IsInt() {
		return Equal
	}

	return helper(other.AsInt(), v.AsInt())
}

func (v IntegerValue) Plus(other Value) (Value, error) {
	if other.DataType().IsFloat() {
		return FloatValue{float64(v.AsInt()) + other.AsFloat()}, nil
	}

	lhs := v.AsInt()
	rhs := other.AsInt()

	if lhs > 0 && rhs > 0 && lhs > math.MaxInt64-rhs {
		return nil, errors.New("integer overflow")
	}

	if lhs < 0 && rhs < 0 && lhs < math.MinInt64-rhs {
		return nil, errors.New("integer underflow")
	}

	return IntegerValue{lhs + rhs}, nil
}

func (v IntegerValue) Minus(other Value) (Value, error) {
	if other.DataType().IsFloat() {
		return FloatValue{float64(v.AsInt()) - other.AsFloat()}, nil
	}

	lhs := v.AsInt()
	rhs := other.AsInt()

	if lhs < 0 && rhs > 0 && lhs < math.MinInt64+rhs {
		return nil, errors.New("integer underflow")
	}

	if lhs > 0 && rhs < 0 && lhs > math.MaxInt64+rhs {
		return nil, errors.New("integer underflow")
	}

	return IntegerValue{lhs - rhs}, nil
}

func (v IntegerValue) Mul(other Value) (Value, error) {
	if !other.DataType().IsInt() && !other.DataType().IsFloat() {
		return nil, errors.New("invalid data type")
	}

	if other.DataType().IsFloat() {
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
	if !other.DataType().IsInt() && !other.DataType().IsFloat() {
		return nil, errors.New("invalid data type")
	}

	if other.DataType().IsInt() {
		if other.AsInt() == 0 {
			return nil, errors.Errorf("Attempt to divide %s by zero", v.Fmt())
		}
	}

	if other.DataType().IsFloat() {
		return FloatValue{float64(v.AsInt()) / other.AsFloat()}, nil
	}

	return IntegerValue{v.AsInt() / other.AsInt()}, nil
}

func (v IntegerValue) Modulus(other Value) (Value, error) {
	if !other.DataType().IsInt() && !other.DataType().IsFloat() {
		return nil, errors.New("invalid data type")
	}

	if other.DataType().IsInt() {
		if other.AsInt() == 0 {
			return nil, errors.Errorf("Attempt to calculate the remainder of %s with a divisor of zero", v.Fmt())
		}
	}

	if other.DataType().IsFloat() {
		return FloatValue{math.Mod(float64(v.AsInt()), other.AsFloat())}, nil
	}

	return IntegerValue{v.AsInt() % other.AsInt()}, nil
}

func (v IntegerValue) AsInt() int64 {
	return v.Value
}

func (v IntegerValue) AsFloat() float64 {
	return float64(v.Value)
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
	Value float64
}

func (v FloatValue) DataType() DataType {
	return Float{}
}

func (v FloatValue) Fmt() string {
	return v.AsText()
}

func (v FloatValue) Equals(other Value) bool {
	if !other.DataType().IsFloat() {
		return false
	}

	return v.AsFloat() == other.AsFloat()
}

func (v FloatValue) Compare(other Value) Ordering {
	helper := func(a, b float64) Ordering {
		if a < b {
			return Less
		} else if a > b {
			return Greater
		}
		return Equal
	}

	if !other.DataType().IsFloat() {
		return Equal
	}

	return helper(other.AsFloat(), v.AsFloat())
}

func (v FloatValue) Plus(other Value) (Value, error) {
	if other.DataType().IsInt() {
		return FloatValue{v.AsFloat() + float64(other.AsInt())}, nil
	}

	return FloatValue{v.AsFloat() + other.AsFloat()}, nil
}

func (v FloatValue) Minus(other Value) (Value, error) {
	if other.DataType().IsInt() {
		return FloatValue{v.AsFloat() - float64(other.AsInt())}, nil
	}

	return FloatValue{v.AsFloat() - other.AsFloat()}, nil
}

func (v FloatValue) Mul(other Value) (Value, error) {
	if !other.DataType().IsFloat() && !other.DataType().IsInt() {
		return nil, errors.New("invalid data type")
	}

	if other.DataType().IsInt() {
		return FloatValue{v.AsFloat() * float64(other.AsInt())}, nil
	}

	return FloatValue{v.AsFloat() * other.AsFloat()}, nil
}

func (v FloatValue) Div(other Value) (Value, error) {
	if !other.DataType().IsFloat() && !other.DataType().IsInt() {
		return nil, errors.New("invalid data type")
	}

	if other.DataType().IsInt() {
		if other.AsInt() == 0 {
			return nil, errors.Errorf("Attempt to divide %s by zero", v.Fmt())
		}
		return FloatValue{v.AsFloat() / float64(other.AsInt())}, nil
	}

	return FloatValue{v.AsFloat() / other.AsFloat()}, nil
}

func (v FloatValue) Modulus(other Value) (Value, error) {
	if !other.DataType().IsInt() && !other.DataType().IsFloat() {
		return nil, errors.New("invalid data type")
	}

	if other.DataType().IsInt() {
		if other.AsInt() == 0 {
			return nil, errors.Errorf("Attempt to calculate the remainder of %s with a divisor of zero", v.Fmt())
		}
		return FloatValue{math.Mod(v.AsFloat(), float64(other.AsInt()))}, nil
	}

	return FloatValue{math.Mod(v.AsFloat(), other.AsFloat())}, nil
}

func (v FloatValue) AsInt() int64 {
	return int64(v.Value)
}

func (v FloatValue) AsFloat() float64 {
	return v.Value
}

func (v FloatValue) AsText() string {
	return strconv.FormatFloat(v.AsFloat(), 'g', -1, 64)
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
	Value string
}

func (v TextValue) DataType() DataType {
	return Text{}
}

func (v TextValue) Fmt() string {
	return v.AsText()
}

func (v TextValue) Equals(other Value) bool {
	if !other.DataType().IsText() {
		return false
	}

	return v.AsText() == other.AsText()
}

func (v TextValue) Compare(other Value) Ordering {
	helper := func(a, b string) Ordering {
		ret := strings.Compare(a, b)
		if ret == -1 {
			return Less
		} else if ret == 1 {
			return Greater
		}
		return Equal
	}

	if !other.DataType().IsText() {
		return Equal
	}

	return helper(other.AsText(), v.AsText())
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

func (v TextValue) AsInt() int64 {
	return 0
}

func (v TextValue) AsFloat() float64 {
	return 0
}

func (v TextValue) AsText() string {
	return v.Value
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
	Value bool
}

func (v BooleanValue) DataType() DataType {
	return Boolean{}
}

func (v BooleanValue) Fmt() string {
	return v.AsText()
}

func (v BooleanValue) Equals(other Value) bool {
	if !other.DataType().IsBool() {
		return false
	}

	return v.AsBool() == other.AsBool()
}

func (v BooleanValue) Compare(other Value) Ordering {
	return Equal
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
	return v.Value
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
	Value int64
}

func (v DateTimeValue) DataType() DataType {
	return DateTime{}
}

func (v DateTimeValue) Fmt() string {
	return v.AsText()
}

func (v DateTimeValue) Equals(other Value) bool {
	if !other.DataType().IsDateTime() {
		return false
	}

	return v.AsDateTime() == other.AsDateTime()
}

func (v DateTimeValue) Compare(other Value) Ordering {
	helper := func(a, b int64) Ordering {
		if a < b {
			return Less
		} else if a > b {
			return Greater
		}
		return Equal
	}

	if !other.DataType().IsDateTime() {
		return Equal
	}

	return helper(other.AsDateTime(), v.AsDateTime())
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
	return v.Value
}

func (v DateTimeValue) AsDate() int64 {
	return 0
}

func (v DateTimeValue) AsTime() string {
	return ""
}

// DateValue implementation

type DateValue struct {
	Value int64
}

func (v DateValue) DataType() DataType {
	return Date{}
}

func (v DateValue) Fmt() string {
	return v.AsText()
}

func (v DateValue) Equals(other Value) bool {
	if !other.DataType().IsDate() {
		return false
	}

	return v.AsDate() == other.AsDate()
}

func (v DateValue) Compare(other Value) Ordering {
	helper := func(a, b int64) Ordering {
		if a < b {
			return Less
		} else if a > b {
			return Greater
		}
		return Equal
	}

	if !other.DataType().IsDate() {
		return Equal
	}

	return helper(other.AsDate(), v.AsDate())
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
	return v.Value
}

func (v DateValue) AsTime() string {
	return ""
}

// TimeValue implementation

type TimeValue struct {
	Value string
}

func (v TimeValue) DataType() DataType {
	return Time{}
}

func (v TimeValue) Fmt() string {
	return v.AsText()
}

func (v TimeValue) Equals(other Value) bool {
	if !other.DataType().IsTime() {
		return false
	}

	return v.AsTime() == other.AsTime()
}

func (v TimeValue) Compare(other Value) Ordering {
	helper := func(a, b string) Ordering {
		ret := strings.Compare(a, b)
		if ret == -1 {
			return Less
		} else if ret == 1 {
			return Greater
		}
		return Equal
	}

	if !other.DataType().IsTime() {
		return Equal
	}

	return helper(other.AsTime(), v.AsTime())
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

func (v TimeValue) AsInt() int64 {
	return 0
}

func (v TimeValue) AsFloat() float64 {
	return 0
}

func (v TimeValue) AsText() string {
	return v.Value
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
	return v.Value
}

// Undefined implementation

type UndefinedValue struct {
	Value interface{}
}

func (v UndefinedValue) DataType() DataType {
	return Undefined{}
}

func (v UndefinedValue) Fmt() string {
	return v.AsText()
}

func (v UndefinedValue) Equals(other Value) bool {
	return true
}

func (v UndefinedValue) Compare(other Value) Ordering {
	return Equal
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

func (v UndefinedValue) AsInt() int64 {
	return 0
}

func (v UndefinedValue) AsFloat() float64 {
	return 0
}

func (v UndefinedValue) AsText() string {
	return "Undefined"
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
	Value interface{}
}

func (v NullValue) DataType() DataType {
	return Null{}
}

func (v NullValue) Fmt() string {
	return v.AsText()
}

func (v NullValue) Equals(other Value) bool {
	return true
}

func (v NullValue) Compare(other Value) Ordering {
	return Equal
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

func (v NullValue) AsInt() int64 {
	return 0
}

func (v NullValue) AsFloat() float64 {
	return 0
}

func (v NullValue) AsText() string {
	return "Null"
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
