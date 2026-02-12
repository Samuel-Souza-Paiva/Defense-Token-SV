package store

import (
	"errors"
	"sync"
)

var ErrKeyNotFound = errors.New("key not found")

type KeyValueStore interface {
	SetKey(key, value string) error
	GetKey(key string) (string, error)
}

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: map[string]string{},
	}
}

func (ms *MemoryStore) SetKey(key, value string) error {
	if key == "" {
		return errors.New("key cannot be empty")
	}
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.data[key] = value
	return nil
}

func (ms *MemoryStore) GetKey(key string) (string, error) {
	if key == "" {
		return "", errors.New("key cannot be empty")
	}
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	value, ok := ms.data[key]
	if !ok {
		return "", ErrKeyNotFound
	}
	return value, nil
}
