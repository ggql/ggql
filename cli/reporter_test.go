package cli

import (
	. "github.com/ggql/ggql/parser"
	"testing"
)

func TestDiagnosticReporter_ReportError(t *testing.T) {
	cstream := &ColoredStream{}
	type fields struct {
		stdout ColoredStream
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{name: "report error", fields: fields{stdout: *cstream}, args: args{message: "report_error"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &DiagnosticReporter{
				stdout: tt.fields.stdout,
			}
			self.ReportError(tt.args.message)
		})
	}
}

func TestDiagnosticReporter_ReportGqlError(t *testing.T) {
	cstream := &ColoredStream{}
	gerr := GQLError{Message: "This is an error message", Location: Location{Start: 20, End: 30}}
	type fields struct {
		stdout ColoredStream
	}
	type args struct {
		err GQLError
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{name: "report_gql_error", fields: fields{stdout: *cstream}, args: args{err: gerr}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &DiagnosticReporter{
				stdout: tt.fields.stdout,
			}
			self.ReportGqlError(tt.args.err)
		})
	}
}

func TestDiagnosticReporter_ReportRuntimeError(t *testing.T) {
	cstream := &ColoredStream{}
	type fields struct {
		stdout ColoredStream
	}
	type args struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{name: "report_runtime_error", fields: fields{stdout: *cstream}, args: args{message: "report_runtime_error"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			self := &DiagnosticReporter{
				stdout: tt.fields.stdout,
			}
			self.ReportRuntimeError(tt.args.message)
		})
	}
}
