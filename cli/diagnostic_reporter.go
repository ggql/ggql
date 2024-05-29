package cli

import (
	"fmt"

	"github.com/ggql/ggql/parser"
)

type DiagnosticReporter struct {
	stdout ColoredStream
}

func (d *DiagnosticReporter) ReportDiagnostic(query string, diagnostic *parser.Diagnostic) {
	d.stdout.SetColor(Red)
	d.stdout.Printlnf(fmt.Sprintf("[%s]: %s", diagnostic.DiaLabel(), diagnostic.DiaMessage()))

	if diagnostic.DiaLocation().Start >= 0 && diagnostic.DiaLocation().End >= 0 {
		d.stdout.Printlnf(fmt.Sprintf("=> Line %d, Column %d,", diagnostic.DiaLocation().Start, diagnostic.DiaLocation().End))
	}

	if query != "" {
		d.stdout.Printlnf("  |")
		d.stdout.Printlnf(fmt.Sprintf("1 | %s\n", query))
		if diagnostic.DiaLocation().Start >= 0 && diagnostic.DiaLocation().End >= 0 {
			d.stdout.Printf("  | ") // NOTE: d.stdout.Printf appends line break
			d.stdout.Printf(d.repeat("-", diagnostic.DiaLocation().Start))
			d.stdout.SetColor(Yellow)
			d.stdout.Printlnf(d.repeat("^", d.max(1, diagnostic.DiaLocation().End-diagnostic.DiaLocation().Start)))
			d.stdout.SetColor(Red)
		}
		d.stdout.Printlnf("  |")
	}

	d.stdout.SetColor(Yellow)

	for _, note := range diagnostic.DiaNotes() {
		d.stdout.Printlnf(fmt.Sprintf(" = Note: %s", note))
	}

	d.stdout.SetColor(Cyan)

	for _, help := range diagnostic.DiaHelps() {
		d.stdout.Printlnf(fmt.Sprintf(" = Help: %s", help))
	}

	d.stdout.SetColor(Blue)

	if diagnostic.DiaDocs() != "" {
		d.stdout.Printlnf(fmt.Sprintf(" = Docs: %s", diagnostic.DiaDocs()))
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
