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

var ctx = context.Background()
var Rdb *redis.Client

type RedisServiceInfo struct {
	StreamName string
	GroupName  string
}

// Init redis client
func Initialize(addr string, password string, db int, broadcastFunc func(string)) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:         addr,     // Address and port of Redis server
		Password:     password, // The password required to connect to the Redis server, or an empty string if no password has been set.
		DB:           db,       // Specify the number of the Redis database to which you want to connect, the default is usually 0
		PoolSize:     40,       // Maximum number of connections in the connection pool
		MinIdleConns: 10,       // Minimum number of idle connections to keep in the connection pool
		PoolTimeout:  2,        // Maximum number of seconds a client waits for an idle connection to be released when all connections are in use
	})

	pong, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.LogError("Redis", fmt.Sprintf("Failed to connect to Redis: %v", err))
		return
	}
	log.LogInfo("Redis", "Redis connected: "+pong)
	FlushAll()
	log.LogInfo("Redis", "Flush all!")

	log.BroadcastFunc = broadcastFunc
}

// FlushAll flushes all data from the Redis database
func FlushAll() error {
	_, err := Rdb.FlushAll(ctx).Result()
	return err
}

func GetClient() *redis.Client {
	return Rdb
}

// Create Consumer Group
func (rsi *RedisServiceInfo) CreateConsumerGroup() error {
	ctx := context.Background()
	// Check if the consumer group already exists
	groups, err := Rdb.XInfoGroups(ctx, rsi.StreamName).Result()
	if err != nil {
		// Try to create consumer groups directly, creating streams automatically (if needed)
		_, err = Rdb.XGroupCreateMkStream(ctx, rsi.StreamName, rsi.GroupName, "$").Result()
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

	// Check if the group already exists
	for _, group := range groups {
		if group.Name == rsi.GroupName {
			log.LogInfo("Redis", fmt.Sprintf("Consumer group '%s' already exists on stream '%s'", rsi.GroupName, rsi.StreamName))
			return nil
		}
	}

	return nil
}

// WriteToStream writes data to the specified Redis Stream.
func (rsi *RedisServiceInfo) WriteToStream(mes message.Message, producerUsername string) error {
	ctx := context.Background()
	_, err := Rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: rsi.StreamName,
		ID:     "*", // Auto-generate ID
		Values: mes.ToMap(),
	}).Result()
	if err != nil {
		log.LogError("Redis", fmt.Sprintf("Failed to write to stream '%s': %v", rsi.StreamName, err))
		return err
	}
	log.LogInfo("Redis", fmt.Sprintf("Data written to stream '%s' successfully", rsi.StreamName))
	log.LogInfo(fmt.Sprintf("Producer: %s", producerUsername), fmt.Sprintf("Send '%s' successfully", mes.Payload))
	return nil
}

// ReadFromStream reads data from the specified Redis Stream for a consumer group.
func (rsi *RedisServiceInfo) ReadFromStream(ctx context.Context, consumerUsername string) ([]redis.XStream, error) {
	streams, err := Rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    rsi.GroupName,
		Consumer: consumerUsername,
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
			log.LogInfo(fmt.Sprintf("Consumer: %s", consumerUsername), fmt.Sprintf("received: %s", mes.ID))
		}
	}

	return streams, nil
}

// ACK calls xack to remove msg from pending list
func (rsi *RedisServiceInfo) XACK(ctx context.Context, messageID string, consumerUsername string) error {
	_, err := Rdb.XAck(ctx, rsi.StreamName, rsi.GroupName, messageID).Result()
	if err != nil {
		log.LogError("Redis", fmt.Sprintf("Failed to acknowledge message '%s' in stream '%s': %v", messageID, rsi.StreamName, err))
		return err
	}
	log.LogInfo("Redis", fmt.Sprintf("Message '%s' acknowledged successfully in stream '%s'", messageID, rsi.StreamName))
	log.LogInfo(fmt.Sprintf("Consumer: %s", consumerUsername), fmt.Sprintf("Acknowledged: %s", messageID))
	return nil
}
