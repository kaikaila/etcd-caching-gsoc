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

// TestNewSnapshotViewPaging verifies that NewSnapshotView correctly orders snapshot data
// and that paging returns the expected subsets.
func TestNewSnapshotViewPaging(t *testing.T) {
    // Initialize cache and insert entries with out-of-order revisions
    wc := NewWatchCache(nil)
    wc.HandlePut("a", "valA", 1)
    wc.HandlePut("b", "valB", 3)
    wc.HandlePut("c", "valC", 2)

    // Build a snapshot view
    sv := wc.NewSnapshotView()

    // Verify Data length and sorting by GlobalRev ascending
    if len(sv.Data) != 3 {
        t.Fatalf("expected 3 items in SnapshotView, got %d", len(sv.Data))
    }
    expectedGlobals := []int64{1, 2, 3}
    for i, obj := range sv.Data {
        if obj.GlobalRev != expectedGlobals[i] {
            t.Fatalf("at index %d expected GlobalRev %d, got %d", i, expectedGlobals[i], obj.GlobalRev)
        }
    }

    // Test paging: page size 2
    page1 := sv.Page(1, 2)
    if len(page1) != 2 {
        t.Fatalf("expected 2 items on page 1, got %d", len(page1))
    }
    if page1[0].GlobalRev != 1 || page1[1].GlobalRev != 2 {
        t.Fatalf("page1 GlobalRevs mismatch: got [%d, %d]", page1[0].GlobalRev, page1[1].GlobalRev)
    }

    page2 := sv.Page(2, 2)
    if len(page2) != 1 {
        t.Fatalf("expected 1 item on page 2, got %d", len(page2))
    }
    if page2[0].GlobalRev != 3 {
        t.Fatalf("page2 GlobalRev mismatch: expected 3, got %d", page2[0].GlobalRev)
    }

    // Out-of-range page should return empty slice
    page3 := sv.Page(3, 2)
    if len(page3) != 0 {
        t.Fatalf("expected 0 items on page 3, got %d", len(page3))
    }
}
