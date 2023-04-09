package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleQuotes(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		route          string
		body           []byte
		expectedStatus int
		expectedBody   string
	}{
		// POST TESTS
		{
			name:           "post no body",
			method:         http.MethodPost,
			body:           nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"JSON body cannot be empty"}`,
		},
		{
			name:           "post no author",
			method:         http.MethodPost,
			body:           []byte(`{"message":"excellent!"}`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"please provide an author"}`,
		},
		{
			name:           "post no message",
			method:         http.MethodPost,
			body:           []byte(`{"author":"ted"}`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"please provide a message"}`,
		},
		{
			name:           "post no values",
			method:         http.MethodPost,
			body:           []byte(`{"author":"","message":""}`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"please provide both an author and a message"}`,
		},
		{
			name:           "post success",
			method:         http.MethodPost,
			body:           []byte(`{"author":"bill","message":"excellent!"}`),
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message":"succesfully created quote"}`,
		},
		// GET TESTS
		{
			name:           "get invalid id",
			method:         http.MethodGet,
			route:          "/one",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"the id must be a positive integer"}`,
		},
		{
			name:           "get non-existant id",
			method:         http.MethodGet,
			route:          "/42",
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"message":"couldn't find quote matching that id"}`,
		},
		{
			name:           "get no route",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"quotes":[{"id":1,"author":"Gandhi","message":"be the change!"},{"id":2,"author":"bill","message":"excellent!"}]}`,
		},
		{
			name:           "get success",
			method:         http.MethodGet,
			route:          "/1",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"id":1,"author":"Gandhi","message":"be the change!"}`,
		},
	}
	// stand up an instance of our app
	log := log.New(io.Discard, "", 0)
	app := newApp(log, NewInMemStore())
	// grab an http server from the testing package
	ts := httptest.NewServer(http.HandlerFunc(app.handleQuotes))
	// build a request
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("%v%v", ts.URL, tc.route)
			req, err := http.NewRequest(tc.method, url, bytes.NewReader(tc.body))
			if err != nil {
				t.Errorf("couldn't build request: %v", err)
			}
			// send the request to the server
			client := ts.Client()
			resp, err := client.Do(req)
			if err != nil {
				t.Errorf("request failed: %v", err)
			}
			defer resp.Body.Close()
			// test the status code
			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("got %v; expected %v", resp.StatusCode, tc.expectedStatus)
			}
			// test the response body
			bs, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Errorf("reading resp body failed: %v", err)
			}
			body := string(bs)
			if body != tc.expectedBody {
				t.Errorf("got %v; expected %v", body, tc.expectedBody)
			}
		})
	}
}

func BenchmarkCreateQuote(b *testing.B) {
	quoteStore := NewInMemStore()
	for n := 0; n < b.N; n++ {
		_ = quoteStore.Create(Quote{Author: "thedude", Message: fmt.Sprintf("something-%d", n)})
	}
}
