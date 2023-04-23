package store

import (
	"database/sql"
	"errors"
	"os"
	"testing"
	"time"

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
	query := `CREATE TABLE IF NOT EXISTS quotes(id SERIAL PRIMARY KEY, author TEXT UNIQUE, message TEXT UNIQUE);`
	stmt, err := db.Prepare(query)
	if err != nil {
		t.Fatal(err)
	}
	_, err = stmt.Exec()
	if err != nil {
		t.Fatal(err)
	}
}

func cleanup(t *testing.T, db *sql.DB) {
	t.Cleanup(func() {
		db.Exec(`DROP TABLE quotes;`)
	})
}

func TestSetup(t *testing.T) {
	db := getDB(t)
	createTable(t, db)
	time.Sleep(20 * time.Second)
	cleanup(t, db)
}
