package main

import (
	"context"
	"fmt"
	"net/http"
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
