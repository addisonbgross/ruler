package storage

// Store represents a generic storage interface that allows
// key-value pair operations, such as retrieving, storing,
// deleting, and iterating over data.
type Store interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string) bool
	Range() map[string]string
}
