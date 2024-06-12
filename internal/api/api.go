package api

import (
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/redis"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

var (
	rsi redis.RedisServiceInfo
	initRedisOnce sync.Once
)

func initRedis() {
	rsi = redis.RedisServiceInfo{
		StreamName: "mystream",
		GroupName:  "",
	}
}

// Broker API handlers

func HandleProduce(w http.ResponseWriter, r *http.Request) {
	initRedisOnce.Do(initRedis)

	var mes message.Message
	err := json.NewDecoder(r.Body).Decode(&mes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = rsi.WriteToStream(mes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	initRedisOnce.Do(initRedis)

	var mes message.Message
	err := json.NewDecoder(r.Body).Decode(&mes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rsi.GroupName = string(mes.Payload)

	err = rsi.CreateConsumerGroup()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func HandleConsume(w http.ResponseWriter, r *http.Request) {
	initRedisOnce.Do(initRedis)

	ctx := r.Context()
	consumerName := r.URL.Query().Get("consumer")

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

func SendMessage(payload string) error {
	mes := message.Message{
		Type:    "produce",
		Payload: []byte(payload),
	}

	data, err := json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("error marshaling message: %v", err)
	}

	resp, err := http.Post("http://localhost:8889/produce", "application/json", bytes.NewBuffer(data))
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

func RegisterConsumer(group string) error {
	mes := message.Message{
		Type:    "registration",
		Payload: []byte(group),
	}

	data, err := json.Marshal(mes)
	if err != nil {
		return fmt.Errorf("error marshaling registration message: %v", err)
	}

	resp, err := http.Post("http://localhost:8889/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error sending registration message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register consumer, status code: %d", resp.StatusCode)
	}

	return nil
}

func ConsumeMessages(consumerName string) ([]message.Message, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8889/consume?consumer=%s", consumerName))
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
