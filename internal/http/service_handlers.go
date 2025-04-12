package http

import (
	"log"
	u "ruler-node/internal/util"

	"net/http"
)

func HandleMembers(w http.ResponseWriter, r *http.Request) {
	sugar, err := u.GetLogger()
	if err != nil {
		log.Print("Failed to get logger for HandleMembers")
	}
	defer r.Body.Close()

	method := r.Method
	if method == "POST" {

	} else if method == "GET" || method == "" {

	} else {
		sugar.Error("Unknown http method used with /members (not GET or POST)")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
	}
}
