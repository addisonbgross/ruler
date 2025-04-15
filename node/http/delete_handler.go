package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	sh "ruler/node/shared"
	t "ruler/node/types"
	u "ruler/node/util"
)

// HandleDelete handles HTTP DELETE requests to remove a stored key-value entry.
// If the key is successfully deleted, the function triggers an asynchronous replication
// process to propagate the delete operation to other nodes.
func HandleDelete(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// only permit the DELETE method
	if r.Method != http.MethodDelete {
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

	dec := json.NewDecoder(r.Body)
	var e t.StoreEntry
	err = dec.Decode(&e)
	if err != nil {
		sugar.Error("Can't decode Delete payload")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
	} else {
		ok := sh.Store.Delete(e.Key)
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
