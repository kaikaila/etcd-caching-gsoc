package cache
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