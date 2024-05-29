package consumer

import (
	"COMP47250-Team-Software-Project/internal/network"
	"fmt"
	"net"
)

func StartConsumer() {
	fmt.Println("Hi, I am consumer!!!")

	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("Consumer has error connecting to broker:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Consumer has connected to the broker")

	tr := &network.Transport{
		Conn: conn,
	}

	for {
		mes, err := tr.ReceiveMessage(conn)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Broker closed the connection")
				return
			}
			fmt.Println("Consumer has error receiving message:", err)
			return
		}

		fmt.Println("Consumer received message:", string(mes.Payload))
	}
}
