package redis

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisServiceInfo struct {
	Client     *redis.ClusterClient // 使用 ClusterClient 替换 Client
	StreamName string
	GroupName  string
}

func NewRedisClusterClient(addrs []string, password string, db int) *RedisServiceInfo {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        addrs,
		Password:     password,
		PoolSize:     40,
		MinIdleConns: 10,
		PoolTimeout:  2 * time.Second,
	})
	return &RedisServiceInfo{
		Client: rdb,
	}
}

func (rsi *RedisServiceInfo) Ping(ctx context.Context, broadcastFunc func(string)) error {
	pong, err := rsi.Client.Ping(ctx).Result()
	if err != nil {
		log.LogError("Redis", fmt.Sprintf("Failed to connect to Redis: %v", err))
		return err
	}
	log.LogInfo("Redis", "Redis connected: "+pong)
	rsi.FlushAll(ctx)
	log.LogInfo("Redis", "Flush all!")

	log.BroadcastFunc = broadcastFunc
	return nil
}

func (rsi *RedisServiceInfo) FlushAll(ctx context.Context) error {
	_, err := rsi.Client.FlushAll(ctx).Result()
	if err != nil {
		log.LogError("Redis", "Failed to flush all data from Redis")
		return err
	}
	log.LogInfo("Redis", "Flush all!")
	return nil
}

func (rsi *RedisServiceInfo) CreateConsumerGroup(ctx context.Context) error {
	groups, err := rsi.Client.XInfoGroups(ctx, rsi.StreamName).Result()
	if err != nil {
		// 尝试直接创建消费者组，必要时自动创建流
		_, err = rsi.Client.XGroupCreateMkStream(ctx, rsi.StreamName, rsi.GroupName, "$").Result()
		if err != nil {
			if strings.Contains(err.Error(), "Consumer Group name already exists") {
				log.LogWarning("Redis", fmt.Sprintf("Consumer group '%s' already exists on stream '%s'", rsi.GroupName, rsi.StreamName))
				return nil
			}
			log.LogError("Redis", fmt.Sprintf("Failed to create consumer group '%s' on stream '%s': %v", rsi.GroupName, rsi.StreamName, err))
			return err
		}

		log.LogInfo("Redis", fmt.Sprintf("Consumer group '%s' created successfully on stream '%s'", rsi.GroupName, rsi.StreamName))
		return nil
	}

	// 检查组是否已经存在
	for _, group := range groups {
		if group.Name == rsi.GroupName {
			log.LogInfo("Redis", fmt.Sprintf("Consumer group '%s' already exists on stream '%s'", rsi.GroupName, rsi.StreamName))
			return nil
		}
	}

	return nil
}

func (rsi *RedisServiceInfo) WriteToStream(ctx context.Context, producerUserName string, mes message.Message) error {
	messageID, err := rsi.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: rsi.StreamName,
		ID:     "*", // 自动生成ID
		Values: mes.ToMap(),
	}).Result()
	if err != nil {
		log.LogError("Redis", fmt.Sprintf("Failed to write to stream '%s': %v", rsi.StreamName, err))
		return err
	}
	log.LogInfo("Redis", fmt.Sprintf("Message from '%s' to stream '%s' successfully", producerUserName, rsi.StreamName))
	log.LogInfo(fmt.Sprintf("Producer: %s", producerUserName), fmt.Sprintf("sent %s: '%s' successfully", messageID, mes.Payload))
	return nil
}

func (rsi *RedisServiceInfo) ReadFromStream(ctx context.Context, consumerUserName string) ([]redis.XStream, error) {
	streams, err := rsi.Client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    rsi.GroupName,
		Consumer: consumerUserName,
		Streams:  []string{rsi.StreamName, ">"},
		Count:    10,
		Block:    30 * time.Second,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			log.LogError("Redis", fmt.Sprintf("Failed to read from stream  '%s' : %v", rsi.StreamName, err))
			return nil, err
		}
	}

	for _, stream := range streams {
		for _, mes := range stream.Messages {
			message, err := message.NewMessageFromMap(mes.Values, mes.ID)
			if err != nil {
				log.LogError("Message", fmt.Sprintf("Failed to convert message from map: %v", err))
				continue
			}
			log.LogInfo(fmt.Sprintf("Consumer: %s", consumerUserName), fmt.Sprintf("received: %s: '%s'", mes.ID, message.Payload))
		}
	}

	return streams, nil
}

func (rsi *RedisServiceInfo) XACK(ctx context.Context, consumerUserName, messageID string) error {
	_, err := rsi.Client.XAck(ctx, rsi.StreamName, rsi.GroupName, messageID).Result()
	if err != nil {
		log.LogError("Redis", fmt.Sprintf("Failed to acknowledge message '%s' in stream '%s': %v", messageID, rsi.StreamName, err))
		return err
	}
	log.LogInfo("Redis", fmt.Sprintf("Message '%s' to '%s' acknowledged successfully in stream '%s'", messageID, consumerUserName, rsi.StreamName))
	log.LogInfo(fmt.Sprintf("Consumer: %s", consumerUserName), fmt.Sprintf("acknowledged: %s", messageID))
	return nil
}
