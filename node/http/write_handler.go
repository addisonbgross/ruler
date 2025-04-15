package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	rs "ruler/node/discovery"
	sh "ruler/node/shared"
	t "ruler/node/types"
	u "ruler/node/util"
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

func replicate(key, value, method string) error {
	sugar, err := u.GetLogger()
	if err != nil {
		return err
	}

	nodeHostname, err := os.Hostname()
	if err != nil {
		return err
	}

	allHostnames, err := rs.GetAllReplicaHostnames()
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
		req, err := http.NewRequest(http.MethodPost, url, reqBody)
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
