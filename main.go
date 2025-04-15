package main

import (
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
	hostname, err := os.Hostname()
	if err != nil {
		log.Default().Printf("Failed to get hostname: %v", err)
		return
	}

	err = rs.RegisterReplica(hostname)
	if err != nil {
		panic(err)
	}

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
