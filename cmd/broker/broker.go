package main

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/database"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/redis"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	err := database.InitializeMongoDB()
	if err != nil {
		log.LogError("Broker", "Failed to initialize database: "+err.Error())
		return
	}
	log.LogInfo("Broker", "Database initialized successfully")

	// Init broadcast here with redis
	redis.Initialize("localhost:6379", "", 0, api.BroadcastMessage)

	port := os.Getenv("BROKER_PORT")
	if port == "" {
		port = "8080" // default port
	}

	log.LogInfo("Broker", "Starting broker on port "+port+"...")

	mux := http.NewServeMux()
	api.RegisterHandlers(mux)

	// Register WebSocket handler
	mux.HandleFunc("/ws", api.HandleConnections)

	// Set CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	handler := c.Handler(mux)

	go func() {
		log.LogInfo("Broker", "Broker listening on port "+port)
		log.LogInfo("Broker", "Broker waiting for connections...")
		err := http.ListenAndServe(":"+port, handler)
		if err != nil {
			log.LogError("Broker", "broker listen error: "+err.Error())
		}
	}()

	select {} // Block forever
}
