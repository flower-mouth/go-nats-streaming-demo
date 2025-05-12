package main

import (
	"context"
	"go-nats-streaming-demo/internal/configuration"
	"go-nats-streaming-demo/internal/database"
	"go-nats-streaming-demo/internal/services"
	"go-nats-streaming-demo/router"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", router.HomePage)
	mux.HandleFunc("/record", router.IdPage)
	mux.HandleFunc("/list/", router.DataListPage)

	err := database.SyncCacheAndDatabase(database.NewClient(context.Background(), 3, configuration.StorageConfig))
	if err != nil {
		log.Print(err)
	}

	go services.Subscriber()

	log.Printf("Starting server...")
	err = http.ListenAndServe(":8181", mux)
	if err != nil {
		log.Printf("Error in lauching server: %v", err)
	}
}
