package cli

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyArguments(t *testing.T) {
	actual := ParseArguments([]string{"ggql"})
	expected := Command{}

	ret := reflect.DeepEqual(actual.ReplMode, expected.ReplMode)
	assert.Equal(t, true, ret)
}

func TestReposArguments(t *testing.T) {
	actual := ParseArguments([]string{"ggql", "--repos", "."})
	expected := Command{
		QueryMode: struct {
			Query     string
			Arguments Arguments
		}{
			Arguments: struct {
				repos        []string
				analysis     bool
				pagination   bool
				pageSize     int
				outputFormat OutputFormat
			}{
				repos: []string{"--repos", "."},
			},
		},
	}

	ret := reflect.DeepEqual(actual.QueryMode.Arguments.repos, expected.QueryMode.Arguments.repos)
	assert.Equal(t, true, ret)
}

func TestQueryArguments(t *testing.T) {
	actual := ParseArguments([]string{"ggql", "--query", "Select * from table"})
	expected := Command{
		QueryMode: struct {
			Query     string
			Arguments Arguments
		}{
			Query: "Select * from table",
		},
	}

	ret := reflect.DeepEqual(actual.QueryMode.Query, expected.QueryMode.Query)
	assert.Equal(t, true, ret)
}

func TestPaginationArguments(t *testing.T) {
	t.Skip("Skipping TestPaginationArguments.")
}

func TestPagesizeArguments(t *testing.T) {
	actual := ParseArguments([]string{"ggql", "--pagesize", "10"})
	expected := Command{
		QueryMode: struct {
			Query     string
			Arguments Arguments
		}{
			Arguments: struct {
				repos        []string
				analysis     bool
				pagination   bool
				pageSize     int
				outputFormat OutputFormat
			}{
				pageSize: 10,
			},
		},
	}

	ret := reflect.DeepEqual(actual.QueryMode.Arguments.pageSize, expected.QueryMode.Arguments.pageSize)
	assert.Equal(t, true, ret)

	actual = ParseArguments([]string{"ggql", "--pagesize", "-"})
	expected.Error = "Argument --pagesize must be followed by the page size"

	ret = reflect.DeepEqual(actual.Error, expected.Error)
	assert.Equal(t, true, ret)
}

func TestOutputArguments(t *testing.T) {
	actual := ParseArguments([]string{"ggql", "--output", "csv"})
	expected := Command{
		QueryMode: struct {
			Query     string
			Arguments Arguments
		}{
			Arguments: struct {
				repos        []string
				analysis     bool
				pagination   bool
				pageSize     int
				outputFormat OutputFormat
			}{
				outputFormat: CSV,
			},
		},
	}

	ret := reflect.DeepEqual(actual.QueryMode.Arguments.outputFormat, expected.QueryMode.Arguments.outputFormat)
	assert.Equal(t, true, ret)

	actual = ParseArguments([]string{"ggql", "--output", "text"})
	expected.Error = "invalid output format"

	ret = reflect.DeepEqual(actual.Error, expected.Error)
	assert.Equal(t, true, ret)
}

func TestAnalysisArguments(t *testing.T) {
	t.Skip("Skipping TestAnalysisArguments.")
}

func TestHelpArguments(t *testing.T) {
	t.Skip("Skipping TestHelpArguments.")
}

func TestVersionArguments(t *testing.T) {
	actual := ParseArguments([]string{"ggql", "--version"})
	expected := Command{
		Version: Version,
	}

	ret := reflect.DeepEqual(actual.Version, expected.Version)
	assert.Equal(t, true, ret)
}
