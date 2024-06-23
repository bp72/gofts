package fts

/* Обратный индекс - хранит ИД документов для конкретного токена/слова */
type RevIndex struct {
	Token string
	Index map[int]bool
}

func NewRevIndex(token string) *RevIndex {
	ri := &RevIndex{Token: token}

	ri.Index = make(map[int]bool)

	return ri
}

func (ri *RevIndex) Freq() int {
	return len(ri.Index)
}

func (ri *RevIndex) Add(id int) {
	ri.Index[id] = true
}
