package api

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/redis"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var jsonSerializer = &serializer.JSONSerializer{}

func writeErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	log.LogError(err)
	http.Error(w, err.Error(), statusCode)
}

// HandleProduce: Handle the request of producer sending message
func HandleProduce(w http.ResponseWriter, r *http.Request) {

	var msg message.Message
	err := jsonSerializer.DeserializeFromReader(r.Body, &msg)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if msg.ConsumerInfo.StreamName == "" {
		writeErrorResponse(w, http.StatusBadRequest, fmt.Errorf("stream name is required"))
		return
	}

	rsi := redis.RedisServiceInfo{
		StreamName: msg.ConsumerInfo.StreamName,
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
		StreamName: msg.ConsumerInfo.StreamName,
		GroupName:  msg.ConsumerInfo.GroupName,
	}

	err = rsi.CreateConsumerGroup()
	if err != nil {
		if strings.Contains(err.Error(), "Consumer Group name already exists") {
			// log.LogWarning("Consumer Group name already exists")
			w.WriteHeader(http.StatusOK) // Return OK status to not block the process
			return
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
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
		if errors.Is(err, context.Canceled) {
			log.LogWarning(fmt.Sprintf("Consumer exited from stream '%s': %v", rsi.StreamName, err))
			return
		} else {
			writeErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	if len(streams) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	var messages []message.Message
	for _, stream := range streams {

		for _, mes := range stream.Messages {
			m, err := message.NewMessageFromMap(mes.Values, mes.ID)
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

// HandleACK: Handle Consumers' ACK.
func HandleACK(w http.ResponseWriter, r *http.Request) {

	var msg message.Message
	err := jsonSerializer.DeserializeFromReader(r.Body, &msg)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	rsi := redis.RedisServiceInfo{
		StreamName: msg.ConsumerInfo.StreamName,
		GroupName:  msg.ConsumerInfo.GroupName,
	}

	msgID := msg.ID

	ctx := r.Context()

	err = rsi.XACK(ctx, msgID)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err)
	}

}

// SendACK: consumer send ack to broker
func SendACK(brokerPort string, msg message.Message) error {
	data, err := jsonSerializer.Serialize(msg)
	if err != nil {
		return fmt.Errorf("error serializing message: %v", err)
	}
	resp, err := http.Post(fmt.Sprintf("http://localhost:%s/ack", brokerPort), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error sending ACK: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send ACK, status code: %d", resp.StatusCode)
	}

	return nil
}

// SendMessage: Send message to broker
func SendMessage(brokerPort string, msg message.Message) error {

	data, err := jsonSerializer.Serialize(msg)
	if err != nil {
		return fmt.Errorf("error serializing message: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://localhost:%s/produce", brokerPort), "application/json", bytes.NewBuffer(data))
	//resp, err := http.Post(fmt.Sprintf("http://broker:%s/produce?stream=%s", brokerPort, streamName), "application/json", bytes.NewBuffer(data))

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
func RegisterConsumer(brokerPort string, msg message.Message) error {

	data, err := jsonSerializer.Serialize(msg)
	if err != nil {
		return fmt.Errorf("error serializing registration message: %v", err)
	}

	resp, err := http.Post(fmt.Sprintf("http://localhost:%s/register", brokerPort), "application/json", bytes.NewBuffer(data))
	//resp, err := http.Post(fmt.Sprintf("http://broker:%s/register?stream=%s", brokerPort, streamName), "application/json", bytes.NewBuffer(data))

	if err != nil {
		return fmt.Errorf("error sending registration message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register consumer, status code: %d", resp.StatusCode)
	}

	log.LogInfo(fmt.Sprintf("Consumer group '%s' registered for stream: %s", msg.ConsumerInfo.GroupName, msg.ConsumerInfo.StreamName))
	return nil
}

// ConsumeMessages: Send message request to API
func ConsumeMessages(brokerPort, streamName, groupName, consumerName string) ([]message.Message, error) {
	//resp, err := http.Get(fmt.Sprintf("http://broker:%s/consume?stream=%s&group=%s&consumer=%s", brokerPort, streamName, groupName, consumerName))
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

	log.LogInfo(fmt.Sprintf("Messages consumed from broker: %d messages", len(messages)))
	return messages, nil
}
