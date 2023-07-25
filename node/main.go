package main

import (
	h "Git/ruler/node/http"
	u "Git/ruler/node/util"

	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("Node not started with a port. It needs a port!")
	}

	port := os.Args[1]

	sugar := u.GetLogger()

	mux := http.NewServeMux()

	mux.HandleFunc("/read/", h.HandleRead)
	mux.HandleFunc("/write", h.HandleWrite)
	mux.HandleFunc("/dump", h.HandleDump)

	sugar.Info(fmt.Sprintf("Node is listening on port %s", port))
	listener := fmt.Sprintf(":%s", port)
	http.ListenAndServe(listener, mux)
}
