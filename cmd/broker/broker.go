package broker

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/redis"
	"fmt"
	"net"
	"sync"
)

var consumers []net.Conn
var consumersMutex sync.Mutex

func goprocess(conn net.Conn) {
	defer conn.Close()

	processor := Processor{
		Conn: conn,
	}

	// processor.handleConnection()

	err := processor.brokerProcessMes()
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Processor error: %v", err))
	}

}

func StartBroker() {
	log.LogMessage("INFO", "Starting broker...")

	// Init redis client (Rdb)
	redis.Initialize("localhost:6379", "", 0)

	// Listen on port 8889
	log.LogMessage("INFO", "Broker listen on port 8889")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Broker listen error: %v", err))
		return
	}
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.LogMessage("ERROR", fmt.Sprintf("Broker listen.Accept() error: %v", err))
			continue
		}

		// Register the new consumer
		consumersMutex.Lock()
		consumers = append(consumers, conn)
		consumersMutex.Unlock()

		// Start a goroutine to keep the communication between broker and the client
		go goprocess(conn)

	}
}
