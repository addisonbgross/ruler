package storage

type Store interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Range() map[string]string
}
