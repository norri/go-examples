package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"app/datasources"
	"app/datasources/database"
	"app/server"
)

const readHeaderTimeout = 10 * time.Second

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := NewConfiguration()
	db := database.NewDatabase()
	defer db.CloseConnections()

	router := server.NewServer(ctx, &datasources.DataSources{DB: db})
	server := &http.Server{
		Addr:              ":" + conf.Port,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	log.Println("Listening on port", conf.Port)
	log.Panic(server.ListenAndServe())
}
