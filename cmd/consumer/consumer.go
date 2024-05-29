package consumer

import (
	"COMP47250-Team-Software-Project/internal/log"
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

	tr := &network.Transport{
		Conn: conn,
	}

	for {
		mes, err := tr.ReceiveMessage(conn)
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
