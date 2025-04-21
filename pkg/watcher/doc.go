/*
Package watcher provides etcd watch wrappers that produce structured change events.

It encapsulates etcd's low-level watch interface and emits high-level Event objects
that can be passed into the WatchCache or other consumers.

Core components:

- WatchKey: a wrapper for watching a specific etcd key, supporting callback or channel-based consumption.
- EventTransformer: a planned utility to normalize etcd responses into unified event models.

This package abstracts away the low-level stream handling, allowing other modules to consume
semantic events without dealing with raw etcd WatchResponses.
*/
package watcher
