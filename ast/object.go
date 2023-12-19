package ast

import (
	"sync"
)

type GQLObject struct {
	Attributes map[string]Value
}

type Value interface{}

var mutex = &sync.Mutex{}

func FlatGQLGroups(groups *[][]GQLObject) error {
	var mainGroup []GQLObject

	for _, group := range *groups {
		mainGroup = append(mainGroup, group...)
	}

	mutex.Lock()
	*groups = (*groups)[:0] // clear the groups slice
	mutex.Unlock()

	*groups = append(*groups, mainGroup)

	return nil
}
