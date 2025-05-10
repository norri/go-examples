package services

import (
	"context"

	"app/server/domain"
)

type MockBooksService struct {
	GetBooksFunc func(ctx context.Context) ([]domain.Book, error)
	SaveBookFunc func(ctx context.Context, newBook domain.Book) error
}

func (m *MockBooksService) GetBooks(ctx context.Context) ([]domain.Book, error) {
	return m.GetBooksFunc(ctx)
}

func (m *MockBooksService) SaveBook(ctx context.Context, newBook domain.Book) error {
	return m.SaveBookFunc(ctx, newBook)
}
