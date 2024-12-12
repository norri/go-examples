package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/server/domain"
	"app/server/services"
)

var booksRoute = "/api/v1/books"

func TestGetBooks(t *testing.T) {
	mockService := &services.BooksServiceMock{
		GetBooksFunc: func(ctx context.Context) ([]domain.Book, error) {
			return []domain.Book{{Title: "Title"}}, nil
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc(booksRoute, GetBooks(mockService))

	req := httptest.NewRequest("GET", booksRoute, nil)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}

	body := bodyFromResponse[domain.BooksResponse](t, resp.Result())
	if len(body.Books) != 1 {
		t.Fatalf("expected 1 book, got %d", len(body.Books))
	}
}

func TestGetBooks_ServiceFails(t *testing.T) {
	mockService := &services.BooksServiceMock{
		GetBooksFunc: func(ctx context.Context) ([]domain.Book, error) {
			return nil, errors.New("error")
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc(booksRoute, GetBooks(mockService))

	req := httptest.NewRequest("GET", booksRoute, nil)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.Code)
	}

	body := bodyFromResponse[domain.ErrorResponse](t, resp.Result())
	if body.Error != "internal error" {
		t.Fatalf("expected error message 'internal error', got '%s'", body.Error)
	}
}

func TestAddBook(t *testing.T) {
	mockService := &services.BooksServiceMock{
		SaveBookFunc: func(ctx context.Context, newBook domain.Book) error {
			return nil
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc(booksRoute, AddBook(mockService))

	req := postRequest(booksRoute, `{"title":"Title"}`)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.Code)
	}
}

func TestAddBook_InvalidRequest(t *testing.T) {
	mockService := &services.BooksServiceMock{}

	mux := http.NewServeMux()
	mux.HandleFunc(booksRoute, AddBook(mockService))

	req := httptest.NewRequest("POST", booksRoute, nil)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}

	body := bodyFromResponse[domain.ErrorResponse](t, resp.Result())
	if body.Error != "invalid request" {
		t.Fatalf("expected error message 'invalid request', got '%s'", body.Error)
	}
}

func TestAddBook_ServiceFails(t *testing.T) {
	mockService := &services.BooksServiceMock{
		SaveBookFunc: func(ctx context.Context, newBook domain.Book) error {
			return errors.New("error")
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc(booksRoute, AddBook(mockService))

	req := postRequest(booksRoute, `{"title":"Title"}`)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, resp.Code)
	}

	body := bodyFromResponse[domain.ErrorResponse](t, resp.Result())
	if body.Error != "internal error" {
		t.Fatalf("expected error message 'internal error', got '%s'", body.Error)
	}
}

func postRequest(url string, body string) *http.Request {
	req := httptest.NewRequest("POST", url, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func bodyFromResponse[T any](t *testing.T, resp *http.Response) T {
	defer resp.Body.Close()
	var body T
	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return body
}