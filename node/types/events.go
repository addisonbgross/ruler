package types

type NodeActionEvent struct {
	Hostname string            `json:"hostname"`
	Type     ActionType        `json:"type"`
	Data     map[string]string `json:"action"`
}

type ActionType int

const (
	NodeStarted ActionType = iota
	NodeStopped
	Write
	WriteReplication
	Read
	Delete
	DeleteReplication
)

func (a ActionType) String() string {
	switch a {
	case NodeStarted:
		return "NodeStarted"
	case NodeStopped:
		return "NodeStopped"
	case Write:
		return "Write"
	case WriteReplication:
		return "WriteReplication"
	case Read:
		return "Read"
	case Delete:
		return "Delete"
	case DeleteReplication:
		return "DeleteReplication"
	default:
		return "Unknown"
	}
}
