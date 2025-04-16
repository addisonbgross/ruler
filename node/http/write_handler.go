package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	rs "node/discovery"
	sh "node/shared"
	t "node/types"
	u "node/util"
	"os"
	"time"
)

// HandleWrite handles incoming requests to store a key-value pair.
// If the entry is not marked as a replication write, it asynchronously replicates
// the key-value pair across other nodes in the cluster.
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
			go func() {
				err := replicate(e.Key, e.Value, "write")
				if err != nil {
					sugar.Error(err)
				}
			}()
			sugar.Info(fmt.Sprintf("Wrote key(%s) - value(%s)", e.Key, e.Value))
		} else {
			sugar.Info(fmt.Sprintf("Wrote key(%s) - value(%s) - Replication", e.Key, e.Value))
		}

	}
}

// replicate can either send a key-value pair to all other nodes in the cluster to ensure data consistency,
// or send a delete request for the provided key. It skips replication to the current node.
func replicate(key, value, method string) error {
	sugar, err := u.GetLogger()
	if err != nil {
		return err
	}

	nodeHostname, err := os.Hostname()
	if err != nil {
		return err
	}

	allHostnames, err := rs.GetAllNodeHostnames()
	if err != nil {
		return err
	}

	for _, nextHostname := range allHostnames {
		if nextHostname == nodeHostname {
			sugar.Info("Skipping replication to self...")
			continue
		}

		sBody := fmt.Sprintf(`{"key":"%s","value":"%s","isreplicate":true}`, key, value)
		jBody := []byte(sBody)
		reqBody := bytes.NewReader(jBody)

		url := fmt.Sprintf("http://%s:8080/%s", nextHostname, method)
		replicationHttpMethod := http.MethodPost // TODO make this cleaner
		if method == "delete" {
			replicationHttpMethod = http.MethodDelete
		}
		req, err := http.NewRequest(replicationHttpMethod, url, reqBody)
		if err != nil {
			sugar.Errorf("(%s) failed to prepare replication request for key(%s) - value(%s)", nodeHostname, key, value)
			continue
		}

		req.Header.Set("Content-Type", "application/json")
		client := http.Client{Timeout: 10 * time.Second}
		_, err = client.Do(req)
		if err != nil {
			sugar.Errorf("(%s) failed to replicate key(%s) - value(%s) to %s", nodeHostname, key, value, url)
			sugar.Error(err)
		}
	}

	return nil
}
