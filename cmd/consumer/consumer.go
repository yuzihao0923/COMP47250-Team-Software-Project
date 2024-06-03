package consumer

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/network"
	"fmt"
	"net"
)

func StartConsumer() {

	log.LogMessage("INFO", "Starting conumer...")

	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Consumer has error connecting to broker: %v", err))
		return
	}
	defer conn.Close()

	log.LogMessage("INFO", "Consumer has connected to the broker")

	// Currently, do not consider re-connected situation!!!    Every new consumer connection will be treated as first time connection

	mes := message.Message{
		Type:    "registration",
		Payload: []byte("mygroup"), // If the Type is registration, Playload will be the consumer group name
	}

	tr := &network.Transport{
		Conn: conn,
	}

	error := tr.SendMessage(mes)
	if error != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Consumer has error sending registration message: %v", error))
	}

	for {
		mes, err := tr.ReceiveMessage()
		if err != nil {
			if err.Error() == "EOF" {
				log.LogMessage("INFO", "Broker closed the connection")
				return
			}
			log.LogMessage("ERROR", fmt.Sprintf("Consumer has error receiving message: %v", err))
			return
		}

		log.LogMessage("INFO", "Consumer received message: "+string(mes.Payload))
	}
}
