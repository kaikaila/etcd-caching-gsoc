package cache

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
)

// Tests for WatchCache.Snapshot()
// Test Points:
// 1. Content correctness: modifying the original cache after Snapshot() should not change the snapshot copy.
// 2. Concurrency safety: concurrently writing to the cache and taking snapshots should not cause panics or data races.

// TestSnapshotIsolation verifies that Snapshot returns an immutable copy of the cache state.
func TestSnapshotIsolation(t *testing.T) {
	wc := NewWatchCache(nil)
	wc.HandlePut("foo", "bar", 1)
	snap1 := wc.Snapshot()

	// Modify the cache after taking the snapshot
	wc.HandlePut("foo", "baz", 2)

	// Verify the snapshot remains unchanged
	obj, ok := snap1["foo"]
	if !ok {
		t.Fatalf("snapshot missing key 'foo'")
	}
	if string(obj.Value) != "bar" {
		t.Fatalf("expected 'bar' in snapshot, got %s", obj.Value)
	}
}

// TestSnapshotConcurrentSafety verifies that taking snapshots concurrently with writes
// does not cause panics or data races.
func TestSnapshotConcurrentSafety(t *testing.T) {
	wc := NewWatchCache(nil)
	var wg sync.WaitGroup
	numWriters := 10
	numSnapshots := 10

	// Start multiple writers that continuously put values
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			for rev := int64(1); rev <= 100; rev++ {
				wc.HandlePut(key, fmt.Sprintf("val-%d", rev), rev)
			}
		}(i)
	}

	// Start multiple snapshotters that repeatedly take snapshots
	for i := 0; i < numSnapshots; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = wc.Snapshot()
			}
		}()
	}

	wg.Wait()
}