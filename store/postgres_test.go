package store

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func getDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := os.Getenv("QUOTES_TEST_DSN")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		t.Fatal(errors.New("couldn't connect to database"))
	}
	return db
}

func createTable(t *testing.T, db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS quotes(id SERIAL PRIMARY KEY, author TEXT UNIQUE NOT NULL, message TEXT UNIQUE NOT NULL);`
	stmt, err := db.Prepare(query)
	if err != nil {
		t.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		t.Fatal(err)
	}
}

func seed(t *testing.T, db *sql.DB) {
	query := `INSERT INTO quotes (author, message) VALUES ('ghandi', 'be the change');`
	stmt, err := db.Prepare(query)
	if err != nil {
		t.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		t.Error(err)
	}
}

func cleanup(t *testing.T, db *sql.DB) {
	t.Cleanup(func() {
		db.Exec(`DROP TABLE quotes;`)
	})
}

func TestCreateQuote(t *testing.T) {
	testCases := []struct {
		desc  string
		quote Quote
		err   error
	}{
		{
			desc:  "create no author",
			quote: Quote{Message: "be the change"},
			err:   ErrNoAuthor,
		},
		{
			desc:  "create no message",
			quote: Quote{Author: "ghandi"},
			err:   ErrNoMessage,
		},
		{
			desc:  "create success",
			quote: Quote{Author: "ghandi", Message: "be the change"},
			err:   nil,
		},
	}
	db := getDB(t)
	createTable(t, db)
	defer cleanup(t, db)
	store := NewPGStore(db)
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err := store.Create(tc.quote)
			if err != tc.err {
				t.Errorf("got %v; want %v", err, tc.err)
			}
		})
	}
}

func TestGetAQuoteByID(t *testing.T) {
	testCases := []struct {
		desc  string
		quote Quote
		id    int
		err   error
	}{
		{
			desc: "get invalid id",
			id:   0,
			err:  sql.ErrNoRows,
		},
		{
			desc:  "get success",
			id:    1,
			quote: Quote{Author: "ghandi", Message: "be the change"},
			err:   nil,
		},
	}
	db := getDB(t)
	createTable(t, db)
	seed(t, db)
	defer cleanup(t, db)
	store := NewPGStore(db)
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			qt, err := store.ByID(tc.id)
			if err != tc.err {
				t.Errorf("got %v; want %v", err, tc.err)
			}
			if err == nil {
				if qt.Author != tc.quote.Author {
					t.Errorf("got %v; want %v", err, tc.err)
				}
			}
		})
	}
}
