package http

import (
	"bytes"
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	s "ruler/internal/storage"
	t "ruler/internal/types"
	u "ruler/internal/util"
	"strconv"
	"time"

	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// data storage

var store s.Store
var info *t.NodeInfo
var members *t.MemberList

func InitHandlers(s s.Store, n *t.NodeInfo, m *t.MemberList) {
	store = s
	info = n
	members = m
}

func HandleRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	sugar, err := u.GetLogger()
	if err != nil {
		log.Print("Failed to get logger for HandleRead")
	}
	defer r.Body.Close()

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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	sugar, err := u.GetLogger()
	if err != nil {
		log.Print("Failed to get logger for HandleWrite")
	}
	sugar.Info("Write request received")
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	var e t.StoreEntry
	err = dec.Decode(&e)
	if err != nil {
		sugar.Error("Can't decode Write payload")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
	} else {
		store.Set(e.Key, e.Value)

		if !e.IsReplicate {
			// TODO add dead-letter-queue for failed replications
			go replicate(e.Key, e.Value, "write")
			sugar.Info(fmt.Sprintf("Wrote key(%s) - value(%s)", e.Key, e.Value))
		} else {
			sugar.Info(fmt.Sprintf("Wrote key(%s) - value(%s) - Replication", e.Key, e.Value))
		}

	}
}

func HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	sugar, err := u.GetLogger()
	if err != nil {
		log.Print("Failed to get logger for HandleDelete")
	}
	defer r.Body.Close()

	dec := json.NewDecoder(r.Body)
	var e t.StoreEntry
	err = dec.Decode(&e)
	if err != nil {
		sugar.Error("Can't decode Delete payload")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
	} else {
		ok := store.Delete(e.Key)
		if !ok {
			sugar.Warnf("No key '%s' to delete", e.Key)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte{})
		}

		if !e.IsReplicate {
			// TODO add dead-letter-queue for failed replications
			go replicate(e.Key, "", "delete")
			sugar.Info(fmt.Sprintf("Deleted key(%s)", e.Key))
		} else {
			sugar.Info(fmt.Sprintf("Deleted key(%s) - Replication", e.Key))
		}
	}
}

func HandleDump(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	sugar, err := u.GetLogger()
	if err != nil {
		log.Print("Failed to get logger for HandleDump")
	}
	defer r.Body.Close()

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

func replicate(key, value, method string) {
	sugar, err := u.GetLogger()
	if err != nil {
		log.Print("Failed to get logger for Replicate")
		return
	}

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
	maxCounter, err := client.Get(ctx, "ruler-node-counter").Result()
	if err != nil {
		sugar.Error("Failed to get max counter")
		return
	}

	var i int
	maxCounterInt, err := strconv.Atoi(maxCounter)
	if err != nil {
		sugar.Errorf("Failed to convert maxCounter to int: %v", err)
		return
	}

	for i = 1; i < maxCounterInt; i++ {
		sBody := fmt.Sprintf(`{"key":"%s","value":"%s","isreplicate":true}`, key, value)
		jBody := []byte(sBody)
		reqBody := bytes.NewReader(jBody)

		url := fmt.Sprintf("http://%s-%d:%s/%s", "ruler-node-ruler-node", i, "8080", method)
		req, err := http.NewRequest(http.MethodPost, url, reqBody)
		if err != nil {
			sugar.Errorf("%s:%s failed to prepare replication '%s' request for key(%s) - value(%s)", info.Ip, info.Port, method, key, value)
		}

		req.Header.Set("Content-Type", "application/json")
		client := http.Client{Timeout: 10 * time.Second}
		_, err = client.Do(req)
		if err != nil {
			sugar.Errorf("%s:%s failed to replicate key(%s) - value(%s) to %s", info.Ip, info.Port, key, value, url)
			sugar.Error(err)
		}
	}
}
