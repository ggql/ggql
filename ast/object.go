package ast

import (
	"encoding/csv"
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
)

type Row struct {
	Values []Value
}

type Group struct {
	Rows []Row
}

func (g *Group) IsEmpty() bool {
	return len(g.Rows) == 0
}

func (g *Group) Len() int {
	return len(g.Rows)
}

type GitQLObject struct {
	Titles []string
	Groups []Group
}

func (g *GitQLObject) Flat() {
	var rows []Row

	for i := range g.Groups {
		rows = append(rows, g.Groups[i].Rows...)
		g.Groups[i].Rows = nil
	}

	g.Groups = g.Groups[:0]
	g.Groups = append(g.Groups, Group{Rows: rows})
}

func (g *GitQLObject) IsEmpty() bool {
	return len(g.Groups) == 0
}

func (g *GitQLObject) Len() int {
	return len(g.Groups)
}

func (g *GitQLObject) AsJson() (string, error) {
	var elements []map[string]interface{}

	if len(g.Groups) == 0 {
		return "", errors.New("invalid groups")
	}

	for _, row := range g.Groups[0].Rows {
		object := make(map[string]interface{})
		for i, value := range row.Values {
			object[g.Titles[i]] = value.AsText()
		}
		elements = append(elements, object)
	}

	bytes, err := json.Marshal(elements)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal json")
	}

	return string(bytes), nil
}

func (g *GitQLObject) AsCsv() (string, error) {
	b := &strings.Builder{}
	writer := csv.NewWriter(b)

	if err := writer.Write(g.Titles); err != nil {
		return "", errors.Wrap(err, "failed to write")
	}

	if len(g.Groups) == 0 {
		return "", errors.New("invalid groups")
	}

	for _, row := range g.Groups[0].Rows {
		var valuesRow []string
		for _, value := range row.Values {
			valuesRow = append(valuesRow, value.AsText())
		}
		if err := writer.Write(valuesRow); err != nil {
			return "", errors.Wrap(err, "failed to write")
		}
	}

	writer.Flush()

	return b.String(), writer.Error()
}
