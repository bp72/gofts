package fts

type SearchResultDocumentContainer struct {
	Doc   *DocumContainer
	Score float64
}

type SearchResult struct {
	Total     int
	Documents []*SearchResultDocumentContainer
	Tokens    []string
}
