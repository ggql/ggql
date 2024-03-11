package ast

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGroupIsEmpty(t *testing.T) {
	group := Group{Rows: []Row{}}

	ret := group.IsEmpty()
	assert.Equal(t, true, ret)
}

func TestGroupLen(t *testing.T) {
	group := Group{Rows: []Row{}}

	ret := group.Len()
	assert.Equal(t, 0, ret)
}

func TestGitQLObjectFlat(t *testing.T) {
	groups := []Group{
		{
			Rows: []Row{
				{
					Values: []Value{},
				},
			},
		},
	}

	object := GitQLObject{
		Titles: []string{},
		Groups: groups,
	}

	object.Flat()
	assert.Equal(t, 1, len(object.Groups))

	for _, item := range object.Groups {
		assert.Equal(t, 1, item.Len())
	}

	object.Groups = []Group{}
	object.Groups = append(object.Groups, Group{Rows: []Row{{Values: []Value{}}}}, Group{Rows: []Row{{Values: []Value{}}}})
	assert.Equal(t, 2, len(object.Groups))

	object.Flat()
	assert.Equal(t, 1, len(object.Groups))

	for _, item := range object.Groups {
		assert.Equal(t, 2, len(item.Rows))
	}
}

func TestGitQLObjectIsEmpty(t *testing.T) {
	object := GitQLObject{
		Titles: []string{},
		Groups: []Group{},
	}

	ret := object.IsEmpty()
	assert.Equal(t, true, ret)
}

func TestGitQLObjectLen(t *testing.T) {
	object := GitQLObject{
		Titles: []string{},
		Groups: []Group{},
	}

	ret := object.Len()
	assert.Equal(t, 0, ret)

	object.Groups = append(object.Groups, Group{Rows: []Row{}})

	ret = object.Len()
	assert.Equal(t, 1, ret)
}

func TestGitQLObjectAsJSON(t *testing.T) {
	object := GitQLObject{
		Titles: []string{"title1", "title2"},
		Groups: []Group{
			{
				Rows: []Row{
					{
						Values: []Value{
							IntegerValue{1},
							IntegerValue{1},
						},
					},
				},
			},
			{
				Rows: []Row{
					{
						Values: []Value{
							TextValue{"hello"},
							TextValue{"world"},
						},
					},
				},
			},
		},
	}

	ret, err := object.AsJson()
	fmt.Println(ret)
	assert.Equal(t, nil, err)

	var result []map[string]string
	err = json.Unmarshal([]byte(ret), &result)
	assert.Equal(t, nil, err)
}

func TestGitQLObjectAsCSV(t *testing.T) {
	object := GitQLObject{
		Titles: []string{"title1", "title2"},
		Groups: []Group{
			{
				Rows: []Row{
					{
						Values: []Value{
							IntegerValue{1},
							IntegerValue{1},
						},
					},
				},
			},
			{
				Rows: []Row{
					{
						Values: []Value{
							TextValue{"hello"},
							TextValue{"hello"},
						},
					},
				},
			},
		},
	}

	ret, err := object.AsCsv()
	fmt.Println(ret)
	assert.Equal(t, nil, err)
}
