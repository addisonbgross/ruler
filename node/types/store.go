package types

type StoreEntry struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	IsReplicate bool   `json:"isreplicate"`
}
