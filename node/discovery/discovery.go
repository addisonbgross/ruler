package discovery

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	rs "node/redis"
	"strconv"
)

// RegisterNode increments the node counter in Redis and
// registers a new hostname with the format "node-hostname:<hostname>".
func RegisterNode(hostname string) error {
	client := rs.GetRedisClient()
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
	client := rs.GetRedisClient()
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
