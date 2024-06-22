package api

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/redis"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"fmt"
	"net/http"
)

// HandleProduce: Handle the request of producer sending message
func HandleProduce(w http.ResponseWriter, r *http.Request) {
	var msg message.Message
	err := serializer.JSONSerializerInstance.DeserializeFromReader(r.Body, &msg)
	if err != nil {
		log.WriteErrorResponse(w, http.StatusBadRequest, err)
		return
	}

	if msg.ConsumerInfo.StreamName == "" {
		log.WriteErrorResponse(w, http.StatusBadRequest, fmt.Errorf("stream name is required"))
		return
	}

	rsi := redis.RedisServiceInfo{
		StreamName: msg.ConsumerInfo.StreamName,
	}
	err = rsi.WriteToStream(msg)
	if err != nil {
		log.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
