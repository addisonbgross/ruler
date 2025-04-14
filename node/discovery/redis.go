package discovery

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
)

var Client *redis.Client

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
