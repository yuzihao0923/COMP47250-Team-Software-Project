package broker

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"fmt"
	"net/http"
	"os"
)

func StartBroker() {
	port := os.Getenv("BROKER_PORT")
	if port == "" {
		port = "8080" // default port
	}

	log.LogMessage("INFO", "Starting broker on port " + port + "...")

	http.HandleFunc("/produce", api.HandleProduce)
	http.HandleFunc("/register", api.HandleRegister)
	http.HandleFunc("/consume", api.HandleConsume)

	log.LogMessage("INFO", "Broker listening on port " + port)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Broker listen error: %v", err))
	}
}
