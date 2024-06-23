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
	}

	index := fts.NewFullTextSearchIndex()

	for _, doc := range docs {
		index.IndexDoc(fts.NewDocumContainer(doc))
	}

	serp := index.Search("fire rabbit", fts.SearchTypeOR, true)
	for pos, found := range serp {
		doc := found.Document.(*Doc)
		fmt.Printf("P=%d T='%s' A='%s' S=%f\n", pos, doc.Title, doc.Author, found.Score)
	}
}
