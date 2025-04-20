package event

// EventType represents the type of operation that occurred on a key.
type EventType int

const (
    // These values must match mvccpb.Event_EventType for safe type casting.
    EventPut EventType = iota   // 类似 Java 中的 enum，用来表示是一次 PUT 操作
    EventDelete                 // 表示是一次 DELETE 操作
)

// Event represents a single operation that occurred in the system.
type Event struct {
    Type        EventType // 操作类型，例如 PUT 或 DELETE
    Key         string    // 被操作的 key
    Value       []byte    // 当前值（如果是 DELETE，则可以为 nil）
    KeyRev int64     // 该 key 的局部版本号（用于快照判断是否为最新）
    GlobalRev   int64     // 全局排序用的 revision（用于 replay 顺序）
}

