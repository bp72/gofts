package fts

import (
	"fmt"
	"reflect"
	"testing"
)

type TestDocument struct {
	ID   int
	Text string
}

func (td *TestDocument) GetID() int {
	return td.ID
}

func (td *TestDocument) AsText() string {
	return td.Text
}

func TestGetTokens(t *testing.T) {
	params := FullTextSearchParams{
		ExcludeBySimhash:   true,
		MaxSearchResults:   150,
		MinSimhashDistance: 4,
		UseStemming:        true,
	}
	fts := NewFullTextSearchIndex(params)
	tokens := fts.GetTokens("You and me toWer")
	expect := []string{"tower"}

	// Check toLower, stop-words
	if !reflect.DeepEqual(tokens, expect) {
		t.Fatalf("Test GetTokens failed. Expected %v. Got: %v", tokens, expect)
	}

	// Check toLower, stop-words and stemmer
	tokens = fts.GetTokens("CHECK STEMMER APPLES You and me toWers")
	expect = []string{"check", "stemmer", "appl", "tower"}

	if !reflect.DeepEqual(tokens, expect) {
		t.Fatalf("Test GetTokens failed. Expected %v. Got: %v", tokens, expect)
	}

	// Check AddStopWord and its effect
	fts.AddStopWord("tower")
	tokens = fts.GetTokens("CHECK STEMMER APPLES You and me toWers")
	expect = []string{"check", "stemmer", "appl"}

	if !reflect.DeepEqual(tokens, expect) {
		t.Fatalf("Test GetTokens failed. Expected %v. Got: %v", tokens, expect)
	}

	// Check GetTokens with Ngrams
	fts.FullTextSearchParams.UseNgrams = true
	tokens = fts.GetTokens("CHECK STEMMER APPLES You and me toWers")
	expect = []string{"check stemmer appl", "check stemmer", "stemmer appl", "check", "stemmer", "appl"}
	if !reflect.DeepEqual(tokens, expect) {
		t.Fatalf("Test GetTokens failed. Expected %v. Got: %v", tokens, expect)
	}

	tokens = fts.GetTokens("CHECK")
	expect = []string{"check"}
	if !reflect.DeepEqual(tokens, expect) {
		t.Fatalf("Test GetTokens failed. Expected %v. Got: %v", tokens, expect)
	}
	tokens = fts.GetTokens("CHECK cook")
	expect = []string{"check cook", "check", "cook"}
	if !reflect.DeepEqual(tokens, expect) {
		t.Fatalf("Test GetTokens failed. Expected %v. Got: %v", tokens, expect)
	}
}

func TestIndexDoc(t *testing.T) {
	params := FullTextSearchParams{
		ExcludeBySimhash:   true,
		MinSimhashDistance: 4,
		MaxSearchResults:   150,
		UseStemming:        true,
	}
	fts := NewFullTextSearchIndex(params)

	DocTexts := []string{
		"You and me toWer",
		"CHECK STEMMER APPLES You and me TOWERS",
	}

	for pos, text := range DocTexts {
		dc := NewDocumContainer(&TestDocument{ID: pos + 1, Text: text})
		fts.IndexDoc(dc)

		if len(fts.Documents) != dc.ID {
			t.Fatalf("Test IndexDoc failed. Expected size %d. Got: %d", dc.ID, fts.DocCount())
		}
	}

	Terms := map[string]int{
		"tower":   2,
		"check":   1,
		"stemmer": 1,
		"appl":    1,
	}

	for term, ri := range fts.Index {
		if expected, exists := Terms[term]; exists {
			if expected != ri.Freq() {
				t.Fatalf("Test 'RevIndex By Term %s Size' failed. Expected %d, got %d", term, expected, ri.Freq())
			}
		} else {
			t.Fatalf("Test 'RevIndex By Term %s Size' failed. Expected to exist.", term)
		}

	}
}

func TestIndexDocNgrams(t *testing.T) {
	params := FullTextSearchParams{
		ExcludeBySimhash:   true,
		MinSimhashDistance: 4,
		MaxSearchResults:   150,
		UseStemming:        true,
		UseNgrams:          true,
	}
	fts := NewFullTextSearchIndex(params)

	DocTexts := []string{
		"You and me toWer",
		"CHECK STEMMER APPLES You and me TOWERS",
	}

	for pos, text := range DocTexts {
		dc := NewDocumContainer(&TestDocument{ID: pos + 1, Text: text})
		fts.IndexDoc(dc)

		if len(fts.Documents) != dc.ID {
			t.Fatalf("Test IndexDoc failed. Expected size %d. Got: %d", dc.ID, fts.DocCount())
		}
	}

	Terms := map[string]int{
		"stemmer appl tower": 1,
		"check stemmer appl": 1,
		"check stemmer":      1,
		"stemmer appl":       1,
		"appl tower":         1,
		"tower":              2,
		"check":              1,
		"stemmer":            1,
		"appl":               1,
	}

	for term, ri := range fts.Index {
		if expected, exists := Terms[term]; exists {
			if expected != ri.Freq() {
				t.Fatalf("Test 'RevIndex By Term '%s' Size' failed. Expected %d, got %d", term, expected, ri.Freq())
			}
		} else {
			t.Fatalf("Test 'RevIndex By Term '%s' Size' failed. Expected to exist.", term)
		}

	}
}

func TestSearchResult(t *testing.T) {
	params := FullTextSearchParams{
		ExcludeBySimhash:   true,
		MinSimhashDistance: 4,
		MaxSearchResults:   150,
		UseStemming:        true,
	}
	fts := NewFullTextSearchIndex(params)

	DocTexts := []string{
		"You and me toWer",
		"CHECK STEMMER APPLES You and me TOWERS",
		"APPLE builds the TOWER and phones tower tower",
		"Empty doc",
	}

	for pos, text := range DocTexts {
		dc := NewDocumContainer(&TestDocument{ID: pos + 1, Text: text})
		fts.IndexDoc(dc)

		if len(fts.Documents) != dc.ID {
			t.Fatalf("Test IndexDoc failed. Expected size %d. Got: %d", dc.ID, fts.DocCount())
		}
	}

	expectedTokens := []string{"appl", "tower"}
	searchResult := fts.Search("apPle the toWer", SearchTypeAND, true)

	if len(searchResult.Tokens) != len(expectedTokens) {
		t.Fatalf("Test 'SearchResult' failed. Tokens expected size %d. Got: %d", len(expectedTokens), len(searchResult.Tokens))
	}

	for i := 0; i < len(expectedTokens); i++ {
		if searchResult.Tokens[i] != expectedTokens[i] {
			t.Fatalf("Test 'SearchResult' failed. Token expected %s. Got: %s", expectedTokens[i], searchResult.Tokens[i])
		}
	}

	if len(searchResult.Documents) != 2 {
		t.Fatalf("Test 'Search' failed. Expected size %d. Got: %d", 2, len(searchResult.Documents))
	}

	searchResult = fts.Search("apple tower", SearchTypeOR, false)

	if len(searchResult.Documents) != 3 {
		t.Fatalf("Test 'Search' failed. Expected size %d. Got: %d", 3, len(searchResult.Documents))
	}

	// for _, docum := range searchResult.Documents {
	// 	fmt.Printf("doc='%s' score=%f\n", docum.Doc.AsText(), docum.Score)
	// }
}

func TestSearchResultNgram(t *testing.T) {
	params := FullTextSearchParams{
		ExcludeBySimhash:   true,
		MinSimhashDistance: 4,
		MaxSearchResults:   150,
		UseStemming:        true,
		UseNgrams:          true,
	}
	fts := NewFullTextSearchIndex(params)

	DocTexts := []string{
		"You and me toWer",
		"CHECK STEMMER APPLES You and me TOWERS",
		"APPLE builds the TOWER and phones tower tower",
		"Empty doc",
	}

	for pos, text := range DocTexts {
		dc := NewDocumContainer(&TestDocument{ID: pos + 1, Text: text})
		fts.IndexDoc(dc)

		if len(fts.Documents) != dc.ID {
			t.Fatalf("Test IndexDoc failed. Expected size %d. Got: %d", dc.ID, fts.DocCount())
		}
	}

	expectedTokens := []string{"appl tower", "appl", "tower"}
	searchResult := fts.Search("apPle the toWer", SearchTypeAND, true)

	if len(searchResult.Tokens) != len(expectedTokens) {
		t.Fatalf("Test 'SearchResult' failed. Tokens expected size %d. Got: %d", len(expectedTokens), len(searchResult.Tokens))
	}

	for i := 0; i < len(expectedTokens); i++ {
		if searchResult.Tokens[i] != expectedTokens[i] {
			t.Fatalf("Test 'SearchResult' failed. Token expected %s. Got: %s", expectedTokens[i], searchResult.Tokens[i])
		}
	}

	expectedDocCount := 1
	if len(searchResult.Documents) != expectedDocCount {
		t.Fatalf("Test 'Search' failed. Expected size %d. Got: %d", expectedDocCount, len(searchResult.Documents))
	}

	searchResult = fts.Search("apple tower", SearchTypeOR, false)

	if len(searchResult.Documents) != 3 {
		t.Fatalf("Test 'Search' failed. Expected size %d. Got: %d", 3, len(searchResult.Documents))
	}

	// for _, docum := range searchResult.Documents {
	// 	fmt.Printf("doc='%s' score=%f\n", docum.Doc.AsText(), docum.Score)
	// }
}

func TestSearchResultNgramStandalone(t *testing.T) {
	params := FullTextSearchParams{
		ExcludeBySimhash:   true,
		MinSimhashDistance: 4,
		MaxSearchResults:   150,
		UseStemming:        true,
		UseNgrams:          true,
	}
	fts := NewFullTextSearchIndex(params)

	DocTexts := []string{
		"You and me toWer",
		"CHECK STEMMER APPLES You and me TOWERS",
		"APPLE builds the TOWER and phones tower tower",
		"Empty doc",
	}

	for pos, text := range DocTexts {
		dc := NewDocumContainer(&TestDocument{ID: pos + 1, Text: text})
		fts.IndexDoc(dc)

		if len(fts.Documents) != dc.ID {
			t.Fatalf("Test IndexDoc failed. Expected size %d. Got: %d", dc.ID, fts.DocCount())
		}
	}

	expectedTokens := []string{"appl tower"}
	searchResult := fts.SearchNgram("apPle the toWer")

	if len(searchResult.Tokens) != len(expectedTokens) {
		t.Fatalf("Test 'SearchResult' failed. Tokens expected size %d. Got: %d", len(expectedTokens), len(searchResult.Tokens))
	}

	for i := 0; i < len(expectedTokens); i++ {
		if searchResult.Tokens[i] != expectedTokens[i] {
			t.Fatalf("Test 'SearchResult' failed. Token expected %s. Got: %s", expectedTokens[i], searchResult.Tokens[i])
		}
	}

	expectedDocCount := 1
	if len(searchResult.Documents) != expectedDocCount {
		t.Fatalf("Test 'Search' failed. Expected size %d. Got: %d", expectedDocCount, len(searchResult.Documents))
	}

}

func TestSearchQuery(t *testing.T) {
	params := FullTextSearchParams{
		ExcludeBySimhash:   true,
		MinSimhashDistance: 4,
		MaxSearchResults:   150,
		UseStemming:        true,
		UseNgrams:          true,
	}
	fts := NewFullTextSearchIndex(params)

	testDocs := []*TestDocument{
		{ID: 1, Text: "stopword-1 aaa stopword-2 bbb stopword-3 ccc stopword-4 ddd stopword-5"},
		{ID: 2, Text: "stopword-1 aaa bbb stopword-3 ccc stopword-4 ddd stopword-5"},
		{ID: 3, Text: "stopword-1 aaa bbb ccc stopword-4 ddd stopword-5"},
		{ID: 4, Text: "stopword-1 aaa stopword-4 bbb ccc ddd stopword-5"},
		{ID: 5, Text: "stopword-1 aaa bbb stopword-4 ccc ddd stopword-5"},
		{ID: 6, Text: "stopword-1 aaa bbb ccc ddd stopword-5"},
	}

	type TestCase struct {
		Query        string
		SearchResult []*TestDocument
	}

	testCases := []TestCase{
		{Query: "aaa stopword-2 bbb", SearchResult: []*TestDocument{}},
		{Query: "aaa bbb stopword-2", SearchResult: []*TestDocument{testDocs[1], testDocs[2], testDocs[4], testDocs[5]}},
		{Query: "aaa bbb", SearchResult: []*TestDocument{testDocs[1], testDocs[2], testDocs[4], testDocs[5]}},
		{Query: "stopword-1 aaa bbb stopword-1", SearchResult: []*TestDocument{testDocs[1], testDocs[2], testDocs[4], testDocs[5]}},
	}

	qp := NewQueryParser("english", true)

	for stopWord, _ := range fts.StopWords {
		qp.stopWords[stopWord] = true
	}

	for i := 0; i < 100; i++ {
		fts.AddStopWord(fmt.Sprintf("stopword-%d", i))
		qp.stopWords[fmt.Sprintf("stopword-%d", i)] = true
	}

	for _, testDoc := range testDocs {
		dc := NewDocumContainer(testDoc)
		// fts.IndexDoc(dc)
		fts.IndexDocQueryParser(dc, qp)

		// if len(fts.Documents) != dc.ID {
		// 	t.Fatalf("Test IndexDoc failed. Expected size %d. Got: %d", dc.ID, fts.DocCount())
		// }
	}

	for testNo, testCase := range testCases {
		q := qp.ParseQuery(testCase.Query)
		searchResult := fts.SearchQuery(q)

		if len(searchResult.Documents) != len(testCase.SearchResult) {
			fmt.Println(q)
			t.Fatalf("Test 'TestSearchQuery' failed %d. Token expected %d. Got: %d", testNo, len(testCase.SearchResult), len(searchResult.Documents))
		}

		got := make(map[*TestDocument]bool)
		expected := make(map[*TestDocument]bool)
		for pos, docContainer := range searchResult.Documents {
			testDoc := docContainer.Doc.Document.(*TestDocument)
			got[testDoc] = true
			expected[testCase.SearchResult[pos]] = true
		}

		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("Test 'TestSearchQuery' failed %d. Expected %v got: %v", testNo, got, testCase.SearchResult)
		}
	}
}
