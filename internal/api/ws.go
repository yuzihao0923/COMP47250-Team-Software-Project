package api

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan string)            // broadcast channel
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all incoming sources
	},
}

func init() {
	go handleMessages()
}

// HandleConnections handles incoming websocket requests from clients
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection: %v", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg string
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
	}
}

// handleMessages sends messages to all clients in the broadcast channel
func handleMessages() {
	for {
		msg := <-broadcast
		// log.Printf("Broadcasting message: %s", msg)
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// BroadcastMessage sends a message to the broadcast channel
func BroadcastMessage(message string) {
	// log.Printf("Received message to broadcast: %s", message)
	broadcast <- message
}
