package main

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"COMP47250-Team-Software-Project/internal/client"
	"COMP47250-Team-Software-Project/internal/database"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"fmt"
	"time"
)

var proxyURL = "http://localhost:8888"

var username string
var password string

// RegisterConsumerGroup: use API to register a consumer group (with groupName) in a stream (with streamName)
func RegisterConsumerGroup(brokerAddr, streamName, groupName, token string) {
	msg := message.Message{
		Type: "registration",
		ConsumerInfo: &message.ConsumerInfo{
			StreamName: streamName,
			GroupName:  groupName,
		},
	}
	err := client.RegisterConsumer(brokerAddr, msg, token)
	if err != nil {
		log.LogError("Consumer", "consumer has error registering: "+err.Error())
		return
	}
	log.LogInfo("Consumer", "Consumer registered to Broker...")
}

func ConsumeMessages(brokerAddr, streamName, groupName, consumerUsername, token string) {
	for {
		messages, err := client.ConsumeMessages(brokerAddr, streamName, groupName, consumerUsername, token)
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
			log.LogInfo("Consumer", "Consumer received message: "+string(msg.Payload))

			AcknowledgeMessage(brokerAddr, msg, token)
		}

		// time.Sleep(time.Second * 1)
	}
}

func AcknowledgeMessage(brokerAddr string, msg message.Message, token string) {
	err := client.SendACK(brokerAddr, msg, token)
	if err != nil {
		log.LogError("Consumer", "consumer has error sending ACK: "+err.Error())
		return
	}
	log.LogInfo("Consumer", "Consumer sending ACK successfully...")
}

func main() {
	fmt.Println("[INFO] [Consumer] Starting consumer...")

	broker, err := client.GetBroker(proxyURL)
	if err != nil {
		log.LogError("Consumer", fmt.Sprintf("Get broker failed, error: %v", err))
		return
	}

	brokerAddr := broker.Address

	err = database.ConnectMongoDB()
	if err != nil {
		fmt.Println("[ERROR] [Consumer] Failed to connect to database:", err)
		return
	}
	fmt.Println("[INFO] [Consumer] Database connected successfully")

	var token, role string
	for {
		username = auth.GetUserInput("\nEnter username: ")
		password = auth.GetPasswordInput("Enter password: ")

		token, role, err = auth.AuthenticateUser(username, password, brokerAddr)
		if err != nil {
			fmt.Println(err)
		} else if role != "consumer" {
			fmt.Println("this user is not a consumer, please try again")
		} else {
			// successfully login
			break
		}
	}

	// Register consumer group
	RegisterConsumerGroup(brokerAddr, "mystream", "mygroup", token)

	ConsumeMessages(brokerAddr, "mystream", "mygroup", username, token)
}
