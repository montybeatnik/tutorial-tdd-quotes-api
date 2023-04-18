package store

import (
	"errors"
	"sync"
)

var ErrQuoteNotFound = errors.New("couldn't find quote matching that id")

// Quote reprsents insightful text and
// it's origin.
type Quote struct {
	ID      int    `json:"id"`
	Author  string `json:"author"`
	Message string `json:"message"`
}

// Validate ensures certain quote conditions are met.
func (q *Quote) Validate() error {
	if q.Author == "" && q.Message == "" {
		return errors.New("please provide both an author and a message")
	}
	if q.Author == "" {
		return errors.New("please provide an author")
	}
	if q.Message == "" {
		return errors.New("please provide a message")
	}
	return nil
}

type Repo interface {
	Create(quote Quote) error
	All() ([]Quote, error)
	ByID(id int) (Quote, error)
}

type InMemStore struct {
	mu    sync.Mutex
	store map[int]Quote
}

func NewInMem() *InMemStore {
	store := make(map[int]Quote)
	count++
	store[count] = Quote{ID: count, Author: "Gandhi", Message: "be the change!"}
	return &InMemStore{store: store}
}

// count serves as our PK/ID for quotes
// in the data store.
var count int

func (ms *InMemStore) Create(quote Quote) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	count++
	quote.ID = count
	ms.store[quote.ID] = quote
	return nil
}

func (ms *InMemStore) All() ([]Quote, error) {
	var quotes []Quote
	for _, qt := range ms.store {
		quotes = append(quotes, Quote{ID: qt.ID, Author: qt.Author, Message: qt.Message})
	}
	return quotes, nil
}

func (m *InMemStore) ByID(id int) (Quote, error) {
	quote, found := m.store[id]
	if !found {
		return quote, ErrQuoteNotFound
	}
	return quote, nil
}
