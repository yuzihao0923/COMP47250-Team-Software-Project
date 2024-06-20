package main

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/redis"
	"fmt"
	"net/http"
	"os"
)

func StartBroker() {

	redis.Initialize("localhost:6379", "", 0)

	port := os.Getenv("BROKER_PORT")
	if port == "" {
		port = "8080" // default port
	}

	log.LogInfo("Starting broker on port " + port + "...")

	api.RegisterHandlers()

	log.LogInfo("Broker listening on port " + port)
	log.LogInfo("Broker waiting for connections...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.LogError(fmt.Errorf("broker listen error: %v", err))
	}

	log.LogInfo("Broker waiting for connections...")
}

func main() {
	StartBroker()
}
