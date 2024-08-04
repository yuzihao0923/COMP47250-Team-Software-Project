package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func shutdownNode(client *redis.ClusterClient, port string) error {
	result := client.Do(context.Background(), "SHUTDOWN", "NOSAVE")
	if result.Err() != nil && result.Err() != redis.Nil {
		return result.Err()
	}
	fmt.Printf("Node on port %s has been shut down.\n", port)
	return nil
}

func main() {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: []string{
			"localhost:6381",
			"localhost:6382",
			"localhost:6383",
			"localhost:6384",
			"localhost:6385",
			"localhost:6386",
		},
	})

	ctx := context.Background()

	// Ping the cluster to ensure it's up and running
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis cluster: %v", err)
	}

	fmt.Println("Connected to Redis cluster.")

	// Get initial cluster state
	clusterInfoBefore, err := rdb.ClusterInfo(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to get initial cluster info: %v", err)
	}

	start := time.Now()
	err = shutdownNode(rdb, "6381")
	if err != nil {
		log.Fatalf("Failed to shut down node: %v", err)
	}

	// Wait for the failover to occur
	failoverWaitTime := 10 * time.Second
	fmt.Printf("Waiting %s for failover to complete...\n", failoverWaitTime)
	time.Sleep(failoverWaitTime)

	// Check for the new master
	newMasterFound := false
	var newMaster string
	for !newMasterFound {
		clusterInfoAfter, err := rdb.ClusterInfo(ctx).Result()
		if err != nil {
			log.Fatalf("Failed to get cluster info: %v", err)
		}

		if clusterInfoBefore != clusterInfoAfter {
			newMaster = extractNewMasterNode(clusterInfoBefore, clusterInfoAfter)
			if newMaster != "" {
				newMasterFound = true
			}
		}

		time.Sleep(1 * time.Second) // Check every second
	}

	// Measure the time taken for the failover
	elapsed := time.Since(start)
	fmt.Printf("Failover took %s\n", elapsed)

	// Output new master node information
	fmt.Printf("New master node after failover: %s\n", newMaster)
}

// Helper function to extract the new master node from cluster info
func extractNewMasterNode(infoBefore, infoAfter string) string {
	// Implement your logic here to extract the new master node from cluster info
	// This can be parsing the output of redis.ClusterInfo() or using other methods
	// Return the node ID or address of the new master node
	return "" // Placeholder
}
