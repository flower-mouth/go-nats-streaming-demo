package services

import (
	"fmt"
	stan "github.com/nats-io/stan.go"
	"sync"
	"wbLab0/internal/database"
	"wbLab0/internal/models"
)

func Subscriber(postgreSQLClient database.Client) {
	fmt.Printf("subscriber started\n")

	sc, _ := stan.Connect("prod", "sub-1")
	defer sc.Close()

	fmt.Printf("subscriber connected\n")

	sc.Subscribe("Test3", func(m *stan.Msg) {
		fmt.Printf("Got: %v \n\n", string(m.Data))
		testtt := models.IntTest{m.Data}

		err := database.AddMessageToDatabase(postgreSQLClient, testtt) // add to database
		if err != nil {
			fmt.Printf("%v", err)
		}
	})

	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}
