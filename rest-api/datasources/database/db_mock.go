package database

import (
	"context"
)

type MockDatabase struct {
	LoadAllBooksFunc func(ctx context.Context) ([]Book, error)
	CreateBookFunc   func(ctx context.Context, newBook NewBook) error
}

func (m *MockDatabase) LoadAllBooks(ctx context.Context) ([]Book, error) {
	return m.LoadAllBooksFunc(ctx)
}

func (m *MockDatabase) CreateBook(ctx context.Context, newBook NewBook) error {
	return m.CreateBookFunc(ctx, newBook)
}

func (m *MockDatabase) CloseConnections() {
}
