package main

import (
	"encoding/json"
	"fmt"
	stan "github.com/nats-io/stan.go"
	"go-nats-streaming-demo/internal/models"
	"log"
	"strconv"
	"time"
)

func main() {
	fmt.Printf("publisher started\n")

	sc, err := stan.Connect(
		"prod",
		"simple-pub",
		stan.NatsURL("nats://localhost:4222"),
	)
	if err != nil {
		log.Fatalf("Ошибка подключения к NATS Streaming: %v", err)
	}
	defer sc.Close()

	item := models.Items{
		ChrtID:      993930,
		TrackNumber: "WBILTESTTRACK",
		Price:       9999,
		Rid:         "utyadfg097656jhfcgvhbtest",
		Name:        "testtesttest",
		Sale:        300,
		Size:        "0",
		TotalPrice:  9699,
		NmID:        864578,
		Brand:       "Vivienne Westwood",
		Status:      202,
	}

	var items []models.Items
	items = append(items, item)

	order := models.Order{
		OrderUID:    "hgf78jan06b2rd8b1test",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "oleg",
			Phone:   "+9156120000",
			Zip:     "3900006",
			City:    "Ryazan",
			Address: "Griboedov St",
			Region:  "Ryazan Oblast",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  "hgf78jan06b2rd8b1test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			PaymentDT:    987,
			Amount:       9876656,
			Bank:         "tinkoff",
			DeliveryCost: 500,
			GoodsTotal:   100,
			CustomFee:    0,
		},
		Items:             items,
		Locale:            "en",
		InternalSignature: "",
		CustomerId:        "test",
		DeliveryService:   "pochta ROSSII",
		Shardkey:          "9",
		SmId:              99,
		DateCreated:       "2023-12-2T06:22:19Z",
		OofShred:          "1",
	}

	for i := 400; ; i++ {

		order.OrderUID = strconv.Itoa(i)            // create unique identifier
		order.Payment.Transaction = strconv.Itoa(i) // create unique identifier
		order.Items[0].Price = int64(100 + i)       // create unique identifier
		order.Payment.Amount = int64(i * 10)
		order.Delivery.Address = "Griboedov St " + strconv.Itoa(i+10)
		record, _ := json.Marshal(order)

		err := sc.Publish("Test4", record)
		fmt.Printf("published successful\n")
		if err != nil {
			log.Printf("Error in publishing message: %v\n", err)
		}

		time.Sleep(10 * time.Second)
	}
}
