package main

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"COMP47250-Team-Software-Project/internal/client"
	"COMP47250-Team-Software-Project/internal/database"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"fmt"
	"os"
	"time"
)

// SendMessage: send a new message to a stream (with streamName)
func SendMessage(brokerPort, streamName string, payload []byte, token string) error {
	msg := message.Message{
		Type: "produce",
		ConsumerInfo: &message.ConsumerInfo{
			StreamName: streamName,
		},
		Payload: payload,
	}
	var err error
	// err := client.SendMessage(brokerPort, msg, token)
	// if err != nil {
	// 	log.LogError("Producer", "producer has error sending message: "+err.Error())
	// 	return
	// }
	// log.LogInfo("Producer", fmt.Sprintf("Producer sent message: %s", msg.Payload))
	for retryCount := 0; retryCount < MaxRetryCount; retryCount++ {
		err = client.SendMessage(brokerPort, msg, token)
		if err == nil {
			log.LogInfo("Producer", fmt.Sprintf("Producer sent message: %s", msg.Payload))
			return nil
		}
		log.LogError("Producer", fmt.Sprintf("Error sending message (attempt %d/%d): %s", retryCount+1, MaxRetryCount, err.Error()))
		time.Sleep(RetryInterval)
	}
	return fmt.Errorf("failed to send message after %d attempts: %w", MaxRetryCount, err)
}

const (
	MaxRetryCount = 1000
	RetryInterval = 5 * time.Second
)

func main() {
	// Ensure logs are printed before prompting user input
	fmt.Println("[INFO] [Producer] Starting producer...")

	brokerPort := os.Getenv("BROKER_PORT")
	if brokerPort == "" {
		brokerPort = "8080" // Default port
	}

	err := database.ConnectMongoDB()
	if err != nil {
		fmt.Println("[ERROR] [Producer] Failed to connect to database:", err)
		return
	}
	fmt.Println("[INFO] [Producer] Database connected successfully")

	var token, role string
	for {
		username := auth.GetUserInput("\nEnter username: ")
		password := auth.GetPasswordInput("Enter password: ")

		token, role, err = auth.AuthenticateUser(username, password)
		if err != nil {
			fmt.Println(err)
		} else if role != "producer" {
			fmt.Println("this user is not a producer, please try again")
		} else {
			break
		}
	}

	for i := 0; i < 5; i++ {
		payload := []byte(fmt.Sprintf("Hello %d", i))
		err := SendMessage(brokerPort, "mystream", payload, token)
		if err != nil {
			// fmt.Println("[ERROR] [Producer] Failed to send message after retries:", err)

			log.LogError("Producer", fmt.Sprintf("Failed to send message after retries: %v", err))
		}
		time.Sleep(time.Millisecond) // Slight delay to prevent overwhelming the broker
	}
}
