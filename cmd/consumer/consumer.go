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

var username string
var password string

// RegisterConsumerGroup: use API to register a consumer group (with groupName) in a stream (with streamName)
func RegisterConsumerGroup(brokerPort, streamName, groupName, token string) {
	msg := message.Message{
		Type: "registration",
		ConsumerInfo: &message.ConsumerInfo{
			StreamName: streamName,
			GroupName:  groupName,
		},
	}
	err := client.RegisterConsumer(brokerPort, msg, token)
	if err != nil {
		log.LogError("Consumer", "consumer has error registering: "+err.Error())
		return
	}
	log.LogInfo("Consumer", "Consumer registered to Broker...")
}

func ConsumeMessages(brokerPort, streamName, groupName, consumerID, token string) {
	for {
		messages, err := client.ConsumeMessages(brokerPort, streamName, groupName, consumerID, token)
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
			time.Sleep(time.Millisecond) // Ensure the order of log between "producer send" & "consumer receive"
			log.LogInfo("Consumer", "Consumer received message: "+string(msg.Payload))

			AcknowledgeMessage(brokerPort, msg, token)
		}

		time.Sleep(time.Second * 1)
	}
}

func AcknowledgeMessage(brokerPort string, msg message.Message, token string) {
	err := client.SendACK(brokerPort, msg, token)
	if err != nil {
		log.LogError("Consumer", "consumer has error sending ACK: "+err.Error())
		return
	}
	log.LogInfo("Consumer", "Consumer sending ACK successfully...")
}

func main() {
	// Ensure logs are printed before prompting user input
	fmt.Println("[INFO] [Consumer] Starting consumer...")

	brokerPort := os.Getenv("BROKER_PORT")
	if brokerPort == "" {
		brokerPort = "8080" // Default port
	}

	err := database.ConnectMongoDB()
	if err != nil {
		fmt.Println("[ERROR] [Consumer] Failed to connect to database:", err)
		return
	}
	fmt.Println("[INFO] [Consumer] Database connected successfully")

	var token, role string
	for {
		username = auth.GetUserInput("\nEnter username: ")
		password = auth.GetPasswordInput("Enter password: ")

		token, role, err = auth.AuthenticateUser(username, password)
		if err != nil {
			fmt.Println(err)
		} else if role != "consumer" {
			fmt.Println("this user is not a consumer, please try again")
		} else {
			break
		}
	}

	// Register consumer group
	RegisterConsumerGroup(brokerPort, "mystream", "mygroup", token)

	ConsumeMessages(brokerPort, "mystream", "mygroup", username, token)
}
