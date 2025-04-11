package types

type NodeInfo struct {
	Rank int    `json:"rank"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type MemberList struct {
	Members []NodeInfo
}
