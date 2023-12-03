package main

import (
	"context"
	"log"
	"wbLab0/internal/configuration"
	"wbLab0/internal/database"
	"wbLab0/internal/services"
)

func main() {
	err := database.SyncCacheAndDatabase(database.NewClient(context.Background(), 3, configuration.StorageConfig))
	if err != nil {
		log.Print(err)
	}

	services.Subscriber()
}
