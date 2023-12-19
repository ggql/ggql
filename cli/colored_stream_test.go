package cli

import (
	"reflect"
	"testing"
)

func TestSetColor(t *testing.T) {
	var tmp = NewColoredStream()
	var setColor = RED
	tmp.SetColor(setColor)
	var getColor = tmp.outColor
	if !reflect.DeepEqual(getColor, setColor) {
		t.Errorf("setColor: %v  getColor: %v", setColor, getColor)
	}
}

func TestReSet(t *testing.T) {
	var tmp = NewColoredStream()
	var wantColor = WHITE
	tmp.Reset()
	var getColor = tmp.outColor
	if !reflect.DeepEqual(getColor, wantColor) {
		t.Errorf("wantColor: %v  getColor: %v", wantColor, getColor)
	}
}
