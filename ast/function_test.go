package ast

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// String functions

func TestTextLowercase(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"HELLO"})
	ret := textLowercase(buf)
	assert.Equal(t, "hello", ret.AsText())
}

func TestTextUppercase(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"hello"})
	ret := textUppercase(buf)
	assert.Equal(t, "HELLO", ret.AsText())
}

func TestTextReverse(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"hello"})
	ret := textReverse(buf)
	assert.Equal(t, "olleh", ret.AsText())
}

func TestTextReplicate(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"a"}, IntegerValue{3})
	ret := textReplicate(buf)
	assert.Equal(t, "aaa", ret.AsText())
}

func TestTextSpace(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{3})
	ret := textSpace(buf)
	assert.Equal(t, "   ", ret.AsText())
}

func TestTextTrim(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{" hello "})
	ret := textTrim(buf)
	assert.Equal(t, "hello", ret.AsText())
}

func TestTextLeftTrim(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{" hello"})
	ret := textLeftTrim(buf)
	assert.Equal(t, "hello", ret.AsText())
}

func TestTextRightTrim(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"hello "})
	ret := textRightTrim(buf)
	assert.Equal(t, "hello", ret.AsText())
}

func TestTextLen(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"hello"})
	ret := textLen(buf)
	assert.Equal(t, int64(len("hello")), ret.AsInt())
}

func TestTextAscii(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{""})
	ret := textAscii(buf)
	assert.Equal(t, int64(0), ret.AsInt())

	buf = nil
	buf = append(buf, TextValue{"a"})
	ret = textAscii(buf)
	assert.Equal(t, int64(97), ret.AsInt())
}

func TestTextLeft(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{""})
	ret := textLeft(buf)
	assert.Equal(t, "", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{"hello"}, IntegerValue{6})
	ret = textLeft(buf)
	assert.Equal(t, "hello", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{"hello"}, IntegerValue{2})
	ret = textLeft(buf)
	assert.Equal(t, "he", ret.AsText())
}

func TestTextDataLength(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{""})
	ret := textDataLength(buf)
	assert.Equal(t, int64(0), ret.AsInt())

	buf = nil
	buf = append(buf, TextValue{"hello"})
	ret = textDataLength(buf)
	assert.Equal(t, int64(5), ret.AsInt())
}

func TestTextChar(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{-1})
	ret := textChar(buf)
	assert.Equal(t, "", ret.AsText())

	buf = nil
	buf = append(buf, IntegerValue{97})
	ret = textChar(buf)
	assert.Equal(t, "a", ret.AsText())
}

func TestTextReplace(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"hello"}, TextValue{"he"}, TextValue{"eh"})
	ret := textReplace(buf)
	assert.Equal(t, "ehllo", ret.AsText())
}

func TestTextSubstring(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"hello"}, IntegerValue{7}, IntegerValue{2})
	ret := textSubstring(buf)
	assert.Equal(t, "hello", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{"hello"}, IntegerValue{1}, IntegerValue{-1})
	ret = textSubstring(buf)
	assert.Equal(t, "", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{"hello"}, IntegerValue{1}, IntegerValue{2})
	ret = textSubstring(buf)
	assert.Equal(t, "he", ret.AsText())
}

func TestTextStuff(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{""}, IntegerValue{1}, IntegerValue{2}, TextValue{"world"})
	ret := textStuff(buf)
	assert.Equal(t, "", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{"hello"}, IntegerValue{7}, IntegerValue{2}, TextValue{"world"})
	ret = textStuff(buf)
	assert.Equal(t, "hello", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{"hello"}, IntegerValue{1}, IntegerValue{2}, TextValue{"aa"})
	ret = textStuff(buf)
	assert.Equal(t, "aallo", ret.AsText())
}

func TestTextRight(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{""}, IntegerValue{1})
	ret := textRight(buf)
	assert.Equal(t, "", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{"hello"}, IntegerValue{6})
	ret = textRight(buf)
	assert.Equal(t, "hello", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{"hello"}, IntegerValue{2})
	ret = textRight(buf)
	assert.Equal(t, "lo", ret.AsText())
}

func TestTextTranslate(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"hello"}, TextValue{"he"}, TextValue{"aaa"})
	ret := textTranslate(buf)
	assert.Equal(t, "", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{"hello"}, TextValue{"he"}, TextValue{"aa"})
	ret = textTranslate(buf)
	assert.Equal(t, "aallo", ret.AsText())
}

func TestTextUnicode(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"aa"})
	ret := textUnicode(buf)
	assert.Equal(t, int64(97), ret.AsInt())
}

func TestTextSoundex(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{""})
	ret := textSoundex(buf)
	assert.Equal(t, "", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{">>>>"})
	ret = textSoundex(buf)
	assert.Equal(t, ">000", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{">>>"})
	ret = textSoundex(buf)
	assert.Equal(t, ">000", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{">>>>>"})
	ret = textSoundex(buf)
	assert.Equal(t, ">000", ret.AsText())

	buf = nil
	buf = append(buf, TextValue{"BFPVC"})
	ret = textSoundex(buf)
	assert.Equal(t, "B111", ret.AsText())
}

func TestTextConcat(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"hello"}, TextValue{"world"})
	ret := textConcat(buf)
	assert.Equal(t, "helloworld", ret.AsText())
}

// Date functions

func TestDateCurrentDate(t *testing.T) {
	var buf []Value

	ret := dateCurrentDate(buf)
	fmt.Printf("date_current_date: %d", ret.AsDate())
	assert.NotEqual(t, 0, ret.AsDate())
}

func TestDateCurrentTime(t *testing.T) {
	var buf []Value

	ret := dateCurrentTime(buf)
	fmt.Printf("date_current_time: %s", ret.AsTime())
	assert.NotEqual(t, "", ret.AsTime())
}

func TestDateCurrentTimestamp(t *testing.T) {
	var buf []Value

	ret := dateCurrentTimestamp(buf)
	fmt.Printf("date_current_timestamp: %d", ret.AsDateTime())
	assert.NotEqual(t, 0, ret.AsDateTime())
}

func TestDateMakeDate(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{2024}, IntegerValue{1})
	ret := dateMakeDate(buf)
	fmt.Printf("date_make_date: %d", ret.AsDate())
	assert.NotEqual(t, 0, ret.AsDate())
}

// Numeric functions

func TestNumericAbs(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{1})
	ret := numericAbs(buf)
	assert.Equal(t, int64(1), ret.AsInt())

	buf = nil
	buf = append(buf, IntegerValue{-1})
	ret = numericAbs(buf)
	assert.Equal(t, int64(1), ret.AsInt())
}

func TestNumericPi(t *testing.T) {
	var buf []Value

	ret := numericPi(buf)
	assert.Equal(t, math.Pi, ret.AsFloat())
}

func TestNumericFloor(t *testing.T) {
	var buf []Value

	buf = append(buf, FloatValue{1.1})
	ret := numericFloor(buf)
	assert.Equal(t, int64(1), ret.AsInt())

	buf = nil
	buf = append(buf, FloatValue{1.5})
	ret = numericFloor(buf)
	assert.Equal(t, int64(1), ret.AsInt())

	buf = nil
	buf = append(buf, FloatValue{1.9})
	ret = numericFloor(buf)
	assert.Equal(t, int64(1), ret.AsInt())
}

func TestNumericRound(t *testing.T) {
	var buf []Value

	buf = append(buf, FloatValue{1.1})
	ret := numericRound(buf)
	assert.Equal(t, int64(1), ret.AsInt())

	buf = nil
	buf = append(buf, FloatValue{1.5})
	ret = numericRound(buf)
	assert.Equal(t, int64(2), ret.AsInt())

	buf = nil
	buf = append(buf, FloatValue{1.9})
	ret = numericRound(buf)
	assert.Equal(t, int64(2), ret.AsInt())
}

func TestNumericSquare(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{2})
	ret := numericSquare(buf)
	assert.Equal(t, int64(4), ret.AsInt())
}

func TestNumericSin(t *testing.T) {
	var buf []Value

	buf = append(buf, FloatValue{0})
	ret := numericSin(buf)
	assert.Equal(t, float64(0), ret.AsFloat())

	buf = nil
	buf = append(buf, FloatValue{90})
	ret = numericSin(buf)
	assert.NotEqual(t, float64(0), ret.AsFloat())
}

func TestNumericAsin(t *testing.T) {
	var buf []Value

	buf = append(buf, FloatValue{0})
	ret := numericAsin(buf)
	assert.Equal(t, float64(0), ret.AsFloat())
}

func TestNumericCos(t *testing.T) {
	var buf []Value

	buf = append(buf, FloatValue{0})
	ret := numericCos(buf)
	assert.Equal(t, float64(1), ret.AsFloat())
}

func TestNumericTan(t *testing.T) {
	var buf []Value

	buf = append(buf, FloatValue{0})
	ret := numericTan(buf)
	assert.Equal(t, float64(0), ret.AsFloat())
}

// General functions

func TestGeneralIsNull(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{1})
	ret := generalIsNull(buf)
	assert.Equal(t, false, ret.AsBool())

	buf = nil
	buf = append(buf, NullValue{})
	ret = generalIsNull(buf)
	assert.Equal(t, true, ret.AsBool())
}

func TestGeneralIsNumeric(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{""})
	ret := generalIsNumeric(buf)
	assert.Equal(t, false, ret.AsBool())

	buf = nil
	buf = append(buf, IntegerValue{1})
	ret = generalIsNumeric(buf)
	assert.Equal(t, true, ret.AsBool())
}

func TestGeneralTypeOf(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{""})
	ret := generalTypeOf(buf)
	assert.Equal(t, "Text", ret.AsText())

	buf = nil
	buf = append(buf, IntegerValue{1})
	ret = generalTypeOf(buf)
	assert.Equal(t, "Integer", ret.AsText())
}
