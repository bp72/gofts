package fts

import (
	"testing"

	sh1 "github.com/mfonda/simhash"
	sh2 "github.com/safeie/simhash"
)

func TestDocumSimhash(t *testing.T) {
	params := FullTextSearchParams{
		MaxSearchResults:   150,
		ExcludeBySimhash:   true,
		MinSimhashDistance: 4,
	}
	idx := NewFullTextSearchIndex(params)

	type TestParams struct {
		TextA      string
		TextB      string
		Distance   int
		Similarity float64
		Compare    uint8
	}

	testParams := []TestParams{
		{"humpty dumpty sat on the wall", "humpty dumpty had a great falt", 8, 0.75, 15},
		{"humpty dumpty sat on the wall", "humpty sat wall on dumpty", 7, 0.78125, 10},
		{"humpty dumpty sat on the wall", "humpty dumpty sat wall", 4, 0.875, 9},
		{"humpty dumpty had a great falt", "humpty sat wall on dumpty", 13, 0.59375, 21},
		{"humpty dumpty had a great falt", "humpty dumpty sat wall", 6, 0.8125, 16},
		{"humpty sat wall on dumpty", "humpty dumpty sat wall", 11, 0.65625, 15},
	}

	texts := []string{
		"humpty dumpty sat on the wall",
		"humpty dumpty had a great falt",
		"humpty sat wall on dumpty",
		"humpty dumpty sat wall",
	}

	for pos, text := range texts {
		dc := NewDocumContainer(&TestDocument{ID: pos, Text: text})
		idx.IndexDoc(dc)
	}

	for pos, testParam := range testParams {
		SimhashA1 := sh2.Simhash(testParam.TextA)
		SimhashA2 := sh1.Simhash(sh1.NewWordFeatureSet([]byte(testParam.TextA)))
		SimhashB1 := sh2.Simhash(testParam.TextB)
		if testParam.Distance != sh2.Distance(SimhashA1, SimhashB1) {
			t.Errorf("test 'Test simhash Distance' TestNo %d. Expected: %d got %d", pos, testParam.Distance, sh2.Distance(SimhashA1, SimhashB1))
		}
		if testParam.Similarity != sh2.Similar(SimhashA1, SimhashB1) {
			t.Errorf("test 'Test simhash Similar' TestNo %d. Expected: %f got %f", pos, testParam.Similarity, sh2.Similar(SimhashA1, SimhashB1))
		}

		SimhashB2 := sh1.Simhash(sh1.NewWordFeatureSet([]byte(testParam.TextB)))

		if testParam.Compare != sh1.Compare(SimhashA2, SimhashB2) {
			t.Errorf("test 'Test simhash Compare' TestNo %d. Expected: %d got %d", pos, testParam.Compare, sh1.Compare(SimhashA2, SimhashB2))

		}
	}
}
