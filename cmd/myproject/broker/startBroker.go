package broker

import (
	"fmt"
	"myproject/internal/message/redis"
	"net"
)

func goprocess(conn net.Conn) {
	defer conn.Close()

	processor := Processor{
		Conn: conn,
	}

	processor.brokerProcessMes()

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
		fmt.Println("Wait connection from clients")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Broker listen.Accept() err=", err)
			return
		}

		//Start a goroutin to keep the communication between broker and the client
		go goprocess(conn)
	}
}
