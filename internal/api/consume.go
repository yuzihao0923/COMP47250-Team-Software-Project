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
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx.Err() != nil {
		log.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("request canceled or the client closed the connection"))
		return
	}

	var msg message.Message
	err := serializer.JSONSerializerInstance.DeserializeFromReader(r.Body, &msg)
	if err != nil {
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
			log.LogWarning("Consumer", fmt.Sprintf("Consumer exited from stream '%s': %v", rsi.StreamName, err))
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
