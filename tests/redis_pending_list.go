package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis 服务器地址
		DB:   0,                // 使用默认数据库
	})

	streamName := "mystream1"
	groupName := "mygroup1"
	consumerName := "myconsumer1"

	// 创建消费者组
	err := rdb.XGroupCreateMkStream(ctx, streamName, groupName, "$").Err()
	if err != nil && err != redis.Nil {
		log.Fatalf("Failed to create group: %v", err)
	}

	// 生产者发送大量消息到流中
	for i := 0; i < 100000; i++ {
		err := rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: streamName,
			Values: map[string]interface{}{"key": fmt.Sprintf("value_%d", i)},
		}).Err()
		if err != nil {
			log.Fatalf("Failed to add message to stream: %v", err)
		}
		if i%1000 == 0 {
			fmt.Printf("Sent %d messages\n", i)
		}
	}

	// 消费者读取消息但不确认，从而使其留在 pending list 中
	for {
		res, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    groupName,
			Consumer: consumerName,
			Streams:  []string{streamName, ">"},
			Count:    10,
			Block:    1000 * time.Millisecond,
		}).Result()
		if err != nil {
			log.Fatalf("Failed to read messages: %v", err)
		}

		if len(res) == 0 {
			break
		}

		for _, msg := range res[0].Messages {
			fmt.Printf("Processing message: %v\n", msg)
			// 这里故意不使用 XAck 确认消息，以增加 pending list
		}

		pendingCount, err := rdb.XPending(ctx, streamName, groupName).Result()
		if err != nil {
			log.Fatalf("Failed to get pending messages: %v", err)
		}
		fmt.Printf("Pending messages: %v\n", pendingCount)

		time.Sleep(1 * time.Second)
	}
}