package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ggql/ggql/parser"
)

func TestReportDiagnostic(t *testing.T) {
	reporter := DiagnosticReporter{}
	err := parser.NewError("error")

	reporter.ReportDiagnostic("keyword", *err)
	assert.Equal(t, nil, nil)
}

func TestRepeat(t *testing.T) {
	reporter := DiagnosticReporter{}

	buf := reporter.repeat("a", 2)
	assert.Equal(t, 2, len(buf))
}

func TestMax(t *testing.T) {
	reporter := DiagnosticReporter{}

	buf := reporter.max(2, 1)
	assert.Equal(t, 2, buf)

	buf = reporter.max(1, 2)
	assert.Equal(t, 2, buf)
}
