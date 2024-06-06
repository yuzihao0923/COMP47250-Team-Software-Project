package main

import (
	"encoding/json"
	"fmt"
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

func main() {
	rand.Seed(time.Now().UnixNano())

	var orders []Order
	for i := 0; i < 1000; i++ { // Generate 1000 orders
		order := generateOrder()
		orders = append(orders, order)
	}

	file, err := os.Create("orders.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(orders)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	fmt.Println("Generated orders.json with 1000 orders")
}
