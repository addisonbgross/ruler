package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	sh "node/shared"
	t "node/types"
	u "node/util"
	"strings"
)

// HandleRead processes HTTP GET requests to retrieve the value associated with a given key
// from the store. The key is extracted from the URL path in the form "/read/{key}".
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

	sugar.Info(fmt.Sprintf("Read key: '%s' with value: '%s'", key, value))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(value))
}

// HandleDump processes HTTP GET requests to return a JSON representation
// of all key-value pairs in the store.
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

	var resp []t.StoreEntry
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
