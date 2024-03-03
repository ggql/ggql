package cli

import (
	"fmt"
)

type DiagnosticReporter struct {
	stdout ColoredStream
}

func (d *DiagnosticReporter) ReportDiagnostic(query string, diagnostic Diagnostic) {
	d.stdout.SetColor(Red)
	d.stdout.Printlnf(fmt.Sprintf("[%s]: %s", diagnostic.Label(), diagnostic.Message()))

	if diagnostic.Location() != nil {
		d.stdout.Printlnf(fmt.Sprintf("=> Line %d, Column %d,", diagnostic.Location()[0], diagnostic.Location()[1]))
	}

	if query != "" {
		d.stdout.Printlnf("  |")
		d.stdout.Printlnf(fmt.Sprintf("1 | %s", query))
		if diagnostic.Location() != nil {
			d.stdout.Printf("  | ")
			d.stdout.Printf(d.repeat("-", diagnostic.Location()[0]))
			d.stdout.SetColor(Yellow)
			d.stdout.Printlnf(d.repeat("^", d.max(1, diagnostic.Location()[1]-diagnostic.Location()[0])))
			d.stdout.SetColor(Red)
		}
		d.stdout.Printlnf("  |")
	}

	d.stdout.SetColor(Yellow)

	for _, note := range diagnostic.Notes() {
		d.stdout.Printlnf(fmt.Sprintf(" = Note: %s", note))
	}

	d.stdout.SetColor(Cyan)

	for _, help := range diagnostic.Helps() {
		d.stdout.Printlnf(fmt.Sprintf(" = Help: %s", help))
	}

	d.stdout.SetColor(Blue)

	if diagnostic.Docs() != "" {
		d.stdout.Printlnf(fmt.Sprintf(" = Docs: %s", diagnostic.Docs()))
	}

	d.stdout.Reset()
}

func (d *DiagnosticReporter) repeat(s string, count int) string {
	str := ""
	for i := 0; i < count; i++ {
		str += s
	}
	return str
}

func (d *DiagnosticReporter) max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
