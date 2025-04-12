package main

import (
	"log"
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
	// start params
	ipFlag := flag.String("ip", "", "IP address of the node")
	portFlag := flag.String("port", "", "Port of the node")
	flag.Parse()

	if *ipFlag == "" || *portFlag == "" {
		log.Fatal("Not enough params, needs: <ip> <port>")
	}

	ip := *ipFlag
	port := *portFlag

	config, err := readConfig()
	if err != nil {
		panic(err)
	}

	// initialize data storage
	info, err := getOwnNodeInfo(&config, ip, port)
	if err != nil {
		panic(err)
	}

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
	mux.HandleFunc("/read/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte{})
			return
		}
		h.HandleRead(w, r)
	})
	mux.HandleFunc("/write", h.HandleWrite)
	mux.HandleFunc("/delete", h.HandleDelete)
	mux.HandleFunc("/dump", h.HandleDump)

	// service discovery
	mux.HandleFunc("/members", h.HandleMembers)

	logger.Info(fmt.Sprintf("Node is listening on port %s", port))
	listener := fmt.Sprintf(":%s", port)
	http.ListenAndServe(listener, mux)
}

func readConfig() (t.NodeConfig, error) {
	//_, b, _, _ := runtime.Caller(0)
	//root := filepath.Join(filepath.Dir(b), "./node-config.json")
	//config, err := os.ReadFile(root)
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

func getOwnNodeInfo(c *t.NodeConfig, ip, port string) (t.NodeInfo, error) {
	for _, node := range c.Nodes {
		if node.Ip == ip && node.Port == port {
			return t.NodeInfo{Ip: ip, Port: port, Rank: node.Rank}, nil
		}
	}

	return t.NodeInfo{}, errors.New(fmt.Sprintf("Unable to find own node info in config file for %s:%s", ip, port))
}
