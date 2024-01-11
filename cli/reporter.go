package cli

import (
	"strconv"
	"strings"

	"github.com/ggql/ggql/parser"
)

const (
	PorpotLength int = 6
)

type DiagnosticReporter struct {
	stdout ColoredStream
}

func (d *DiagnosticReporter) ReportError(message string) {
	d.stdout.SetColor(RED)
	d.stdout.Printf("ERROR: ")
	d.stdout.Printlnf(message)
	d.stdout.Reset()
}

func (d *DiagnosticReporter) ReportGqlError(err parser.GQLError) {
	d.stdout.SetColor(RED)
	start := err.Location.Start

	d.stdout.Printf(strings.Repeat("-", PorpotLength+start))
	d.stdout.Printlnf("^")
	d.stdout.Printf("Compiletime ERROR: ")

	end := err.Location.End
	message := err.Message

	d.stdout.Printf("[")
	d.stdout.Printf(strconv.Itoa(start))
	d.stdout.Printf(" - ")
	d.stdout.Printf(strconv.Itoa(end))
	d.stdout.Printf("] -> ")
	d.stdout.Printlnf(message)
	d.stdout.Reset()
}

func (d *DiagnosticReporter) ReportRuntimeError(message string) {
	d.stdout.SetColor(RED)
	d.stdout.Printf("RUNTIME EXCEPTION: ")
	d.stdout.Printlnf(message)
	d.stdout.Reset()
}
