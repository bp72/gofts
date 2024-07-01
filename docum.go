package fts

import (
	"github.com/mfonda/simhash"
)

type Document interface {
	AsText() string
	GetID() int
}

type DocumContainer struct {
	ID       int
	Terms    map[string]int
	Tokens   []string
	Score    float64
	Document Document
	Simhash  uint64
}

func (dc *DocumContainer) AsText() string {
	return dc.Document.AsText()
}

func (dc *DocumContainer) TermFreq(term string) int {
	if freq, exists := dc.Terms[term]; exists {
		return freq
	}
	return 0
}

func NewDocumContainer(docum Document) *DocumContainer {
	dc := &DocumContainer{}

	dc.Terms = make(map[string]int)
	dc.Tokens = make([]string, 0)
	dc.Document = docum
	dc.ID = docum.GetID()
	dc.Score = 0.0

	return dc
}

func (dc *DocumContainer) AddTerm(term string) {
	dc.Terms[term]++
	dc.Tokens = append(dc.Tokens, term)
}

func (dc *DocumContainer) SetSimhash() {
	tokens := make([][]byte, 0)

	for _, token := range dc.Tokens {
		tokens = append(tokens, []byte(token))
	}

	shingles := simhash.Shingle(2, tokens)
	dc.Simhash = simhash.SimhashBytes(shingles)

	fs := simhash.NewWordFeatureSet([]byte(dc.Document.AsText()))
	dc.Simhash = simhash.Simhash(fs)
}
