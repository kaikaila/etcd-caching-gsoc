package proxy

import "testing"

// RunBasicCacheTests runs a suite of unit tests against any implementation of the Cache interface.
// It verifies expected behavior for basic operations like Set, Get, key overwrite, and handling of missing keys.
// This test can be reused across multiple Cache implementations to ensure compliance with the interface contract.
func RunBasicCacheTests(t *testing.T, c Cache) {
    
    // Basic Set and Get test 
    c.Set("foo", "bar")

    val, ok := c.Get("foo")
    if !ok {
        t.Fatal("Expected key 'foo' to be found")
    }
    if val != "bar" {
        t.Fatalf("Expected value 'bar', got '%s'", val)
    }

    // Overwriting existing key
    c.Set("foo", "baz")
    val, ok = c.Get("foo")
    if !ok || val != "baz" {
        t.Fatalf("Expected overwritten value 'baz', got '%s'", val)
    }

    // Getting non-existent key
    val, ok = c.Get("notfound")
    if ok {
        t.Fatalf("Expected key 'notfound' to be missing, got '%s'", val)
    }

    // TODO: use multiple goroutines to concurrently Set and Get, and verify no data race or inconsistency occurs
}

func TestMemoryCache_GetSet(t *testing.T) {
    RunBasicCacheTests(t, NewMemoryCache())
}