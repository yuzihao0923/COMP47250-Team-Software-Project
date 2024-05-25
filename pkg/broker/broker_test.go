// pkg/broker/broker_test.go
package broker

import (
	"testing"
	"time"
)

func TestBroker(t *testing.T) {
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

	// 测试 Exists 方法
	exists, err := br.Exists("testkey")
	if err != nil {
		t.Fatalf("Failed to check existence of key: %v", err)
	}
	if !exists {
		t.Fatalf("Expected key 'testkey' to exist")
	}

	// 测试 Expire 方法
	err = br.Expire("testkey", 1)
	if err != nil {
		t.Fatalf("Failed to set expiration: %v", err)
	}

	// 等待 2 秒，确保键过期
	time.Sleep(2 * time.Second)

	// 测试 Expire 后 Get 方法，键应该不存在
	value, err = br.Get("testkey")
	if err != nil {
		t.Fatalf("Failed to get key after expiration: %v", err)
	}
	if value != "" {
		t.Fatalf("Expected key 'testkey' to be empty after expiration, got %s", value)
	}

	// 再次测试 Exists 方法，键应该不存在
	exists, err = br.Exists("testkey")
	if err != nil {
		t.Fatalf("Failed to check existence of key after expiration: %v", err)
	}
	if exists {
		t.Fatalf("Expected key 'testkey' to not exist after expiration")
	}

	// 测试 Del 方法
	err = br.Set("testkey", "testvalue")
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	err = br.Del("testkey")
	if err != nil {
		t.Fatalf("Failed to delete key: %v", err)
	}

	// 测试 Del 后 Get 方法，键应该不存在
	value, err = br.Get("testkey")
	if err != nil {
		t.Fatalf("Failed to get key after deletion: %v", err)
	}
	if value != "" {
		t.Fatalf("Expected key 'testkey' to be empty after deletion, got %s", value)
	}

	// 再次测试 Exists 方法，键应该不存在
	exists, err = br.Exists("testkey")
	if err != nil {
		t.Fatalf("Failed to check existence of key after deletion: %v", err)
	}
	if exists {
		t.Fatalf("Expected key 'testkey' to not exist after deletion")
	}
}
