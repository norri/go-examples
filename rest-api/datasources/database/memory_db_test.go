package database

import (
	"context"
	"testing"
)

func TestMemoryDB_LoadBooks(t *testing.T) {
	db := newMemoryDB()
	books, err := db.LoadAllBooks(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(books) != 0 {
		t.Fatalf("expected 0 books, got %d", len(books))
	}
}

func TestMemoryDB_SaveBook(t *testing.T) {
	db := newMemoryDB()
	newBook := NewBook{Title: "Title"}
	err := db.CreateBook(context.Background(), newBook)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	books, err := db.LoadAllBooks(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(books))
	}
	assertBook(t, books[0], 0, newBook)
}

func TestMemoryDB_SaveBookMultiple(t *testing.T) {
	db := newMemoryDB()
	newBook1 := NewBook{Title: "Title1"}
	err := db.CreateBook(context.Background(), newBook1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	newBook2 := NewBook{Title: "Title2"}
	err = db.CreateBook(context.Background(), newBook2)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	books, err := db.LoadAllBooks(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(books) != 2 {
		t.Fatalf("expected 2 books, got %d", len(books))
	}
	assertBook(t, books[0], 0, newBook1)
	assertBook(t, books[1], 1, newBook2)
}

func assertBook(t *testing.T, book Book, expectedID int, expected NewBook) {
	t.Helper()
	if book.ID != expectedID {
		t.Fatalf("expected ID %d, got %d", expectedID, book.ID)
	}
	if book.Title != expected.Title {
		t.Fatalf("expected Title %s, got %s", expected.Title, book.Title)
	}
}
