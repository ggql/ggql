package ast

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
