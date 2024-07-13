package main

import (
	"COMP47250-Team-Software-Project/internal/api"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/register", api.HandleRegisterBroker).Methods("POST")
	// r.HandleFunc("/unregister", api.HandleUnRegisterBroker).Methods("DELETE")
	r.HandleFunc("/heartbeat", api.HandleHeartbeat).Methods("POST")
	r.HandleFunc("/get-broker", api.HandleGetBroker).Methods("GET")

	go api.CheckHeartbeat()

	log.Println("api server started on :8888")
	log.Fatal(http.ListenAndServe(":8888", r))
}
