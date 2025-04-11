package main

import (
	h "Git/ruler/node/http"
	s "Git/ruler/node/storage"
	t "Git/ruler/node/types"
	u "Git/ruler/node/util"
	"encoding/json"
	"path/filepath"
	"runtime"

	"fmt"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		panic("Not enough params, needs: <ip> <port>")
	}

	// start params
	ip := os.Args[1]
	port := os.Args[2]

	config := readConfig()

	// initialize data storage
	info, ok := getOwnNodeInfo(&config, ip, port)
	if !ok {
		panic("Unable to get own node info from config file")
	}

	members := t.MemberList{Members: config.Nodes}
	store := s.InMemoryStore{}
	h.InitHandlers(&store, &info, &members)

	u.InitLogger(&info)
	sugar := u.GetLogger()
	mux := http.NewServeMux()

	// data storage
	mux.HandleFunc("/read/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte{})
			return
		}
		h.HandleRead(w, r)
	})
	mux.HandleFunc("/write", h.HandleWrite)
	mux.HandleFunc("/delete", h.HandleDelete)
	mux.HandleFunc("/dump", h.HandleDump)

	// service discovery
	mux.HandleFunc("/members", h.HandleMembers)

	sugar.Info(fmt.Sprintf("Node is listening on port %s", port))
	listener := fmt.Sprintf(":%s", port)
	http.ListenAndServe(listener, mux)
}

func readConfig() t.NodeConfig {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "./node-config.json")
	// config, err := os.ReadFile("./node-config.json")
	config, err := os.ReadFile(root)
	if err != nil {
		panic("Failed to read node-config.json")
	}

	nodeConfig := t.NodeConfig{}
	err = json.Unmarshal(config, &nodeConfig)
	if err != nil {
		panic("Failed to parse node-config.json")
	}

	return nodeConfig
}

func getOwnNodeInfo(c *t.NodeConfig, ip, port string) (t.NodeInfo, bool) {
	for _, node := range c.Nodes {
		if node.Ip == ip && node.Port == port {
			return t.NodeInfo{Ip: ip, Port: port, Rank: node.Rank}, true
		}
	}

	panic(fmt.Sprintf("Unable to find own node info in config file for %s:%s", ip, port))
}
