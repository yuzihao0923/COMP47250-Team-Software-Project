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

// HandleACK: Handle Consumers' ACK.
func HandleACK(w http.ResponseWriter, r *http.Request) {

	var msg message.Message
	err := serializer.JSONSerializerInstance.DeserializeFromReader(r.Body, &msg)
	if err != nil {
		log.WriteErrorResponse(w, http.StatusBadRequest, err)
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
		log.WriteErrorResponse(w, http.StatusInternalServerError, err)
	}

}

// SendACK: consumer send ack to broker
func SendACK(brokerPort string, msg message.Message) error {
	data, err := serializer.JSONSerializerInstance.Serialize(msg)
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