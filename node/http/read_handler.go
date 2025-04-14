package http

import (
	"encoding/json"
	"log"
	"net/http"
	sh "ruler/node/shared"
	t "ruler/node/types"
	u "ruler/node/util"
	"strings"
)

func HandleRead(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	sugar, err := u.GetLogger()
	if err != nil {
		log.Print("Failed to get logger for HandleRead")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}

	const readPath = "/read/"
	var key string
	if strings.HasPrefix(r.URL.Path, readPath) {
		key = r.URL.Path[len(readPath):]
	} else {
		sugar.Error("Missing key in Read request. Needs -> 'ip:port/read/myKey'")
		return
	}

	value, ok := sh.Store.Get(key)
	if !ok {
		sugar.Info("Missing value for key: ", key)
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})
		return
	}

	sugar.Info("Read key: ", key)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

func HandleDump(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	sugar, err := u.GetLogger()
	if err != nil {
		log.Print("Failed to get logger for HandleRead")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}

	resp := []t.StoreEntry{}
	for k, v := range sh.Store.Range() {
		resp = append(resp, t.StoreEntry{Key: k, Value: v})
	}

	jResp, err := json.Marshal(resp)
	if err != nil {
		sugar.Error("Failed to marshall store data into /dump response")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jResp)
}
