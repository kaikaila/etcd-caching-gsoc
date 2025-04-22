// This file defines high-level interfaces based on the original Mermaid architecture diagram.
// It serves as the top-down blueprint to align all implementation efforts.

package api

// ======================================================
//                  COMPONENT: GENERIC PROXY
// ======================================================

// ========== Layer: Request Handling Layer ==========

// RequestProcessor receives incoming client requests and delegates work.
type RequestProcessor interface {
    ProcessWatchRequest(key string, fromRevision int64) (<-chan Event, error)
    ProcessListRequest(prefix string) ([]KV, error)
}

// RequestMerger merges concurrent watch requests on the same key.
type RequestMerger interface {
    MergeWatch(key string, fromRev int64, ch <-chan Event) <-chan Event
}

// RequestRouter decides whether to serve request from cache or etcd.
type RequestRouter interface {
    RouteWatch(key string, fromRev int64) (<-chan Event, error)
    RouteList(prefix string) ([]KV, error)
}

// ========== Layer: Cache Layer ==========

// CacheBackfiller fills cache on miss from etcd.
type CacheBackfiller interface {
    Backfill(key string) (KV, error)
    BackfillRange(prefix string) ([]KV, error)
}

// BTreeStore provides indexable range storage.
type BTreeStore interface {
    Get(key string) (KV, bool)
    List(prefix string) ([]KV, error)
    Insert(key string, value []byte, revision int64)
    Delete(key string)
}

// RevisionTracker tracks current revision of the cache.
type RevisionTracker interface {
    CurrentRevision() int64
    IsStale(key string, newRev int64) bool
}

// VersionValidator checks for outdated/stale updates.
type VersionValidator interface {
    IsValidUpdate(key string, newRev int64) bool
}

// ========== Layer: Observability & Metrics Layer ==========

// MetricsCollector observes internal cache components.
type MetricsCollector interface {
    ObserveStoreSize(size int)
    ObserveRequestRate(endpoint string, count int)
}

// MetricsExporter exposes metrics to Prometheus or others.
type MetricsExporter interface {
    Export() error
}

// ======================================================
//                 COMPONENT: CLIENT LIBRARY
// ======================================================

// ClientSession manages the lifecycle of one connected client.
type ClientSession interface {
    // --- Lifecycle management ---
    Start() error     // Starts session: init resources, register, etc.
    Stop() error      // Terminates session, releases any goroutine or channel
    ID() string       // Returns session ID for tracking

    // --- View and watch capabilities ---
    CacheView() SnapshotView
    // WatchSingle subscribes to changes on a single key
    Watch(key string, fromRev int64) (<-chan Event, error)
    // WatchPrefix subscribes to changes on a key prefix
    WatchPrefix(prefix string, fromRev int64) (<-chan Event, error)
}

// ClientLibrary provides an interface for SDK-level usage.
type ClientLibrary interface {
    NewSession(clientID string) (ClientSession, error)
    BroadcastUpdate(ev Event)
}

// ======================================================
//                 COMPONENT: SPECIFIC ADAPTER
// ======================================================

// K8sAdapter wraps the proxy for Kubernetes API semantics.
type K8sAdapter interface {
    ServeList(kind string, namespace string, opts map[string]string) ([]byte, error)
    ServeWatch(kind string, namespace string, fromResourceVersion string) (<-chan []byte, error)
    Encode(obj interface{}) ([]byte, error)
    Decode(data []byte) (interface{}, error)
}

// EtcdAdapter bridges raw etcd events to internal event format.
type EtcdAdapter interface {
    TranslateEtcdEvent(kv EtcdKV) (Event, error)
    IsWatchableKey(key string) bool
}

// EtcdKV is a simplified placeholder for etcd's key-value event.
type EtcdKV struct {
    Key   string
    Value []byte
    ModRevision int64
}

// ======================================================
//                       Shared Types
// ======================================================
// pkg/api/cacheview.go
type SnapshotView interface {
    Get(key string) (KV, bool)
    List(prefix string) ([]KV, error)
    Page(page, size int) ([]KV, error)
    Revision() int64
}

type KV struct {
    Key      string
    Value    []byte
    Revision int64
}

// Event struct
// EventType represents the type of operation that occurred on a key.
type EventType int

const (
    // These values must match mvccpb.Event_EventType for safe type casting.
    EventPut EventType = iota   // 类似 Java 中的 enum，用来表示是一次 PUT 操作
    EventDelete                 // 表示是一次 DELETE 操作
)

// Event represents a single operation that occurred in the system.
type Event struct {
    Type      EventType // Type of operation: PUT or DELETE
    Key       string    // The key that was operated on
    Value     []byte    // The new value (nil if DELETE)
    Revision int64     // Monotonic revision assigned by the watch cache, used for local event ordering
    ModRev    int64     // etcd's original ModRevision for this key
}