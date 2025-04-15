package types

// StoreEntry represents a key-value pair with an optional flag for replication status.
type StoreEntry struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	IsReplicate bool   `json:"isreplicate"`
}
