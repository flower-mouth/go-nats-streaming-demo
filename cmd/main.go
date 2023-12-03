package main

import (
	"context"
	"fmt"
	"wbLab0/internal/configuration"
	"wbLab0/internal/database"
	"wbLab0/internal/services"
)

func main() {
	postgreSQLClient, err := database.NewClient(context.Background(), 3, configuration.StorageConfig)
	if err != nil {
		fmt.Printf("%v", err)
	}
	services.Subscriber(postgreSQLClient)
}
