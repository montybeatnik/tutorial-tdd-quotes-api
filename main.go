package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

// validate ensures certain quote conditions are met.
func (q *Quote) validate() error {
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

func NewInMemStore() *InMemStore {
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

// application holds app dependencies
type application struct {
	log   *log.Logger
	store Repo
}

// newApp spins up a new app, factoring in dependencies.
func newApp(log *log.Logger, store Repo) *application {
	return &application{store: store, log: log}
}

// handleQuotes performs validation and interacts with the quotes store.
func (app *application) handleQuotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var quote Quote
		defer r.Body.Close()
		// deserialize the request body
		if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"JSON body cannot be empty"}`))
			return
		}
		// run validations againgst the deserialized data
		if err := quote.validate(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := map[string]string{"error": err.Error()}
			bs, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(bs)
			return
		}
		// add the quote to the store.
		if err := app.store.Create(quote); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := map[string]string{"error": err.Error()}
			bs, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(bs)
			return
		}
		resp := map[string]string{"message": "succesfully created quote"}
		bs, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(bs)
	case http.MethodGet:
		// Get all logic
		if strings.Split(r.URL.Path, "/")[1] == "" {
			quotes, err := app.store.All()
			if err != nil {
				// deal with error
			}
			resp := map[string][]Quote{"quotes": quotes}
			bs, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(bs)
			return
		}
		idStr := strings.Split(r.URL.Path, "/")[1]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			resp := map[string]string{"error": "the id must be a positive integer"}
			bs, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(bs)
			return
		}
		quote, err := app.store.ByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			resp := map[string]string{"message": err.Error()}
			bs, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(bs)
			return
		}
		bs, err := json.Marshal(quote)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(bs)

	default:
		resp := map[string]string{"error": "allowed methods [POST, GET]"}
		bs, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write(bs)
	}
}

func main() {
	log := log.New(os.Stdout, "[quotes] ", log.Ldate|log.Ltime|log.Lshortfile)
	app := newApp(log, NewInMemStore())
	appAddr := "localhost:8000"
	app.log.Println("quotes app listening on ", appAddr)
	http.Handle("/", http.HandlerFunc(app.handleQuotes))
	app.log.Fatal(http.ListenAndServe(appAddr, nil))
}
