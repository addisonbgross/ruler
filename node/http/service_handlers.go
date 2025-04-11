package http

import (
	u "Git/ruler/node/util"

	"net/http"
)

func HandleMembers(w http.ResponseWriter, r *http.Request) {
	sugar := u.GetLogger()

	method := r.Method
	if method == "POST" {

	} else if method == "GET" || method == "" {

	} else {
		sugar.Error("Unknown http method used with /members (not GET or POST)")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte{})
	}
}
