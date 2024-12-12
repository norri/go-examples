package database

import (
	"context"
)

type DatabaseMock struct {
	LoadAllBooksFunc func(ctx context.Context) ([]Book, error)
	CreateBookFunc   func(ctx context.Context, newBook NewBook) error
}

func (m *DatabaseMock) LoadAllBooks(ctx context.Context) ([]Book, error) {
	return m.LoadAllBooksFunc(ctx)
}

func (m *DatabaseMock) CreateBook(ctx context.Context, newBook NewBook) error {
	return m.CreateBookFunc(ctx, newBook)
}

func (m *DatabaseMock) CloseConnections() {
}
