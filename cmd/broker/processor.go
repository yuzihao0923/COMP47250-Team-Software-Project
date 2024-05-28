package main

import (
	"COMP47250-Team-Software-Project/internal/network"
	"fmt"
	"net"
)

type Processor struct {
	Conn net.Conn
}

func (p *Processor) brokerProcessMes() (err error) {
	tr := &network.Transport{
		Conn: p.Conn,
	}

	mes, err := tr.ReceiveMessage(tr.Conn)
	if err != nil {
		fmt.Println("Broker can not receive message successfully!!")
		return
	}

	fmt.Println("Broker received a message: ", mes.Payload)
	// fmt.Println("The message detail:")
	// fmt.Println("ID: ", mes.ID)
	// fmt.Println("Timestamp: ", mes.Timestamp)
	// fmt.Println("Type: ", mes.Type)

	return
}
