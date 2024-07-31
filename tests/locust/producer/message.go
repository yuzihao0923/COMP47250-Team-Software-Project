package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// Message represents the structure of the message
type Message struct {
	ConsumerInfo ConsumerInfo `json:"consumer_info"`
}

// ConsumerInfo represents the consumer information
type ConsumerInfo struct {
	StreamName string `json:"stream_name"`
	GroupName  string `json:"group_name"`
}

// generateRandomString generates a random string of a given length
func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	// Open a file for writing
	file, err := os.Create("messages.json")
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	// Create a JSON encoder
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	// Start a JSON array
	file.WriteString("[\n")

	// Generate 50000 messages
	for i := 0; i < 50000; i++ {
		message := Message{
			ConsumerInfo: ConsumerInfo{
				StreamName: generateRandomString(10),
				GroupName:  generateRandomString(10),
			},
		}

		// Encode the message as JSON
		err := encoder.Encode(message)
		if err != nil {
			log.Fatalf("Failed to write message to file: %v", err)
		}

		// Add a comma between objects
		if i < 49999 {
			file.WriteString(",\n")
		}
	}

	// End the JSON array
	file.WriteString("\n]")

	fmt.Println("Successfully generated 50000 messages.")
}
