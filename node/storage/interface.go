package storage

type Store interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string) bool
	Range() map[string]string
}
