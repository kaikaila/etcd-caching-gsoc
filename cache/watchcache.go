package cache

import (
	"sort"
	"sync"
)

type WatchCache struct {
	mu        sync.RWMutex
	store     map[string]*storeObj // The current latest key-value state snapshot
	globalRevision  int64             // The latest global revision, ensures order consistency
	eventSink EventSink         // Downstream sink (observer pattern)
}

func NewWatchCache(sink EventSink) *WatchCache {
	return &WatchCache{
		store:     make(map[string]*storeObj),
		eventSink: sink,
	}
}

// HandlePut is a convenience wrapper that accepts string values.
// For performance-sensitive scenarios, use HandlePutBytes instead.
func (w *WatchCache) HandlePut(key, val string, rev int64) {
    w.HandlePutBytes(key, []byte(val), rev)
}

// HandlePutBytes is the high-throughput version of HandlePut that accepts raw byte slices.
// It avoids extra string<->[]byte conversions for high-frequency workloads.
func (w *WatchCache) HandlePutBytes(key string, valBytes []byte, rev int64) {
    w.mu.Lock()
    defer w.mu.Unlock()
    existing, ok := w.store[key]
	if ok && rev <= existing.Revision {
		return
	}
    w.store[key] = &storeObj{
        Key:      key,
        Value:    valBytes,
        Revision: rev,
    }

	if rev > w.globalRevision {
		w.globalRevision = rev
	}

    if w.eventSink != nil {
        w.eventSink.HandlePut(key, string(valBytes))
    }
}

// HandleDeleteBytes is the high-throughput version of HandleDelete that accepts raw data.
// It avoids extra overhead in delete operations.
func (w *WatchCache) HandleDeleteBytes(key string, rev int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	existing, ok := w.store[key]
	if ok && rev <= existing.Revision {
		return
	}

	delete(w.store, key)
	if rev > w.globalRevision {
		w.globalRevision = rev
	}
	if w.eventSink != nil {
		w.eventSink.HandleDelete(key)
	}
}

// HandleDelete is a convenience wrapper for deletion.
// For performance-sensitive scenarios, use HandleDeleteBytes instead.
func (w *WatchCache) HandleDelete(key string, rev int64) {
    w.HandleDeleteBytes(key, rev)
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
	return w.globalRevision
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

// NewSnapshotView builds a view from the cache snapshot.
func (wc *WatchCache) NewSnapshotView() *SnapshotView {
	snap := wc.Snapshot()
	items := make([]*storeObj, 0, len(snap))
	for _, obj := range snap {
	  items = append(items, obj)
	}
	// This is to sort by revision.  Optional for implementation: sort by keys
	sort.Slice(items, func(i, j int) bool {
	  return items[i].Revision < items[j].Revision
	})
	return &SnapshotView{Data: items}
  }