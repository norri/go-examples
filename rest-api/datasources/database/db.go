package database

import (
	"context"
)

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

func NewDatabase() *DB {
	return &DB{impl: newMemoryDB()}
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
