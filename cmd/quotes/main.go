package main

import (
	"log"
	"net/http"
	"os"

	"github.com/montybeatnik/tutorials/tdd-quotes-api/store"
)

func main() {
	log := log.New(os.Stdout, "[quotes] ", log.Ldate|log.Ltime|log.Lshortfile)
	app := newApp(log, store.NewInMem())
	appAddr := "localhost:8000"
	app.log.Println("quotes app listening on ", appAddr)
	http.Handle("/", http.HandlerFunc(app.handleQuotes))
	app.log.Fatal(http.ListenAndServe(appAddr, nil))
}
