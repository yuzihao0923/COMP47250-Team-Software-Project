package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/IBM/sarama"
)

func kafka_consumermain() {
	brokerList := []string{"localhost:9092"} // Kafka broker 地址

	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokerList, config)
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer consumer.Close()

	topic := "test-topic" // Kafka 主题名

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatalf("Error consuming partition: %v", err)
	}
	defer partitionConsumer.Close()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case msg := <-partitionConsumer.Messages():
				fmt.Printf("Received message from Kafka: %s\n", string(msg.Value))
			case err := <-partitionConsumer.Errors():
				log.Printf("Error from partition consumer: %v\n", err)
			}
		}
	}()

	// 等待中断信号，优雅地关闭消费者
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan

	fmt.Println("Interrupt signal received, shutting down consumer...")
	wg.Wait()
}
