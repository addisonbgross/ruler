package main

import (
	"data/events"
	h "data/http"
	u "data/util"
	"net/http"
)

func main() {
	logger, err := u.GetLogger()
	if err != nil {
		panic(err)
	}

	defer events.CloseEventQueue()

	mux := http.NewServeMux()
	mux.HandleFunc("/events/", h.HandleEvents)

	logger.Info("Data service is listening on port 8081")
	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		panic(err)
	}
}
