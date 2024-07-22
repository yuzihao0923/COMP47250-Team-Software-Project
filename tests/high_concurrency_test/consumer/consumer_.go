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

var proxyURL = "http://localhost:8888"

// User 结构用于保存用户信息
type User struct {
	Username string
	Password string
	Token    string
	Role     string
}

// RegisterConsumerGroup 注册消费者组
func RegisterConsumerGroup(brokerAddr, streamName, groupName, token string) error {
	msg := message.Message{
		Type: "registration",
		ConsumerInfo: &message.ConsumerInfo{
			StreamName: streamName,
			GroupName:  groupName,
		},
	}
	return client.RegisterConsumer(brokerAddr, msg, token)
}

// ConsumeMessages 消费消息
func ConsumeMessages(brokerAddr, streamName, groupName string, user User, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		messages, err := client.ConsumeMessages(brokerAddr, streamName, groupName, user.Username, user.Token)
		if err != nil {
			if err.Error() == "no new messages" {
				log.LogWarning("Consumer", "No new messages, retrying...")
			} else {
				log.LogError("Consumer", "consumer has error receiving message: "+err.Error())
			}
			time.Sleep(time.Second * 1)
			continue
		}

		for _, msg := range messages {
			log.LogInfo("Consumer", fmt.Sprintf("Consumer %s received message: %s", user.Username, msg.Payload))
			AcknowledgeMessage(brokerAddr, msg, user.Token)
		}
	}
}

// AcknowledgeMessage 发送消息确认
func AcknowledgeMessage(brokerAddr string, msg message.Message, token string) {
	if err := client.SendACK(brokerAddr, msg, token); err != nil {
		log.LogError("Consumer", "consumer has error sending ACK: "+err.Error())
	}
}

func main() {
	fmt.Println("[INFO] [Consumer] Starting consumer...")

	broker, err := client.GetBroker(proxyURL)
	if err != nil {
		log.LogError("Consumer", fmt.Sprintf("Get broker failed, error: %v", err))
		return
	}

	brokerAddr := broker.Address

	_, err = database.ConnectMongoDB("mongodb://localhost:27017", "comp47250", "users")
	if err != nil {
		fmt.Println("[ERROR] [Consumer] Failed to connect to database:", err)
		return
	}
	fmt.Println("[INFO] [Consumer] Database connected successfully")

	var wg sync.WaitGroup
	users := []User{
		{Username: "c1", Password: "123"},
		{Username: "c2", Password: "123"},
		{Username: "c3", Password: "123"},
		{Username: "c4", Password: "123"},
	}

	for i := range users {
		wg.Add(1)
		go func(u int) {
			user, err := authenticateUser(brokerAddr, users[u].Username, users[u].Password)
			if err != nil {
				fmt.Println("Authentication failed for", users[u].Username, ":", err)
				wg.Done()
				return
			}
			ConsumeMessages(brokerAddr, "mystream", "mygroup", user, &wg)
		}(i)
	}

	wg.Wait() // 等待所有goroutine完成
	fmt.Println("[INFO] [Consumer] All consumers have stopped.")
}

func authenticateUser(brokerAddr, username, password string) (User, error) {
	token, role, err := auth.AuthenticateUser(username, password, brokerAddr)
	if err != nil {
		return User{}, err
	}
	if role != "consumer" {
		return User{}, fmt.Errorf("this user is not a consumer, please try again")
	}
	return User{Username: username, Password: password, Token: token, Role: role}, nil
}
