package ast

type Ordering int

const (
	Less Ordering = iota
	Equal
	Greater
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

	return v.value == other.(IntegerValue).value
}

func (v IntegerValue) Compare(other Value) (Ordering, error) {
	// TODO: FIXME
	return Less, nil
}

func (v IntegerValue) Plus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v IntegerValue) Minus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v IntegerValue) Mul(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v IntegerValue) Div(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v IntegerValue) Modulus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v IntegerValue) Literal() string {
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

	return v.value == other.(FloatValue).value
}

func (v FloatValue) Compare(other Value) (Ordering, error) {
	// TODO: FIXME
	return Less, nil
}

func (v FloatValue) Plus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v FloatValue) Minus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v FloatValue) Mul(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v FloatValue) Div(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v FloatValue) Modulus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v FloatValue) Literal() string {
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

	return v.value == other.(TextValue).value
}

func (v TextValue) Compare(other Value) (Ordering, error) {
	// TODO: FIXME
	return Less, nil
}

func (v TextValue) Plus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v TextValue) Minus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v TextValue) Mul(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v TextValue) Div(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v TextValue) Modulus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v TextValue) Literal() string {
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

	return v.value == other.(BooleanValue).value
}

func (v BooleanValue) Compare(other Value) (Ordering, error) {
	// TODO: FIXME
	return Less, nil
}

func (v BooleanValue) Plus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v BooleanValue) Minus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v BooleanValue) Mul(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v BooleanValue) Div(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v BooleanValue) Modulus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v BooleanValue) Literal() string {
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

	return v.value == other.(DateTimeValue).value
}

func (v DateTimeValue) Compare(other Value) (Ordering, error) {
	// TODO: FIXME
	return Less, nil
}

func (v DateTimeValue) Plus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v DateTimeValue) Minus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v DateTimeValue) Mul(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v DateTimeValue) Div(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v DateTimeValue) Modulus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v DateTimeValue) Literal() string {
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

	return v.value == other.(DateValue).value
}

func (v DateValue) Compare(other Value) (Ordering, error) {
	// TODO: FIXME
	return Less, nil
}

func (v DateValue) Plus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v DateValue) Minus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v DateValue) Mul(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v DateValue) Div(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v DateValue) Modulus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v DateValue) Literal() string {
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

	return v.value == other.(TimeValue).value
}

func (v TimeValue) Compare(other Value) (Ordering, error) {
	// TODO: FIXME
	return Less, nil
}

func (v TimeValue) Plus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v TimeValue) Minus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v TimeValue) Mul(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v TimeValue) Div(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v TimeValue) Modulus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v TimeValue) Literal() string {
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
	if other.DataType() != Null {
		return false
	}

	return v.value == other.(NullValue).value
}

func (v NullValue) Compare(other Value) (Ordering, error) {
	// TODO: FIXME
	return Less, nil
}

func (v NullValue) Plus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v NullValue) Minus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v NullValue) Mul(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v NullValue) Div(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v NullValue) Modulus(other Value) (Value, error) {
	// TODO: FIXME
	return other, nil
}

func (v NullValue) Literal() string {
	return ""
}
