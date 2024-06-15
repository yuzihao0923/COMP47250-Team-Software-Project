package api

import (
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/redis"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HandleProduce: Handle the request of produecer sending message
func HandleProduce(w http.ResponseWriter, r *http.Request) {
	streamName := r.URL.Query().Get("stream")
	if streamName == "" {
		http.Error(w, "Stream name is required", http.StatusBadRequest)
		return
	}

	var msg message.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if msg.ConsumerInfo == nil || msg.ConsumerInfo.StreamName == "" {
		msg.ConsumerInfo = &message.ConsumerInfo{
			StreamName: streamName,
		}
	}

	rsi := redis.RedisServiceInfo{
		StreamName: streamName,
	}
	err = rsi.WriteToStream(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleRegister: Handle the register request of consumer group
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	streamName := r.URL.Query().Get("stream")
	if streamName == "" {
		http.Error(w, "Stream name is required", http.StatusBadRequest)
		return
	}

	var msg message.Message
	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if msg.ConsumerInfo == nil {
		http.Error(w, "Consumer info is required", http.StatusBadRequest)
		return
	}

	rsi := redis.RedisServiceInfo{
		StreamName: streamName,
		GroupName:  msg.ConsumerInfo.GroupName,
	}
	err = rsi.CreateConsumerGroup()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleConsume: Handle comsuners' request to consume message
func HandleConsume(w http.ResponseWriter, r *http.Request) {
	streamName := r.URL.Query().Get("stream")
	if streamName == "" {
		http.Error(w, "Stream name is required", http.StatusBadRequest)
		return
	}

	groupName := r.URL.Query().Get("group")
	if groupName == "" {
		http.Error(w, "Group name is required", http.StatusBadRequest)
		return
	}

	consumerName := r.URL.Query().Get("consumer")
	if consumerName == "" {
		http.Error(w, "Consumer name is required", http.StatusBadRequest)
		return
	}

	rsi := redis.RedisServiceInfo{
		StreamName: streamName,
		GroupName:  groupName,
	}

	ctx := r.Context()
	streams, err := rsi.ReadFromStream(ctx, consumerName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var messages []message.Message
	for _, stream := range streams {
		for _, mes := range stream.Messages {
			m, err := message.NewMessageFromMap(mes.Values)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			messages = append(messages, *m)
		}
	}

	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// SendMessage: Send message to broker
func SendMessage(brokerPort string, msg message.Message) error {
	streamName := msg.ConsumerInfo.StreamName
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://localhost:%s/produce?stream=%s", brokerPort, streamName), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}

// RegisterConsumer: Send request of registering consumer to API
func RegisterConsumer(brokerPort, streamName, group string) error {
	msg := message.Message{
		Type: "registration",
		ConsumerInfo: &message.ConsumerInfo{
			StreamName: streamName,
			GroupName:  group,
		},
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling registration message: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://localhost:%s/register?stream=%s", brokerPort, streamName), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error sending registration message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register consumer, status code: %d", resp.StatusCode)
	}

	return nil
}

// ConsumeMessages: Send message request to API
func ConsumeMessages(brokerPort, streamName, groupName, consumerName string) ([]message.Message, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/consume?stream=%s&group=%s&consumer=%s", brokerPort, streamName, groupName, consumerName))
	if err != nil {
		return nil, fmt.Errorf("error receiving messages: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to receive messages, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var messages []message.Message
	err = json.Unmarshal(body, &messages)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %v", err)
	}

	return messages, nil
}
