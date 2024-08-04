package main

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"COMP47250-Team-Software-Project/internal/client"
	"COMP47250-Team-Software-Project/internal/database"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"fmt"
	"sync"
	"time"
)

const (
	MaxRetryCount = 3               // 最大重试次数
	RetryInterval = 2 * time.Second // 重试间隔
)

var (
	proxyURL = "http://localhost:8888"
	users    = []struct {
		username string
		password string
	}{
		{"p1", "123"},
		{"p2", "123"},
		{"p3", "123"},
		{"p4", "123"},
	}
)

func SendMessage(brokerAddr, streamName string, payload []byte, token string) error {
	msg := message.Message{
		Type: "produce",
		ConsumerInfo: &message.ConsumerInfo{
			StreamName: streamName,
		},
		Payload: payload,
	}

	var err error
	for retryCount := 0; retryCount < MaxRetryCount; retryCount++ {
		err = client.SendMessage(brokerAddr, msg, token)
		if err == nil {
			log.LogInfo("Producer", fmt.Sprintf("Producer sent message: %s", msg.Payload))
			return nil
		}
		log.LogError("Producer", fmt.Sprintf("Error sending message (attempt %d/%d): %s", retryCount+1, MaxRetryCount, err.Error()))
		time.Sleep(RetryInterval)
	}
	return fmt.Errorf("failed to send message after %d attempts: %w", MaxRetryCount, err)
}

func main() {
	fmt.Println("[INFO] [Producer] Starting producer...")

	broker, err := client.GetBroker(proxyURL)
	if err != nil {
		log.LogError("Producer", fmt.Sprintf("Get broker failed, error: %v", err))
		return
	}
	brokerAddr := broker.Address

	_, err = database.ConnectMongoDB("mongodb://localhost:27017", "comp47250", "users")
	if err != nil {
		fmt.Println("[ERROR] [Consumer] Failed to connect to database:", err)
		return
	}
	fmt.Println("[INFO] [Producer] Database connected successfully")

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(user struct{ username, password string }) {
			defer wg.Done()
			token, role, err := auth.AuthenticateUser(user.username, user.password, brokerAddr)
			if err != nil {
				fmt.Println(err)
				return
			}
			if role != "producer" {
				fmt.Println("this user is not a producer, please try again")
				return
			}

			for i := 0; i < 2500; i++ {
				payload := []byte(fmt.Sprintf("Hello %d from %s", i, user.username))
				err := SendMessage(brokerAddr, "mystream", payload, token)
				if err != nil {
					log.LogError("Producer", fmt.Sprintf("Failed to send message: %v", err))
				}
				// time.Sleep(time.Millisecond * 100) // Slight delay to prevent overwhelming the broker
			}
		}(user)
	}
	wg.Wait() // Wait for all goroutines to finish
	fmt.Println("[INFO] [Producer] All messages sent.")
}
