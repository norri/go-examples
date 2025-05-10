package services

import (
	"context"
	"errors"
	"testing"

	"app/datasources/database"
	"app/server/domain"
)

func TestGetBooks(t *testing.T) {
	mockDB := &database.MockDatabase{
		LoadAllBooksFunc: func(_ context.Context) ([]database.Book, error) {
			return []database.Book{{Title: "Title"}}, nil
		},
	}

	service := NewBooksService(mockDB)
	books, err := service.GetBooks(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}
}

func TestGetBooks_Fails(t *testing.T) {
	mockDB := &database.MockDatabase{
		LoadAllBooksFunc: func(_ context.Context) ([]database.Book, error) {
			return nil, errors.New("error")
		},
	}

	service := NewBooksService(mockDB)
	_, err := service.GetBooks(context.Background())
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestSaveBook(t *testing.T) {
	mockDB := &database.MockDatabase{
		CreateBookFunc: func(_ context.Context, _ database.NewBook) error {
			return nil
		},
	}

	service := NewBooksService(mockDB)
	err := service.SaveBook(context.Background(), domain.Book{Title: "Title"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSaveBook_Fails(t *testing.T) {
	mockDB := &database.MockDatabase{
		CreateBookFunc: func(_ context.Context, _ database.NewBook) error {
			return errors.New("error")
		},
	}

	service := NewBooksService(mockDB)
	err := service.SaveBook(context.Background(), domain.Book{Title: "Title"})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}
