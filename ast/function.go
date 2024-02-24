package ast

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type Function func([]Value) Value

type Prototype struct {
	Parameters []DataType
	Result     DataType
}

var Functions = map[string]Function{
	// String functions
	"lower":      textLowercase,
	"upper":      textUppercase,
	"reverse":    textReverse,
	"replicate":  textReplicate,
	"space":      textSpace,
	"trim":       textTrim,
	"ltrim":      textLeftTrim,
	"rtrim":      textRightTrim,
	"len":        textLen,
	"ascii":      textAscii,
	"left":       textLeft,
	"datalength": textDataLength,
	"char":       textChar,
	"nchar":      textChar,
	"charindex":  textCharIndex,
	"replace":    textReplace,
	"substring":  textSubstring,
	"stuff":      textStuff,
	"right":      textRight,
	"translate":  textTranslate,
	"soundex":    textSoundex,
	"concat":     textConcat,
	"concat_ws":  textConcatWs,
	"unicode":    textUnicode,
	"strcmp":     textStrcmp,

	// Date functions
	"current_date":      dateCurrentDate,
	"current_time":      dateCurrentTime,
	"current_timestamp": dateCurrentTimestamp,
	"now":               dateCurrentTimestamp,
	"makedate":          dateMakeDate,
	"maketime":          dateMakeTime,
	"day":               dateDay,
	"dayname":           dateDayname,
	"monthname":         dateMonthname,
	"hour":              dateHour,
	"isdate":            dateIsDate,

	// Numeric functions
	"abs":    numericAbs,
	"pi":     numericPi,
	"floor":  numericFloor,
	"round":  numericRound,
	"square": numericSquare,
	"sin":    numericSin,
	"asin":   numericAsin,
	"cos":    numericCos,
	"acos":   numericAcos,
	"tan":    numericTan,
	"atan":   numericAtan,
	"atn2":   numericAtn2,
	"sign":   numericSign,

	// General functions
	"isnull":    generalIsNull,
	"isnumeric": generalIsNumeric,
	"typeof":    generalTypeOf,
	"greatest":  generalGreatest,
	"least":     generalLeast,
}

var Prototypes = map[string]Prototype{
	// String functions
	"lower":      {Parameters: []DataType{Text{}}, Result: Text{}},
	"upper":      {Parameters: []DataType{Text{}}, Result: Text{}},
	"reverse":    {Parameters: []DataType{Text{}}, Result: Text{}},
	"replicate":  {Parameters: []DataType{Text{}, Integer{}}, Result: Text{}},
	"space":      {Parameters: []DataType{Integer{}}, Result: Text{}},
	"trim":       {Parameters: []DataType{Text{}}, Result: Text{}},
	"ltrim":      {Parameters: []DataType{Text{}}, Result: Text{}},
	"rtrim":      {Parameters: []DataType{Text{}}, Result: Text{}},
	"len":        {Parameters: []DataType{Text{}}, Result: Integer{}},
	"ascii":      {Parameters: []DataType{Text{}}, Result: Integer{}},
	"left":       {Parameters: []DataType{Text{}, Integer{}}, Result: Text{}},
	"datalength": {Parameters: []DataType{Text{}}, Result: Integer{}},
	"char":       {Parameters: []DataType{Integer{}}, Result: Text{}},
	"nchar":      {Parameters: []DataType{Integer{}}, Result: Text{}},
	"charindex":  {Parameters: []DataType{Text{}, Text{}}, Result: Integer{}},
	"replace":    {Parameters: []DataType{Text{}, Text{}, Text{}}, Result: Text{}},
	"substring":  {Parameters: []DataType{Text{}, Integer{}, Integer{}}, Result: Text{}},
	"stuff":      {Parameters: []DataType{Text{}, Integer{}, Integer{}, Text{}}, Result: Text{}},
	"right":      {Parameters: []DataType{Text{}, Integer{}}, Result: Text{}},
	"translate":  {Parameters: []DataType{Text{}, Text{}, Text{}}, Result: Text{}},
	"soundex":    {Parameters: []DataType{Text{}}, Result: Text{}},
	"concat":     {Parameters: []DataType{Any{}, Any{}, Varargs{Any{}}}, Result: Text{}},
	"concat_ws":  {Parameters: []DataType{Text{}, Any{}, Any{}, Varargs{Any{}}}, Result: Text{}},
	"unicode":    {Parameters: []DataType{Text{}}, Result: Integer{}},
	"strcmp":     {Parameters: []DataType{Text{}, Text{}}, Result: Integer{}},

	// Date functions
	"current_date":      {Parameters: []DataType{}, Result: Date{}},
	"current_time":      {Parameters: []DataType{}, Result: Time{}},
	"current_timestamp": {Parameters: []DataType{}, Result: DateTime{}},
	"now":               {Parameters: []DataType{}, Result: DateTime{}},
	"makedate":          {Parameters: []DataType{Integer{}, Integer{}}, Result: Date{}},
	"maketime":          {Parameters: []DataType{Integer{}, Integer{}, Integer{}}, Result: Time{}},
	"day":               {Parameters: []DataType{Date{}}, Result: Integer{}},
	"dayname":           {Parameters: []DataType{Date{}}, Result: Text{}},
	"monthname":         {Parameters: []DataType{Date{}}, Result: Text{}},
	"hour":              {Parameters: []DataType{DateTime{}}, Result: Integer{}},
	"isdate":            {Parameters: []DataType{Any{}}, Result: Boolean{}},

	// Numeric functions
	"abs":    {Parameters: []DataType{Integer{}}, Result: Integer{}},
	"pi":     {Parameters: []DataType{}, Result: Float{}},
	"floor":  {Parameters: []DataType{Float{}}, Result: Integer{}},
	"round":  {Parameters: []DataType{Float{}}, Result: Integer{}},
	"square": {Parameters: []DataType{Integer{}}, Result: Integer{}},
	"sin":    {Parameters: []DataType{Float{}}, Result: Float{}},
	"asin":   {Parameters: []DataType{Float{}}, Result: Float{}},
	"cos":    {Parameters: []DataType{Float{}}, Result: Float{}},
	"acos":   {Parameters: []DataType{Float{}}, Result: Float{}},
	"tan":    {Parameters: []DataType{Float{}}, Result: Float{}},
	"atan":   {Parameters: []DataType{Float{}}, Result: Float{}},
	"atn2":   {Parameters: []DataType{Float{}, Float{}}, Result: Float{}},
	"sign":   {Parameters: []DataType{Variant{Integer{}, Float{}}}, Result: Integer{}},

	// General functions
	"isnull":    {Parameters: []DataType{Any{}}, Result: Boolean{}},
	"isnumeric": {Parameters: []DataType{Any{}}, Result: Boolean{}},
	"typeof":    {Parameters: []DataType{Any{}}, Result: Text{}},
	"greatest":  {Parameters: []DataType{Any{}, Any{}, Varargs{Any{}}}, Result: Any{}},
	"least":     {Parameters: []DataType{Any{}, Any{}, Varargs{Any{}}}, Result: Any{}},
}

// String functions

func textLowercase(inputs []Value) Value {
	return TextValue{strings.ToLower(inputs[0].AsText())}
}

func textUppercase(inputs []Value) Value {
	return TextValue{strings.ToUpper(inputs[0].AsText())}
}

func textReverse(inputs []Value) Value {
	runes := []rune(inputs[0].AsText())

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return TextValue{string(runes)}
}

func textReplicate(inputs []Value) Value {
	return TextValue{strings.Repeat(inputs[0].AsText(), int(inputs[1].AsInt()))}
}

func textSpace(inputs []Value) Value {
	return TextValue{strings.Repeat(" ", int(inputs[0].AsInt()))}
}

func textTrim(inputs []Value) Value {
	return TextValue{strings.TrimSpace(inputs[0].AsText())}
}

func textLeftTrim(inputs []Value) Value {
	return TextValue{strings.TrimLeftFunc(inputs[0].AsText(), unicode.IsSpace)}
}

func textRightTrim(inputs []Value) Value {
	return TextValue{strings.TrimRightFunc(inputs[0].AsText(), unicode.IsSpace)}
}

func textLen(inputs []Value) Value {
	return IntegerValue{int64(len(inputs[0].AsText()))}
}

func textAscii(inputs []Value) Value {
	if inputs[0].AsText() == "" {
		return IntegerValue{0}
	}

	return IntegerValue{int64(inputs[0].AsText()[0])}
}

func textLeft(inputs []Value) Value {
	text := inputs[0].AsText()
	if text == "" {
		return TextValue{""}
	}

	numberOfChars := int(inputs[1].AsInt())
	if numberOfChars > len(text) {
		return TextValue{text}
	}

	substring := text[:numberOfChars]

	return TextValue{substring}
}

func textDataLength(inputs []Value) Value {
	return IntegerValue{int64(len([]byte(inputs[0].AsText())))}
}

func textChar(inputs []Value) Value {
	code := inputs[0].AsInt()
	if code >= 0 && code <= 0x10FFFF {
		return TextValue{string(rune(code))}
	}

	return TextValue{""}
}

func textCharIndex(inputs []Value) Value {
	substr := inputs[0].AsText()
	input := inputs[1].AsText()

	index := strings.Index(strings.ToLower(input), strings.ToLower(substr))
	if index == -1 {
		return IntegerValue{0}
	}

	return IntegerValue{int64(index + 1)}
}

func textReplace(inputs []Value) Value {
	text := inputs[0].AsText()
	oldString := inputs[1].AsText()
	newString := inputs[2].AsText()

	result := strings.ReplaceAll(strings.ToLower(text), strings.ToLower(oldString), newString)

	return TextValue{result}
}

func textSubstring(inputs []Value) Value {
	text := inputs[0].AsText()
	start := int(inputs[1].AsInt()) - 1
	length := int(inputs[2].AsInt())

	if start > len(text) || length > len(text) {
		return TextValue{text}
	}

	if length < 0 {
		return TextValue{""}
	}

	end := start + length
	if end > len(text) {
		end = len(text)
	}

	return TextValue{text[start:end]}
}

func textStuff(inputs []Value) Value {
	text := inputs[0].AsText()
	start := inputs[1].AsInt() - 1
	length := inputs[2].AsInt()
	newString := inputs[3].AsText()

	if text == "" {
		return TextValue{text}
	}

	if start > int64(len(text)) || length > int64(len(text)) {
		return TextValue{text}
	}

	textRunes := []rune(text)
	newStringRunes := []rune(newString)
	textRunes = append(textRunes[:start], append(newStringRunes, textRunes[start+length:]...)...)

	return TextValue{string(textRunes)}
}

func textRight(inputs []Value) Value {
	text := inputs[0].AsText()
	if text == "" {
		return TextValue{""}
	}

	numberOfChars := inputs[1].AsInt()
	if numberOfChars > int64(len(text)) {
		return TextValue{text}
	}

	return TextValue{text[len(text)-int(numberOfChars):]}
}

func textTranslate(inputs []Value) Value {
	text := inputs[0].AsText()
	characters := inputs[1].AsText()
	translations := inputs[2].AsText()

	if len(translations) != len(characters) {
		return TextValue{""}
	}

	for i, letter := range characters {
		text = strings.ReplaceAll(text, string(letter), string(translations[i]))
	}

	return TextValue{text}
}

func textUnicode(inputs []Value) Value {
	text := inputs[0].AsText()
	if text == "" {
		return IntegerValue{0}
	}

	r, _ := utf8.DecodeRuneInString(text)

	return IntegerValue{int64(r)}
}

// nolint: gomnd
func textSoundex(inputs []Value) Value {
	text := inputs[0].AsText()
	if text == "" {
		return TextValue{""}
	}

	result := string(text[0])

	for idx, letter := range text {
		if idx != 0 {
			letter = unicode.ToUpper(letter)
			if !strings.ContainsRune("AEIOUHWY", letter) {
				var intVal int64
				switch letter {
				case 'B', 'F', 'P', 'V':
					intVal = 1
				case 'C', 'G', 'J', 'K', 'Q', 'S', 'X', 'Z':
					intVal = 2
				case 'D', 'T':
					intVal = 3
				case 'L':
					intVal = 4
				case 'M', 'N':
					intVal = 5
				case 'R':
					intVal = 6
				default:
					intVal = 0
				}
				result += strconv.FormatInt(intVal, 10)

				if len(result) == 4 {
					return TextValue{result}
				}
			}
		}
	}

	if len(result) < 4 {
		diff := 4 - len(result)
		for i := 0; i < diff; i++ {
			result += "0"
		}
	}

	return TextValue{result}
}

func textConcat(inputs []Value) Value {
	var text []string

	for _, v := range inputs {
		text = append(text, v.AsText())
	}

	return TextValue{strings.Join(text, "")}
}

func textConcatWs(inputs []Value) Value {
	var text []string

	for _, v := range inputs[1:] {
		text = append(text, v.AsText())
	}

	separator := inputs[0].AsText()

	return TextValue{strings.Join(text, separator)}
}

func textStrcmp(inputs []Value) Value {
	comparison := strings.Compare(inputs[0].AsText(), inputs[1].AsText())
	switch {
	case comparison < 0:
		return IntegerValue{1}
	case comparison == 0:
		return IntegerValue{2}
	case comparison > 0:
		return IntegerValue{0}
	default:
		return IntegerValue{}
	}
}

// Date functions

func dateCurrentDate(inputs []Value) Value {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	return DateValue{timestamp}
}

func dateCurrentTime(inputs []Value) Value {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	t := time.Unix(0, timestamp*int64(time.Millisecond))

	return TimeValue{t.String()}
}

func dateCurrentTimestamp(inputs []Value) Value {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	return DateTimeValue{timestamp}
}

func dateMakeDate(inputs []Value) Value {
	year := int(inputs[0].AsInt())
	dayOfYear := int(inputs[1].AsInt())

	t := time.Date(year, time.January, dayOfYear, 0, 0, 0, 0, time.UTC)
	timestamp := t.UnixNano() / int64(time.Millisecond)

	return DateValue{timestamp}
}

func dateMakeTime(inputs []Value) Value {
	hour := inputs[0].AsInt()
	minute := inputs[1].AsInt()
	second := inputs[2].AsInt()

	return TimeValue{fmt.Sprintf("%d:%02d:%02d", hour, minute, second)}
}

func dateDay(inputs []Value) Value {
	date := inputs[0].AsDate()
	dateNum := DateToDayNumberInMonth(date)

	return IntegerValue{int64(dateNum)}
}

func dateDayname(inputs []Value) Value {
	date := inputs[0].AsDate()
	dateStr := DateToDayName(date)

	return TextValue{dateStr}
}

func dateMonthname(inputs []Value) Value {
	date := inputs[0].AsDate()
	monthStr := DateToMonthName(date)

	return TextValue{monthStr}
}

func dateHour(inputs []Value) Value {
	date := inputs[0].AsDateTime()
	hour := DateTimeToHour(date)

	return IntegerValue{hour}
}

func dateIsDate(inputs []Value) Value {
	return BooleanValue{inputs[0].DataType().IsDate()}
}

// Numeric functions

func numericAbs(inputs []Value) Value {
	value := inputs[0].AsInt()

	return IntegerValue{int64(math.Abs(float64(value)))}
}

func numericPi(inputs []Value) Value {
	pi := math.Pi

	return FloatValue{pi}
}

func numericFloor(inputs []Value) Value {
	floatValue := inputs[0].AsFloat()

	return IntegerValue{int64(math.Floor(floatValue))}
}

func numericRound(inputs []Value) Value {
	floatValue := inputs[0].AsFloat()

	return IntegerValue{int64(math.Round(floatValue))}
}

func numericSquare(inputs []Value) Value {
	intValue := inputs[0].AsInt()

	return IntegerValue{intValue * intValue}
}

func numericSin(inputs []Value) Value {
	floatValue := inputs[0].AsFloat()

	return FloatValue{math.Sin(floatValue)}
}

func numericAsin(inputs []Value) Value {
	floatValue := inputs[0].AsFloat()

	return FloatValue{math.Asin(floatValue)}
}

func numericCos(inputs []Value) Value {
	floatValue := inputs[0].AsFloat()

	return FloatValue{math.Cos(floatValue)}
}

func numericAcos(inputs []Value) Value {
	floatValue := inputs[0].AsFloat()
	return FloatValue{math.Acos(floatValue)}
}

func numericTan(inputs []Value) Value {
	floatValue := inputs[0].AsFloat()

	return FloatValue{math.Tan(floatValue)}
}

func numericAtan(inputs []Value) Value {
	floatValue := inputs[0].AsFloat()

	return FloatValue{math.Atan(floatValue)}
}

func numericAtn2(inputs []Value) Value {
	first := inputs[0].AsFloat()
	other := inputs[1].AsFloat()

	return FloatValue{math.Atan2(first, other)}
}

func numericSign(inputs []Value) Value {
	helper := func(x int64) int64 {
		switch {
		case x < 0:
			return -1
		case x > 0:
			return 1
		}
		return 0
	}

	value := inputs[0]

	if value.DataType().IsInt() {
		intValue := value.AsInt()
		return IntegerValue{helper(intValue)}
	}

	floatValue := value.AsFloat()

	if floatValue == 0.0 {
		return IntegerValue{0}
	} else if floatValue > 0.0 {
		return IntegerValue{1}
	}

	return IntegerValue{-1}
}

// General functions

func generalIsNull(inputs []Value) Value {
	return BooleanValue{inputs[0].DataType().IsNull()}
}

func generalIsNumeric(inputs []Value) Value {
	inputType := inputs[0].DataType()
	isNumber := inputType.IsInt() || inputType.IsFloat()

	return BooleanValue{isNumber}
}

func generalTypeOf(inputs []Value) Value {
	inputType := inputs[0].DataType()

	return TextValue{inputType.Fmt()}
}

func generalGreatest(inputs []Value) Value {
	maxValue := inputs[0]

	for _, value := range inputs[1:] {
		if ret := maxValue.Compare(value); ret == Greater {
			maxValue = value
		}
	}

	return maxValue
}

func generalLeast(inputs []Value) Value {
	minValue := inputs[0]

	for _, value := range inputs[1:] {
		if ret := minValue.Compare(value); ret == Less {
			minValue = value
		}
	}

	return minValue
}
