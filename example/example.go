package main

import (
	"fmt"

	fts "github.com/bp72/gofts"
)

type Doc struct {
	ID     int
	Title  string
	Author string
	Descr  string
}

func (td *Doc) GetID() int {
	return td.ID
}

func (td *Doc) AsText() string {
	return td.Title + " " + td.Descr + " " + td.Author
}

func main() {
	fmt.Println("example indexing and search")

	docs := []*Doc{
		{ID: 1, Title: "The Poky Little Puppy", Author: "Janette Sebring Lowrey", Descr: "Puppy is slower than other, bigger animals."},
		{ID: 2, Title: "The Tale of Peter Rabbit", Author: "Beatrix Potter", Descr: "Rabbit eats some vegetables."},
		{ID: 3, Title: "Tootle", Author: "Gertrude Crampton", Descr: "Little toy train has big dreams."},
		{ID: 4, Title: "Green Eggs and Ham", Author: "Dr. Seuss", Descr: "Sam has changing food preferences and eats unusually colored food."},
		{ID: 5, Title: "Harry Potter and the Goblet of Fire", Author: "J.K. Rowling", Descr: "Fourth year of school starts, big drama ensues."},
		{ID: 6, Title: "Pride and Prejudice", Author: "Jane Austen", Descr: ""},
		{ID: 7, Title: "Jane Eyre", Author: "Charlotte Brontë", Descr: ""},
		{ID: 8, Title: "The Adventures of Sherlock Holmes (Sherlock Holmes, #3)", Author: "Arthur Conan Doyle", Descr: ""},
		{ID: 9, Title: "The Importance of Being Earnest", Author: "Oscar Wilde", Descr: ""},
		{ID: 10, Title: "The Adventures of Tom Sawyer", Author: "Mark Twain", Descr: ""},
		{ID: 11, Title: "When God Laughs and Other Stories", Author: "Jack London", Descr: "This etext was prepared by Les Bowler, St. Ives, Dorset from the 1911 Mills and Boon edition. WHEN GOD LAUGHS, AND OTHER STORIES CONTENTS WHEN GOD LAUGHS THE APOSTATE A WICKED WOMAN JUST MEAT CREATED HE THEM THE CHINAGO MAKE WESTING SEMPER IDEM A NOSE FOR THE KING THE “FRANCIS SPAIGHT” A CURIOUS FRAGMENT A"},
		{ID: 12, Title: "Travels with a Donkey in the Cévennes", Author: "Robert Louis Stevenson", Descr: "Travels with a Donkey in the Cevennes by Robert Louis Stevenson. Scanned and proofed by David Price, ccx074@coventry.ac.uk Travels with a Donkey in the Cevennes My Dear Sidney Colvin, The journey which this little book is to describe was very agreeable and fortunate for me. After an uncouth beginning, I had the best of luck"},
		{ID: 13, Title: "The Master of Ballantrae", Author: "Robert Louis Stevenson", Descr: "The Master of Ballantrae by Robert Louis Stevenson Scanned and proofed by David Price ccx074@coventry.ac.uk The Master of Ballantrae A Winter’s Tale To Sir Percy Florence and Lady Shelley Here is a tale which extends over many years and travels into many countries. By a peculiar fitness of circumstance the writer began, continued it, and"},
		{ID: 14, Title: "Crime and Punishment", Author: "Fyodor Dostoevsky", Descr: "On an exceptionally hot evening early in July a young man came out of the garret in which he lodged in S. Place and walked slowly, as though in hesitation, towards K. bridge."},
		{ID: 15, Title: "The Adventures of Sherlock Holmes (Sherlock Holmes, #4)", Author: "Arthur Conan Doyle", Descr: ""},
	}

	index := fts.NewFullTextSearchIndex(fts.FullTextSearchParams{ExcludeBySimhash: true, MinSimhashDistance: 1})

	for _, doc := range docs {
		index.IndexDoc(fts.NewDocumContainer(doc))
	}

	serp := index.Search("fire rabbit adventure", fts.SearchTypeOR, true)
	for pos, found := range serp.Documents {
		doc := found.Doc.Document.(*Doc)
		fmt.Printf("P=%d T='%s' A='%s' S=%f H=%d\n", pos, doc.Title, doc.Author, found.Score, found.Doc.Simhash)
	}
}
