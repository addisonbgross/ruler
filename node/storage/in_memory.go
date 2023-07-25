package storage

var storage = map[string]string{}

type InMemoryStore struct{}

func (s InMemoryStore) Get(key string) (string, bool) {
	v, ok := storage[key]
	return v, ok
}

func (s InMemoryStore) Set(key string, value string) {
	storage[key] = value
}

func (s InMemoryStore) Range() map[string]string {
	copy := map[string]string{}
	for k, v := range storage {
		copy[k] = v
	}
	return copy
}
