package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	rs "ruler/node/discovery"
	h "ruler/node/http"
	sh "ruler/node/shared"
	s "ruler/node/storage"
	u "ruler/node/util"
)

func main() {
	id, err := getNodeReplicaIdentifier()
	if err != nil {
		panic(err)
	}
	sh.NodeID = fmt.Sprintf("ruler-node-%d", id)
	setPublicHostname(id)

	sh.Store = s.InMemoryStore{}

	logger, err := u.GetLogger()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// data storage
	mux.HandleFunc("/read/", h.HandleRead)
	mux.HandleFunc("/write", h.HandleWrite)
	mux.HandleFunc("/delete", h.HandleDelete)
	mux.HandleFunc("/dump", h.HandleDump)

	// service discovery
	mux.HandleFunc("/health", h.HandleHealth)

	logger.Info("New node is listening on port 8080")
	http.ListenAndServe(":8080", mux)

	// close shared Redis connection
	client := rs.GetRedisClient()
	err = client.Close()
	if err != nil {
		logger.Error("Failed to close Redis client")
	}
}

func getNodeReplicaIdentifier() (int64, error) {
	client := rs.GetRedisClient()
	ctx := context.Background()

	// Try to increment a counter and use that as our replica number
	replicaNumber, err := client.Incr(ctx, "ruler-node-counter").Result()
	if err != nil {
		return -1, err
	}

	return replicaNumber, nil
}

func setPublicHostname(id int64) {
	client := rs.GetRedisClient()
	ctx := context.Background()

	key := fmt.Sprintf("node-hostname:ruler-node-%d", id)
	log.Default().Printf("Setting hostname to %s", key)

	hostname, err := os.Hostname()
	if err != nil {
		log.Default().Printf("Failed to get hostname: %v", err)
		return
	}

	client.Set(ctx, key, hostname, 0)
}
