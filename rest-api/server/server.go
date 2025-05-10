package server

import (
	"log/slog"
	"net/http"

	"app/datasources"
	"app/server/handlers"
	"app/server/services"
)

func NewServer(dataSources *datasources.DataSources) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/status", func(w http.ResponseWriter, _ *http.Request) {
		_, err := w.Write([]byte("ok"))
		if err != nil {
			slog.Error("failed to write status response", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	mux.HandleFunc("GET /api/v1/books", handlers.GetBooks(services.NewBooksService(dataSources.DB)))
	mux.HandleFunc("POST /api/v1/books", handlers.AddBook(services.NewBooksService(dataSources.DB)))

	return mux
}
