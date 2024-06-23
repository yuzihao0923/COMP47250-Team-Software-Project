package api

import (
	"COMP47250-Team-Software-Project/internal/auth"
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/redis"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"fmt"
	"net/http"
)

// HandleACK: Handle Consumers' ACK.
func HandleACK(r *http.Request) HandlerResult {
	consumerUsername := r.Context().Value(auth.UsernameKey).(string)
	var msg message.Message
	err := serializer.JSONSerializerInstance.DeserializeFromReader(r.Body, &msg)
	if err != nil {
		return HandlerResult{Error: fmt.Errorf("failed to deserialize message: %v", err)}
	}

	rsi := redis.RedisServiceInfo{
		StreamName: msg.ConsumerInfo.StreamName,
		GroupName:  msg.ConsumerInfo.GroupName,
	}

	msgID := msg.ID
	ctx := r.Context()

	err = rsi.XACK(ctx, msgID, consumerUsername)
	if err != nil {
		return HandlerResult{Error: fmt.Errorf("failed to acknowledge message: %v", err)}
	}

	return HandlerResult{Data: "Acknowledged successfully"}
}
