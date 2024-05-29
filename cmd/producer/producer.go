package producer

import (
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/network"
	"fmt"
	"net"
	"time"
)

/*
*

	Send messages to MQ

*
*/

func StartProducer() {

	fmt.Println("Hi, I am producer!!!")

	// Connect to the Q
	conn, err := net.Dial("tcp", "localhost:8889")

	if err != nil {
		fmt.Println(err)
		return
	}

	// Close conn
	defer conn.Close()

	mes := message.Message{
		ID:        "0",
		Type:      "test",
		Payload:   []byte("Hello"),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	fmt.Println(mes)

	tr := &network.Transport{
		Conn: conn,
	}

	tr.SendMessage(conn, mes)

}
