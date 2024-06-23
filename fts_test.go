package fts

import (
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
	fts := NewFullTextSearchIndex()
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
}

func TestIndexDoc(t *testing.T) {
	fts := NewFullTextSearchIndex()

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

func TestSearch(t *testing.T) {
	fts := NewFullTextSearchIndex()

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

	docums := fts.Search("apple tower", SearchTypeAND, true)

	if len(docums) != 2 {
		t.Fatalf("Test 'Search' failed. Expected size %d. Got: %d", 2, len(docums))
	}

	docums = fts.Search("apple tower", SearchTypeOR, true)

	if len(docums) != 3 {
		t.Fatalf("Test 'Search' failed. Expected size %d. Got: %d", 3, len(docums))
	}

	// for _, docum := range docums {
	// 	fmt.Printf("doc='%s' score=%f\n", docum.AsText(), docum.Score)
	// }
}
