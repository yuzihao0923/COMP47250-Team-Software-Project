package main

import (
	"COMP47250-Team-Software-Project/cmd/broker"
	"COMP47250-Team-Software-Project/cmd/consumer"
	"COMP47250-Team-Software-Project/cmd/producer"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// Start broker
	wg.Add(1)
	go func() {
		defer wg.Done()
		broker.StartBroker()
	}()

	// Waiting for broker to connect
	time.Sleep(2 * time.Second)

	// start consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		consumer.StartConsumer()
	}()

	// Waiting for consumer to connect
	time.Sleep(2 * time.Second)

	// Start producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		producer.StartProducer()
	}()

	// Waiting for all goroutines to finish
	wg.Wait()

	fmt.Println("All services have been started and executed.")
}
