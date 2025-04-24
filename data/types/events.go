package types

import "time"

type NodeActionEventResponse struct {
	Hostname  string            `json:"hostname"`
	Type      ActionType        `json:"type"`
	Data      map[string]string `json:"data"`
	CreatedAt time.Time         `json:"created_at"`
}

type ActionType int

const (
	NodeStarted ActionType = iota
	NodeStopped
	Write
	Read
	Delete
	Replication
)

func (a ActionType) String() string {
	switch a {
	case NodeStarted:
		return "NodeStarted"
	case NodeStopped:
		return "NodeStopped"
	case Write:
		return "Write"
	case Read:
		return "Read"
	case Delete:
		return "Delete"
	case Replication:
		return "Replication"
	default:
		return "Unknown"
	}
}
