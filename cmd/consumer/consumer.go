package consumer

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"fmt"
	"time"
)

func StartConsumer() {
	log.LogMessage("INFO", "Starting consumer...")

	streamName := "mystream"
	groupName := "mygroup"

	err := api.RegisterConsumer(streamName, groupName)
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Consumer has error registering: %v", err))
		return
	}

	for {
		messages, err := api.ConsumeMessages(streamName, groupName, "myconsumer")
		if err != nil {
			log.LogMessage("ERROR", fmt.Sprintf("Consumer has error receiving message: %v", err))
			return
		}

		for _, mes := range messages {
			time.Sleep(time.Millisecond) // insure the order of log between "producer send" & "consumer receive"
			log.LogMessage("INFO", "Consumer received message: "+string(mes.Payload))
		}

		time.Sleep(time.Second * 2)
	}
}
