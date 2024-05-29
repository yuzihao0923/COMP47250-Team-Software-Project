package broker

import (
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

	err := processor.brokerProcessMes()
	if err != nil {
		fmt.Println("Processor error:", err)
	}
}

func StartBroker() {
	//	Init redis client(Rdb)
	redis.Initialize("localhost:6379", "", 0)

	// Listen port 8889
	fmt.Println("Broker listen to port 8889")
	listen, err := net.Listen("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("Broker listen err= ", err)
		return
	}
	defer listen.Close()

	for {
		// fmt.Println("Waiting for connections from clients")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Broker listen.Accept() error:", err)
			continue
		}

		fmt.Println("Successfully accepted connection from client")

		// Register the new consumer
		consumersMutex.Lock()
		consumers = append(consumers, conn)
		consumersMutex.Unlock()

		// Start a goroutine to keep the communication between broker and the client
		go goprocess(conn)
	}
}
