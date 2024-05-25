package broker

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// Init redis
func initRedis() {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379", // Address and port of Redis server
		Password:     "",               // The password required to connect to the Redis server, or an empty string if no password has been set.
		DB:           0,                //Specify the number of the Redis database to which you want to connect, the default is usually 0
		PoolSize:     40,               //Maximum number of connections in the connection pool
		MinIdleConns: 10,               //Minimum number of idle connections to keep in the connection pool
		PoolTimeout:  2,                //Maximum number of seconds a client waits for an idle connection to be released when all connections are in use
	})

	pong, err := rdb.Ping(ctx).Result()
	fmt.Println(pong, err)

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

