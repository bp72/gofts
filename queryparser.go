package fts

import (
	"fmt"
	"strings"

	stemmer "github.com/rjohnsondev/golibstemmer"
)

// type QueryParser struct {
// 	idx  *FullTextSearchIndex
// 	Type SearchType
// }

// func NewQueryParser(idx *FullTextSearchIndex) *QueryParser {
// 	return &QueryParser{Type: SearchTypeAND, idx: idx}
// }

// func (q *QueryParser) ParseQuery(query string) Query {
// 	tokens := q.idx.GetTokens(query)

// 	// for _, token := range tokens {
// 	// 	fmt.Println(token)
// 	// }

// 	return Query{Tokens: tokens}
// }

type QueryParser struct {
	stem      *stemmer.Stemmer
	useStem   bool
	stopWords map[string]bool
}

func NewQueryParser(lang string, useStem bool) *QueryParser {
	stem, err := stemmer.NewStemmer(lang)
	if err != nil {
		panic(err)
	}
	qp := &QueryParser{
		stem:      stem,
		useStem:   useStem,
		stopWords: make(map[string]bool),
	}

	return qp
}

func (qp *QueryParser) String() string {
	return "fu"
}

func (qp *QueryParser) ParseQuery(query string) Query {
	words := strings.Fields(strings.ToLower(query))

	q := Query{Text: query}

	for pos, word := range words {
		if qp.useStem {
			word = qp.stem.StemWord(word)
		}
		if _, exists := qp.stopWords[word]; exists {
			continue
		}
		q.Terms = append(q.Terms, Term{Text: word, Pos: pos})
	}

	for i := 0; i < len(q.Terms)-2; i++ {
		a := q.Terms[i]
		b := q.Terms[i+1]
		c := q.Terms[i+2]
		if a.Pos+1 == b.Pos && b.Pos+1 == c.Pos {
			q.Ngrams = append(q.Ngrams, fmt.Sprintf("%s %s %s", a.Text, b.Text, c.Text))
		}
	}

	for i := 0; i < len(q.Terms)-1; i++ {
		a := q.Terms[i]
		b := q.Terms[i+1]
		if a.Pos+1 == b.Pos {
			q.Ngrams = append(q.Ngrams, fmt.Sprintf("%s %s", a.Text, b.Text))
		}
	}

	return q
}
