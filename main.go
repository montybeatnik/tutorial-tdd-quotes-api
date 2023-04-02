package main

import (
	"encoding/json"
	"errors"
	"net/http"
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
	store map[int]Quote
}

// count serves as our PK/ID for quotes
// in the data store.
var count int

// newApp spins up a new app, factoring in dependencies.
func newApp() *application {
	store := make(map[int]Quote)
	count++
	store[count] = Quote{Author: "Gahndi", Message: "be the change!"}
	return &application{store: store}
}

// handleQuotes performs validation and interacts with the quotes store.
func (app *application) handleQuotes(w http.ResponseWriter, r *http.Request) {
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
}

func main() {

}
