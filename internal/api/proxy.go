package api

import (
	"COMP47250-Team-Software-Project/internal/log"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Broker struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	LastCheck time.Time `json:"-"`
}

var brokers []Broker
var lock sync.RWMutex
var pos int32

func HandleRegisterBroker(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var broker Broker

	if err := json.NewDecoder(r.Body).Decode(&broker); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	broker.LastCheck = time.Now()
	lock.Lock()
	brokers = append(brokers, broker)

	lock.Unlock()

	w.WriteHeader(http.StatusOK)
}

func HandleUnRegisterBroker(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		http.Error(w, "Only DELETE method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var broker Broker
	if err := json.NewDecoder(r.Body).Decode(&broker); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lock.Lock()
	for i, b := range brokers {
		if b.Address == broker.Address {
			brokers = append(brokers[:i], brokers[i+1:]...)
			break
		}
	}
	lock.Unlock()

	log.LogInfo("ProxyServer", fmt.Sprintf("Current available numeber of brokers: %v", len(brokers)))

	w.WriteHeader(http.StatusOK)
}

func HandleGetBroker(w http.ResponseWriter, r *http.Request) {
	lock.RLock()
	defer lock.RUnlock()

	if len(brokers) == 0 {
		http.Error(w, "No available brokers", http.StatusNotFound)
		return
	}

	// Simple round-robin load balancing
	atomic.AddInt32(&pos, 1) // 原子地递增 pos
	index := int(pos) % len(brokers)
	json.NewEncoder(w).Encode(brokers[index])
}

func HandleHeartbeat(w http.ResponseWriter, r *http.Request) {
	var heartbeatData struct {
		ID      string `json:"id"`
		Address string `json:"address"`
	}
	if err := json.NewDecoder(r.Body).Decode(&heartbeatData); err != nil {
		http.Error(w, "Invalid heartbeat data", http.StatusBadRequest)
		return
	}
	lock.Lock()
	for i, broker := range brokers {
		if broker.ID == heartbeatData.ID {
			brokers[i].LastCheck = time.Now()
			break
		}
	}
	lock.Unlock()

	w.WriteHeader(http.StatusOK)
}

func CheckHeartbeat() {
	for {
		time.Sleep(5 * time.Second)
		lock.Lock()
		currentTime := time.Now()
		var newBrokers []Broker
		for i, broker := range brokers {
			if currentTime.Sub(brokers[i].LastCheck) <= 10*time.Second {
				newBrokers = append(newBrokers, broker)
			}
		}
		brokers = newBrokers
		log.LogInfo("ProxyServer", fmt.Sprintf("Current available numeber of brokers: %v", len(brokers)))
		lock.Unlock()
	}
}
