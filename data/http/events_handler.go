package http

import (
	e "data/events"
	u "data/util"
	"encoding/json"
	"net/http"
)

func HandleEvents(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	logger, err := u.GetLogger()
	if err != nil {
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	actions, err := e.Read()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to read events from database"))
		return
	}

	data, err := json.Marshal(actions)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to serialize events"))
		return
	}

	logger.Info("Read all events")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
