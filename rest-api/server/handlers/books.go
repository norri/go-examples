package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"app/server/domain"
	"app/server/services"
)

func GetBooks(service services.BooksService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		books, err := service.GetBooks(r.Context())
		if err != nil {
			slog.Error("GetBooks failed", "error", err)
			sendJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
				Error: "internal error",
			})
			return
		}

		sendJSON(w, http.StatusOK, domain.BooksResponse{
			Books: books,
		})
	}
}

func AddBook(service services.BooksService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var book domain.Book
		err := json.NewDecoder(r.Body).Decode(&book)
		if err != nil {
			slog.Warn("AddBook request parsing failed", "error", err)
			sendJSON(w, http.StatusBadRequest, domain.ErrorResponse{
				Error: "invalid request",
			})
			return
		}
		// For production use add proper validation here

		err = service.SaveBook(r.Context(), book)
		if err != nil {
			slog.Error("AddBook failed", "error", err)
			sendJSON(w, http.StatusInternalServerError, domain.ErrorResponse{
				Error: "internal error",
			})
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func sendJSON(w http.ResponseWriter, code int, response any) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		slog.Error("sendResponse json.Marshal failed", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(jsonResponse)
	if err != nil {
		slog.Error("sendResponse w.Write failed", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
