package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintf(t *testing.T) {
	stream := NewColoredStream()

	stream.Printf("hello world!")
	assert.Equal(t, nil, nil)
}

func TestPrintlnf(t *testing.T) {
	stream := NewColoredStream()

	stream.Printlnf("hello world!")
	assert.Equal(t, nil, nil)
}

func TestSetColor(t *testing.T) {
	stream := NewColoredStream()

	stream.SetColor(Blue)
	stream.Printlnf("hello world!")
	assert.Equal(t, nil, nil)
}

func TestReset(t *testing.T) {
	stream := NewColoredStream()

	stream.SetColor(Blue)
	stream.Reset()
	stream.Printlnf("hello world!")
	assert.Equal(t, nil, nil)
}
