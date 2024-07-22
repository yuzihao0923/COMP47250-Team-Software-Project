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

	for i := 0; i < 100; i++ {
		msg := Message{
			ID:        i,
			Data:      "Hello from Producer",
			Timestamp: time.Now(),
		}
		jsonData, err := json.Marshal(msg)
		if err != nil {
			log.Fatalf("Error marshalling JSON: %v", err)
		}

		err = rsi.Set(ctx, fmt.Sprintf("message:%d", i), jsonData, 10*time.Second).Err()
		if err != nil {
			log.Fatalf("Error setting key in Redis: %v", err)
		}
		log.Printf("Sent: %s", jsonData)
	}
}
