package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var (
	ErrNoAuthor  = errors.New("must provide author")
	ErrNoMessage = errors.New("must provide message")
)

type PGStore struct{ db *sql.DB }

func NewPGStore(db *sql.DB) *PGStore {
	return &PGStore{db: db}
}

func GetDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("couldn't connect to RDBMS: %v", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("couldn't connect to database: %v", err)
	}
	return db, nil
}

func (ps *PGStore) Create(quote Quote) error {
	if quote.Author == "" {
		return ErrNoAuthor
	}
	if quote.Message == "" {
		return ErrNoMessage
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	query := `INSERT INTO quotes (author, message) VALUES ($1, $2);`
	stmt, err := ps.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("couldn't prepare statement: %v", err)
	}
	_, err = stmt.ExecContext(ctx, quote.Author, quote.Message)
	if err != nil {
		return fmt.Errorf("couldn't add quote to DB: %v", err)
	}
	return nil
}

func (ps *PGStore) ByID(id int) (Quote, error) {
	var quote Quote
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	query := `SELECT id, author, message FROM quotes WHERE id = $1;`
	stmt, err := ps.db.PrepareContext(ctx, query)
	if err != nil {
		return quote, fmt.Errorf("couldn't prepare statement: %v", err)
	}
	if err := stmt.QueryRow(id).Scan(&quote.ID, &quote.Author, &quote.Message); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return quote, sql.ErrNoRows
		}
		return quote, fmt.Errorf("couldn't scan row into quote: %v", err)
	}
	return quote, nil
}

func (ps *PGStore) All() ([]Quote, error) {
	var quotes []Quote
	query := `SELECT id, author, message FROM quotes;`
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stmt, err := ps.db.PrepareContext(ctx, query)
	if err != nil {
		return quotes, fmt.Errorf("couldn't prepare statement: %v", err)
	}
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return quotes, fmt.Errorf("couldn't get rows: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var q Quote
		if err := rows.Scan(&q.ID, &q.Author, &q.Message); err != nil {
			log.Println("scan failed:", err)
		}
		quotes = append(quotes, q)
	}
	return quotes, nil
}
