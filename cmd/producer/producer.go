package producer

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/api"
	"fmt"
)

func StartProducer() {
	log.LogMessage("INFO", "Starting producer...")

	messages := []string{"Hello 0"}

	for _, msg := range messages {
		err := api.SendMessage(msg)
		if err != nil {
			log.LogMessage("ERROR", fmt.Sprintf("Producer has error sending message: %v", err))
			return
		}
		log.LogMessage("INFO", fmt.Sprintf("Producer sent message: %s", msg))
	}
}
