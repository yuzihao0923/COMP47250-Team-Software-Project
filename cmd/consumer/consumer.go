package main

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"fmt"
	"os"
	"time"
)

// RegisterConsumerGroup: use api to register a consumer group (with groupName) in stream (with streamName)
func RegisterConsumerGroup(brokerPort, streamName, groupName, token string) {
	msg := message.Message{
		Type: "registration",
		ConsumerInfo: &message.ConsumerInfo{
			StreamName: streamName,
			GroupName:  groupName,
		},
	}
	err := api.RegisterConsumer(brokerPort, msg, token)
	if err != nil {
		log.LogError("Consumer", "consumer has error registering: "+err.Error())
		return
	}
	log.LogInfo("Consumer", "Consumer registered to Broker...")
}

func ConsumeMessages(brokerPort, streamName, groupName, consumerID, token string) {
	for {
		messages, err := api.ConsumeMessages(brokerPort, streamName, groupName, consumerID, token)
		if err != nil {
			if err.Error() == "no new messages" {
				log.LogWarning("Consumer", "No new messages, retrying...")
				time.Sleep(time.Second * 1)
				continue
			} else {
				log.LogError("Consumer", "consumer has error receiving message: "+err.Error())
				time.Sleep(time.Second * 1)
				continue
			}
		}

		for _, msg := range messages {
			time.Sleep(time.Millisecond) // ensure the order of log between "producer send" & "consumer receive"
			log.LogInfo("Consumer", "Consumer received message: "+string(msg.Payload))

			AcknowledgeMessage(brokerPort, msg, token)
		}

		time.Sleep(time.Second * 1)
	}
}

func AcknowledgeMessage(brokerPort string, msg message.Message, token string) {
	err := api.SendACK(brokerPort, msg, token)
	if err != nil {
		log.LogError("Consumer", "consumer has error sending ACK: "+err.Error())
		return
	}
	log.LogInfo("Consumer", "Consumer sending ACK successfully...")
}

func StartConsumer() {
	log.LogInfo("Consumer", "Starting consumer...")

	brokerPort := os.Getenv("BROKER_PORT")
	if brokerPort == "" {
		brokerPort = "8080" // default port
	}

	// Initialize BroadcastFunc for logging
	log.BroadcastFunc = api.BroadcastMessage

	token, err := api.GetJWTToken("consumer", "123")
	if err != nil {
		log.LogError("Consumer", fmt.Sprintf("Failed to get JWT token: %v", err))
		return
	}

	// register consumer group
	RegisterConsumerGroup(brokerPort, "mystream", "mygroup", token)

	ConsumeMessages(brokerPort, "mystream", "mygroup", "myconsumer", token)
}

func main() {
	StartConsumer()
}
