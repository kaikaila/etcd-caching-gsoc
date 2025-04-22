package eventlog

// MemoryEventLog is the default in-memory implementation of EventLog.
// It uses a slice as a ring buffer to store recent events.
type MemoryEventLog struct {
    events      []Event
    capacity    int
    startIndex  int
    count       int
    latestRev   int64
}

// NewMemoryEventLog initializes a new MemoryEventLog with a fixed capacity.
func NewMemoryEventLog(capacity int) *MemoryEventLog {
    return &MemoryEventLog{
        events:   make([]Event, capacity),
        capacity: capacity,
    }
}

// Append adds a new event to the log, maintaining a fixed-size ring buffer.
func (l *MemoryEventLog) Append(ev Event) error {
    l.latestRev = ev.Revision
    pos := (l.startIndex + l.count) % l.capacity
    l.events[pos] = ev
    if l.count < l.capacity {
        l.count++
    } else {
        l.startIndex = (l.startIndex + 1) % l.capacity
    }
    return nil
}

// ListSince returns all events with Revision >= fromRev.
func (l *MemoryEventLog) ListSince(fromRev int64) ([]Event, error) {
    result := []Event{}
    for i := 0; i < l.count; i++ {
        idx := (l.startIndex + i) % l.capacity
        ev := l.events[idx]
        if ev.Revision >= fromRev {
            result = append(result, ev)
        }
    }
    return result, nil
}

// LatestRevision returns the highest Revision seen so far.
func (l *MemoryEventLog) LatestRevision() int64 {
    return l.latestRev
}

// Compact removes all events with Revision <= rev and returns the count of removed events.
func (l *MemoryEventLog) Compact(rev int64) int {
    removed := 0
    for l.count > 0 {
        ev := l.events[l.startIndex]
        if ev.Revision > rev {
            break
        }
        l.startIndex = (l.startIndex + 1) % l.capacity
        l.count--
        removed++
    }
    return removed
}