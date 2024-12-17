package services

import (
	"context"

	"app/server/domain"

	"github.com/stretchr/testify/mock"
)

type BooksServiceMock struct {
	mock.Mock
}

func (m *BooksServiceMock) GetBooks(ctx context.Context) ([]domain.Book, error) {
	args := m.Called(ctx)
	books, ok := args.Get(0).([]domain.Book)
	if !ok {
		return nil, args.Error(1)
	}
	return books, args.Error(1)
}

func (m *BooksServiceMock) SaveBook(ctx context.Context, newBook domain.Book) error {
	args := m.Called(ctx, newBook)
	return args.Error(0)
}
