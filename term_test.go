package fts

import "testing"

func TestTermBasic(t *testing.T) {
	term := &Term{Text: "aaa", Pos: 0}

	expected := "aaa(P0)"
	if term.String() != expected {
		t.Errorf("test 'TestTermBasic' failed. Expected %s got %s", expected, term.String())
	}
}
