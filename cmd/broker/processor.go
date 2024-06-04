package broker

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"COMP47250-Team-Software-Project/internal/network"
	"COMP47250-Team-Software-Project/internal/redis"
	"context"
	"fmt"
	"net"
	"time"
)

var rsi redis.RedisServiceInfo

type Processor struct {
	Conn net.Conn
}

func (p *Processor) brokerProcessMes() (err error) {
	tr := &network.Transport{
		Conn: p.Conn,
	}

	for {
		mes, err := tr.ReceiveMessage()
		if err != nil {
			if err.Error() == "EOF" {
				log.LogMessage("INFO", "Client closed the connection")
				// p.removeConsumer(p.Conn)
				return nil
			}
			// p.removeConsumer(p.Conn)
			return err
		}

		log.LogMessage("INFO", "Broker received a message: "+string(mes.Payload))

		// if mes.Type == "registration" {
		// 	handleConnection(mes.Payload)
		// }
		handleConnection(mes)

		p.brokerReadMesFromStream()

		// Broadcast the message to all consumers
		// consumersMutex.Lock()
		// var activeConsumers []net.Conn
		// for _, consumer := range consumers {
		// 	tr := &network.Transport{
		// 		Conn: consumer,
		// 	}
		// 	err := tr.SendMessage(consumer, mes)
		// 	if err != nil {
		// 		// fmt.Printf("Failed to send message to consumer: %v\n", err)
		// 		log.LogMessage("ERROR", fmt.Sprintf("Failed to send message to consumer: %v", err))
		// 		consumer.Close() // Close the connection if sending fails
		// 	} else {
		// 		activeConsumers = append(activeConsumers, consumer)
		// 	}
		// }
		// consumers = activeConsumers
		// consumersMutex.Unlock()
	}
}

// func (p *Processor) removeConsumer(conn net.Conn) {
// 	consumersMutex.Lock()
// 	defer consumersMutex.Unlock()

// 	for i, consumer := range consumers {
// 		if consumer == conn {
// 			consumers = append(consumers[:i], consumers[i+1:]...)
// 			break
// 		}
// 	}
// }

func handleConnection(mes message.Message) {

	// rsi := redis.RedisServiceInfo{
	// 	StreamName: "",
	// 	GroupName:  "",
	// }

	if mes.Type == "registration" {
		rsi.StreamName = "mystream"
		rsi.GroupName = string(mes.Payload)

		// make sure that consumer group exist
		err := rsi.CreateConsumerGroup() // Every consumer will create a new consumer group.
		if err != nil {
			log.LogMessage("ERROR", fmt.Sprintf("Failed to create or check consumer group: %v", err))
			return
		}
	} else {
		// if the message type is not 'registration', the message is from producer
		err := rsi.WriteToStream(mes)
		if err != nil {
			log.LogMessage("ERROR", fmt.Sprintf("Failed to write message to stream: %v", err))
			return
		}
	}

}

// ReadMessagesFromStream reads messages for a specific consumer group and processes them.
func (p *Processor) brokerReadMesFromStream() {
	ctx := context.Background()
	consumerName := "myconsumer"

	for {
		// read from stream
		streams, err := rsi.ReadFromStream(ctx, consumerName)
		if err != nil {
			log.LogMessage("ERROR", fmt.Sprintf("Failed to read from stream: %v", err))
			time.Sleep(time.Second * 10) // retry
			continue
		}

		// process messages from each returned stream
		for _, stream := range streams {
			for _, mes := range stream.Messages {
				m, err := message.NewMessageFromMap(mes.Values)
				if err != nil {
					log.LogMessage("ERROR", fmt.Sprintf("Failed to convert message from map[string]interface{} to mes: %v", err))
					return
				}

				p.brokerSendMesToConsumer(*m)

				// _, err := redisService.AcknowledgeMessage(ctx, stream.Stream, message.ID)
				// if err != nil {
				// 	log.LogMessage("ERROR", fmt.Sprintf("Failed to acknowledge message %s: %v", message.ID, err))
				// }
			}
		}

		// Adjust wait times or use appropriate spacing strategies to balance performance and resource consumption
		time.Sleep(time.Millisecond * 500)
	}
}


func (p *Processor) brokerSendMesToConsumer(mes message.Message) {
	tr := &network.Transport{
		Conn: p.Conn,
	}

	err := tr.SendMessage(mes)
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Broker has error sending message to consumer: %v", err))
		return
	}
	log.LogMessage("INFO", fmt.Sprintf("Broker sent message to consumer: %s", mes.Payload))

}

