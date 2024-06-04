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

// SendMessage serializes the message and writes it to the connection
func (tr *Transport) SendMessage(mes message.Message) error {
	serializer := &serializer.JSONSerializer{}
	data, err := serializer.Serialize(mes)
	if err != nil {
		fmt.Println("Error serializing message:", err)
		return err
	}


	_, err = tr.Conn.Write(append(data, '\n'))

	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}
	return nil

}

/*
*

	1.create a reader to read from conn
	2.read conn
	3.deserialize

*
*/

// ReceiveMessage reads from the connection and deserializes the message

func (tr *Transport) ReceiveMessage() (message.Message, error) {
	reader := bufio.NewReader(tr.Conn)
	Data, err := reader.ReadBytes('\n')
	if err != nil {

		if err.Error() == "EOF" {
			return message.Message{}, err
		}

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
