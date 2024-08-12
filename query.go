package fts

import "strings"

type Query struct {
	Text   string
	Terms  []Term
	Ngrams []string
}

func (q *Query) String() string {
	var sb strings.Builder

	for i, term := range q.Terms {
		sb.WriteString(term.String())
		if i < len(q.Terms)-1 {
			sb.WriteRune('|')
		}
	}

	if len(q.Ngrams) > 0 {
		sb.WriteRune('|')
	}

	for i, ngram := range q.Ngrams {
		sb.WriteString(ngram)
		if i < len(ngram)-1 {
			sb.WriteRune('|')
		}
	}

	return sb.String()
}
