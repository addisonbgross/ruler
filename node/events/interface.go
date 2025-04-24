package events

// EventQueue represents a generic event queue interface
// used for recording the actions taken by the node.
type EventQueue interface {
	Push(data string) error
}
