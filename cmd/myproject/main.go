// cmd/myproject/main.go
package main

import (
	"fmt"
	"myproject/pkg/broker"
)

func main() {
	// 创建一个 Broker 实例并连接到 Redis
	br := broker.NewBroker("localhost:6379", "", 0)

	// 设置一个键值对
	err := br.Set("mykey", "myvalue")
	if err != nil {
		fmt.Println("Error setting key:", err)
		return
	}

	// 获取键值对
	value, err := br.Get("mykey")
	if err != nil {
		fmt.Println("Error getting key:", err)
		return
	}

	fmt.Println("mykey:", value)
}
