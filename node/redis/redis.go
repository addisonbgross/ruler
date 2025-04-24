package redis

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"node/util"
	"os"
)

var client *redis.Client

// GetRedisClient initializes or retrieves a Redis client instance.
// If the client is already initialized, it returns the existing client.
func GetRedisClient() *redis.Client {
	if client != nil {
		return client
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})

	return client
}

func CloseClient() {
	if client != nil {
		err := client.Close()
		if err != nil {
			logger, _ := util.GetLogger()
			logger.Error("Failed to close Redis client")
		}
	}
}
