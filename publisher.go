package main

import (
	"fmt"
	stan "github.com/nats-io/stan.go"
	"log"
	"strconv"
	"time"
)

// .\nats-streaming-server -cid prod -store file -dir store
func main() {
	fmt.Printf("producer working\n")
	sc, _ := stan.Connect("prod", "simple-pub")
	defer sc.Close()

	for i := 1; ; i++ {

		err := sc.Publish("Test3", []byte("Test "+strconv.Itoa(i)))

		if err != nil {
			log.Printf("Error in publishing message: %v\n", err)
		}
		fmt.Printf("POSTED %v \n", strconv.Itoa(i))
		time.Sleep(5 * time.Second)

	}
}
