package http

import (
	"net/http"
)

func HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}
	defer r.Body.Close()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
