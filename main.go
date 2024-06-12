package main

import (
	"COMP47250-Team-Software-Project/cmd/broker"
	"COMP47250-Team-Software-Project/cmd/consumer"
	"COMP47250-Team-Software-Project/cmd/producer"
	"COMP47250-Team-Software-Project/internal/redis"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// Initialize Redis and flush all data
	redis.Initialize("localhost:6379", "", 0)
	err := redis.FlushAll()
	if err != nil {
		fmt.Printf("Failed to flush Redis: %v\n", err)
		return
	}
	fmt.Println("Redis database has been flushed")

	// Start broker
	wg.Add(1)
	go func() {
		defer wg.Done()
		broker.StartBroker()
	}()

	// Waiting for broker to connect
	time.Sleep(2 * time.Second)

	// Start consumer
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
