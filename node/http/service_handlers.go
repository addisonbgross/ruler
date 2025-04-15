package http

import (
	"net/http"
)

// HandleHealth responds to a health check HTTP request.
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
