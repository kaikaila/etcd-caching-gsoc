package proxy

import "github.com/kaikaila/etcd-caching-gsoc/pkg/api"
type WatchCacheInterface interface {
	// this interface cannot be defined by building upon CacheWithSink, because the methods have different signatures(i.e. rev)
	HandlePut(key, val string, rev int64)  
	HandleDelete(key string, rev int64)
	Get(key string) (string, bool)
	Revision() int64
	Snapshot() api.SnapshotView
}

// CacheWithSink represents a cache implementation that can also handle etcd watch events.
// Implementations must satisfy both the Cache and EventSink interfaces.
type CacheWithSink interface {
	Cache
	EventSink
}
type Cache interface {
	Get(key string)(string, bool)
	Set(key string, value string)
}

// EventSink defines the interface for receiving etcd events.
// This interface is currently implemented by memoryCache to apply Watch updates.
// TODO: If more components implement EventSink in the future, consider moving this interface to a dedicated eventsink.go file.
type EventSink interface {
	HandlePut(key, val string)
	HandleDelete(key string)
}