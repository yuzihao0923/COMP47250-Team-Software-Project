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

// HandleRegister: Handle the register request of consumer group
func HandleRegister(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	if ctx.Err() != nil {
		log.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("request canceled or the client closed the connection"))
		return
	}

	var msg message.Message
	err := serializer.JSONSerializerInstance.DeserializeFromReader(r.Body, &msg)
	if err != nil {
		fmt.Println(err.Error())
		log.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if msg.ConsumerInfo == nil {
		log.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("consumer info is required"))
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
			log.WriteErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// HandleConsume: Handle consumers' request to consume message
func HandleConsume(w http.ResponseWriter, r *http.Request) {
	streamName := r.URL.Query().Get("stream")
	if streamName == "" {
		log.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("stream name is required"))
		return
	}

	groupName := r.URL.Query().Get("group")
	if groupName == "" {
		log.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("group name is required"))
		return
	}

	consumerName := r.URL.Query().Get("consumer")
	if consumerName == "" {
		log.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("consumer name is required"))
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
			log.WriteErrorResponse(w, http.StatusInternalServerError, err)
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
				log.WriteErrorResponse(w, http.StatusInternalServerError, err)
				return
			}
			messages = append(messages, *m)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	err = serializer.JSONSerializerInstance.SerializeToWriter(messages, w)
	if err != nil {
		log.WriteErrorResponse(w, http.StatusInternalServerError, err)
	}
}

// RegisterConsumer: Send request of registering consumer to API
func RegisterConsumer(brokerPort string, msg message.Message) error {

	data, err := serializer.JSONSerializerInstance.Serialize(msg)
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

	return nil
}

func ConsumeMessages(brokerPort, streamName, groupName, consumerName string) ([]message.Message, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/consume?stream=%s&group=%s&consumer=%s", brokerPort, streamName, groupName, consumerName))
	// resp, err := http.Get(fmt.Sprintf("http://broker:%s/consume?stream=%s&group=%s&consumer=%s", brokerPort, streamName, groupName, consumerName))

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
	err = serializer.JSONSerializerInstance.Deserialize(body, &messages)
	if err != nil {
		return nil, fmt.Errorf("error deserializing response body: %v", err)
	}

	log.LogInfo(fmt.Sprintf("Messages consumed from broker: %d messages", len(messages)))
	return messages, nil
}