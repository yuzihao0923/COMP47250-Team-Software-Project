package consumer

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"fmt"
	"time"
)

// use api to register a consumer group
func RegisterConsumerGroup(streamName, groupName string) {
	err := api.RegisterConsumer(streamName, groupName)
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Consumer has error registering: %v", err))
		return
	}
}

// consumer(consumerID) gets messages from a group of a stream
func ConsumeMessages(consumerID, streamName, groupName string) {
	msg := message.Message{
		Type: "consumer",
		ConsumerInfo: &message.ConsumerInfo{
			ConsumerID: "myconsumer",
			StreamName: "mystream",
			GroupName:  "mygroup",
		},
		Payload: []byte(""),
	}
	for {
		messages, err := api.ConsumeMessages(msg.ConsumerInfo.StreamName, msg.ConsumerInfo.GroupName, msg.ConsumerInfo.ConsumerID)
		if err != nil {
			log.LogMessage("ERROR", fmt.Sprintf("Consumer has error receiving message: %v", err))
			return
		}

		for _, mes := range messages {
			time.Sleep(time.Millisecond) // insure the order of log between "producer send" & "consumer receive"
			log.LogMessage("INFO", "Consumer received message: "+string(mes.Payload))
		}

		time.Sleep(time.Second * 1)
	}
}

func StartConsumer() {
	log.LogMessage("INFO", "Starting consumer...")

	// register consumer group
	RegisterConsumerGroup("mystream", "mygroup")

	// get messages
	ConsumeMessages("myconsumer", "mystream", "mygroup")
}
