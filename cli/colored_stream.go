package cli

import (
	"github.com/fatih/color"
	"strings"
)

const (
	Blue   = color.FgBlue
	Cyan   = color.FgCyan
	Green  = color.FgGreen
	Red    = color.FgRed
	Yellow = color.FgYellow
	White  = color.FgWhite
)

type ColoredStream struct {
	outColor color.Attribute
}

func NewColoredStream() *ColoredStream {
	return &ColoredStream{
		outColor: color.FgWhite,
	}
}

func (cs *ColoredStream) Printf(format string, value ...interface{}) {
	switch cs.outColor {
	case Green:
		color.Green(format, value...)
	case Red:
		color.Red(format, value...)
	case Yellow:
		color.Yellow(format, value...)
	case White:
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
