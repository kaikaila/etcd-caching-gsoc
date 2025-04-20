package main

import (
	"fmt"

	"github.com/kaikaila/etcd-caching-gsoc/cache"
	"github.com/kaikaila/etcd-caching-gsoc/cache/event"
)

func testWatcherWithLogOutput() {
	fmt.Println("== Logging Demo ==")

	// Create a memory-backed event log with capacity 5
	log := event.NewMemoryEventLog(5)
	cache := cache.NewWatchCacheWithLog(nil, log)

	// Add a few events
	cache.AddEvent(event.Event{
		Key:       "foo",
		Value:     []byte("bar"),
		KeyRev:    1,
		GlobalRev: 100,
		Type:      event.EventPut,
		ModRev:    100,
	})
	cache.AddEvent(event.Event{
		Key:       "baz",
		Value:     []byte("qux"),
		KeyRev:    1,
		GlobalRev: 101,
		Type:      event.EventPut,
		ModRev:    101,
	})

	// Replay and print
	fmt.Println("== Replayed Events ==")
	events, err := log.Replay(0)
	if err != nil {
		fmt.Println("Replay error:", err)
		return
	}
	for _, ev := range events {
		fmt.Printf("Key=%s Type=%v GlobalRev=%d Value=%s\n", ev.Key, ev.Type, ev.GlobalRev, string(ev.Value))
	}
}
