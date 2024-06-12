package broker

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"fmt"
	"net/http"
)

func StartBroker() {
	log.LogMessage("INFO", "Starting broker...")

	http.HandleFunc("/produce", api.HandleProduce)
	http.HandleFunc("/register", api.HandleRegister)
	http.HandleFunc("/consume", api.HandleConsume)

	log.LogMessage("INFO", "Broker listen on port 8889")
	err := http.ListenAndServe(":8889", nil)
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Broker listen error: %v", err))
	}
}
