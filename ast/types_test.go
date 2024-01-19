package ast

import (
	"testing"
)

func TestDataType(t *testing.T) {
	if Integer.Clone() != Integer {
		t.Error("Integer.Clone() should be Integer")
	}

	if !Any.IsType(Any) {
		t.Error("Any.IsType(Any) should be true")
	}

	if Integer.IsType(Float) {
		t.Error("Integer.IsType(Float) should be false")
	}

	if !Integer.IsInt() {
		t.Error("Integer.IsInt should be true")
	}

	if !Float.IsFloat() {
		t.Error("Float.IsFloat should be true")
	}

	if !Integer.IsNumber() {
		t.Error("Integer.IsNumber should be true")
	}

	if !Float.IsNumber() {
		t.Error("Float.IsNumber should be true")
	}

	if !Text.IsText() {
		t.Error("Text.IsText should be true")
	}

	if !Time.IsTime() {
		t.Error("Time.IsTime should be true")
	}

	if !Date.IsDate() {
		t.Error("Date.IsDate should be true")
	}

	if !DateTime.IsDateTime() {
		t.Error("DateTime.IsDateTime should be true")
	}

	if !Undefined.IsUndefined() {
		t.Error("Undefined.IsUndefined should be true")
	}
}

func TestLiteral(t *testing.T) {
	if Any.Literal() != "Any" {
		t.Error("Any.Literal() should be 'Any'")
	}

	if Text.Literal() != "Text" {
		t.Error("Text.Literal() should be 'Text'")
	}

	if Integer.Literal() != "Integer" {
		t.Error("Integer.Literal() should be 'Integer'")
	}

	if Float.Literal() != "Float" {
		t.Error("Float.Literal() should be 'Float'")
	}

	if Boolean.Literal() != "Boolean" {
		t.Error("Boolean.Literal() should be 'Boolean'")
	}

	if Date.Literal() != "Date" {
		t.Error("Date.Literal() should be 'Date'")
	}

	if Time.Literal() != "Time" {
		t.Error("Time.Literal() should be 'Time'")
	}

	if DateTime.Literal() != "DateTime" {
		t.Error("DateTime.Literal() should be 'DateTime'")
	}

	if Undefined.Literal() != "Undefined" {
		t.Error("Undefined.Literal() should be 'Undefined'")
	}

	if Null.Literal() != "Null" {
		t.Error("Null.Literal() should be 'Null'")
	}
}
