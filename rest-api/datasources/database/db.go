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

func NewDatabase() Database {
	return newMemoryDB()
}
