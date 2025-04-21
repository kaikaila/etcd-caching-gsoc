package main

import (
	"fmt"

	"github.com/kaikaila/etcd-caching-gsoc/pkg/eventlog"
	"github.com/kaikaila/etcd-caching-gsoc/pkg/proxy"
)

func testWatcherWithLogOutput() {
	fmt.Println("== Logging Demo ==")

	// Create a memory-backed event log with capacity 5
	log := eventlog.NewMemoryEventLog(5)
	cache := proxy.NewWatchCacheWithLog(nil, log)

	// Add a few events
	cache.AddEvent(eventlog.Event{
		Key:       "foo",
		Value:     []byte("bar"),
		GlobalRev: 100,
		Type:      eventlog.EventPut,
		ModRev:    100,
	})
	cache.AddEvent(eventlog.Event{
		Key:       "baz",
		Value:     []byte("qux"),
		GlobalRev: 101,
		Type:      eventlog.EventPut,
		ModRev:    101,
	})

	// Replay and print
	fmt.Println("== Replayed Events ==")
	events, err := log.ListSince(0)
	if err != nil {
		fmt.Println("ListSince (Replay) error:", err)
		return
	}
	for _, ev := range events {
		fmt.Printf("Key=%s Type=%v GlobalRev=%d Value=%s\n", ev.Key, ev.Type, ev.GlobalRev, string(ev.Value))
	}
}
