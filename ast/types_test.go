package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestPartialeqEq(t *testing.T) {
	partialAny := Any{}
	otherAny := Any{}

	ret := partialAny.Equal(otherAny)
	assert.Equal(t, true, ret)

	partialVariant := Variant{Text{}, Integer{}}
	otherText := Text{}

	ret = partialVariant.Equal(otherText)
	assert.Equal(t, true, ret)

	partialText := Text{}
	otherVariant := Variant{Text{}, Integer{}}

	ret = partialText.Equal(otherVariant)
	assert.Equal(t, true, ret)

	partialOptional := Optional{Text{}}
	otherText = Text{}

	ret = partialOptional.Equal(otherText)
	assert.Equal(t, true, ret)

	partialText = Text{}
	otherOptional := Optional{Text{}}

	ret = partialText.Equal(otherOptional)
	assert.Equal(t, true, ret)

	partialVarargs := Varargs{Text{}}
	otherText = Text{}

	ret = partialVarargs.Equal(otherText)
	assert.Equal(t, true, ret)

	partialText = Text{}
	otherVarargs := Varargs{Text{}}

	ret = partialText.Equal(otherVarargs)
	assert.Equal(t, true, ret)

	partialBoolean := Boolean{}
	otherBoolean := Boolean{}

	ret = partialBoolean.Equal(otherBoolean)
	assert.Equal(t, true, ret)

	partialInteger := Integer{}
	otherInteger := Integer{}

	ret = partialInteger.Equal(otherInteger)
	assert.Equal(t, true, ret)

	partialFloat := Float{}
	otherFloat := Float{}

	ret = partialFloat.Equal(otherFloat)
	assert.Equal(t, true, ret)

	partialInteger = Integer{}
	otherInteger = Integer{}

	ret = partialInteger.Equal(otherInteger)
	assert.Equal(t, true, ret)

	partialText = Text{}
	otherText = Text{}

	ret = partialText.Equal(otherText)
	assert.Equal(t, true, ret)

	partialDate := Date{}
	otherDate := Date{}

	ret = partialDate.Equal(otherDate)
	assert.Equal(t, true, ret)

	partialTime := Time{}
	otherTime := Time{}

	ret = partialTime.Equal(otherTime)
	assert.Equal(t, true, ret)

	partialDateTime := DateTime{}
	otherDateTime := DateTime{}

	ret = partialDateTime.Equal(otherDateTime)
	assert.Equal(t, true, ret)

	partialNull := Null{}
	otherNull := Null{}

	ret = partialNull.Equal(otherNull)
	assert.Equal(t, true, ret)

	partialUndefined := Undefined{}
	otherUndefined := Undefined{}

	ret = partialUndefined.Equal(otherUndefined)
	assert.Equal(t, true, ret)
}

func TestDatatypeFmt(t *testing.T) {
	datatypeAny := Any{}
	assert.Equal(t, "Any", datatypeAny.Fmt())

	datatypeText := Text{}
	assert.Equal(t, "Text", datatypeText.Fmt())

	datatypeInteger := Integer{}
	assert.Equal(t, "Integer", datatypeInteger.Fmt())

	datatypeFloat := Float{}
	assert.Equal(t, "Float", datatypeFloat.Fmt())

	datatypeBoolean := Boolean{}
	assert.Equal(t, "Boolean", datatypeBoolean.Fmt())

	datatypeDate := Date{}
	assert.Equal(t, "Date", datatypeDate.Fmt())

	datatypeTime := Time{}
	assert.Equal(t, "Time", datatypeTime.Fmt())

	datatypeDateTime := DateTime{}
	assert.Equal(t, "DateTime", datatypeDateTime.Fmt())

	datatypeUndefined := Undefined{}
	assert.Equal(t, typesUndefined, datatypeUndefined.Fmt())

	datatypeNull := Null{}
	assert.Equal(t, typesNull, datatypeNull.Fmt())

	datatypeVariant := Variant{Text{}, Integer{}}
	assert.Equal(t, "[Text | Integer]", datatypeVariant.Fmt())

	datatypeOptional := Optional{Text{}}
	assert.Equal(t, "Text?", datatypeOptional.Fmt())

	datatypeVarargs := Varargs{Text{}}
	assert.Equal(t, "...Text", datatypeVarargs.Fmt())
}

func TestDatatypeIsAny(t *testing.T) {
	datatype := Any{}

	ret := datatype.IsAny()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsBool(t *testing.T) {
	datatype := Boolean{}

	ret := datatype.IsBool()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsInt(t *testing.T) {
	datatype := Integer{}

	ret := datatype.IsInt()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsFloat(t *testing.T) {
	datatype := Float{}

	ret := datatype.IsFloat()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsNumber(t *testing.T) {
	datatypeInteger := Integer{}

	ret := datatypeInteger.IsNumber()
	assert.Equal(t, true, ret)

	datatypeFloat := Float{}

	ret = datatypeFloat.IsNumber()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsText(t *testing.T) {
	datatype := Text{}

	ret := datatype.IsText()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsTime(t *testing.T) {
	datatype := Time{}

	ret := datatype.IsTime()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsDate(t *testing.T) {
	datatype := Date{}

	ret := datatype.IsDate()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsDatetime(t *testing.T) {
	datatype := DateTime{}

	ret := datatype.IsDateTime()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsNull(t *testing.T) {
	datatype := Null{}

	ret := datatype.IsNull()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsUndefined(t *testing.T) {
	datatype := Undefined{}

	ret := datatype.IsUndefined()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsVariant(t *testing.T) {
	datatype := Variant{Text{}, Integer{}}

	ret := datatype.IsVariant()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsOptional(t *testing.T) {
	datatype := Optional{Text{}}

	ret := datatype.IsOptional()
	assert.Equal(t, true, ret)
}

func TestDatatypeIsVarargs(t *testing.T) {
	datatype := Varargs{Text{}}

	ret := datatype.IsVarargs()
	assert.Equal(t, true, ret)
}
