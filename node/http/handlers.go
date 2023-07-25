package http

import (
	s "Git/ruler/node/storage"
	t "Git/ruler/node/types"
	u "Git/ruler/node/util"

	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var store s.InMemoryStore

func HandleRead(w http.ResponseWriter, r *http.Request) {
	sugar := u.GetLogger()

	const readPath = "/read/"
	var key string
	if strings.HasPrefix(r.URL.Path, readPath) {
		key = r.URL.Path[len(readPath):]
	} else {
		sugar.Error("Missing key in Read request. Needs -> 'ip:port/read/myKey'")
		return
	}

	value, ok := store.Get(key)
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

func HandleWrite(w http.ResponseWriter, r *http.Request) {
	sugar := u.GetLogger()

	dec := json.NewDecoder(r.Body)
	var e t.StoreEntry
	err := dec.Decode(&e)
	if err != nil {
		sugar.Error("Can't decode Write payload")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
	} else {
		store.Set(e.Key, e.Value)
		sugar.Info(fmt.Sprintf("Wrote key(%s) - value(%s)", e.Key, e.Value))
	}
}

func HandleDump(w http.ResponseWriter, r *http.Request) {
	sugar := u.GetLogger()

	resp := []t.StoreEntry{}
	for k, v := range store.Range() {
		resp = append(resp, t.StoreEntry{Key: k, Value: v})
	}

	jResp, err := json.Marshal(resp)
	if err != nil {
		sugar.Error("Failed to marshall store data into /dump response")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jResp)
}
