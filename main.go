package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Quote reprsents insightful text and
// it's origin.
type Quote struct {
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

// application holds app dependencies
type application struct {
	log   *log.Logger
	store map[int]Quote
}

// count serves as our PK/ID for quotes
// in the data store.
var count int

// newApp spins up a new app, factoring in dependencies.
func newApp(log *log.Logger) *application {
	store := make(map[int]Quote)
	count++
	store[count] = Quote{Author: "Gandhi", Message: "be the change!"}
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
		count++
		app.store[count] = quote
		resp := map[string]string{"message": "succesfully created quote"}
		bs, err := json.Marshal(resp)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write(bs)
	case http.MethodGet:
		if strings.Split(r.URL.Path, "/")[1] == "" {
			var quotes []Quote
			for _, qt := range app.store {
				quotes = append(quotes, qt)
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
		quote, found := app.store[id]
		if !found {
			resp := map[string]string{"message": "couldn't find quote matching that id"}
			bs, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNotFound)
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
	app := newApp(log)
	appAddr := "localhost:8000"
	app.log.Println("quotes app listening on ", appAddr)
	http.Handle("/", http.HandlerFunc(app.handleQuotes))
	app.log.Fatal(http.ListenAndServe(appAddr, nil))
}
