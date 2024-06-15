package message

import (
	"fmt"
)

type Message struct {
	Type         string        `json:"type"` // Used to determine producer or consumer
	ConsumerInfo *ConsumerInfo `json:"consumer_info,omitempty"`
	Payload      []byte        `json:"payload,omitempty"`
	// Timestamp string `json:"timestamp"`
}

type ConsumerInfo struct {
	ConsumerID string `json:"consumer_id"` // It can be used for identifing different consumers in the same group
	StreamName string `json:"stream_name"`
	GroupName  string `json:"group_name"`
}

func (m Message) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"Type": m.Type,
	}

	if m.ConsumerInfo != nil {
		if m.ConsumerInfo.ConsumerID != "" {
			result["ConsumerID"] = m.ConsumerInfo.ConsumerID
		}
		if m.ConsumerInfo.StreamName != "" {
			result["StreamName"] = m.ConsumerInfo.StreamName
		}
		if m.ConsumerInfo.GroupName != "" {
			result["GroupName"] = m.ConsumerInfo.GroupName
		}
	}

	if m.Payload != nil {
		result["Content"] = string(m.Payload)
	}

	return result
	// return map[string]interface{}{
	// 	"ID":        m.ID,
	// 	"Type":      m.Type,
	// 	"Content":   m.Payload,
	// 	"Timestamp": m.Timestamp, // Redis prefers strings or numbers
	// }
}

// NewMessageFromMap creates a Message struct from a map
func NewMessageFromMap(data map[string]interface{}) (*Message, error) {
	msg := &Message{}

	if typ, ok := data["Type"].(string); ok {
		msg.Type = typ
	} else {
		return nil, fmt.Errorf("Type missing or not a string")
	}

	if content, ok := data["Content"].(string); ok {
		msg.Payload = []byte(content)
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
