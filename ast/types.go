package ast

import "strings"

const (
	typesNull      = "Null"
	typesUndefined = "Undefined"
)

type DataType interface {
	Equal(DataType) bool
	Fmt() string
	IsAny() bool
	IsBool() bool
	IsInt() bool
	IsFloat() bool
	IsNumber() bool
	IsText() bool
	IsTime() bool
	IsDate() bool
	IsDateTime() bool
	IsNull() bool
	IsUndefined() bool
	IsVariant() bool
	IsOptional() bool
	IsVarargs() bool
}

type Any struct{}
type Text struct{}
type Integer struct{}
type Float struct{}
type Boolean struct{}
type Date struct{}
type Time struct{}
type DateTime struct{}
type Undefined struct{}
type Null struct{}
type Variant []DataType
type Optional struct{ DataType }
type Varargs struct{ DataType }

var TablesFieldsTypes = map[string]DataType{
	"commit_id":     Text{},
	"title":         Text{},
	"message":       Text{},
	"name":          Text{},
	"full_name":     Text{},
	"insertions":    Integer{},
	"deletions":     Integer{},
	"files_changed": Integer{},
	"email":         Text{},
	"type":          Text{},
	"datetime":      DateTime{},
	"is_head":       Boolean{},
	"is_remote":     Boolean{},
	"commit_count":  Integer{},
	"repo":          Text{},
}

// Any implementation

func (a Any) Equal(other DataType) bool {
	return other.IsAny()
}

func (a Any) Fmt() string {
	return "Any"
}

func (a Any) IsAny() bool {
	return true
}

func (a Any) IsBool() bool {
	return true
}

func (a Any) IsInt() bool {
	return true
}

func (a Any) IsFloat() bool {
	return true
}

func (a Any) IsNumber() bool {
	return true
}

func (a Any) IsText() bool {
	return true
}

func (a Any) IsTime() bool {
	return true
}

func (a Any) IsDate() bool {
	return true
}

func (a Any) IsDateTime() bool {
	return true
}

func (a Any) IsNull() bool {
	return true
}

func (a Any) IsUndefined() bool {
	return true
}

func (a Any) IsVariant() bool {
	return true
}

func (a Any) IsOptional() bool {
	return true
}

func (a Any) IsVarargs() bool {
	return true
}

// Text implementation

func (t Text) Equal(other DataType) bool {
	return other.IsText()
}

func (t Text) Fmt() string {
	return "Text"
}

func (t Text) IsAny() bool {
	return false
}

func (t Text) IsBool() bool {
	return false
}

func (t Text) IsInt() bool {
	return false
}

func (t Text) IsFloat() bool {
	return false
}

func (t Text) IsNumber() bool {
	return false
}

func (t Text) IsText() bool {
	return true
}

func (t Text) IsTime() bool {
	return false
}

func (t Text) IsDate() bool {
	return false
}

func (t Text) IsDateTime() bool {
	return false
}

func (t Text) IsNull() bool {
	return false
}

func (t Text) IsUndefined() bool {
	return false
}

func (t Text) IsVariant() bool {
	return false
}

func (t Text) IsOptional() bool {
	return false
}

func (t Text) IsVarargs() bool {
	return false
}

// Integer implementation

func (i Integer) Equal(other DataType) bool {
	return other.IsInt()
}

func (i Integer) Fmt() string {
	return "Integer"
}

func (i Integer) IsAny() bool {
	return false
}

func (i Integer) IsBool() bool {
	return false
}

func (i Integer) IsInt() bool {
	return true
}

func (i Integer) IsFloat() bool {
	return false
}

func (i Integer) IsNumber() bool {
	return false
}

func (i Integer) IsText() bool {
	return false
}

func (i Integer) IsTime() bool {
	return false
}

func (i Integer) IsDate() bool {
	return false
}

func (i Integer) IsDateTime() bool {
	return false
}

func (i Integer) IsNull() bool {
	return false
}

func (i Integer) IsUndefined() bool {
	return false
}

func (i Integer) IsVariant() bool {
	return false
}

func (i Integer) IsOptional() bool {
	return false
}

func (i Integer) IsVarargs() bool {
	return false
}

// Float implementation

func (f Float) Equal(other DataType) bool {
	return other.IsFloat()
}

func (f Float) Fmt() string {
	return "Float"
}

func (f Float) IsAny() bool {
	return false
}

func (f Float) IsBool() bool {
	return false
}

func (f Float) IsInt() bool {
	return false
}

func (f Float) IsFloat() bool {
	return true
}

func (f Float) IsNumber() bool {
	return false
}

func (f Float) IsText() bool {
	return false
}

func (f Float) IsTime() bool {
	return false
}

func (f Float) IsDate() bool {
	return false
}

func (f Float) IsDateTime() bool {
	return false
}

func (f Float) IsNull() bool {
	return false
}

func (f Float) IsUndefined() bool {
	return false
}

func (f Float) IsVariant() bool {
	return false
}

func (f Float) IsOptional() bool {
	return false
}

func (f Float) IsVarargs() bool {
	return false
}

// Boolean implementation

func (b Boolean) Equal(other DataType) bool {
	return other.IsBool()
}

func (b Boolean) Fmt() string {
	return "Boolean"
}

func (b Boolean) IsAny() bool {
	return false
}

func (b Boolean) IsBool() bool {
	return true
}

func (b Boolean) IsInt() bool {
	return false
}

func (b Boolean) IsFloat() bool {
	return false
}

func (b Boolean) IsNumber() bool {
	return false
}

func (b Boolean) IsText() bool {
	return false
}

func (b Boolean) IsTime() bool {
	return false
}

func (b Boolean) IsDate() bool {
	return false
}

func (b Boolean) IsDateTime() bool {
	return false
}

func (b Boolean) IsNull() bool {
	return false
}

func (b Boolean) IsUndefined() bool {
	return false
}

func (b Boolean) IsVariant() bool {
	return false
}

func (b Boolean) IsOptional() bool {
	return false
}

func (b Boolean) IsVarargs() bool {
	return false
}

// Date implementation

func (d Date) Equal(other DataType) bool {
	return other.IsDate()
}

func (d Date) Fmt() string {
	return "Date"
}

func (d Date) IsAny() bool {
	return false
}

func (d Date) IsBool() bool {
	return false
}

func (d Date) IsInt() bool {
	return false
}

func (d Date) IsFloat() bool {
	return false
}

func (d Date) IsNumber() bool {
	return false
}

func (d Date) IsText() bool {
	return false
}

func (d Date) IsTime() bool {
	return false
}

func (d Date) IsDate() bool {
	return true
}

func (d Date) IsDateTime() bool {
	return false
}

func (d Date) IsNull() bool {
	return false
}

func (d Date) IsUndefined() bool {
	return false
}

func (d Date) IsVariant() bool {
	return false
}

func (d Date) IsOptional() bool {
	return false
}

func (d Date) IsVarargs() bool {
	return false
}

// Time implementation

func (t Time) Equal(other DataType) bool {
	return other.IsTime()
}

func (t Time) Fmt() string {
	return "Time"
}

func (t Time) IsAny() bool {
	return false
}

func (t Time) IsBool() bool {
	return false
}

func (t Time) IsInt() bool {
	return false
}

func (t Time) IsFloat() bool {
	return false
}

func (t Time) IsNumber() bool {
	return false
}

func (t Time) IsText() bool {
	return false
}

func (t Time) IsTime() bool {
	return true
}

func (t Time) IsDate() bool {
	return false
}

func (t Time) IsDateTime() bool {
	return false
}

func (t Time) IsNull() bool {
	return false
}

func (t Time) IsUndefined() bool {
	return false
}

func (t Time) IsVariant() bool {
	return false
}

func (t Time) IsOptional() bool {
	return false
}

func (t Time) IsVarargs() bool {
	return false
}

// DateTime implementation

func (dt DateTime) Equal(other DataType) bool {
	return other.IsDateTime()
}

func (dt DateTime) Fmt() string {
	return "DateTime"
}

func (dt DateTime) IsAny() bool {
	return false
}

func (dt DateTime) IsBool() bool {
	return false
}

func (dt DateTime) IsInt() bool {
	return false
}

func (dt DateTime) IsFloat() bool {
	return false
}

func (dt DateTime) IsNumber() bool {
	return false
}

func (dt DateTime) IsText() bool {
	return false
}

func (dt DateTime) IsTime() bool {
	return false
}

func (dt DateTime) IsDate() bool {
	return false
}

func (dt DateTime) IsDateTime() bool {
	return true
}

func (dt DateTime) IsNull() bool {
	return false
}

func (dt DateTime) IsUndefined() bool {
	return false
}

func (dt DateTime) IsVariant() bool {
	return false
}

func (dt DateTime) IsOptional() bool {
	return false
}

func (dt DateTime) IsVarargs() bool {
	return false
}

// Undefined implementation

func (u Undefined) Equal(other DataType) bool {
	return other.IsUndefined()
}

func (u Undefined) Fmt() string {
	return typesUndefined
}

func (u Undefined) IsAny() bool {
	return false
}

func (u Undefined) IsBool() bool {
	return false
}

func (u Undefined) IsInt() bool {
	return false
}

func (u Undefined) IsFloat() bool {
	return false
}

func (u Undefined) IsNumber() bool {
	return false
}

func (u Undefined) IsText() bool {
	return false
}

func (u Undefined) IsTime() bool {
	return false
}

func (u Undefined) IsDate() bool {
	return false
}

func (u Undefined) IsDateTime() bool {
	return false
}

func (u Undefined) IsNull() bool {
	return false
}

func (u Undefined) IsUndefined() bool {
	return true
}

func (u Undefined) IsVariant() bool {
	return false
}

func (u Undefined) IsOptional() bool {
	return false
}

func (u Undefined) IsVarargs() bool {
	return false
}

// Null implementation

func (n Null) Equal(other DataType) bool {
	return other.IsNull()
}

func (n Null) Fmt() string {
	return typesNull
}

func (n Null) IsAny() bool {
	return false
}

func (n Null) IsBool() bool {
	return false
}

func (n Null) IsInt() bool {
	return false
}

func (n Null) IsFloat() bool {
	return false
}

func (n Null) IsNumber() bool {
	return false
}

func (n Null) IsText() bool {
	return false
}

func (n Null) IsTime() bool {
	return false
}

func (n Null) IsDate() bool {
	return false
}

func (n Null) IsDateTime() bool {
	return false
}

func (n Null) IsNull() bool {
	return true
}

func (n Null) IsUndefined() bool {
	return false
}

func (n Null) IsVariant() bool {
	return false
}

func (n Null) IsOptional() bool {
	return false
}

func (n Null) IsVarargs() bool {
	return false
}

// Variant implementation

func (v Variant) Equal(other DataType) bool {
	for _, item := range v {
		if item == other {
			return true
		}
	}

	if other.IsVariant() {
		for _, o := range other.(Variant) {
			for _, item := range v {
				if item == o {
					return true
				}
			}
		}
	}

	return false
}

func (v Variant) Fmt() string {
	var types []string

	for _, item := range v {
		types = append(types, item.Fmt())
	}

	return "[" + strings.Join(types, " | ") + "]"
}

func (v Variant) IsAny() bool {
	return false
}

func (v Variant) IsBool() bool {
	return false
}

func (v Variant) IsInt() bool {
	return false
}

func (v Variant) IsFloat() bool {
	return false
}

func (v Variant) IsNumber() bool {
	return false
}

func (v Variant) IsText() bool {
	return false
}

func (v Variant) IsTime() bool {
	return false
}

func (v Variant) IsDate() bool {
	return false
}

func (v Variant) IsDateTime() bool {
	return false
}

func (v Variant) IsNull() bool {
	return false
}

func (v Variant) IsUndefined() bool {
	return false
}

func (v Variant) IsVariant() bool {
	return true
}

func (v Variant) IsOptional() bool {
	return false
}

func (v Variant) IsVarargs() bool {
	return false
}

// Optional implementation

func (o Optional) Equal(other DataType) bool {
	if o.DataType == other {
		return true
	}

	if other.IsOptional() {
		if o.DataType == other.(Optional).DataType {
			return true
		}
	}

	return false
}

func (o Optional) Fmt() string {
	return o.DataType.Fmt()
}

func (o Optional) IsAny() bool {
	return false
}

func (o Optional) IsBool() bool {
	return false
}

func (o Optional) IsInt() bool {
	return false
}

func (o Optional) IsFloat() bool {
	return false
}

func (o Optional) IsNumber() bool {
	return false
}

func (o Optional) IsText() bool {
	return false
}

func (o Optional) IsTime() bool {
	return false
}

func (o Optional) IsDate() bool {
	return false
}

func (o Optional) IsDateTime() bool {
	return false
}

func (o Optional) IsNull() bool {
	return false
}

func (o Optional) IsUndefined() bool {
	return false
}

func (o Optional) IsVariant() bool {
	return false
}

func (o Optional) IsOptional() bool {
	return true
}

func (o Optional) IsVarargs() bool {
	return false
}

// Varargs implementation

func (va Varargs) Equal(other DataType) bool {
	if va.DataType == other {
		return true
	}

	if other.IsVarargs() {
		if va.DataType == other.(Varargs).DataType {
			return true
		}
	}

	return false
}

func (va Varargs) Fmt() string {
	return "..." + va.DataType.Fmt()
}

func (va Varargs) IsAny() bool {
	return false
}

func (va Varargs) IsBool() bool {
	return false
}

func (va Varargs) IsInt() bool {
	return false
}

func (va Varargs) IsFloat() bool {
	return false
}

func (va Varargs) IsNumber() bool {
	return false
}

func (va Varargs) IsText() bool {
	return false
}

func (va Varargs) IsTime() bool {
	return false
}

func (va Varargs) IsDate() bool {
	return false
}

func (va Varargs) IsDateTime() bool {
	return false
}

func (va Varargs) IsNull() bool {
	return false
}

func (va Varargs) IsUndefined() bool {
	return false
}

func (va Varargs) IsVariant() bool {
	return false
}

func (va Varargs) IsOptional() bool {
	return false
}

func (va Varargs) IsVarargs() bool {
	return true
}
