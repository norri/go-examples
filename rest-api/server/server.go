package server

import (
	"context"
	"net/http"

	"app/datasources"
	"app/server/handlers"
	"app/server/services"
)

func NewServer(ctx context.Context, dataSources *datasources.DataSources) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.HandleFunc("GET /api/v1/books", handlers.GetBooks(services.NewBooksService(dataSources.DB)))
	mux.HandleFunc("POST /api/v1/books", handlers.AddBook(services.NewBooksService(dataSources.DB)))

	return mux
}
