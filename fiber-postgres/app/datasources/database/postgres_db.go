package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPool interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Close()
}

func newPostgresDB(ctx context.Context, databaseURL string) (Database, error) {
	// For production use set connection pool settings and validate connection with ping
	dbpool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %v", err)
	}
	return &postgresDB{
		pool: dbpool,
	}, nil
}

type postgresDB struct {
	pool PostgresPool
}

func (db *postgresDB) LoadAllBooks(ctx context.Context) ([]Book, error) {
	rows, err := db.pool.Query(ctx, "SELECT id, title FROM books")
	if err != nil {
		return nil, fmt.Errorf("failed to query books table: %w", err)
	}
	defer rows.Close()

	books, err := pgx.CollectRows(rows, pgx.RowToStructByName[Book])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows: %w", err)
	}
	return books, nil
}

func (db *postgresDB) CreateBook(ctx context.Context, newBook NewBook) error {
	_, err := db.pool.Exec(ctx, "INSERT INTO books (title) VALUES ($1)", newBook.Title)
	if err != nil {
		return fmt.Errorf("failed to insert book: %w", err)
	}
	return nil
}

func (db *postgresDB) CloseConnections() {
	db.pool.Close()
}
