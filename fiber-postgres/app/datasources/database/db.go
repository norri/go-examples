package database

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
)

type Book struct {
	ID    int
	Title string
}

type NewBook struct {
	Title string
}

type Database interface {
	// LoadAllBooks retrieves all books from the database.
	LoadAllBooks(ctx context.Context) ([]Book, error)

	// CreateBook adds a new book to the database.
	CreateBook(ctx context.Context, newBook NewBook) error

	// CloseConnections closes all open connections to the database.
	CloseConnections()
}

func NewDatabase(ctx context.Context, databaseURL string) (Database, error) {
	if databaseURL == "" {
		slog.Info("Using in-memory database implementation")
		return newMemoryDB(), nil
	}

	if strings.HasPrefix(databaseURL, "postgres://") {
		db, err := newPostgresDB(ctx, databaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize PostgreSQL database connection: %w", err)
		}
		slog.Info("Using PostgreSQL database implementation")
		return db, nil
	}

	return nil, fmt.Errorf("unsupported database URL scheme: %s", databaseURL)
}
