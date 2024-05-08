package engine

import (
	"os"

	"github.com/go-git/go-git/v5"
)

func TestNewRepo(path string) error {
	// Clone the given repository to the given path
	_, err := git.PlainInit(path, true)
	return err
}

func TestDeleteRepo(path string) error {
	err := os.RemoveAll(path)
	return err
}
