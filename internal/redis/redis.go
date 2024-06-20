package redis

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"context"
	"fmt"
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
func Initialize(addr string, password string, db int) {
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
		log.LogMessage("ERROR", fmt.Sprintf("Failed to connect to Redis: %v", err))
		return
	}
	log.LogMessage("INFO", "Redis connected: "+pong)
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
			log.LogMessage("ERROR", fmt.Sprintf("Failed to create consumer group '%s' on stream '%s': %v", rsi.GroupName, rsi.StreamName, err))
			return err
		}

		log.LogMessage("INFO", fmt.Sprintf("Consumer group '%s' created successfully on stream '%s'", rsi.GroupName, rsi.StreamName))
		return nil
	}

	// Check if the group already exists
	for _, group := range groups {
		if group.Name == rsi.GroupName {
			log.LogMessage("INFO", fmt.Sprintf("Consumer group '%s' already exists on stream '%s'", rsi.GroupName, rsi.StreamName))
			return nil
		}
	}

	return nil
}

// WriteToStream writes data to the specified Redis Stream.
func (rsi *RedisServiceInfo) WriteToStream(mes message.Message) error {
	ctx := context.Background()
	_, err := Rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: rsi.StreamName,
		ID:     "*", // Auto-generate ID
		Values: mes.ToMap(),
	}).Result()
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Failed to write to stream '%s': %v", rsi.StreamName, err))
		return err
	}
	log.LogMessage("INFO", fmt.Sprintf("Data written to stream '%s' successfully", rsi.StreamName))
	return nil
}

// ReadFromStream reads data from the specified Redis Stream for a consumer group.
func (rsi *RedisServiceInfo) ReadFromStream(ctx context.Context, consumerName string) ([]redis.XStream, error) {
	streams, err := Rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    rsi.GroupName,
		Consumer: consumerName,
		Streams:  []string{rsi.StreamName, ">"},
		Count:    10,
		Block:    30 * time.Second,
	}).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		} else {
			log.LogMessage("ERROR", fmt.Sprintf("Failed to read from stream  '%s' : %v", rsi.StreamName, err))
			return nil, err
		}
	}
	return streams, nil
}

// ACK calls xack to remove msg from peding list
func (rsi *RedisServiceInfo) XACK(ctx context.Context, messageID string) error {
	_, err := Rdb.XAck(ctx, rsi.StreamName, rsi.GroupName, messageID).Result()
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Failed to acknowledge message '%s' in stream '%s': %v", messageID, rsi.StreamName, err))
		return err
	}
	log.LogMessage("INFO", fmt.Sprintf("Message '%s' acknowledged successfully in stream '%s'", messageID, rsi.StreamName))
	return nil
}
