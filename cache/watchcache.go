package cache

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/kaikaila/etcd-caching-gsoc/cache/event"
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
	eventLog      event.EventLog
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
func NewWatchCacheWithLog(sink EventSink, log event.EventLog) *WatchCache {
	return &WatchCache{
		store:     make(map[string]*StoreObj),
		eventSink: sink,
		eventLog:  log,
	}
}

// HandlePut is a convenience wrapper that accepts string values.
// For performance-sensitive scenarios, use HandlePutBytes instead.
func (w *WatchCache) HandlePut(key, val string, globalRev int64) {
	w.HandlePutBytes(key, []byte(val), globalRev)
}

// HandlePutBytes is the high-throughput version of HandlePut that accepts raw byte slices.
// It avoids extra string<->[]byte conversions for high-frequency workloads.
func (w *WatchCache) HandlePutBytes(key string, valBytes []byte, globalRev int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	existing, ok := w.store[key]
	if ok && globalRev <= existing.GlobalRev {
		return
	}

	w.store[key] = &StoreObj{
		Key:      key,
		Value:    valBytes,
		GlobalRev: globalRev,
	}

	if globalRev > w.revision {
		w.revision = globalRev
	}

	if w.eventSink != nil {
		w.eventSink.HandlePut(key, string(valBytes))
	}
}

// HandleDeleteBytes is the high-throughput version of HandleDelete that accepts raw data.
// It avoids extra overhead in delete operations.
func (w *WatchCache) HandleDeleteBytes(key string, globalRev int64) {
	w.mu.Lock()
	defer w.mu.Unlock()

	existing, ok := w.store[key]
	if ok && globalRev <= existing.GlobalRev {
		return
	}

	delete(w.store, key)

	if globalRev > w.revision {
		w.revision = globalRev
	}
	if w.eventSink != nil {
		w.eventSink.HandleDelete(key)
	}
}

// HandleDelete is a convenience wrapper for deletion.
// For performance-sensitive scenarios, use HandleDeleteBytes instead.
func (w *WatchCache) HandleDelete(key string, globalRev int64) {
	w.HandleDeleteBytes(key, globalRev)
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

func (w *WatchCache) AddEvent(ev event.Event) error {
	switch ev.Type {
	case event.EventPut:
		w.HandlePutBytes(ev.Key, ev.Value, ev.GlobalRev)
	case event.EventDelete:
		w.HandleDeleteBytes(ev.Key, ev.GlobalRev)
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

// Snapshot returns a deep copy of the cache, in case external modification
func (wc *WatchCache) Snapshot() map[string]*StoreObj {
	wc.mu.RLock()
	defer wc.mu.RUnlock()

	snapshot := make(map[string]*StoreObj, len(wc.store))
	for k, v := range wc.store {
		snapshot[k] = v.DeepCopy()
	}
	return snapshot
}

// NewSnapshotView builds a view from the cache snapshot.
func (wc *WatchCache) NewSnapshotView() *SnapshotView {
	snap := wc.Snapshot()
	items := make([]*StoreObj, 0, len(snap))
	for _, obj := range snap {
		items = append(items, obj)
	}
	// This is to sort by revision.  Optional for implementation: sort by keys
	sort.Slice(items, func(i, j int) bool {
		return items[i].GlobalRev < items[j].GlobalRev
	})
	return &SnapshotView{Data: items}
}
