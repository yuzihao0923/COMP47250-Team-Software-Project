package api

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/redis"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// HandleRegister: Handle the register request of consumer group
func HandleRegister(r *http.Request) HandlerResult {
	ctx := r.Context()
	if ctx.Err() != nil {
		return HandlerResult{Error: fmt.Errorf("request canceled or the client closed the connection")}
	}

	var msg message.Message
	err := serializer.JSONSerializerInstance.DeserializeFromReader(r.Body, &msg)
	if err != nil {
		return HandlerResult{nil, fmt.Errorf("failed to deserialize message: %v", err)}
	}

	if msg.ConsumerInfo == nil {
		return HandlerResult{Error: fmt.Errorf("consumer info is required")}
	}

	rsi := redis.RedisServiceInfo{
		StreamName: msg.ConsumerInfo.StreamName,
		GroupName:  msg.ConsumerInfo.GroupName,
	}

	err = rsi.CreateConsumerGroup()
	if err != nil {
		if strings.Contains(err.Error(), "Consumer Group name already exists") {
			return HandlerResult{Data: "Consumer Group name already exists"} // Return OK status to not block the process
		} else {
			return HandlerResult{Error: err}
		}
	}

	return HandlerResult{Data: "Consumer group created successfully"}
}

// HandleConsume: Handle consumers' request to consume message
func HandleConsume(r *http.Request) HandlerResult {
	streamName := r.URL.Query().Get("stream")
	if streamName == "" {
		return HandlerResult{Error: fmt.Errorf("stream name is required")}
	}

	groupName := r.URL.Query().Get("group")
	if groupName == "" {
		return HandlerResult{Error: fmt.Errorf("group name is required")}
	}

	consumerUsername := r.URL.Query().Get("consumer")
	if consumerUsername == "" {
		return HandlerResult{Error: fmt.Errorf("consumer username is required")}
	}

	rsi := redis.RedisServiceInfo{
		StreamName: streamName,
		GroupName:  groupName,
	}

	ctx := r.Context()
	streams, err := rsi.ReadFromStream(ctx, consumerUsername)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.LogWarning("Consumer", fmt.Sprintf("Consumer exited from stream '%s': %v", rsi.StreamName, err))
			return HandlerResult{Error: err}
		} else {
			return HandlerResult{Error: err}
		}
	}

	if len(streams) == 0 {
		return HandlerResult{Data: nil}
	}

	var messages []message.Message
	for _, stream := range streams {
		for _, mes := range stream.Messages {
			m, err := message.NewMessageFromMap(mes.Values, mes.ID)
			if err != nil {
				return HandlerResult{Error: err}
			}
			messages = append(messages, *m)
		}
	}

	return HandlerResult{Data: messages}
}
