package eventlog

// EventLog defines the interface for an append-only event history.
// Different implementations may include: in-memory ring buffer, WAL file, etcd historical API, etc.
// For now, only the in-memory version is implemented; WAL and etcd-backed versions may be added later.

import "context"
type EventLog interface {
    Append(ev Event) error                       // Appends an event to the log
    ListSince(fromRev int64) ([]Event, error)    // Returns all events with Revision > fromRev
    Compact(rev int64) int 
    LatestRevision() int64                       // Returns the current max Revision in the log
    // Watch returns a channel streaming events with Revision > sinceRev.
    Watch(ctx context.Context, sinceRev int64) (<-chan Event, error)
}
