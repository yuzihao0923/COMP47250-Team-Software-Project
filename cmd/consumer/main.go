package main

import (
	"COMP47250-Team-Software-Project/internal/network"
	"fmt"
	"net"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("Error connecting to broker:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to broker")

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
			fmt.Println("Error receiving message:", err)
			return
		}

		fmt.Println("Received message:", string(mes.Payload))
	}
}


