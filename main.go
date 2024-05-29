package main

import (
	"fmt"
	"time"

	"COMP47250-Team-Software-Project/cmd/broker"
	"COMP47250-Team-Software-Project/cmd/producer"
)

func main() {
	// Start the broker in a goroutine
	go func() {
		fmt.Println("Starting broker...")
		broker.StartBroker()
	}()

	// Wait a bit to ensure the broker is up and running
	time.Sleep(2 * time.Second)

	// Start the producer
	fmt.Println("Starting producer...")
	producer.StartProducer()

	// Keep the main function running to allow broker and producer to communicate
	select {}
}
