package api

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/redis"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"fmt"
	"net/http"
)

// HandleProduce: Handle the request of producer sending message
func HandleProduce(r *http.Request) HandlerResult {
	producerUsername := r.Context().Value(auth.UsernameKey).(string)
	var msg message.Message
	err := serializer.JSONSerializerInstance.DeserializeFromReader(r.Body, &msg)
	if err != nil {
		return HandlerResult{nil, fmt.Errorf("failed to deserialize message: %v", err)}
	}

	if msg.ConsumerInfo.StreamName == "" {
		return HandlerResult{nil, fmt.Errorf("stream name is required")}
	}

	rsi := redis.RedisServiceInfo{
		StreamName: msg.ConsumerInfo.StreamName,
	}
	err = rsi.WriteToStream(msg, producerUsername)
	if err != nil {
		return HandlerResult{nil, fmt.Errorf("failed to write to stream: %v", err)}
	}

	return HandlerResult{Data: "Message produced successfully", Error: nil}
}
