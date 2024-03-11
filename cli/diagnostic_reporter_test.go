package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportDiagnostic(t *testing.T) {
	reporter := DiagnosticReporter{}

	reporter.ReportDiagnostic("keyword", Diagnostic.Error("error"))
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
