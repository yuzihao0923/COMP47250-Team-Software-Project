package network

import (
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/pkg/serializer"
	"bufio"
	"fmt"
	"net"
)

type Transport struct {
	Conn net.Conn
}

/*
*

	1.serialize message
	2.write to conn

*
*/
func (tr *Transport) SendMessage(conn net.Conn, mes message.Message) error {
	serializer := &serializer.JSONSerializer{}
	Data, err := serializer.Serialize(mes)
	if err != nil {
		fmt.Println("Error serializing message:", err)
		return err
	}

	_, err = conn.Write(append(Data, '\n'))
	if err != nil {
		fmt.Println("Error sending message:", err)
	}
	return err
}

/*
*

	1.create a reader to read from conn
	2.read conn
	3.deserialize

*
*/
func (tr *Transport) ReceiveMessage(conn net.Conn) (message.Message, error) {
	reader := bufio.NewReader(conn)
	Data, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("Error receiving message:", err)
		return message.Message{}, err
	}

	deserializer := &serializer.JSONSerializer{}
	var mes message.Message

	// Use the pointer of mes!!!!!
	err = deserializer.Deserialize(Data, &mes)
	if err != nil {
		fmt.Println("Error deserialzing message:", err)
		return message.Message{}, err
	}
	return mes, err
}
