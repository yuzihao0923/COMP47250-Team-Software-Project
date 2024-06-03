package message

import (
	"fmt"
)

type Message struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Payload   []byte `json:"payload"`
	Timestamp string `json:"timestamp"`
}

func (m Message) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"ID":        m.ID,
		"Type":      m.Type,
		"Content":   m.Payload,
		"Timestamp": m.Timestamp, // Redis prefers strings or numbers
	}
}

// NewMessageFromMap creates a Message struct from a map
func NewMessageFromMap(data map[string]interface{}) (*Message, error) {
	msg := &Message{}

	// Perform security checks before using type assertions
	if id, ok := data["ID"].(string); ok {
		msg.ID = id
	} else {
		return nil, fmt.Errorf("ID missing or not a string")
	}

	if typ, ok := data["Type"].(string); ok {
		msg.Type = typ
	} else {
		return nil, fmt.Errorf("type missing or not a string")
	}

	if timestamp, ok := data["Timestamp"].(string); ok {
		msg.Timestamp = timestamp
	} else {
		msg.Timestamp = ""
	}

	if payload, ok := data["Content"].(string); ok {
		msg.Payload = []byte(payload)
	} else {
		return nil, fmt.Errorf("content missing or not a string")
	}

	return msg, nil
}

// func NewMessageFromMap(data map[string]interface{}) (*Message, error) {
// 	msg := &Message{
// 		// ID:        data["id"].(string),
// 		Type:      data["type"].(string),
// 		Timestamp: data["timestamp"].(string),
// 	}
// 	// Assuming Payload is stored as string and needs conversion to []byte
// 	if payload, ok := data["payload"].(string); ok {
// 		msg.Payload = []byte(payload)
// 	}
// 	return msg, nil
// }
