package cli

import (
	"strings"

	"github.com/fatih/color"
)

const (
	GREEN  = color.FgGreen
	RED    = color.FgRed
	YELLOW = color.FgYellow
	WHITE  = color.FgWhite
)

type ColoredStream struct {
	outColor color.Attribute
}

// The Default implementation is replaced with a New function
func NewColoredStream() *ColoredStream {
	// Default color is white
	return &ColoredStream{
		outColor: color.FgWhite,
	}
}

func (cs *ColoredStream) Printf(format string, value ...interface{}) {
	switch cs.outColor {
	case GREEN:
		color.Green(format, value...)
	case RED:
		color.Red(format, value...)
	case YELLOW:
		color.Yellow(format, value...)
	case WHITE:
		color.White(format, value...)
	}
}
func (cs *ColoredStream) Printlnf(format string, value ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	cs.Printf(format, value...)
}
func (cs *ColoredStream) SetColor(c color.Attribute) {
	cs.outColor = c
}

func (cs *ColoredStream) Reset() {
	cs.outColor = color.FgWhite
}
