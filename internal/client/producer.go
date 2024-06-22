package client

import (
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"bytes"
	"fmt"
	"net/http"
)

// // SendMessage: Send message to broker
func SendMessage(brokerPort string, msg message.Message, token string) error {
	client := GetClientWithToken(token)

	data, err := serializer.JSONSerializerInstance.Serialize(msg)
	if err != nil {
		return fmt.Errorf("error serializing message: %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%s/produce", brokerPort), bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("error creating produce request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}
