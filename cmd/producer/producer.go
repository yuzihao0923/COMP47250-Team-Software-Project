package producer

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/network"
	"fmt"
	"net"
)

/*
*

	Send messages to MQ

*
*/

func StartProducer() {

	log.LogMessage("INFO", "Starting producer...")

	// Connect to the Q
	conn, err := net.Dial("tcp", "localhost:8889")

	if err != nil {
		fmt.Println(err)
		return
	}

	// Close conn
	defer conn.Close()

	log.LogMessage("INFO", "Producer has connected to the broker")

	tr := &network.Transport{
		Conn: conn,
	}

	// Send messages to the broker
	// change ur times
	// 目前的版本情况只可以发送 1 次
	for i := 0; i < 1; i++ {
		mes := message.Message{
			Payload: []byte(fmt.Sprintf("Hello %d", i)),
		}
		err := tr.SendMessage(conn, mes)
		if err != nil {
			log.LogMessage("ERROR", fmt.Sprintf("Producer has error sending message %d: %v", i, err))
			return
		}
		log.LogMessage("INFO", fmt.Sprintf("Producer sent message %d: %s", i, mes.Payload))
	}

}
