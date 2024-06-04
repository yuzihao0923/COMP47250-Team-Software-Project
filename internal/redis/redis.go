package redis

import (
	"COMP47250-Team-Software-Project/internal/log"
	"COMP47250-Team-Software-Project/internal/message"
	"context"
	"fmt"

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
		DB:           db,       //Specify the number of the Redis database to which you want to connect, the default is usually 0
		PoolSize:     40,       //Maximum number of connections in the connection pool
		MinIdleConns: 10,       //Minimum number of idle connections to keep in the connection pool
		PoolTimeout:  2,        //Maximum number of seconds a client waits for an idle connection to be released when all connections are in use
	})

	pong, err := Rdb.Ping(ctx).Result()
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Failed to connect to Redis: %v", err))
		return
	}
	log.LogMessage("INFO", "Redis connected: "+pong)

	// err := rdb.Set(ctx, "key", "value", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	// val, err := rdb.Get(ctx, "key").Result()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("key", val)
}

func GetClient() *redis.Client {
	return Rdb
}

// Create Consumer Group
// func CreateConsumerGroup(streamName, groupName string) error {
// 	ctx := context.Background()

// 	//Try to create a consumer group, if stream is not exist, create stream first
// 	_, err := Rdb.XGroupCreateMkStream(ctx, streamName, groupName, "$").Result() // The link for the definition of XGroupCreateMkStream: https://github.com/redis/go-redis/blob/v6.15.9/commands.go#L1429
// 	if err != nil {
// 		// Check the error is that do not exist stream
// 		if strings.Contains(err.Error(), "NOGROUP") {
// 			// create a empty stream
// 			_, err := Rdb.XAdd(ctx, &redis.XAddArgs{
// 				Stream: streamName,
// 				ID:     "*", // Auto ID
// 				// Values: map[string]interface{}{"init": "init"},
// 			}).Result()
// 			if err != nil {
// 				log.LogMessage("ERROR", fmt.Sprintf("Failed to create stream '%s':%v", streamName, err))
// 				return err
// 			}

// 			// create a consumer group again
// 			_, err = Rdb.XGroupCreateMkStream(ctx, streamName, groupName, "$").Result()
// 			if err != nil {
// 				log.LogMessage("ERROR", fmt.Sprintf("Failed to create consumer group '%s' after stream creation:%v", groupName, err))
// 				return err
// 			}
// 		} else {
// 			log.LogMessage("ERROR", fmt.Sprintf("Failed to create consumer group '%s': %v", groupName, err))
// 			return err
// 		}
// 	}
// 	log.LogMessage("INFO", fmt.Sprintf("Consumer group '%s' created successfully on stream '%s'", groupName, streamName))
// 	return nil
// }

// Create Consumer Group
func (rsi *RedisServiceInfo) CreateConsumerGroup() error {
	ctx := context.Background()
	// Try to create consumer groups directly, creating streams automatically (if needed)
	_, err := Rdb.XGroupCreateMkStream(ctx, rsi.StreamName, rsi.GroupName, "$").Result()
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Failed to create consumer group '%s' on stream '%s': %v", rsi.GroupName, rsi.StreamName, err))
		return err
	}

	log.LogMessage("INFO", fmt.Sprintf("Consumer group '%s' created successfully on stream '%s'", rsi.GroupName, rsi.StreamName))
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
		Block:    0,
	}).Result()
	if err != nil {
		log.LogMessage("ERROR", fmt.Sprintf("Failed to read from stream '%s': %v", rsi.StreamName, err))
		return nil, err
	}
	return streams, nil
}
