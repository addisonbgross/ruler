package main

import (
	"context"
	"net/http"
	d "node/discovery"
	e "node/events"
	h "node/http"
	r "node/redis"
	sh "node/shared"
	s "node/storage"
	t "node/types"
	u "node/util"
	w "node/workers"
	"os"
	"os/signal"
	"syscall"
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
		logger.Error("No connection to Redis. This node will not be discoverable!")
	}

	wp, err := w.GetWorkerPool()
	if err != nil {
		panic(err)
	}

	wp.Submit(t.NodeActionEvent{
		Hostname: hostname,
		Type:     t.NodeStarted,
		Data:     map[string]string{"test": "worker pool data"},
	})

	// Context for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// ensure that the Redis & Postgres connections, and the worker pool,
	// are closed when the program exits
	defer r.CloseClient()
	defer e.CloseEventQueue()
	defer wp.Close()

	// initialize the key/value store
	// TODO: enable other storage mediums
	sh.Store = s.InMemoryStore{}

	mux := http.NewServeMux()
	mux.HandleFunc("/read/", h.HandleRead)
	mux.HandleFunc("/write", h.HandleWrite)
	mux.HandleFunc("/delete", h.HandleDelete)
	mux.HandleFunc("/dump", h.HandleDump)
	mux.HandleFunc("/health", h.HandleHealth)

	go func() {
		logger.Info("New node is listening on port 8080")
		err = http.ListenAndServe(":8080", mux)
		if err != nil {
			panic(err)
		}
	}()

	// Wait for a shutdown signal
	<-ctx.Done()
	stop()
	logger.Info("Shutting down gracefully...")
}
