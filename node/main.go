package main

import (
	"log"
	"net/http"
	d "node/discovery"
	e "node/events"
	h "node/http"
	r "node/redis"
	sh "node/shared"
	s "node/storage"
	t "node/types"
	u "node/util"
	"os"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	logger, err := u.GetLogger()
	if err != nil {
		panic(err)
	}

	// register this node to the Redis discovery service
	err = d.RegisterNode(hostname)
	if err != nil {
		log.Default().Printf("No connection to Redis. This node will not be discoverable!")
	}

	err = e.Push(t.NodeActionEvent{
		Hostname: hostname,
		Type:     t.NodeStarted,
		Data:     map[string]string{},
	})
	if err != nil {
		panic(err)
	}

	// ensure that the Redis + Postgres connections are closed when the program exits
	defer r.CloseClient()
	defer e.CloseEventQueue()

	// initialize the key/value store
	// TODO: enable other storage mediums
	sh.Store = s.InMemoryStore{}

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
