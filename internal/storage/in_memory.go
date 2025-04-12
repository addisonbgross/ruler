package storage

import t "ruler-node/internal/types"

var storage = map[string]string{}

type InMemoryStore struct {
	Members t.MemberList
}

func (s InMemoryStore) Get(key string) (string, bool) {
	v, ok := storage[key]
	return v, ok
}

func (s InMemoryStore) Set(key string, value string) {
	storage[key] = value
}

func (s InMemoryStore) Delete(key string) bool {
	_, ok := storage[key]
	delete(storage, key)
	return ok
}

func (s InMemoryStore) Range() map[string]string {
	copy := map[string]string{}
	for k, v := range storage {
		copy[k] = v
	}
	return copy
}
