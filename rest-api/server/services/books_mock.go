package services

import (
	"context"

	"app/server/domain"
)

type BooksServiceMock struct {
	GetBooksFunc func(ctx context.Context) ([]domain.Book, error)
	SaveBookFunc func(ctx context.Context, newBook domain.Book) error
}

func (m *BooksServiceMock) GetBooks(ctx context.Context) ([]domain.Book, error) {
	return m.GetBooksFunc(ctx)
}

func (m *BooksServiceMock) SaveBook(ctx context.Context, newBook domain.Book) error {
	return m.SaveBookFunc(ctx, newBook)
}
