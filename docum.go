package fts

type Document interface {
	AsText() string
	GetID() int
}

type DocumContainer struct {
	ID       int
	Terms    map[string]int
	Score    float64
	Document Document
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
	dc.Document = docum
	dc.ID = docum.GetID()

	return dc
}
