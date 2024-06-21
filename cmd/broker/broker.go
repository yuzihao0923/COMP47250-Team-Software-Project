package main

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/redis"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func StartBroker() {
	redis.Initialize("localhost:6379", "", 0)

	port := os.Getenv("BROKER_PORT")
	if port == "" {
		port = "8080" // default port
	}

	log.LogInfo("Broker", "Starting broker on port "+port+"...")

	// Initialize BroadcastFunc for logging
	log.BroadcastFunc = api.BroadcastMessage

	mux := http.NewServeMux()
	api.RegisterHandlers(mux)

	// 注册 WebSocket 处理程序
	mux.HandleFunc("/ws", api.HandleConnections)

	// 设置 CORS 选项
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	handler := c.Handler(mux)

	log.LogInfo("Broker", "Broker listening on port "+port)
	log.LogInfo("Broker", "Broker waiting for connections...")
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.LogError("Broker", "broker listen error: "+err.Error())
	}
}

func main() {
	StartBroker()
}
