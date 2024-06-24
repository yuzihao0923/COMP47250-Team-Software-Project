package client

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// RegisterConsumer: Send request of registering consumer to API
func RegisterConsumer(brokerPort string, msg message.Message, token string) error {
	client := GetClientWithToken(token)

	data, err := serializer.JSONSerializerInstance.Serialize(msg)
	if err != nil {
		return fmt.Errorf("error serializing registration message: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%s/register", brokerPort), bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating registration request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending registration message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register consumer, status code: %d", resp.StatusCode)
	}

	return nil
}

func ConsumeMessages(brokerPort, streamName, groupName, consumerUsername, token string) ([]message.Message, error) {
	client := GetClientWithToken(token)

	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%s/consume?stream=%s&group=%s&consumer=%s", brokerPort, streamName, groupName, consumerUsername), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating consume request: %v", err)
	}

	resp, err := client.Do(req)
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

	// log.LogInfo("Consumer", fmt.Sprintf("Messages consumed from broker: %d messages", len(messages)))
	if len(messages) == 0 {
		log.LogWarning("Consumer", "No new message now, please wait.")
	}
	return messages, nil
}

// SendACK: consumer send ack to broker
func SendACK(brokerPort string, msg message.Message, token string) error {
	client := GetClientWithToken(token)

	data, err := serializer.JSONSerializerInstance.Serialize(msg)
	if err != nil {
		return fmt.Errorf("error serializing message: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%s/ack", brokerPort), bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating ACK request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending ACK: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send ACK, status code: %d", resp.StatusCode)
	}

	return nil
}
