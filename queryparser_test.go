package fts

import "testing"

func TestQueryParsrw(t *testing.T) {
	params := FullTextSearchParams{
		ExcludeBySimhash:   true,
		MaxSearchResults:   150,
		MinSimhashDistance: 4,
		UseStemming:        true,
		UseNgrams:          true,
	}
	idx := NewFullTextSearchIndex(params)

	qp := NewQueryParser(idx)
	qp.ParseQuery("word1    Word2       WORD3 wOrD4 WorD5")

}
