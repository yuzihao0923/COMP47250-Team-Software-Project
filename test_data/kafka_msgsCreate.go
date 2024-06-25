package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/IBM/sarama"
)

// Order_kafka represents an e-commerce order
type Order_kafka struct {
	OrderID     string  `json:"order_id"`
	CustomerID  string  `json:"customer_id"`
	ProductID   string  `json:"product_id"`
	Quantity    int     `json:"quantity"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}

func kafka_generateOrder() Order_kafka {
	orderID := fmt.Sprintf("ORD%08d", rand.Intn(100000000))
	customerID := fmt.Sprintf("CUST%06d", rand.Intn(1000000))
	productID := fmt.Sprintf("PROD%05d", rand.Intn(100000))
	quantity := rand.Intn(10) + 1
	totalAmount := float64(quantity) * (rand.Float64()*100 + 1)
	totalAmount = float64(int(totalAmount*100)) / 100
	statuses := []string{"pending", "processed", "failed"}
	status := statuses[rand.Intn(len(statuses))]

	return Order_kafka{
		OrderID:     orderID,
		CustomerID:  customerID,
		ProductID:   productID,
		Quantity:    quantity,
		TotalAmount: totalAmount,
		Status:      status,
	}
}

func kafkaMsg_createmain() {
	rand.Seed(time.Now().UnixNano())

	brokerList := []string{"localhost:9092"} // Kafka broker 地址
	topic := "test-topic"                    // Kafka 主题名

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal       // 等待本地确认
	config.Producer.Compression = sarama.CompressionSnappy   // 使用 Snappy 压缩消息
	config.Producer.Flush.Frequency = 500 * time.Millisecond // 设置刷新时间间隔

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer producer.Close()

	var orders []Order_kafka
	for i := 0; i < 1000; i++ { // Generate 1000 orders
		order := kafka_generateOrder()
		orders = append(orders, order)

		// 将订单编码为 JSON
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error encoding order to JSON: %v", err)
			continue
		}

		// 发送订单消息到 Kafka 主题
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(orderJSON),
		})
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			continue
		}
	}

	fmt.Println("Generated and sent 1000 orders to Kafka")
}

func kafka_msgcreatemain() {
	kafkaMsg_createmain()
}
