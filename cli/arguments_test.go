package cli

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyArguments(t *testing.T) {
	wd, _ := os.Getwd()

	actual := ParseArguments([]string{"ggql"})
	expected := Command{
		ReplMode: Arguments{
			repos:        []string{wd},
			analysis:     false,
			pagination:   false,
			pageSize:     10,
			outputFormat: 0,
		},
	}

	ret := reflect.DeepEqual(actual.ReplMode, expected.ReplMode)
	assert.Equal(t, true, ret)
}

func TestReposArguments(t *testing.T) {
	actual := ParseArguments([]string{"ggql", "--repos", "."})
	expected := Command{
		ReplMode: Arguments{
			repos:        []string{"."},
			analysis:     false,
			pagination:   false,
			pageSize:     10,
			outputFormat: 0,
		},
	}

	ret := reflect.DeepEqual(actual.ReplMode.repos, expected.ReplMode.repos)
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
		ReplMode: Arguments{
			repos:        []string{},
			analysis:     false,
			pagination:   false,
			pageSize:     10,
			outputFormat: 0,
		},
	}

	ret := reflect.DeepEqual(actual.QueryMode.Arguments.pageSize, expected.QueryMode.Arguments.pageSize)
	assert.Equal(t, true, ret)

	actual = ParseArguments([]string{"ggql", "--pagesize", "-"})
	expected.Error = "Invalid page size"

	ret = reflect.DeepEqual(actual.Error, expected.Error)
	assert.Equal(t, true, ret)
}

func TestOutputArguments(t *testing.T) {
	actual := ParseArguments([]string{"ggql", "--output", "csv"})
	expected := Command{
		ReplMode: Arguments{
			repos:        []string{},
			analysis:     false,
			pagination:   false,
			pageSize:     10,
			outputFormat: CSV,
		},
	}

	ret := reflect.DeepEqual(actual.ReplMode.outputFormat, expected.ReplMode.outputFormat)
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
