package producer

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"fmt"
)

// SetMessage: set new message's stream name and payload
func SendMessage(streamName string, payload []byte){
	msg := message.Message{
		Type: "produce",
		ConsumerInfo: &message.ConsumerInfo{
			StreamName: streamName,
		},
		Payload: payload,
	}
	err := api.SendMessage(msg)
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Producer has error sending message: %v", err))
		return
	}
	log.LogMessage("INFO", fmt.Sprintf("Producer sent message: %s", msg.Payload))
}

func StartProducer() {
	log.LogMessage("INFO", "Starting producer...")

	// test
	payloads := [][]byte{
		[]byte("Hello 0"),
		[]byte("Hello 1"),
		[]byte("Hello 2"),
	}

	for _, payload := range payloads {
		SendMessage("mystream", payload)
	}
}
