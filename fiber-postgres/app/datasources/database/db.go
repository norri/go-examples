package database

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

var ErrUnsupportedDatabase = errors.New("unsupported database URL")

type Book struct {
	ID    int
	Title string
}

type NewBook struct {
	Title string
}

type Database interface {
	LoadAllBooks(ctx context.Context) ([]Book, error)
	CreateBook(ctx context.Context, newBook NewBook) error
	CloseConnections()
}

func NewDatabase(ctx context.Context, databaseURL string) (*DB, error) {
	switch {
	case databaseURL == "":
		db := newMemoryDB()
		slog.Info("Using in-memory database implementation")
		return &DB{impl: db}, nil

	case strings.HasPrefix(databaseURL, "postgres://"):
		db, err := newPostgresDB(ctx, databaseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize PostgreSQL database connection: %w", err)
		}
		slog.Info("Using PostgreSQL database implementation")
		return &DB{impl: db}, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrUnsupportedDatabase, databaseURL)
}

type DB struct {
	impl Database
}

func (db *DB) LoadAllBooks(ctx context.Context) ([]Book, error) {
	return db.impl.LoadAllBooks(ctx)
}

func (db *DB) CreateBook(ctx context.Context, newBook NewBook) error {
	return db.impl.CreateBook(ctx, newBook)
}

func (db *DB) CloseConnections() {
	db.impl.CloseConnections()
}
