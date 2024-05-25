// pkg/broker/broker_test.go
package broker

import (
	"testing"
)

func TestBroker(t *testing.T) {
	// 假设 Redis 在本地运行
	br := NewBroker("localhost:6379", "", 0)

	// 测试 Set 方法
	err := br.Set("testkey", "testvalue")
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// 测试 Get 方法
	value, err := br.Get("testkey")
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}

	if value != "testvalue" {
		t.Fatalf("Expected 'testvalue', got %s", value)
	}
}
