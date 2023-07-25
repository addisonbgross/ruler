package storage

type Store interface {
	Get() (string, bool)
	Set(key string, value string)
	Range() map[string]string
}
