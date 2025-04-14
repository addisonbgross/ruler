package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	sh "ruler/node/shared"
	t "ruler/node/types"
	u "ruler/node/util"
	"strconv"
	"time"
)

func HandleWrite(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	sugar, err := u.GetLogger()
	if err != nil {
		log.Print("Failed to get logger for HandleWrite")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return
	}

	dec := json.NewDecoder(r.Body)
	var e t.StoreEntry
	err = dec.Decode(&e)
	if err != nil {
		sugar.Error("Can't decode Write payload")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
	} else {
		sh.Store.Set(e.Key, e.Value)

		if !e.IsReplicate {
			// TODO add dead-letter-queue for failed replications
			go replicate(e.Key, e.Value, "write")
			sugar.Info(fmt.Sprintf("Wrote key(%s) - value(%s) to", e.Key, e.Value))
		} else {
			sugar.Info(fmt.Sprintf("Wrote key(%s) - value(%s) - Replication to", e.Key, e.Value))
		}

	}
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

	for i = 1; i <= maxCounterInt; i++ {
		sBody := fmt.Sprintf(`{"key":"%s","value":"%s","isreplicate":true}`, key, value)
		jBody := []byte(sBody)
		reqBody := bytes.NewReader(jBody)

		nodeUrl := fmt.Sprintf("%s-%d", "ruler-node", i)
		url := fmt.Sprintf("http://%s-%d:%s/%s", "ruler-node", i, "8080", method)
		if nodeUrl == sh.NodeID {
			sugar.Info(fmt.Sprintf("Skipping replication to self: %s", url))
			continue
		}

		req, err := http.NewRequest(http.MethodPost, url, reqBody)
		if err != nil {
			sugar.Errorf("(%s) failed to prepare replication request for key(%s) - value(%s)", sh.NodeID, key, value)
		}

		req.Header.Set("Content-Type", "application/json")
		client := http.Client{Timeout: 10 * time.Second}
		_, err = client.Do(req)
		if err != nil {
			sugar.Errorf("(%s) failed to replicate key(%s) - value(%s) to %s", sh.NodeID, key, value, url)
			sugar.Error(err)
		}
	}
}
