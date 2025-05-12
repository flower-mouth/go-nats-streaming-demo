package services

import (
	"context"
	"encoding/json"
	"fmt"
	stan "github.com/nats-io/stan.go"
	"go-nats-streaming-demo/internal/configuration"
	"go-nats-streaming-demo/internal/database"
	"go-nats-streaming-demo/internal/models"
	"log"
	"sync"
)

func unmarshalMessage(m []byte) (models.Order, error) {
	var order models.Order

	err := json.Unmarshal(m, &order)
	if err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func Subscriber() {
	fmt.Printf("subscriber started\n")

	sc, err := stan.Connect("prod", "sub-1")
	if err != nil {
		log.Fatalf("Ошибка подключения к NATS Streaming: %v", err)
	}
	if sc != nil {
		defer sc.Close()
	}

	fmt.Printf("subscriber connected\n")

	_, err = sc.Subscribe("Test4", func(m *stan.Msg) {
		order, err := unmarshalMessage(m.Data)
		if err != nil {
			log.Printf("Marshaling failed (incorrect message type): %v\n", err)
		} else {
			err = database.AddMessageToDatabase(database.NewClient(context.Background(), 3, configuration.StorageConfig), order)
			fmt.Printf("Stored in database\n")
			if err != nil {
				log.Print(err)
			} else {
				models.Cache[order.OrderUID] = order
			}
		}
	})
	if err != nil {
		log.Printf("Error in subscription: %v\n", err)
	}

	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}
