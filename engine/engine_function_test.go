package engine

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
)

func newFunctionRepo(path string) error {
	_, err := git.PlainInit(path, true)
	return err
}

func deleteFunctionRepo(path string) error {
	return os.RemoveAll(path)
}

func TestSelectGQLObjects(t *testing.T) {
}

func TestSelectReferences(t *testing.T) {
}

func TestSelectCommits(t *testing.T) {
}

func TestSelectBranches(t *testing.T) {
}

func TestSelectDiffs(t *testing.T) {
}

func TestSelectTags(t *testing.T) {
}

func TestSelectValues(t *testing.T) {
}

func TestGetColumnName(t *testing.T) {
}
