package cli

import (
	. "github.com/ggql/ggql/parser"
	"strconv"
	"strings"
)

const (
	PORPOT_LENGTH int = 6
)

type DiagnosticReporter struct {
	stdout ColoredStream
}

func (self *DiagnosticReporter) ReportError(message string) {
	self.stdout.SetColor(RED)
	self.stdout.Printf("ERROR: ")
	self.stdout.Printlnf(message)
	self.stdout.Reset()
}

func (self *DiagnosticReporter) ReportGqlError(err GQLError) {
	self.stdout.SetColor(RED)
	start := err.Location.Start
	self.stdout.Printf(strings.Repeat("-", PORPOT_LENGTH+start))
	self.stdout.Printlnf("^")
	self.stdout.Printf("Compiletime ERROR: ")
	end := err.Location.End
	message := err.Message
	self.stdout.Printf("[")
	self.stdout.Printf(strconv.Itoa(start))
	self.stdout.Printf(" - ")
	self.stdout.Printf(strconv.Itoa(end))
	self.stdout.Printf("] -> ")
	self.stdout.Printlnf(message)
	self.stdout.Reset()
}

func (self *DiagnosticReporter) ReportRuntimeError(message string) {
	self.stdout.SetColor(RED)
	self.stdout.Printf("RUNTIME EXCEPTION: ")
	self.stdout.Printlnf(message)
	self.stdout.Reset()
}
