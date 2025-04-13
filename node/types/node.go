package types

type NodeInfo struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type MemberList struct {
	Members []NodeInfo
}
