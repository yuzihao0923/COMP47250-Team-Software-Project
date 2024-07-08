package main

import (
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/database"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/redis"
	"COMP47250-Team-Software-Project/pkg/pool"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	db, err := database.NewMongoDB("mongodb://localhost:27017", "userdb", "users")
	if err != nil {
		log.LogError("Broker", "Failed to initialize database: "+err.Error())
		return
	}
	defer func() {
		ctx := context.Background()
		if err := db.Close(ctx); err != nil {
			log.LogError("Broker", "Failed to close MongoDB connection: "+err.Error())
		}
	}()
	// err := database.InitializeMongoDB()
	// if err != nil {
	// 	log.LogError("Broker", "Failed to initialize database: "+err.Error())
	// 	return
	// }
	// log.LogInfo("Broker", "Database initialized successfully")

	// Init broadcast here with redis
	// redis.Initialize("localhost:6379", "", 0, api.BroadcastMessage)

	// Create redis client instance
	rsi := redis.NewRedisClient("localhost:6379", "", 0)
	ctx := context.Background()

	// Check connection, Ping func will flush all data in redis
	if err := rsi.Ping(ctx, api.BroadcastMessage); err != nil {
		log.LogError("Broker", fmt.Sprintf("Failed to connect to Redis: %v", err))
	}

	port := os.Getenv("BROKER_PORT")
	if port == "" {
		port = "8080" // default port
	}

	log.LogInfo("Broker", "Starting broker on port "+port+"...")

	pool := pool.NewWorkerPool(10, 100) // 10 workers, JobQueueSize 100
	pool.Start()

	mux := http.NewServeMux()
	api.RegisterHandlers(mux, pool, db, rsi)

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
