package message

import (
	"fmt"
)

type Message struct {
	ID           string        `json:"id,omitempty"` // message id
	Type         string        `json:"type"`         // Used to determine producer or consumer
	ConsumerInfo *ConsumerInfo `json:"consumer_info,omitempty"`
	Payload      []byte        `json:"payload,omitempty"`
	// Timestamp string `json:"timestamp"`
}

type ConsumerInfo struct {
	ConsumerUsername string `json:"consumer_username"` // It can be used for identifing different consumers in the same group
	StreamName       string `json:"stream_name"`
	GroupName        string `json:"group_name"`
}

func (m Message) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"Type": m.Type,
	}

	if m.ConsumerInfo != nil {
		if m.ConsumerInfo.ConsumerUsername != "" {
			result["ConsumerUsername"] = m.ConsumerInfo.ConsumerUsername
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
func NewMessageFromMap(data map[string]interface{}, mesID string) (*Message, error) {
	msg := &Message{}
	msg.ID = mesID

	// Required: Type
	typ, ok := data["Type"].(string)
	if !ok {
		return nil, fmt.Errorf("type missing or not a string")
	}
	msg.Type = typ

	// Optional: Payload (Content)
	if content, ok := data["Content"].(string); ok {
		msg.Payload = []byte(content)
	}

	// Optional: ConsumerInfo
	consumerInfo := ConsumerInfo{}
	anyConsumerInfo := false // Flag to check if any consumer info is provided

	if ConsumerUsername, ok := data["ConsumerUsername"].(string); ok && ConsumerUsername != "" {
		consumerInfo.ConsumerUsername = ConsumerUsername
		anyConsumerInfo = true
	}
	if streamName, ok := data["StreamName"].(string); ok && streamName != "" {
		consumerInfo.StreamName = streamName
		anyConsumerInfo = true
	}
	if groupName, ok := data["GroupName"].(string); ok && groupName != "" {
		consumerInfo.GroupName = groupName
		anyConsumerInfo = true
	}

	// Attach consumerInfo to msg if any consumerInfo field was provided
	if anyConsumerInfo {
		msg.ConsumerInfo = &consumerInfo
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
