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
	Client     *redis.Client
	StreamName string
	GroupName  string
}

func NewRedisClient(addr, password string, db int) *RedisServiceInfo {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,     // Address and port of Redis server
		Password:     password, // The password required to connect to the Redis server, or an empty string if no password has been set.
		DB:           db,       // Specify the number of the Redis database to which you want to connect, the default is usually 0
		PoolSize:     40,       // Maximum number of connections in the connection pool
		MinIdleConns: 10,       // Minimum number of idle connections to keep in the connection pool
		PoolTimeout:  2,        // Maximum number of seconds a client waits for an idle connection to be released when all connections are in use
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
		// Try to create consumer groups directly, creating streams automatically (if needed)
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

	// Check if the group already exists
	for _, group := range groups {
		if group.Name == rsi.GroupName {
			log.LogInfo("Redis", fmt.Sprintf("Consumer group '%s' already exists on stream '%s'", rsi.GroupName, rsi.StreamName))
			return nil
		}
	}

	return nil
}

func (rsi *RedisServiceInfo) WriteToStream(ctx context.Context, producerUserName string, mes message.Message) error {
	// 模拟redis没有打开
	// rsi.Client.Close()

	// 模拟超时
	// 假设超时时间是1毫秒，但是操作需要2毫秒才能完成，就会超时 context deadline exceeded
	// ctx, cancel := context.WithTimeout(ctx, 1*time.Millisecond)
	// defer cancel()
	// time.Sleep(2 * time.Millisecond)

	messageID, err := rsi.Client.XAdd(ctx, &redis.XAddArgs{
		Stream: rsi.StreamName,
		ID:     "*", // Auto-generate ID
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

// // var ctx = context.Background()
// // var Rdb *redis.Client

// // type RedisServiceInfo struct {
// // 	StreamName string
// // 	GroupName  string
// // }

// // Init redis client
// // func Initialize(addr string, password string, db int, broadcastFunc func(string)) {
// // 	Rdb = redis.NewClient(&redis.Options{
// // 		Addr:         addr,     // Address and port of Redis server
// // 		Password:     password, // The password required to connect to the Redis server, or an empty string if no password has been set.
// // 		DB:           db,       // Specify the number of the Redis database to which you want to connect, the default is usually 0
// // 		PoolSize:     40,       // Maximum number of connections in the connection pool
// // 		MinIdleConns: 10,       // Minimum number of idle connections to keep in the connection pool
// // 		PoolTimeout:  2,        // Maximum number of seconds a client waits for an idle connection to be released when all connections are in use
// // 	})

// // 	pong, err := Rdb.Ping(ctx).Result()
// // 	if err != nil {
// // 		log.LogError("Redis", fmt.Sprintf("Failed to connect to Redis: %v", err))
// // 		return
// // 	}
// // 	log.LogInfo("Redis", "Redis connected: "+pong)
// // 	FlushAll()
// // 	log.LogInfo("Redis", "Flush all!")

// // 	log.BroadcastFunc = broadcastFunc
// // }

// // FlushAll flushes all data from the Redis database
// // func FlushAll() error {
// // 	_, err := Rdb.FlushAll(ctx).Result()
// // 	return err
// // }

// // func GetClient() *redis.Client {
// // 	return Rdb
// // }

// // Create Consumer Group
// // func (rsi *RedisServiceInfo) CreateConsumerGroup() error {
// // 	ctx := context.Background()
// // 	// Check if the consumer group already exists
// // 	groups, err := Rdb.XInfoGroups(ctx, rsi.StreamName).Result()
// // 	if err != nil {
// // 		// Try to create consumer groups directly, creating streams automatically (if needed)
// // 		_, err = Rdb.XGroupCreateMkStream(ctx, rsi.StreamName, rsi.GroupName, "$").Result()
// // 		if err != nil {
// // 			if strings.Contains(err.Error(), "Consumer Group name already exists") {
// // 				log.LogWarning("Redis", fmt.Sprintf("Consumer group '%s' already exists on stream '%s'", rsi.GroupName, rsi.StreamName))
// // 				return nil
// // 			}
// // 			log.LogError("Redis", fmt.Sprintf("Failed to create consumer group '%s' on stream '%s': %v", rsi.GroupName, rsi.StreamName, err))
// // 			return err
// // 		}

// // 		log.LogInfo("Redis", fmt.Sprintf("Consumer group '%s' created successfully on stream '%s'", rsi.GroupName, rsi.StreamName))
// // 		return nil
// // 	}

// // 	// Check if the group already exists
// // 	for _, group := range groups {
// // 		if group.Name == rsi.GroupName {
// // 			log.LogInfo("Redis", fmt.Sprintf("Consumer group '%s' already exists on stream '%s'", rsi.GroupName, rsi.StreamName))
// // 			return nil
// // 		}
// // 	}

// // 	return nil
// // }

// // WriteToStream writes data to the specified Redis Stream.
// // func (rsi *RedisServiceInfo) WriteToStream(mes message.Message, producerUsername string) error {
// // 	ctx := context.Background()
// // 	messageID, err := Rdb.XAdd(ctx, &redis.XAddArgs{
// // 		Stream: rsi.StreamName,
// // 		ID:     "*", // Auto-generate ID
// // 		Values: mes.ToMap(),
// // 	}).Result()
// // 	if err != nil {
// // 		log.LogError("Redis", fmt.Sprintf("Failed to write to stream '%s': %v", rsi.StreamName, err))
// // 		return err
// // 	}
// // 	log.LogInfo("Redis", fmt.Sprintf("Message from '%s' to stream '%s' successfully", producerUsername, rsi.StreamName))
// // 	log.LogInfo(fmt.Sprintf("Producer: %s", producerUsername), fmt.Sprintf("sent %s: '%s' successfully", messageID, mes.Payload))
// // 	return nil
// // }

// // ReadFromStream reads data from the specified Redis Stream for a consumer group.
// // func (rsi *RedisServiceInfo) ReadFromStream(ctx context.Context, consumerUsername string) ([]redis.XStream, error) {
// // 	streams, err := Rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
// // 		Group:    rsi.GroupName,
// // 		Consumer: consumerUsername,
// // 		Streams:  []string{rsi.StreamName, ">"},
// // 		Count:    10,
// // 		Block:    30 * time.Second,
// // 	}).Result()
// // 	if err != nil {
// // 		if err == redis.Nil {
// // 			return nil, nil
// // 		} else {
// // 			log.LogError("Redis", fmt.Sprintf("Failed to read from stream  '%s' : %v", rsi.StreamName, err))
// // 			return nil, err
// // 		}
// // 	}

// // 	for _, stream := range streams {
// // 		for _, mes := range stream.Messages {
// // 			message, err := message.NewMessageFromMap(mes.Values, mes.ID)
// // 			if err != nil {
// // 				log.LogError("Message", fmt.Sprintf("Failed to convert message from map: %v", err))
// // 				continue
// // 			}
// // 			log.LogInfo(fmt.Sprintf("Consumer: %s", consumerUsername), fmt.Sprintf("received: %s: '%s'", mes.ID, message.Payload))
// // 		}
// // 	}

// // 	return streams, nil
// // }

// // ACK calls xack to remove msg from pending list
// // func (rsi *RedisServiceInfo) XACK(ctx context.Context, messageID string, consumerUsername string) error {
// // 	_, err := Rdb.XAck(ctx, rsi.StreamName, rsi.GroupName, messageID).Result()
// // 	if err != nil {
// // 		log.LogError("Redis", fmt.Sprintf("Failed to acknowledge message '%s' in stream '%s': %v", messageID, rsi.StreamName, err))
// // 		return err
// // 	}
// // 	log.LogInfo("Redis", fmt.Sprintf("Message '%s' to '%s' acknowledged successfully in stream '%s'", messageID, consumerUsername, rsi.StreamName))
// // 	log.LogInfo(fmt.Sprintf("Consumer: %s", consumerUsername), fmt.Sprintf("acknowledged: %s", messageID))
// // 	return nil
// // }
