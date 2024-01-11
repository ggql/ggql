package ast

import (
	"testing"
)

func TestDataType(t *testing.T) {
	// Test Clone method
	if Integer.Clone() != Integer {
		t.Error("Integer.Clone() should be Integer")
	}

	// Test isType method
	if !Any.isType(Any) {
		t.Error("Any.isType(Any) should be true")
	}
	if Integer.isType(Float) {
		t.Error("Integer.isType(Float) should be false")
	}

	// Test isInt method
	if !Integer.isInt() {
		t.Error("Integer.isInt should be true")
	}

	// Test isFloat method
	if !Float.isFloat() {
		t.Error("Float.isFloat should be true")
	}

	// Test isNumber method
	if !Integer.isNumber() {
		t.Error("Integer.isNumber should be true")
	}
	if !Float.isNumber() {
		t.Error("Float.isNumber should be true")
	}

	// Test isText method
	if !Text.isText() {
		t.Error("Text.isText should be true")
	}

	// Test isTime method
	if !Time.isTime() {
		t.Error("Time.isTime should be true")
	}

	// Test isDate method
	if !Date.isDate() {
		t.Error("Date.isDate should be true")
	}

	// Test idDateTime method
	if !DateTime.isDateTime() {
		t.Error("DateTime.isDateTime should be true")
	}

	// Test idDateTime method
	if !Undefined.isUndefined() {
		t.Error("Undefined.isUndefined should be true")
	}
}

func TestLiteral(t *testing.T) {
	// Test literal method
	if Any.literal() != "Any" {
		t.Error("Any.literal() should be 'Any'")
	}
	if Text.literal() != "Text" {
		t.Error("Text.literal() should be 'Text'")
	}
	if Integer.literal() != "Integer" {
		t.Error("Integer.literal() should be 'Integer'")
	}
	if Float.literal() != "Float" {
		t.Error("Float.literal() should be 'Float'")
	}
	if Boolean.literal() != "Boolean" {
		t.Error("Boolean.literal() should be 'Boolean'")
	}
	if Date.literal() != "Date" {
		t.Error("Date.literal() should be 'Date'")
	}
	if Time.literal() != "Time" {
		t.Error("Time.literal() should be 'Time'")
	}
	if DateTime.literal() != "DateTime" {
		t.Error("DateTime.literal() should be 'DateTime'")
	}
	if Undefined.literal() != "Undefined" {
		t.Error("Undefined.literal() should be 'Undefined'")
	}
	if Null.literal() != "Null" {
		t.Error("Null.literal() should be 'Null'")
	}
}
