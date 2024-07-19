package main

import (
	"COMP47250-Team-Software-Project/configs/configloader"
	"COMP47250-Team-Software-Project/internal/api"
	"COMP47250-Team-Software-Project/internal/database"
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/redis"
	"COMP47250-Team-Software-Project/pkg/pool"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/rs/cors"
)

var proxyURL = "http://localhost:8888"

type Broker struct {
	ID           string
	Address      string // Record its own IP address for registering with the proxy server
	dbManager    *database.MongoDB
	redisService *redis.RedisServiceInfo
	pool         *pool.WorkerPool
	mux          *http.ServeMux
}

func NewBroker(id, address string, db *database.MongoDB, redis *redis.RedisServiceInfo, pool *pool.WorkerPool) *Broker {
	mux := http.NewServeMux()
	broker := &Broker{
		ID:           id,
		Address:      address,
		dbManager:    db,
		redisService: redis,
		pool:         pool,
		mux:          mux,
	}
	api.RegisterHandlers(mux, pool, db, redis)
	broker.setupWebSocketHandler()
	return broker
}

func (b *Broker) setupWebSocketHandler() {
	b.mux.HandleFunc("/ws", api.HandleConnections)
}

func (b *Broker) Start() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
	handler := c.Handler(b.mux)

	server := &http.Server{
		Addr:    b.Address,
		Handler: handler,
	}

	// Capture system interrupt signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	// Start the server
	go func() {
		log.LogInfo("Broker", "Broker listening on "+b.Address+"...")
		log.LogInfo("Broker", "Broker waiting for connections...")
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			log.LogError("Broker", "broker listen error: "+err.Error())
		}
	}()

	// Send heartbeat periodically
	ticker := time.NewTicker(5 * time.Second) // Every 5 seconds
	go func() {
		for {
			select {
			case <-ticker.C:
				log.LogInfo("Broker", "broker sends a heartbeat to proxyURL")
				b.sendHeartbeat(proxyURL)
			case <-stop:
				ticker.Stop()
				log.LogInfo("Broker", "Stopping heartbeat...")
				return
			}
		}
	}()

	// Wait for system interrupt signal
	<-stop
	log.LogInfo("Broker", "Shutdown signal received, shutting down server...")

	// Send unregister signal
	// b.UnregisterFromProxy(proxyURL)

	// Create a timeout context for shutting down the broker server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.LogError("Broker", "Server shutdown error: "+err.Error())
	} else {
		log.LogInfo("Broker", "Server shutdown successfully")
	}

}

func (b *Broker) register2Proxy(proxyURL string) error {
	brokerInfo := struct {
		ID      string `json:"id"`
		Address string `json:"address"`
	}{
		ID:      b.ID,
		Address: b.Address,
	}
	data, err := json.Marshal(brokerInfo)
	if err != nil {
		return err
	}

	resp, err := http.Post(proxyURL+"/register", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register with proxy: status %d", resp.StatusCode)
	}

	return nil
}

func (b *Broker) UnregisterFromProxy(proxyURL string) {
	unregisterData := struct {
		ID string `json:"id"`
	}{
		ID: b.ID,
	}
	jsonData, err := json.Marshal(unregisterData)
	if err != nil {
		log.LogError("Broker", "Error encoding unregister data: "+err.Error())
		return
	}

	req, err := http.NewRequest(http.MethodDelete, proxyURL+"/unregister", bytes.NewBuffer(jsonData))
	if err != nil {
		log.LogError("Broker", "Error creating unregister request: "+err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.LogError("Broker", "Failed to send unregister request: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.LogError("Broker", fmt.Sprintf("Unregister request failed with status: %v", resp.StatusCode))
	} else {
		log.LogInfo("Broker", "Successfully unregistered from proxy")
	}
}

func (b *Broker) sendHeartbeat(proxyURL string) {
	heartbeatData := struct {
		ID      string `json:"id"`
		Address string `json:"address"`
	}{
		ID:      b.ID,
		Address: b.Address,
	}

	jsonData, err := json.Marshal(heartbeatData)
	if err != nil {
		log.LogError("Broker", "Error encoding heartbeat data: "+err.Error())
		return
	}

	resp, err := http.Post(proxyURL+"/heartbeat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.LogError("Broker", "Failed to send heartbeat: "+err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.LogError("Broker", fmt.Sprintf("Heartbeat failed with status: %d", resp.StatusCode))
	}
}

func initDB() (*database.MongoDB, error) {
	db, err := database.NewMongoDB("mongodb://localhost:27017", "userdb", "users")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	ctx := context.Background()
	err = db.InitializeMongoDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	log.LogInfo("Broker", "Database initialized successfully")
	return db, nil
}

func startBroker(brokerConfig configloader.BrokerConfig, db *database.MongoDB, rsi *redis.RedisServiceInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	pool := pool.NewWorkerPool(10, 100) // 10 workers, JobQueueSize 100
	pool.Start()

	broker := NewBroker(brokerConfig.ID, brokerConfig.Address, db, rsi, pool)

	err := broker.register2Proxy(proxyURL)
	if err != nil {
		log.LogError("Broker", fmt.Sprintf("Failed to register broker %s: %v", brokerConfig.ID, err))
		return
	}

	broker.Start()
}

func main() {
	fmt.Println("Starting Broker cluster...")

	configPath := "../../configs/configloader/brokers.yaml"
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.LogError("Broker", fmt.Sprintf("Configuration file does not exist: %s", configPath))
		time.Sleep(1 * time.Second) // Wait 1 second before exit
		os.Exit(1)
	}
	configLoader := configloader.NewYAMLConfigLoader(configPath)
	conf, err := configLoader.LoadConfig()
	if err != nil {
		log.LogError("Broker", "Failed to load configuration: "+err.Error())
		return
	} else {
		fmt.Println("Load Config success..")
	}

	db, err := initDB()
	if err != nil {
		log.LogError("Broker", err.Error())
		return
	} else {
		fmt.Println("Database connected...")
	}
	defer func() {
		ctx := context.Background()
		if err := db.Close(ctx); err != nil {
			log.LogError("Broker", "Failed to close MongoDB connection: "+err.Error())
		}
	}()

	// Create Redis cluster client instance
	rsi := redis.NewRedisClusterClient([]string{
		"localhost:6381",
		"localhost:6382",
		"localhost:6383",
		"localhost:6384",
		"localhost:6385",
		"localhost:6386",
	}, "", 0, api.BroadcastMessage)
	ctx := context.Background()

	// Check connection, Ping function will flush all data in Redis
	if err := rsi.Ping(ctx, api.BroadcastMessage); err != nil {
		log.LogError("Broker", fmt.Sprintf("Failed to connect to Redis: %v", err))
	} else {
		fmt.Println("Redis connected...")
	}

	var wg sync.WaitGroup
	for _, brokerConfig := range conf.Brokers {
		wg.Add(1)
		go startBroker(brokerConfig, db, rsi, &wg)
	}

	wg.Wait()
	fmt.Println("All brokers started.")
}
