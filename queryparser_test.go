package fts

import (
	"fmt"
	"reflect"
	"testing"
)

func TestQueryParser(t *testing.T) {
	qp := NewQueryParser("english", true)
	q := qp.ParseQuery("word1    Word2       WORD3 wOrD4 WorD5")

	expectedNgrams := []string{
		"word1 word2 word3",
		"word2 word3 word4",
		"word3 word4 word5",
		"word1 word2",
		"word2 word3",
		"word3 word4",
		"word4 word5",
	}

	if !reflect.DeepEqual(q.Ngrams, expectedNgrams) {
		t.Fatalf("Test 'TestQueryParser' Ngrams failed. Expected %v. Got: %v", expectedNgrams, q.Ngrams)
	}

	expectedTerms := []Term{
		{Text: "word1", Pos: 0},
		{Text: "word2", Pos: 1},
		{Text: "word3", Pos: 2},
		{Text: "word4", Pos: 3},
		{Text: "word5", Pos: 4},
	}

	if !reflect.DeepEqual(q.Terms, expectedTerms) {
		t.Fatalf("Test 'TestQueryParser' Terms failed. Expected %v. Got: %v", expectedTerms, q.Terms)
	}
}

func TestQueryParserWithStopWords(t *testing.T) {
	qp := NewQueryParser("english", true)

	qp.stopWords["stop-word-1"] = true
	qp.stopWords["stop-word-2"] = true
	qp.stopWords["stop-word-3"] = true

	q := qp.ParseQuery("word1     stop-word-1 Word2       WORD3  stop-word-2 wOrD4 WorD5 stop-word-3")

	expectedNgrams := []string{
		"word2 word3",
		"word4 word5",
	}

	if !reflect.DeepEqual(q.Ngrams, expectedNgrams) {
		t.Fatalf("Test 'TestQueryParserWithStopWords' Ngrams failed. Expected %v. Got: %v", expectedNgrams, q.Ngrams)
	}

	expectedTerms := []Term{
		{Text: "word1", Pos: 0},
		{Text: "word2", Pos: 2},
		{Text: "word3", Pos: 3},
		{Text: "word4", Pos: 5},
		{Text: "word5", Pos: 6},
	}

	if !reflect.DeepEqual(q.Terms, expectedTerms) {
		t.Fatalf("Test 'TestQueryParser' Terms failed. Expected %v. Got: %v", expectedTerms, q.Terms)
	}

	type TestCase struct {
		Text   string
		Terms  []Term
		Ngrams []string
	}

	testCases := []TestCase{
		{
			Text: "stopword-1 aaa stopword-2 bbb stopword-3 ccc stopword-4 ddd stopword-5",
			Terms: []Term{
				{Text: "aaa", Pos: 1},
				{Text: "bbb", Pos: 3},
				{Text: "ccc", Pos: 5},
				{Text: "ddd", Pos: 7},
			},
			Ngrams: []string{},
		},
		{
			Text: "stopword-1 aaa bbb stopword-3 ccc stopword-4 ddd stopword-5",
			Terms: []Term{
				{Text: "aaa", Pos: 1},
				{Text: "bbb", Pos: 2},
				{Text: "ccc", Pos: 4},
				{Text: "ddd", Pos: 6},
			},
			Ngrams: []string{
				"aaa bbb",
			},
		},
		{
			Text: "stopword-1 aaa bbb ccc stopword-4 ddd stopword-5",
			Terms: []Term{
				{Text: "aaa", Pos: 1},
				{Text: "bbb", Pos: 2},
				{Text: "ccc", Pos: 3},
				{Text: "ddd", Pos: 5},
			},
			Ngrams: []string{
				"aaa bbb ccc",
				"aaa bbb",
				"bbb ccc",
			},
		},
		{
			Text: "stopword-1 aaa stopword-4 bbb ccc ddd stopword-5",
			Terms: []Term{
				{Text: "aaa", Pos: 1},
				{Text: "bbb", Pos: 3},
				{Text: "ccc", Pos: 4},
				{Text: "ddd", Pos: 5},
			},
			Ngrams: []string{
				"bbb ccc ddd",
				"bbb ccc",
				"ccc ddd",
			},
		},
		{
			Text: "stopword-1 aaa bbb stopword-4 ccc ddd stopword-5",
			Terms: []Term{
				{Text: "aaa", Pos: 1},
				{Text: "bbb", Pos: 2},
				{Text: "ccc", Pos: 4},
				{Text: "ddd", Pos: 5},
			},
			Ngrams: []string{
				"aaa bbb",
				"ccc ddd",
			},
		},
		{
			Text: "stopword-1 aaa bbb ccc ddd stopword-5",
			Terms: []Term{
				{Text: "aaa", Pos: 1},
				{Text: "bbb", Pos: 2},
				{Text: "ccc", Pos: 3},
				{Text: "ddd", Pos: 4},
			},
			Ngrams: []string{
				"aaa bbb ccc",
				"bbb ccc ddd",
				"aaa bbb",
				"bbb ccc",
				"ccc ddd",
				"ccc ddd",
			},
		},
	}

	for i := 0; i < 100; i++ {
		qp.stopWords[fmt.Sprintf("stopword-%d", i)] = true
	}

	for testNo, testCase := range testCases {
		q := qp.ParseQuery(testCase.Text)
		if len(q.Terms) != len(testCase.Terms) {
			t.Fatalf("Test 'TestQueryParserWithStopWords' testCase %d term check failed. Expected size %d. Got: %d", testNo, len(testCase.Terms), len(q.Terms))
		}
		if !reflect.DeepEqual(q.Terms, testCase.Terms) {
			t.Fatalf("Test 'TestQueryParserWithStopWords' testCase %d term check failed. Expected size %d. Got: %d", testNo, len(testCase.Terms), len(q.Terms))
		}

	}

}
