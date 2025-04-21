package eventlog

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
    GlobalRev int64     // Monotonic revision assigned by the watch cache, used for local event ordering
    ModRev    int64     // etcd's original ModRevision for this key
}
