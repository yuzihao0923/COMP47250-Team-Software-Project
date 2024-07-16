package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

// Order represents an e-commerce order
type Order struct {
	OrderID     string  `json:"order_id"`
	CustomerID  string  `json:"customer_id"`
	ProductID   string  `json:"product_id"`
	Quantity    int     `json:"quantity"`
	TotalAmount float64 `json:"total_amount"`
	Status      string  `json:"status"`
}

func generateOrder() Order {
	orderID := fmt.Sprintf("ORD%08d", rand.Intn(100000000))
	customerID := fmt.Sprintf("CUST%06d", rand.Intn(1000000))
	productID := fmt.Sprintf("PROD%05d", rand.Intn(100000))
	quantity := rand.Intn(10) + 1
	totalAmount := float64(quantity) * (rand.Float64()*100 + 1)
	totalAmount = float64(int(totalAmount*100)) / 100
	statuses := []string{"pending", "processed", "failed"}
	status := statuses[rand.Intn(len(statuses))]

	return Order{
		OrderID:     orderID,
		CustomerID:  customerID,
		ProductID:   productID,
		Quantity:    quantity,
		TotalAmount: totalAmount,
		Status:      status,
	}
}

func kafka_performancemain() {
	rand.Seed(time.Now().UnixNano())

	filePath := "performance_results.csv"
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	err = writer.Write([]string{"Messages Sent", "Producer Start Time", "Consumer Finish Time", "Total Test Time (seconds)"})
	if err != nil {
		log.Fatalf("Error writing CSV header: %v", err)
	}

	// Start testing time
	startTime := time.Now()

	initialMessages := 100
	maxMessages := 10000

	for numMessages := initialMessages; numMessages <= maxMessages; numMessages += 200 {
		var orders []Order
		for j := 0; j < numMessages; j++ {
			order := generateOrder()
			orders = append(orders, order)
		}

		// Measure producer start time
		producerStartTime := time.Now()

		// Simulate producer sending messages (replace with your Kafka producer logic)
		// Example:
		// kafkaProducer.ProduceMessages(orders)

		// Simulate consumer receiving messages (replace with your Kafka consumer logic)
		time.Sleep(time.Second * time.Duration(numMessages/100)) // Simulating consumer processing time

		// Measure consumer finish time
		consumerFinishTime := time.Now()

		// Calculate total test time for this iteration
		totalTestTime := consumerFinishTime.Sub(startTime).Seconds()

		// Write results to CSV
		err := writer.Write([]string{
			fmt.Sprintf("%d", numMessages),
			producerStartTime.Format(time.RFC3339Nano),
			consumerFinishTime.Format(time.RFC3339Nano),
			fmt.Sprintf("%.2f", totalTestTime), // Total test time in seconds for this iteration
		})
		if err != nil {
			log.Fatalf("Error writing CSV record: %v", err)
		}

		fmt.Printf("Sent %d messages\n", numMessages)
		fmt.Println(orders) // Example usage of orders slice
	}

	// End testing time
	endTime := time.Now()
	totalTime := endTime.Sub(startTime).Seconds()

	// Output total test time
	fmt.Printf("Total test time: %.2f seconds\n", totalTime)

	// Write total test time to CSV
	err = writer.Write([]string{"Total Test Time", "", "", fmt.Sprintf("%.2f", totalTime)})
	if err != nil {
		log.Fatalf("Error writing total test time to CSV: %v", err)
	}

	fmt.Printf("Performance results saved to %s\n", filePath)
}

func main() {
	kafka_performancemain()
}
