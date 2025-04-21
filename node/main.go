package main

import (
	"log"
	"net/http"
	rs "node/discovery"
	h "node/http"
	sh "node/shared"
	s "node/storage"
	u "node/util"
	"os"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	// register this node to the Redis discovery service
	err = rs.RegisterNode(hostname)
	if err != nil {
		log.Default().Printf("No connection to Redis. This node will not be discoverable!")
	}
	defer rs.CloseClient()

	// initialize the key/value store
	// TODO: enable other storage mediums
	sh.Store = s.InMemoryStore{}

	logger, err := u.GetLogger()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/read/", h.HandleRead)
	mux.HandleFunc("/write", h.HandleWrite)
	mux.HandleFunc("/delete", h.HandleDelete)
	mux.HandleFunc("/dump", h.HandleDump)
	mux.HandleFunc("/health", h.HandleHealth)

	logger.Info("New node is listening on port 8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
