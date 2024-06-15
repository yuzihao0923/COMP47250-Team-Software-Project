package main

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"fmt"
	"os"
	"time"
)

// RegisterConsumerGroup: use api to register a consumer group (with groupName) in stream (with streamName)
func RegisterConsumerGroup(brokerPort, streamName, groupName string) {
	err := api.RegisterConsumer(brokerPort, streamName, groupName)
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Consumer has error registering: %v", err))
		return
	}
	log.LogMessage("INFO", "Consumer register to Broker...")

}

// ConsumeMessages: consumer (with consumerID) gets messages from a group (with groupName) of a stream (with streamName)
func ConsumeMessages(brokerPort, streamName, groupName, consumerID string) {
	for {
		messages, err := api.ConsumeMessages(brokerPort, streamName, groupName, consumerID)
		if err != nil {
			log.LogMessage("ERROR", fmt.Sprintf("Consumer has error receiving message: %v", err))
			time.Sleep(time.Second * 1)
			continue
		}

		for _, msg := range messages {
			time.Sleep(time.Millisecond) // insure the order of log between "producer send" & "consumer receive"
			log.LogMessage("INFO", "Consumer received message: "+string(msg.Payload))
		}

		time.Sleep(time.Second * 1)
	}
}

func StartConsumer() {
	log.LogMessage("INFO", "Starting consumer...")

	brokerPort := os.Getenv("BROKER_PORT")
	if brokerPort == "" {
		brokerPort = "8080" // default port
	}

	// register consumer group
	RegisterConsumerGroup(brokerPort, "mystream", "mygroup")

	ConsumeMessages(brokerPort, "mystream", "mygroup", "myconsumer")
}

func main() {
	StartConsumer()
}
