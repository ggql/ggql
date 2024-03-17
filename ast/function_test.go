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

func TestTextCharIndex(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"h"}, TextValue{"hello"})
	ret := textCharIndex(buf)
	assert.Equal(t, int64(1), ret.AsInt())

	buf = nil
	buf = append(buf, TextValue{"w"}, TextValue{"hello"})
	ret = textCharIndex(buf)
	assert.Equal(t, int64(0), ret.AsInt())
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

func TestTextConcatWs(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{" "}, TextValue{"hello"}, TextValue{"world"})
	ret := textConcatWs(buf)
	assert.Equal(t, "hello world", ret.AsText())
}

func TestTextStrcmp(t *testing.T) {
	var buf []Value

	buf = append(buf, TextValue{"hello"}, TextValue{"hello"})
	ret := textStrcmp(buf)
	assert.Equal(t, int64(2), ret.AsInt())

	buf = nil
	buf = append(buf, TextValue{"hello"}, TextValue{"world"})
	ret = textStrcmp(buf)
	assert.Equal(t, int64(1), ret.AsInt())
}

// Date functions

func TestDateCurrentDate(t *testing.T) {
	var buf []Value

	ret := dateCurrentDate(buf)
	fmt.Printf("dateCurrentDate: %d", ret.AsDate())
	assert.NotEqual(t, 0, ret.AsDate())
}

func TestDateCurrentTime(t *testing.T) {
	var buf []Value

	ret := dateCurrentTime(buf)
	fmt.Printf("dateCurrentTime: %s", ret.AsTime())
	assert.NotEqual(t, "", ret.AsTime())
}

func TestDateCurrentTimestamp(t *testing.T) {
	var buf []Value

	ret := dateCurrentTimestamp(buf)
	fmt.Printf("dateCurrentTimestamp: %d", ret.AsDateTime())
	assert.NotEqual(t, 0, ret.AsDateTime())
}

func TestDateMakeDate(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{2024}, IntegerValue{1})
	ret := dateMakeDate(buf)
	fmt.Printf("dateMakeDate: %d", ret.AsDate())
	assert.NotEqual(t, 0, ret.AsDate())
}

func TestDateMakeTime(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{23}, IntegerValue{59}, IntegerValue{59})
	ret := dateMakeTime(buf)
	fmt.Printf("dateMakeTime: %s", ret.AsTime())
	assert.NotEqual(t, "", ret.AsTime())
}

func TestDateDay(t *testing.T) {
	var buf []Value

	buf = append(buf, DateValue{1705117592})
	ret := dateDay(buf)
	fmt.Printf("dateDay: %d", ret.AsInt())
	assert.NotEqual(t, 0, ret.AsInt())
}

func TestDateDayname(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{1705117592})
	ret := dateDayname(buf)
	fmt.Printf("dateDayname: %s", ret.AsText())
	assert.NotEqual(t, "", ret.AsText())
}

func TestDateMonthname(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{1705117592})
	ret := dateMonthname(buf)
	fmt.Printf("dateMonthname: %s", ret.AsText())
	assert.NotEqual(t, "", ret.AsText())
}

func TestDateHour(t *testing.T) {
	var buf []Value

	buf = append(buf, DateTimeValue{1705117592})
	ret := dateHour(buf)
	fmt.Printf("dateHour: %d", ret.AsInt())
	assert.NotEqual(t, 0, ret.AsInt())
}

func TestDateIsDate(t *testing.T) {
	var buf []Value

	buf = append(buf, DateValue{1705117592})
	ret := dateIsDate(buf)
	fmt.Printf("dateIsDate: %t", ret.AsBool())
	assert.Equal(t, true, ret.AsBool())
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

func TestNumericAcos(t *testing.T) {
	var buf []Value

	buf = append(buf, FloatValue{0})
	ret := numericAcos(buf)
	assert.NotEqual(t, float64(0), ret.AsFloat())
}

func TestNumericTan(t *testing.T) {
	var buf []Value

	buf = append(buf, FloatValue{0})
	ret := numericTan(buf)
	assert.Equal(t, float64(0), ret.AsFloat())
}

func TestNumericAtan(t *testing.T) {
	var buf []Value

	buf = append(buf, FloatValue{0})
	ret := numericAtan(buf)
	assert.Equal(t, float64(0), ret.AsFloat())
}

func TestNumericAtn2(t *testing.T) {
	var buf []Value

	buf = append(buf, FloatValue{0}, FloatValue{0})
	ret := numericAtn2(buf)
	assert.Equal(t, float64(0), ret.AsFloat())
}

func TestNumericSign(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{0})
	ret := numericSign(buf)
	assert.Equal(t, int64(0), ret.AsInt())

	buf = nil
	buf = append(buf, IntegerValue{1})
	ret = numericSign(buf)
	assert.Equal(t, int64(1), ret.AsInt())

	buf = nil
	buf = append(buf, IntegerValue{-1})
	ret = numericSign(buf)
	assert.Equal(t, int64(-1), ret.AsInt())

	buf = nil
	buf = append(buf, FloatValue{0})
	ret = numericSign(buf)
	assert.Equal(t, int64(0), ret.AsInt())

	buf = nil
	buf = append(buf, FloatValue{1})
	ret = numericSign(buf)
	assert.Equal(t, int64(1), ret.AsInt())

	buf = nil
	buf = append(buf, FloatValue{-1})
	ret = numericSign(buf)
	assert.Equal(t, int64(-1), ret.AsInt())
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

func TestGeneralGreatest(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{1}, IntegerValue{2}, IntegerValue{3})
	ret := generalGreatest(buf)
	assert.Equal(t, int64(3), ret.AsInt())
}

func TestGeneralLeast(t *testing.T) {
	var buf []Value

	buf = append(buf, IntegerValue{1}, IntegerValue{2}, IntegerValue{3})
	ret := generalLeast(buf)
	assert.Equal(t, int64(1), ret.AsInt())
}
