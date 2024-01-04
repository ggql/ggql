package ast

// DataType represents the data types.
type DataType int

const (
	Any DataType = iota
	Text
	Integer
	Float
	Boolean
	Date
	Time
	DateTime
	Undefined
	Null
)

const (
	AnyStr       = "Any"
	TextStr      = "Text"
	IntegerStr   = "Integer"
	FloatStr     = "Float"
	BooleanStr   = "Boolean"
	DateStr      = "Date"
	TimeStr      = "Time"
	DateTimeStr  = "DateTime"
	UndefinedStr = "Undefined"
	NullStr      = "Null"
	UnknownStr   = "Unknown"
)

// isType checks if the data type matches the given type.
func (dt DataType) isType(dataType DataType) bool {
	return dt == dataType
}

// isInt checks if the data type is Integer.
func (dt DataType) isInt() bool {
	return dt.isType(Integer)
}

// isFloat checks if the data type is Float.
func (dt DataType) isFloat() bool {
	return dt.isType(Float)
}

// isNumber checks if the data type is Integer or Float.
func (dt DataType) isNumber() bool {
	return dt.isInt() || dt.isFloat()
}

// isText checks if the data type is Text.
func (dt DataType) isText() bool {
	return dt.isType(Text)
}

// isTime checks if the data type is Time.
func (dt DataType) isTime() bool {
	return dt.isType(Time)
}

// isDate checks if the data type is Date.
func (dt DataType) isDate() bool {
	return dt.isType(Date)
}

// isDateTime checks if the data type is DateTime.
func (dt DataType) isDateTime() bool {
	return dt.isType(DateTime)
}

// isUndefined checks if the data type is Undefined.
func (dt DataType) isUndefined() bool {
	return dt.isType(Undefined)
}

// literal returns the string representation of the data type.
func (dt DataType) literal() string {
	switch dt {
	case Any:
		return AnyStr
	case Text:
		return TextStr
	case Integer:
		return IntegerStr
	case Float:
		return FloatStr
	case Boolean:
		return BooleanStr
	case Date:
		return DateStr
	case Time:
		return TimeStr
	case DateTime:
		return DateTimeStr
	case Undefined:
		return UndefinedStr
	case Null:
		return NullStr
	default:
		return UnknownStr
	}
}

var tablesFieldsTypes = map[string]DataType{
	"commit_id":     Text,
	"title":         Text,
	"message":       Text,
	"name":          Text,
	"full_name":     Text,
	"insertions":    Integer,
	"deletions":     Integer,
	"files_changed": Integer,
	"email":         Text,
	"type":          Text,
	"datetime":      DateTime,
	"is_head":       Boolean,
	"is_remote":     Boolean,
	"commit_count":  Integer,
	"repo":          Text,
}
