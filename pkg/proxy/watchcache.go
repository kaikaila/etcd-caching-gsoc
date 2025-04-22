package proxy

import (
	"errors"
	"fmt"
	"sync"

	"github.com/kaikaila/etcd-caching-gsoc/pkg/api"
	"github.com/kaikaila/etcd-caching-gsoc/pkg/eventlog"
)

var (
	ErrKeyNotFound     = errors.New("key not found in WatchCache")
	ErrInvalidRevision = errors.New("revision too old or newer than current")
)

type WatchCache struct {
	mu            sync.RWMutex
	store         map[string]*StoreObj // The current latest key-value state snapshot
	revision      int64                 // revision tracks the total number of write operations across all keys.
	eventSink     EventSink             // Downstream sink (observer pattern)
	eventLog      eventlog.EventLog
	// Optional: If we need to analyze key write frequency, enable eviction policies,
	// or track the most updated key, consider adding:
	// MaxPerKeyRevision int64 // highest key-local revision among all keys
}

func NewWatchCache(sink EventSink) *WatchCache {
	return &WatchCache{
		store:     make(map[string]*StoreObj),
		eventSink: sink,
	}
}

// NewWatchCacheWithLog creates a WatchCache with an optional event log sink.
func NewWatchCacheWithLog(sink EventSink, log eventlog.EventLog) *WatchCache {
	return &WatchCache{
		store:     make(map[string]*StoreObj),
		eventSink: sink,
		eventLog:  log,
	}
}

// HandlePut is a convenience wrapper that accepts string values.
// For performance-sensitive scenarios, use HandlePutBytes instead.
func (w *WatchCache) HandlePut(key, val string, Revision int64) {
	w.HandlePutBytes(key, []byte(val), Revision)
}

// HandlePutBytes is the high-throughput version of HandlePut that accepts raw byte slices.
// It avoids extra string<->[]byte conversions for high-frequency workloads.
func (w *WatchCache) HandlePutBytes(key string, valBytes []byte, Revision int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	existing, ok := w.store[key]
	if ok && Revision <= existing.Revision {
		return
	}

	w.store[key] = &StoreObj{
		Key:      key,
		Value:    valBytes,
		Revision: Revision,
	}

	if Revision > w.revision {
		w.revision = Revision
	}

	if w.eventSink != nil {
		w.eventSink.HandlePut(key, string(valBytes))
	}
}

// HandleDeleteBytes is the high-throughput version of HandleDelete that accepts raw data.
// It avoids extra overhead in delete operations.
func (w *WatchCache) HandleDeleteBytes(key string, Revision int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	existing, ok := w.store[key]
	if ok && Revision <= existing.Revision {
		return
	}

	delete(w.store, key)

	if Revision > w.revision {
		w.revision = Revision
	}
	if w.eventSink != nil {
		w.eventSink.HandleDelete(key)
	}
}

// HandleDelete is a convenience wrapper for deletion.
// For performance-sensitive scenarios, use HandleDeleteBytes instead.
func (w *WatchCache) HandleDelete(key string, Revision int64) {
	w.HandleDeleteBytes(key, Revision)
}

// Get returns a deep copy of the StoreObj associated with the key
func (w *WatchCache) Get(key string) (*StoreObj, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	obj, ok := w.store[key]
	if !ok {
		return nil, false
	}
	// return a deep copy so caller cannot mutate internal state
	return obj.DeepCopy(), true
}

func (w *WatchCache) AddEvent(ev api.Event) error {
	switch ev.Type {
	case api.EventPut:
		w.HandlePutBytes(ev.Key, ev.Value, ev.Revision)
	case api.EventDelete:
		w.HandleDeleteBytes(ev.Key, ev.Revision)
	default:
		return fmt.Errorf("unsupported event type: %v", ev.Type)
	}
	if w.eventLog != nil {
		return w.eventLog.Append(ev)
	}
	return nil
}

func (w *WatchCache) Revision() int64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.revision
}

// Snapshot returns a SnapshotView over the current cache state.
func (wc *WatchCache) Snapshot() api.SnapshotView {
	wc.mu.RLock()
	defer wc.mu.RUnlock()

	// Build map for cacheView
	m := make(map[string]*StoreObj, len(wc.store))
	for key, obj := range wc.store {
		m[key] = obj.DeepCopy()
	}
	// Delegate to clientlibrary.NewCacheView for consistent view
	return clientlibrary.cacheView{data:m, revision:wc.revision}
}
