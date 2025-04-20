package cache

import (
	"fmt"
	"sync"
	"testing"

	"github.com/kaikaila/etcd-caching-gsoc/cache/event"
	"github.com/stretchr/testify/assert"
)

type dummySink struct {
	puts    map[string]string
	deletes []string
}

func newDummySink() *dummySink {
	return &dummySink{
		puts:    make(map[string]string),
		deletes: []string{},
	}
}

func (d *dummySink) HandlePut(key, val string) {
	d.puts[key] = val
}

func (d *dummySink) HandleDelete(key string) {
	d.deletes = append(d.deletes, key)
}

func TestWatchCache_HandlePutAndGet(t *testing.T) {
	sink := newDummySink()
	cache := NewWatchCache(sink)

	cache.HandlePut("foo", "bar", 1)
	val, ok := cache.Get("foo")

	assert.True(t, ok)
	assert.Equal(t, "bar", val)
	assert.Equal(t, int64(1), cache.Revision())
	assert.Equal(t, "bar", sink.puts["foo"])
}

func TestWatchCache_HandleDelete(t *testing.T) {
	sink := newDummySink()
	cache := NewWatchCache(sink)

	cache.HandlePut("foo", "bar", 1)
	cache.HandleDelete("foo", 2)

	_, ok := cache.Get("foo")
	assert.False(t, ok)
	assert.Equal(t, int64(2), cache.Revision())
	assert.Contains(t, sink.deletes, "foo")
}

func TestWatchCache_IgnoreStaleRevision(t *testing.T) {
	sink := newDummySink()
	cache := NewWatchCache(sink)

	cache.HandlePut("foo", "bar", 5)
	cache.HandlePut("foo", "baz", 3)

	val, ok := cache.Get("foo")
	assert.True(t, ok)
	assert.Equal(t, "bar", val)
	assert.Equal(t, int64(5), cache.Revision())
}

func TestWatchCache_ConcurrentAccess(t *testing.T) {
	cache := NewWatchCache(nil)
	var wg sync.WaitGroup

	// write goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			cache.HandlePut(fmt.Sprintf("k%d", i), "val", int64(i+1))
		}
	}()

	// read goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_, _ = cache.Get(fmt.Sprintf("k%d", i))
		}
	}()

	wg.Wait()
}

func TestWatchCache_AddEvent_Put(t *testing.T) {
	sink := newDummySink()
	log := event.NewMemoryEventLog(10)
	cache := NewWatchCacheWithLog(sink, log)

	ev := event.Event{
		Key:        "alpha",
		Value:      []byte("value1"),
		KeyRev:     1,
		GlobalRev:  100,
		Type:       event.EventPut,
		ModRev:     100,
	}

	err := cache.AddEvent(ev)
	assert.NoError(t, err)

	val, ok := cache.Get("alpha")
	assert.True(t, ok)
	assert.Equal(t, "value1", string(val.Value))

	// Check event log recorded it
	events, err := log.Replay(0)
	if err != nil {
		fmt.Println("Replay error:", err)
		return
	}
	assert.Len(t, events, 1)
	assert.Equal(t, "alpha", events[0].Key)
	assert.Equal(t, event.EventPut, events[0].Type)
}
