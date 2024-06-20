package api

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/redis"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"bytes"
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

// SendMessage: Send message to broker
func SendMessage(brokerPort string, msg message.Message) error {

	data, err := serializer.JSONSerializerInstance.Serialize(msg)
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