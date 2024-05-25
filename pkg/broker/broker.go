// pkg/broker/broker.go
package broker

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type Broker struct {
	client *redis.Client
	ctx    context.Context
}

// NewBroker 创建一个新的 Broker 实例并连接到 Redis
func NewBroker(addr string, password string, db int) *Broker {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Broker{
		client: rdb,
		ctx:    context.Background(),
	}
}

// Set 设置一个键值对
func (b *Broker) Set(key string, value interface{}) error {
	err := b.client.Set(b.ctx, key, value, 0).Err()
	if err != nil {
		log.Printf("Failed to set key %s: %v", key, err)
		return err
	}
	return nil
}

// Get 获取一个键的值
func (b *Broker) Get(key string) (string, error) {
	val, err := b.client.Get(b.ctx, key).Result()
	if err == redis.Nil {
		log.Printf("Key %s does not exist", key)
		return "", nil
	} else if err != nil {
		log.Printf("Failed to get key %s: %v", key, err)
		return "", err
	}
	return val, nil
}

// Delete 删除一个键的值
func (b *Broker) Del(key string) error {
	err := b.client.Del(b.ctx, key).Err()
	if err != nil {
		log.Printf("Failed to delete key %s: %v", key, err)
		return err
	}
	return nil
}

// Exists 检查一个键是否存在
func (b *Broker) Exists(key string) (bool, error) {
	exists, err := b.client.Exists(b.ctx, key).Result()
	if err != nil {
		log.Printf("Failed to check existence of key %s: %v", key, err)
		return false, err
	}
	return exists > 0, nil
}

// Expire 设置一个键的过期时间
func (b *Broker) Expire(key string, seconds int64) error {
	err := b.client.Expire(b.ctx, key, time.Duration(seconds)*time.Second).Err()
	if err != nil {
		log.Printf("Failed to set expiration for key %s: %v", key, err)
		return err
	}
	return nil
}
