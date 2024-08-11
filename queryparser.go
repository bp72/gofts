package fts

type QueryParser struct {
	idx  *FullTextSearchIndex
	Type SearchType
}

func NewQueryParser(idx *FullTextSearchIndex) *QueryParser {
	return &QueryParser{Type: SearchTypeAND, idx: idx}
}

func (q *QueryParser) ParseQuery(query string) Query {
	tokens := q.idx.GetTokens(query)

	// for _, token := range tokens {
	// 	fmt.Println(token)
	// }

	return Query{Tokens: tokens}
}
