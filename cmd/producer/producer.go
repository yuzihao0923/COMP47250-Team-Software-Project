package main

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"fmt"
	"os"
	"time"
)

// SendMessage: send a new message to a stream (with streamName)
func SendMessage(brokerPort, streamName string, payload []byte, token string) {
	msg := message.Message{
		Type: "produce",
		ConsumerInfo: &message.ConsumerInfo{
			StreamName: streamName,
		},
		Payload: payload,
	}
	err := api.SendMessage(brokerPort, msg, token)
	if err != nil {
		log.LogError("Producer", "producer has error sending message: "+err.Error())
		return
	}
	log.LogInfo("Producer", fmt.Sprintf("Producer sent message: %s", msg.Payload))
}

func StartProducer() {
	log.LogInfo("Producer", "Starting producer...")

	brokerPort := os.Getenv("BROKER_PORT")
	if brokerPort == "" {
		brokerPort = "8080" // default port
	}

	// Initialize BroadcastFunc for logging
	log.BroadcastFunc = api.BroadcastMessage

	token, err := api.GetJWTToken("producer", "123")
	if err != nil {
		log.LogError("Producer", fmt.Sprintf("Failed to get JWT token: %v", err))
		return
	}

	for i := 0; i < 10; i++ {
		payload := []byte(fmt.Sprintf("Hello %d", i))
		SendMessage(brokerPort, "mystream", payload, token)
		time.Sleep(time.Millisecond) // slight delay to prevent overwhelming the broker
	}
}

func main() {
	StartProducer()
}
