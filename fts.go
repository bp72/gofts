package fts

import (
	"cmp"
	"math"
	"slices"
	"strings"

	stemmer "github.com/rjohnsondev/golibstemmer"
)

type SearchType int

const (
	SearchTypeAND SearchType = 0
	SearchTypeOR  SearchType = 1
)

type FullTextSearchIndex struct {
	Documents map[int]*DocumContainer
	Index     map[string]*RevIndex
	Stemmer   *stemmer.Stemmer
	StopWords map[string]bool
}

func NewFullTextSearchIndex() *FullTextSearchIndex {
	idx := &FullTextSearchIndex{}

	idx.Documents = make(map[int]*DocumContainer)
	idx.Index = make(map[string]*RevIndex)
	idx.StopWords = make(map[string]bool)

	for _, stopWord := range STOPWORDS {
		idx.AddStopWord(stopWord)
	}

	var err error
	idx.Stemmer, err = stemmer.NewStemmer("english")

	if err != nil {
		panic(err)
	}

	return idx
}

func (idx *FullTextSearchIndex) AddStopWord(word string) {
	idx.StopWords[word] = true
}

func (idx *FullTextSearchIndex) DocCount() int {
	return len(idx.Documents)
}

func (idx *FullTextSearchIndex) IndexDoc(dc *DocumContainer) {
	if _, exists := idx.Documents[dc.ID]; !exists {
		idx.Documents[dc.ID] = dc
	}

	Tokens := idx.GetTokens(dc.AsText())

	for _, token := range Tokens {
		if _, exists := idx.Index[token]; !exists {
			idx.Index[token] = NewRevIndex(token)
		}
		idx.Index[token].Add(dc.ID)
		dc.Terms[token]++
	}
}

func (idx *FullTextSearchIndex) Search(query string, searchType SearchType, rank bool) []*DocumContainer {
	res := make([]*DocumContainer, 0)

	if searchType != SearchTypeAND && searchType != SearchTypeOR {
		return res
	}

	tokens := idx.GetTokens(query)

	raw := make(map[int]int)
	for _, token := range tokens {
		if ri, exists := idx.Index[token]; exists {
			for doc := range ri.Index {
				raw[doc]++
			}
		}
	}

	for docId, occurance := range raw {
		if occurance < len(tokens) && searchType == SearchTypeAND {
			continue
		}
		if doc, exists := idx.Documents[docId]; exists {
			res = append(res, doc)
		}
	}

	if rank {
		idx.Rank(tokens, res)
	}

	return res
}

func (idx *FullTextSearchIndex) Rank(tokens []string, docs []*DocumContainer) {
	if len(docs) == 0 {
		return
	}

	for _, doc := range docs {
		for _, token := range tokens {
			// termFreq := float64(doc.TermFreq(token))
			// invDocFreq := idx.InvDocFreq(token)
			// fmt.Println(termFreq)
			// fmt.Println(invDocFreq)
			doc.Score += float64(doc.TermFreq(token)) * idx.InvDocFreq(token)
		}
	}

	slices.SortFunc(docs,
		func(a, b *DocumContainer) int {
			return cmp.Compare(b.Score, a.Score)
		})
}

func (idx *FullTextSearchIndex) DocFreq(token string) float64 {
	if ri, exists := idx.Index[token]; exists {
		return float64(ri.Freq())
	}
	return 0.0
}

func (idx *FullTextSearchIndex) InvDocFreq(token string) float64 {
	/*
		# Manning, Hinrich and Schütze use log10, so we do too, even though it
		# doesn't really matter which log we use anyway
		# https://nlp.stanford.edu/IR-book/html/htmledition/inverse-document-frequency-1.html
	*/

	return math.Log10(float64(idx.DocCount()) / idx.DocFreq(token))
}

func (fts *FullTextSearchIndex) GetTokens(text string) []string {
	words := strings.Fields(strings.ToLower(text))
	res := make([]string, 0)

	for _, word := range words {
		word := fts.Stemmer.StemWord(word)
		if _, exists := fts.StopWords[word]; exists {
			continue
		}
		res = append(res, word)
	}

	return res
}