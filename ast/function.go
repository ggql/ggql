package ast

import (
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

var functions = map[string]Function{
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
	"replace":    textReplace,
	"substring":  textSubstring,
	"stuff":      textStuff,
	"right":      textRight,
	"translate":  textTranslate,
	"soundex":    textSoundex,
	"concat":     textConcat,
	"unicode":    textUnicode,

	// Date functions
	"current_date":      dateCurrentDate,
	"current_time":      dateCurrentTime,
	"current_timestamp": dateCurrentTimestamp,
	"now":               dateCurrentTimestamp,
	"makedate":          dateMakeDate,

	// Numeric functions
	"abs":    numericAbs,
	"pi":     numericPi,
	"floor":  numericFloor,
	"round":  numericRound,
	"square": numericSquare,
	"sin":    numericSin,
	"asin":   numericAsin,
	"cos":    numericCos,
	"tan":    numericTan,

	// General functions
	"isnull":    generalIsNull,
	"isnumeric": generalIsNumeric,
	"typeof":    generalTypeOf,
}

var prototypes = map[string]Prototype{
	// String functions
	"lower":      {},
	"upper":      {},
	"reverse":    {},
	"replicate":  {},
	"space":      {},
	"trim":       {},
	"ltrim":      {},
	"rtrim":      {},
	"len":        {},
	"ascii":      {},
	"left":       {},
	"datalength": {},
	"char":       {},
	"nchar":      {},
	"replace":    {},
	"substring":  {},
	"stuff":      {},
	"right":      {},
	"translate":  {},
	"soundex":    {},
	"concat":     {},
	"unicode":    {},

	// Date functions
	"current_date":      {},
	"current_time":      {},
	"current_timestamp": {},
	"now":               {},
	"makedate":          {},

	// Numeric functions
	"abs":    {},
	"pi":     {},
	"floor":  {},
	"round":  {},
	"square": {},
	"sin":    {},
	"asin":   {},
	"cos":    {},
	"tan":    {},

	// General functions
	"isnull":    {},
	"isnumeric": {},
	"typeof":    {},
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
	if len(inputs[0].AsText()) == 0 {
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

	return IntergerValue{intValue * intValue}
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

func numericTan(inputs []Value) Value {
	floatValue := inputs[0].AsFloat()

	return FloatValue{math.Tan(floatValue)}
}

// General functions

func generalIsNull(inputs []Value) Value {
	return BooleanValue{inputs[0].DataType().isType(Null)}
}

func generalIsNumeric(inputs []Value) Value {
	inputType := inputs[0].DataType()
	isNumber := inputType.isInt() || inputType.isFloat()

	return BooleanValue{isNumber}
}

func generalTypeOf(inputs []Value) Value {
	inputType := inputs[0].DataType()

	return TextValue{inputType.literal()}
}
