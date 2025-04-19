package cache

import "sync"

type WatchCache struct {
	mu        sync.RWMutex
	store     map[string]*storeObj // The current latest key-value state snapshot
	revision  int64             // The latest revision, ensures order consistency
	eventSink EventSink         // Downstream sink (observer pattern)
}

func NewWatchCache(sink EventSink) *WatchCache {
	return &WatchCache{
		store:     make(map[string]*storeObj),
		eventSink: sink,
	}
}

func (w *WatchCache) HandlePut(key, val string, rev int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if rev <= w.revision {
		return // Ignore stale events
	}

	w.store[key] = &storeObj{
		Key:      key,
		Value:    []byte(val),    // store the value as bytes
		Revision: rev,
	}
	w.revision = rev
	if w.eventSink != nil {
		w.eventSink.HandlePut(key, val)
	}
}


func (w *WatchCache) HandlePutBytes(key string, valBytes []byte, rev int64) {
    w.mu.Lock()
    defer w.mu.Unlock()
    // …和 HandlePut 一样，只是参数类型不同…
    if w.eventSink != nil {
        w.eventSink.HandlePut(key, string(valBytes))
    }
}

func (w *WatchCache) HandleDelete(key string, rev int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if rev <= w.revision {
		return
	}

	delete(w.store, key)
	w.revision = rev
	if w.eventSink != nil {
		w.eventSink.HandleDelete(key)
	}
}

// Get returns a deep copy of the storeObj associated with the key
func (w *WatchCache) Get(key string) (*storeObj, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	obj, ok := w.store[key]
	if !ok {
		return nil, false
	}
	// return a deep copy so caller cannot mutate internal state
	return obj.DeepCopy(), true
}

func (w *WatchCache) Revision() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.revision
}

// Snapshot returns a deep copy of the cache, in case external modification
func (wc *WatchCache) Snapshot() map[string]*storeObj {
    wc.mu.RLock()
    defer wc.mu.RUnlock()

    snapshot := make(map[string]*storeObj, len(wc.store))
    for k, v := range wc.store {
        snapshot[k] = v.DeepCopy()
    }
    return snapshot
}