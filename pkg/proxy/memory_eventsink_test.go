package proxy

import "testing"

func TestMemoryCache_EventSink(t *testing.T) {
	s := NewMemoryCache()
	
	c := s.(CacheWithSink)
	c.HandlePut("foo", "bar")
	val, ok := c.Get("foo")
	if !ok || val != "bar" {
		t.Errorf("expected key 'foo' to have value 'bar', got val=%q, ok=%v", val, ok)
	}

	c.HandleDelete("foo")
	val, ok = c.Get("foo")
	if !ok || val != "" {
		t.Errorf("expected key 'foo' to be cleared, got val=%q, ok=%v", val, ok)
	}
}
