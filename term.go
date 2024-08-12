package fts

import "fmt"

type Term struct {
	Text string
	Pos  int
}

func (t *Term) String() string {
	return fmt.Sprintf("%s(P%d)", t.Text, t.Pos)
}
