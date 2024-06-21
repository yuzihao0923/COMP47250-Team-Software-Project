package main

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"fmt"
	"os"
	"time"
)

// SendMessage: send a new message to a stream (with streamName)
func SendMessage(brokerPort, streamName string, payload []byte) {
	msg := message.Message{
		Type: "produce",
		ConsumerInfo: &message.ConsumerInfo{
			StreamName: streamName,
		},
		Payload: payload,
	}
	err := api.SendMessage(brokerPort, msg)
	if err != nil {
		log.LogError("Producer", "producer has error sending message: "+err.Error())
		return
	}
	log.LogInfo("Producer", fmt.Sprintf("Producer sent message: %s", msg.Payload))
}

func StartProducer() {
	log.LogInfo("Producer", "Starting producer...")

	brokerPort := os.Getenv("BROKER_PORT")
	if brokerPort == "" {
		brokerPort = "8080" // default port
	}

	// // payloads for test
	// payloads := [][]byte{
	// 	[]byte("Hello 0"),
	// 	[]byte("Hello 1"),
	// 	[]byte("Hello 2"),
	// }

	// // send all payload to the stream (with streamName)
	// for _, payload := range payloads {
	// 	SendMessage(brokerPort, "mystream", payload)
	// }
	// time.Sleep(time.Millisecond)

	// send 1000 payloads to the stream (with streamName)
	for i := 0; i < 10; i++ {
		payload := []byte(fmt.Sprintf("Hello %d", i))
		SendMessage(brokerPort, "mystream", payload)
		time.Sleep(time.Millisecond) // slight delay to prevent overwhelming the broker
	}
}

func main() {
	StartProducer()
}
