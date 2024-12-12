package main

import (
	"context"
	"log"
	"net/http"

	"app/datasources"
	"app/datasources/database"
	"app/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := NewConfiguration()
	db := database.NewDatabase()
	defer db.CloseConnections()

	router := server.NewServer(ctx, &datasources.DataSources{DB: db})
	server := &http.Server{
		Addr:    ":" + conf.Port,
		Handler: router,
	}
	log.Println("Listening on port", conf.Port)
	log.Fatal(server.ListenAndServe())
}
