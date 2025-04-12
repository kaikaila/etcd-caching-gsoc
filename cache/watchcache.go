package cache

import "sync"

type WatchCache struct {
	mu        sync.RWMutex
	state     map[string]string // The current latest key-value state snapshot
	revision  int64             // The latest revision, ensures order consistency
	eventSink EventSink         // Downstream sink (observer pattern)
}

func NewWatchCache(sink EventSink) *WatchCache {
	return &WatchCache{
		state:     make(map[string]string),
		eventSink: sink,
	}
}

func (w *WatchCache) HandlePut(key, val string, rev int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if rev <= w.revision {
		return // Ignore stale events
	}

	w.state[key] = val
	w.revision = rev
	if w.eventSink != nil {
		w.eventSink.HandlePut(key, val)
	}
}

func (w *WatchCache) HandleDelete(key string, rev int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if rev <= w.revision {
		return
	}

	delete(w.state, key)
	w.revision = rev
	if w.eventSink != nil {
		w.eventSink.HandleDelete(key)
	}
}

func (w *WatchCache) Get(key string) (string, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	val, ok := w.state[key]
	return val, ok
}

func (w *WatchCache) Revision() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.revision
}