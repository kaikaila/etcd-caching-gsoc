package event_test

import (
	"testing"

	"github.com/kaikaila/etcd-caching-gsoc/cache/event"
	"github.com/stretchr/testify/assert"
)

func TestEventLogInterface(t *testing.T) {
    log := event.NewMemoryEventLog(5)  // change this to apply the test to other eventlog implementations

    // Append events
    err := log.Append(event.Event{Key: "foo", Value: []byte("v1"), GlobalRev: 100})
    assert.NoError(t, err)
    err = log.Append(event.Event{Key: "bar", Value: []byte("v2"), GlobalRev: 101})
    assert.NoError(t, err)

    // Test ListSince
    events, err := log.ListSince(100)
    assert.NoError(t, err)
    assert.Len(t, events, 2)

    // Test LatestRevision
    rev := log.LatestRevision()
    assert.Equal(t, int64(101), rev)

    // Test Compact
    removed := log.Compact(100)
    assert.Equal(t, 1, removed)

    // After compact, only one event should remain
    events, _ = log.ListSince(0)
    assert.Len(t, events, 1)
    assert.Equal(t, "bar", events[0].Key)

    // Future edge cases to consider:
    // - Compact with threshold == 0: should retain all events
    // - Compact with threshold >> all event revisions: should remove all events
    // - Compact with threshold == highest GlobalRev: should remove all but latest
    // - Compact when multiple events have the same GlobalRev:
    //     Although GlobalRev is expected to be unique and monotonically increasing,
    //     if duplicates occur (e.g. from replayed events or WAL bugs), all matching entries should be evicted.
}