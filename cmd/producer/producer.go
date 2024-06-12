package producer

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"fmt"
)

func StartProducer() {
	log.LogMessage("INFO", "Starting producer...")

	// test self-defined stream name and multiple messages
	streamName := "mystream"
	messages := []string{"Hello 0", "Hello 1", "Hello 2"}

	for _, msg := range messages {
		err := api.SendMessage(streamName, msg)
		if err != nil {
			log.LogMessage("ERROR", fmt.Sprintf("Producer has error sending message: %v", err))
			return
		}
		log.LogMessage("INFO", fmt.Sprintf("Producer sent message: %s", msg))
	}
}
