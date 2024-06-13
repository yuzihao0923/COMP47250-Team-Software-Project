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

	rsi := redis.RedisServiceInfo{
		StreamName: streamName,
		GroupName:  string(msg.Payload),
	}
	err = rsi.CreateConsumerGroup()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

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

// Producer API client functions

func SendMessage(streamName string, msg message.Message) error {

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://localhost:8889/produce?stream=%s", streamName), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}

// Consumer API client functions

func RegisterConsumer(streamName, group string) error {
	msg := message.Message{
		Type:    "registration",
		Payload: []byte(group),
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling registration message: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://localhost:8889/register?stream=%s", streamName), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error sending registration message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register consumer, status code: %d", resp.StatusCode)
	}

	return nil
}

func ConsumeMessages(streamName, groupName, consumerName string) ([]message.Message, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8889/consume?stream=%s&group=%s&consumer=%s", streamName, groupName, consumerName))
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
