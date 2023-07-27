package main

import (
	h "Git/ruler/node/http"
	s "Git/ruler/node/storage"
	t "Git/ruler/node/types"
	u "Git/ruler/node/util"
	"encoding/json"

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

	config, ok := readConfig()
	if !ok {
		panic("Invalid config file")
	}

	// initialize data storage
	info, ok := getOwnNodeInfo(&config, ip, port)
	if !ok {
		panic("Unable to get own node info from config file")
	}
	members := t.MemberList{Members: config.Nodes}
	store := s.InMemoryStore{}
	h.InitHandlers(&store, &info, &members)

	sugar := u.GetLogger()
	mux := http.NewServeMux()

	// data storage
	mux.HandleFunc("/read/", h.HandleRead)
	mux.HandleFunc("/write", h.HandleWrite)
	mux.HandleFunc("/dump", h.HandleDump)

	// service discovery
	mux.HandleFunc("/members", h.HandleMembers)

	sugar.Info(fmt.Sprintf("Node is listening on port %s", port))
	listener := fmt.Sprintf(":%s", port)
	http.ListenAndServe(listener, mux)
}

func readConfig() (t.NodeConfig, bool) {
	sugar := u.GetLogger()

	config, err := os.ReadFile("./node-config.json")
	if err != nil {
		sugar.Error(("Failed to read node-config.json"))
		return t.NodeConfig{}, false
	}

	nodeConfig := t.NodeConfig{}
	err = json.Unmarshal(config, &nodeConfig)
	if err != nil {
		sugar.Error("Failed to parse node-config.json")
		sugar.Error(err)
		return t.NodeConfig{}, false
	}

	return nodeConfig, true
}

func getOwnNodeInfo(c *t.NodeConfig, ip, port string) (t.NodeInfo, bool) {
	sugar := u.GetLogger()

	for _, node := range c.Nodes {
		if node.Ip == ip && node.Port == port {
			sugar.Infof("Node got own info: rank->%d", node.Rank)
			return t.NodeInfo{Ip: ip, Port: port, Rank: node.Rank}, true
		}
	}

	return t.NodeInfo{}, false
}
