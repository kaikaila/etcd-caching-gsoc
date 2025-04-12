package cache

import "sync"

// memoryCache is a thread-safe in-memory implementation of the Cache interface.
type memoryCache struct {
	store map[string]string
	mu    sync.Mutex
}

// NewMemoryCache creates a new memoryCache instance.
func NewMemoryCache() Cache {
	return &memoryCache{
		store: make(map[string]string),
	}
}

func (m *memoryCache) Get(key string) (string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	val, ok := m.store[key]
	return val, ok
}

func (m *memoryCache) Set(key string, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store[key] = value
}

func (m *memoryCache) HandlePut(k, v string) {
	m.Set(k, v)
}

func (m *memoryCache) HandleDelete(k string) {
	m.Set(k, "") // TODO: support actual deletion
}