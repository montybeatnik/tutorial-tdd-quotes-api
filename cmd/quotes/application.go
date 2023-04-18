package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/montybeatnik/tutorials/tdd-quotes-api/store"
)

// application holds app dependencies
type application struct {
	log   *log.Logger
	store store.Repo
}

// newApp spins up a new app, factoring in dependencies.
func newApp(log *log.Logger, store store.Repo) *application {
	return &application{store: store, log: log}
}

// handleQuotes performs validation and interacts with the quotes store.
func (app *application) handleQuotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var quote store.Quote
		defer r.Body.Close()
		// deserialize the request body
		if err := json.NewDecoder(r.Body).Decode(&quote); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error":"JSON body cannot be empty"}`))
			return
		}
		// run validations againgst the deserialized data
		if err := quote.Validate(); err != nil {
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
			resp := map[string][]store.Quote{"quotes": quotes}
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
