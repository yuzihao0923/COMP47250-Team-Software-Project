package main

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/redis"
	"fmt"
	"net/http"
	"os"

	"github.com/panjf2000/ants/v2"
)

func StartBroker() {

	redis.Initialize("localhost:6379", "", 0)
	// redis.Initialize("redis1:6379", "", 0)

	port := os.Getenv("BROKER_PORT")
	if port == "" {
		port = "8080" // default port
	}

	log.LogMessage("INFO", "Starting broker on port "+port+"...")

	// Create a goroutine pool
	pool, _ := ants.NewPool(10)
	// go http.HandleFunc("/produce", api.HandleProduce)

	register_task := func() {
		http.HandleFunc("/register", api.HandleRegister)
	}

	producer_task := func() {
		http.HandleFunc("/produce", api.HandleProduce)
	}

	consumer_task := func() {
		http.HandleFunc("/consume", api.HandleConsume)
	}

	pool.Submit(register_task)
	pool.Submit(producer_task)
	pool.Submit(consumer_task)
	// go http.HandleFunc("/consume", api.HandleConsume)

	// Create a goroutine pool
	// pool, _ := ants.NewPool(10)

	// defer pool.Release()

	// http.HandleFunc("/produce", func(w http.ResponseWriter, r *http.Request) {
	// 	pool.Submit(func() {
	// 		api.HandleProduce(w, r)
	// 	})
	// })
	// http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
	// 	pool.Submit(func() {
	// 		api.HandleRegister
	// 	})
	// })
	// http.HandleFunc("/consume", func(w http.ResponseWriter, r *http.Request) {
	// 	pool.Submit(func() {
	// 		api.HandleConsume(w, r)
	// 	})
	// })

	log.LogMessage("INFO", "Broker listening on port "+port)
	log.LogMessage("INFO", "Broker waiting for connections...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Broker listen error: %v", err))
	}

	log.LogMessage("INFO", "Broker waiting for connections...")
}

func main() {
	StartBroker()
}
