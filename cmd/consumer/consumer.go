package consumer

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"fmt"
	"time"
)

func StartConsumer() {
	log.LogMessage("INFO", "Starting consumer...")

	// streamName := "mystream"
	// groupName := "mygroup"
	mes := message.Message{
		Type: "consumer",
		ConsumerInfo: &message.ConsumerInfo{
			ConsumerID: "myconsumer",
			StreamName: "mystream",
			GroupName:  "mygroup",
		},
		Payload: []byte(""),
	}
	// err := api.RegisterConsumer(streamName, groupName)
	err := api.RegisterConsumer(mes)
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Consumer has error registering: %v", err))
		return
	}

	for {
		messages, err := api.ConsumeMessages(mes.ConsumerInfo.StreamName, mes.ConsumerInfo.GroupName, mes.ConsumerInfo.ConsumerID)
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
