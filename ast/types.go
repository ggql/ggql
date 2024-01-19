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

func (dt DataType) Clone() DataType {
	return dt
}

func (dt DataType) IsType(dataType DataType) bool {
	return dt == dataType
}

func (dt DataType) IsInt() bool {
	return dt.IsType(Integer)
}

func (dt DataType) IsFloat() bool {
	return dt.IsType(Float)
}

func (dt DataType) IsNumber() bool {
	return dt.IsInt() || dt.IsFloat()
}

func (dt DataType) IsText() bool {
	return dt.IsType(Text)
}

func (dt DataType) IsTime() bool {
	return dt.IsType(Time)
}

func (dt DataType) IsDate() bool {
	return dt.IsType(Date)
}

func (dt DataType) IsDateTime() bool {
	return dt.IsType(DateTime)
}

func (dt DataType) IsUndefined() bool {
	return dt.IsType(Undefined)
}

func (dt DataType) Literal() string {
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
