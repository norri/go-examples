package main

import (
	"log"
	"net/http"
	"time"

	"app/datasources"
	"app/datasources/database"
	"app/server"
)

const readHeaderTimeout = 10 * time.Second

func main() {
	conf := NewConfiguration()
	db := database.NewDatabase()
	defer db.CloseConnections()

	router := server.NewServer(&datasources.DataSources{DB: db})
	server := &http.Server{
		Addr:              ":" + conf.Port,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
	}
	log.Println("Listening on port", conf.Port)
	log.Panic(server.ListenAndServe())
}
