package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Message struct {
	ID        int       `json:"id"`
	Data      string    `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	rsi := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"localhost:6381", "localhost:6382", "localhost:6383",
			"localhost:6384", "localhost:6385", "localhost:6386",
		},
	})
	defer rsi.Close()

	numMessages := 100
	log.Println("Starting to consume messages...")

	var totalStart time.Time
	firstMessageReceived := false
	var totalLatency time.Duration

	for i := 0; i < numMessages; i++ {
		for {
			result, err := rsi.Get(ctx, fmt.Sprintf("message:%d", i)).Result()
			if err == redis.Nil {
				time.Sleep(100 * time.Millisecond) // Adjust sleep time as needed
				continue
			} else if err != nil {
				log.Fatalf("Error getting message %d: %v", i, err)
			}

			if !firstMessageReceived {
				totalStart = time.Now() // Start the timer on first message retrieval
				firstMessageReceived = true
			}

			var msg Message
			err = json.Unmarshal([]byte(result), &msg)
			if err != nil {
				log.Fatalf("Error unmarshalling message %d: %v", i, err)
			}

			latency := time.Since(msg.Timestamp)
			totalLatency += latency
			retrieveTime := time.Since(msg.Timestamp) // Time taken to retrieve and process this message
			log.Printf("Received: %+v, Latency: %.3f seconds, Retrieval Time: %.3f seconds", msg, latency.Seconds(), retrieveTime.Seconds())
			break
		}
	}

	if firstMessageReceived {
		totalTime := time.Since(totalStart)
		averageLatency := totalLatency / time.Duration(numMessages)
		log.Printf("Total time from first sent to last received: %.3f seconds", totalTime.Seconds())
		log.Printf("Average latency: %.3f seconds", averageLatency.Seconds())
	} else {
		log.Println("No messages were received.")
	}
}
