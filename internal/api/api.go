package api

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/redis"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var jsonSerializer = &serializer.JSONSerializer{}

func writeErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	log.LogMessage("ERROR", err.Error())
	http.Error(w, err.Error(), statusCode)
}

// HandleProduce: Handle the request of producer sending message
func HandleProduce(w http.ResponseWriter, r *http.Request) {
	streamName := r.URL.Query().Get("stream")
	if streamName == "" {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("stream name is required"))
		return
	}

	var msg message.Message
	err := jsonSerializer.DeserializeFromReader(r.Body, &msg)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
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
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleRegister: Handle the register request of consumer group
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	streamName := r.URL.Query().Get("stream")
	fmt.Println(streamName)
	if streamName == "" {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("stream name is required"))
		return
	}

	ctx := r.Context()
	if ctx.Err() != nil {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("request canceled or the client closed the connection"))
		return
	}

	var msg message.Message
	err := jsonSerializer.DeserializeFromReader(r.Body, &msg)
	if err != nil {
		fmt.Println(err.Error())
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if msg.ConsumerInfo == nil {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("consumer info is required"))
		return
	}

	rsi := redis.RedisServiceInfo{
		StreamName: streamName,
		GroupName:  msg.ConsumerInfo.GroupName,
	}
	err = rsi.CreateConsumerGroup()
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Consumer group created successfully")
}

// HandleConsume: Handle consumers' request to consume message
func HandleConsume(w http.ResponseWriter, r *http.Request) {
	streamName := r.URL.Query().Get("stream")
	if streamName == "" {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("stream name is required"))
		return
	}

	groupName := r.URL.Query().Get("group")
	if groupName == "" {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("group name is required"))
		return
	}

	consumerName := r.URL.Query().Get("consumer")
	if consumerName == "" {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("consumer name is required"))
		return
	}

	rsi := redis.RedisServiceInfo{
		StreamName: streamName,
		GroupName:  groupName,
	}

	ctx := r.Context()
	streams, err := rsi.ReadFromStream(ctx, consumerName)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	if len(streams) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var messages []message.Message
	for _, stream := range streams {
		for _, mes := range stream.Messages {
			m, err := message.NewMessageFromMap(mes.Values)
			if err != nil {
				writeErrorResponse(w, http.StatusInternalServerError, err)
				return
			}
			messages = append(messages, *m)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = jsonSerializer.SerializeToWriter(messages, w)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err)
	}
}

// SendMessage: Send message to broker
func SendMessage(brokerPort string, msg message.Message) error {
	streamName := msg.ConsumerInfo.StreamName
	data, err := jsonSerializer.Serialize(msg)
	if err != nil {
		return fmt.Errorf("error serializing message: %v", err)
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

	data, err := jsonSerializer.Serialize(msg)
	if err != nil {
		return fmt.Errorf("error serializing registration message: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://localhost:%s/register?stream=%s", brokerPort, streamName), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error sending registration message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register consumer, status code: %d", resp.StatusCode)
	}

	log.LogMessage("INFO", fmt.Sprintf("Consumer group '%s' registered for stream: %s", group, streamName))
	return nil
}

// ConsumeMessages: Send message request to API
func ConsumeMessages(brokerPort, streamName, groupName, consumerName string) ([]message.Message, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/consume?stream=%s&group=%s&consumer=%s", brokerPort, streamName, groupName, consumerName))
	if err != nil {
		return nil, fmt.Errorf("error receiving messages: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, fmt.Errorf("no new messages")
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to receive messages, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	var messages []message.Message
	err = jsonSerializer.Deserialize(body, &messages)
	if err != nil {
		return nil, fmt.Errorf("error deserializing response body: %v", err)
	}

	return messages, nil
}
