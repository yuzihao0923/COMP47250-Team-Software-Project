// cmd/myproject/main.go
package main

import (
	"COMP47250-Team-Software-Project/pkg/broker"
	"fmt"
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

	// 检查键是否存在
	exists, err := br.Exists("mykey")
	if err != nil {
		fmt.Println("Error checking existence:", err)
		return
	}
	fmt.Printf("Does mykey exist? %v\n", exists)

	// 设置键的过期时间
	err = br.Expire("mykey", 60)
	if err != nil {
		fmt.Println("Error setting expiration:", err)
		return
	}

	// 删除键
	err = br.Del("mykey")
	if err != nil {
		fmt.Println("Error deleting key:", err)
		return
	}

	// 再次获取键值对，应该返回空
	value, err = br.Get("mykey")
	if err != nil {
		fmt.Println("Error getting key:", err)
	} else {
		fmt.Println("mykey:", value)
	}
}
