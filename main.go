package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/redis/go-redis/v9"
	"net/http"
	"os"
	h "ruler/node/http"
	s "ruler/node/storage"
	t "ruler/node/types"
	u "ruler/node/util"
)

func main() {
	portFlag := flag.String("port", "", "Exposed port of the node (optional)")
	flag.Parse()

	port := *portFlag
	if port == "" {
		port = "8080"
	}

	config, err := readConfig()
	if err != nil {
		panic(err)
	}

	// initialize data storage
	hostname := "node"
	info := t.NodeInfo{Ip: hostname, Port: port}

	members := t.MemberList{Members: config.Nodes}
	store := s.InMemoryStore{}
	h.InitHandlers(&store, &info, &members)

	u.InitLogger(&info)
	logger, err := u.GetLogger()
	if err != nil {
		panic(err)
	}

	u.InitLogger(&info)

	ress := getNodeIdentifier()
	logger.Info(ress)

	mux := http.NewServeMux()

	// data storage
	mux.HandleFunc("/read/", h.HandleRead)
	mux.HandleFunc("/write", h.HandleWrite)
	mux.HandleFunc("/delete", h.HandleDelete)
	mux.HandleFunc("/dump", h.HandleDump)

	// service discovery
	mux.HandleFunc("/health", h.HandleHealth)

	logger.Info(fmt.Sprintf("Node is listening on port %s", info.Port))
	listener := fmt.Sprintf(":%s", info.Port)
	http.ListenAndServe(listener, mux)
}

func readConfig() (t.NodeConfig, error) {
	config, err := os.ReadFile("./node-config.json")
	if err != nil {
		return t.NodeConfig{}, errors.New(fmt.Sprintf("Failed to read node-config.json: %s", err))
	}

	nodeConfig := t.NodeConfig{}
	err = json.Unmarshal(config, &nodeConfig)
	if err != nil {
		return t.NodeConfig{}, errors.New(fmt.Sprintf("Failed to parse node-config.json: %s", err))
	}

	return nodeConfig, nil
}

func getNodeIdentifier() string {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "redis"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
	})
	defer client.Close()

	ctx := context.Background()

	// Get our container ID for uniqueness
	hostname, _ := os.Hostname()

	// Try to increment a counter and use that as our replica number
	replicaNumber, err := client.Incr(ctx, "ruler-node-counter").Result()
	if err != nil {
		return fmt.Sprintf("ruler-node-unknown-%s", hostname[:8])
	}

	return fmt.Sprintf("Got the result -> ruler-node-%d", replicaNumber)
}
