package consumer

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/api"
	"fmt"
	"time"
)

func StartConsumer() {
	log.LogMessage("INFO", "Starting consumer...")

	err := api.RegisterConsumer("mygroup")
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Consumer has error registering: %v", err))
		return
	}

	for {
		messages, err := api.ConsumeMessages("myconsumer")
		if err != nil {
			log.LogMessage("ERROR", fmt.Sprintf("Consumer has error receiving message: %v", err))
			return
		}

		for _, mes := range messages {
			log.LogMessage("INFO", "Consumer received message: "+string(mes.Payload))
		}

		time.Sleep(time.Second * 2)
	}
}
