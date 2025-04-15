package discovery

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

var Client *redis.Client

// RegisterNode increments the node counter in Redis and
// registers a new hostname with the format "node-hostname:<hostname>".
func RegisterNode(hostname string) error {
	client := GetRedisClient()
	ctx := context.Background()

	_, err := client.Incr(ctx, "ruler-node-counter").Result()
	if err != nil {
		return err
	}

	key := fmt.Sprintf("node-hostname:%s", hostname)
	client.Set(ctx, key, hostname, 0)

	return nil
}

// GetAllNodeHostnames retrieves all registered node hostnames
// from Redis. It scans for keys with the pattern "node-hostname:*".
func GetAllNodeHostnames() ([]string, error) {
	client := GetRedisClient()
	ctx := context.Background()
	numRulerNodes, err := client.Get(ctx, "ruler-node-counter").Result()
	if err != nil {
		return nil, err
	}

	numRulerNodesInt, err := strconv.Atoi(numRulerNodes)
	if err != nil {
		return nil, err
	}

	allHostnameKeys, _, err := client.Scan(ctx, 0, "node-hostname:*", int64(numRulerNodesInt+1)).Result()
	if err != nil {
		return nil, err
	}

	pipe := client.Pipeline()
	allHostnames := make([]*redis.StringCmd, 0, len(allHostnameKeys))

	for _, key := range allHostnameKeys {
		allHostnames = append(allHostnames, pipe.Get(ctx, key))
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	results := make([]string, 0, len(allHostnames))
	for _, hostnamePipeResult := range allHostnames {
		nextHostname, err := hostnamePipeResult.Result()
		if err != nil {
			return nil, err
		}

		results = append(results, nextHostname)
	}

	return results, nil
}

// CloseClient releases the Redis client connection if it is initialized.
func CloseClient() error {
	if Client != nil {
		err := Client.Close()
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

// GetRedisClient initializes or retrieves a Redis client instance.
// If the client is already initialized, it returns the existing client.
func GetRedisClient() *redis.Client {
	if Client != nil {
		return Client
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	Client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	return Client
}
