package main

import (
	h "ruler-node/internal/http"
	s "ruler-node/internal/storage"
	t "ruler-node/internal/types"
	u "ruler-node/internal/util"

	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
)

func main() {
	portFlag := flag.String("port", "", "Exposed port of the node (optional)")
	flag.Parse()

	port := *portFlag
	if port == "" {
		port = "8080"
	}

	config, err := readConfig()
	if err != nil {
		panic(err)
	}

	// initialize data storage
	info := t.NodeInfo{Ip: "0.0.0.0", Port: port}

	members := t.MemberList{Members: config.Nodes}
	store := s.InMemoryStore{}
	h.InitHandlers(&store, &info, &members)

	u.InitLogger(&info)
	logger, err := u.GetLogger()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// data storage
	mux.HandleFunc("/read/", h.HandleRead)
	mux.HandleFunc("/write", h.HandleWrite)
	mux.HandleFunc("/delete", h.HandleDelete)
	mux.HandleFunc("/dump", h.HandleDump)

	// service discovery
	mux.HandleFunc("/members", h.HandleMembers)

	logger.Info(fmt.Sprintf("Node is listening on port %s", info.Port))
	listener := fmt.Sprintf(":%s", info.Port)
	http.ListenAndServe(listener, mux)
}

func readConfig() (t.NodeConfig, error) {
	config, err := os.ReadFile("./node-config.json")
	if err != nil {
		return t.NodeConfig{}, errors.New(fmt.Sprintf("Failed to read node-config.json: %s", err))
	}

	nodeConfig := t.NodeConfig{}
	err = json.Unmarshal(config, &nodeConfig)
	if err != nil {
		return t.NodeConfig{}, errors.New(fmt.Sprintf("Failed to parse node-config.json: %s", err))
	}

	return nodeConfig, nil
}

func getOwnNodeInfo(port string) t.NodeInfo {
	return t.NodeInfo{Ip: "0.0.0.0", Port: "8080"}
}
