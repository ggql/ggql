package cli

import (
	"reflect"
	"testing"
)

func TestParseArguments(t *testing.T) {
	arguments1 := []string{"ggql", "--help"}
	got1 := ParseArguments(arguments1)
	var want1 Command
	want1.Help = true
	if !reflect.DeepEqual(got1, want1) {
		t.Errorf("want: %v  got: %v", want1, got1)
	}
	if got1.Help {
		PrintHelpList()
	}

	arguments2 := []string{"ggql", "--version"}
	got2 := ParseArguments(arguments2)
	var want2 Command
	want2.Version = "1.0.0"
	if !reflect.DeepEqual(got2, want2) {
		t.Errorf("want: %v  got: %v", want2, got2)
	}

	arguments3 := []string{"ggql", "--repos", "repo1", "repo2"}
	got3 := ParseArguments(arguments3)
	var want3 Command
	if !reflect.DeepEqual(got3.QueryMode.Arguments.Repos, want3.QueryMode.Arguments.Repos) {
		t.Errorf("want: %v  got: %v", want3, got3)
	}

	arguments4 := []string{"ggql", "--query", "commitid"}
	got4 := ParseArguments(arguments4)
	var want4 Command
	want4.QueryMode.Query = "commitid"
	if !reflect.DeepEqual(got4.QueryMode.Query, want4.QueryMode.Query) {
		t.Errorf("want: %v  got: %v", want4, got4)
	}

	arguments5 := []string{"ggql", "--analysis"}
	got5 := ParseArguments(arguments5)
	var want5 Command
	want5.ReplMode.Analysis = true
	if !reflect.DeepEqual(got5.ReplMode.Analysis, want5.ReplMode.Analysis) {
		t.Errorf("want: %v  got: %v", want5, got5)
	}

	arguments6 := []string{"ggql", "--pagination"}
	got6 := ParseArguments(arguments6)
	var want6 Command
	want6.ReplMode.Pagination = true
	if !reflect.DeepEqual(got6.ReplMode.Pagination, want6.ReplMode.Pagination) {
		t.Errorf("want: %v  got: %v", want6, got6)
	}

	arguments7 := []string{"ggql", "--pagesize", "5"}
	got7 := ParseArguments(arguments7)
	var want7 Command
	want7.ReplMode.PageSize = 5
	if !reflect.DeepEqual(got7.ReplMode.PageSize, want7.ReplMode.PageSize) {
		t.Errorf("want: %v  got: %v", want7, got7)
	}
}
